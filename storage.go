package main

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type mistakeEntry struct {
	Word        string    `json:"word"`
	Language    string    `json:"language,omitempty"`
	Correction  string    `json:"correction"`
	Explanation string    `json:"explanation"`
	AddedAt     time.Time `json:"added_at"`
}

type learnedWordEntry struct {
	ID                   string    `json:"id"`
	Language             string    `json:"language,omitempty"`
	Russian              string    `json:"russian"`
	English              string    `json:"english"`
	Context              string    `json:"context,omitempty"`
	ReviewCorrectCount   int       `json:"review_correct_count,omitempty"`
	SpellingCorrectCount int       `json:"spelling_correct_count,omitempty"`
	LearnedAt            time.Time `json:"learned_at"`
}

type phrasebookEntry struct {
	ID          string    `json:"id"`
	Phrase      string    `json:"phrase"`
	Translation string    `json:"translation,omitempty"`
	Note        string    `json:"note,omitempty"`
	Source      string    `json:"source,omitempty"`
	Language    string    `json:"language,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type habitDayEntry struct {
	Date      string    `json:"date,omitempty"`
	Login     bool      `json:"login,omitempty"`
	Complete  bool      `json:"complete,omitempty"`
	Claimed   bool      `json:"claimed,omitempty"`
	ClaimedAt time.Time `json:"claimed_at,omitempty"`
	Lessons   int       `json:"lessons,omitempty"`
	Practice  int       `json:"practice,omitempty"`
	Voice     int       `json:"voice,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type navigationLayout struct {
	FunctionRibbon []string `json:"function_ribbon,omitempty"`
	MobilePinned   []string `json:"mobile_pinned,omitempty"`
	MobileMore     []string `json:"mobile_more,omitempty"`
	MobileRail     []string `json:"mobile_rail,omitempty"`
}

type referralInviteeEntry struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Level         string    `json:"level"`
	XPLevel       int       `json:"xp_level"`
	XP            int       `json:"xp"`
	ReachedLevel3 bool      `json:"reached_level_3"`
	LevelRewarded bool      `json:"level_rewarded"`
	Earned        string    `json:"earned"`
	EarnedUSDT    string    `json:"earned_usdt"`
	JoinedAt      time.Time `json:"joined_at"`
}

type userState struct {
	TelegramID              int64                    `json:"telegram_id"`
	FirstName               string                   `json:"first_name"`
	Plan                    string                   `json:"plan"`
	PremiumUntil            time.Time                `json:"premium_until"`
	InvitedBy               int64                    `json:"invited_by"`
	ReferralCode            string                   `json:"referral_code,omitempty"`
	ReferralCount           int                      `json:"referral_count"`
	ReferralBalanceKopecks  int64                    `json:"referral_balance_kopecks"`
	ReferralLevelRewarded   bool                     `json:"referral_level_rewarded,omitempty"`
	ReferralRewardChargeIDs []string                 `json:"referral_reward_charge_ids,omitempty"`
	Mode                    string                   `json:"mode"`
	InterfaceLanguage       string                   `json:"interface_language,omitempty"`
	InterfaceSelected       bool                     `json:"interface_selected,omitempty"`
	LearningLanguage        string                   `json:"learning_language,omitempty"`
	LanguageSelected        bool                     `json:"language_selected,omitempty"`
	Level                   string                   `json:"level"`
	LearningFocus           string                   `json:"learning_focus,omitempty"`
	LessonCount             int                      `json:"lesson_count"`
	PracticeCount           int                      `json:"practice_count"`
	VoiceCount              int                      `json:"voice_count"`
	WordLessonCount         int                      `json:"word_lesson_count"`
	WordGameCount           int                      `json:"word_game_count"`
	XP                      int                      `json:"xp"`
	DailyDate               string                   `json:"daily_date"`
	LessonsToday            int                      `json:"lessons_today"`
	PracticeToday           int                      `json:"practice_today"`
	VoiceToday              int                      `json:"voice_today"`
	LessonHistory           []string                 `json:"lesson_history,omitempty"`
	PracticeHistory         []string                 `json:"practice_history,omitempty"`
	DailyBonusClaims        []string                 `json:"daily_bonus_claims,omitempty"`
	DailyBonusLastClaimedAt time.Time                `json:"daily_bonus_last_claimed_at,omitempty"`
	HabitLog                map[string]habitDayEntry `json:"habit_log,omitempty"`
	NavigationLayout        navigationLayout         `json:"navigation_layout,omitempty"`
	LastPaymentChargeID     string                   `json:"last_payment_charge_id"`
	LastLessonPrompt        string                   `json:"last_lesson_prompt"`
	ReminderEnabled         bool                     `json:"reminder_enabled"`
	ReminderUTCOffset       int                      `json:"reminder_utc_offset_minutes"`
	ReminderHour            int                      `json:"reminder_hour"`
	TimezoneSelected        bool                     `json:"timezone_selected"`
	LastReminderDate        string                   `json:"last_reminder_date"`
	Mistakes                []mistakeEntry           `json:"mistakes,omitempty"`
	LearnedWords            []learnedWordEntry       `json:"learned_words,omitempty"`
	Phrasebook              []phrasebookEntry        `json:"phrasebook,omitempty"`
	CreatedAt               time.Time                `json:"created_at"`
	UpdatedAt               time.Time                `json:"updated_at"`
}

type jsonStore struct {
	path                 string
	mu                   sync.Mutex
	users                map[int64]userState
	tutorLessons         map[string]tutorLesson
	tutorUserLessons     map[int64]map[string]time.Time
	aiTutorLessons       map[string]aiTutorLessonRecord
	aiTutorSessions      map[string]aiTutorSessionRecord
	aiTutorAnswers       map[string]aiTutorAnswerRecord
	aiTutorReviews       map[string]aiTutorReviewRecord
	aiTutorQualityChecks map[string]aiTutorQualityCheckRecord
	aiTutorWordReports   map[string]aiTutorWordReportRecord
}

type leaderboardEntry struct {
	FirstName string
	Score     int
	XP        int
	Words     int
	Mistakes  int
	Level     int
	Title     string
	Languages []string
}

type languageLeaderboardEntry struct {
	FirstName string
	Score     int
	Words     int
	Mistakes  int
	Level     string
}

type reminderTarget struct {
	TelegramID        int64
	FirstName         string
	InterfaceLanguage string
	LocalDate         string
}

type tutorLessonFactory func(sequence int) (tutorLesson, error)

func totalLeaderboardScore(user userState) int {
	return leaderboardActivityScore(masteredWordCount(user), len(user.Mistakes))
}

func leaderboardActivityScore(words int, mistakes int) int {
	return words*15 + mistakes*3
}

type store interface {
	getOrCreateUser(telegramID int64, firstName string) (userState, error)
	applyReferral(newUserID int64, inviterID int64) (referralApplication, error)
	rewardReferralLevel(telegramID int64) (referralLevelReward, error)
	creditReferralPurchase(telegramID int64, chargeID string, paidKopecks int64) (referralPurchaseReward, error)
	setMode(telegramID int64, mode string) error
	setInterfaceLanguage(telegramID int64, language string) error
	setLearningLanguage(telegramID int64, language string) error
	setUserLevel(telegramID int64, level string) error
	setLearningFocus(telegramID int64, focus string) error
	setNavigationLayout(telegramID int64, layout navigationLayout) error
	addXP(telegramID int64, amount int) error
	recordAITutorCompletion(telegramID int64, amount int) error
	recordHabitLogin(telegramID int64) error
	claimDailyBonus(telegramID int64, date string, amount int) (bool, error)
	saveLesson(telegramID int64, prompt string) error
	incrementPractice(telegramID int64) error
	savePracticeHistory(telegramID int64, history []string) error
	incrementVoice(telegramID int64) error
	addLearnedWord(telegramID int64, word vocabWord) (bool, int, error)
	markWordGameCorrect(telegramID int64, wordID string) error
	markSpellingCorrect(telegramID int64, wordID string) error
	addPhrasebookEntry(telegramID int64, entry phrasebookEntry) ([]phrasebookEntry, error)
	removePhrasebookEntry(telegramID int64, id string) ([]phrasebookEntry, error)
	extendPremium(telegramID int64, chargeID string, duration time.Duration, tier string) (time.Time, error)
	addMistakes(telegramID int64, language string, entries []mistakeEntry) error
	removeMistake(telegramID int64, language string, word string, correction string) error
	clearMistakes(telegramID int64, language string) error
	nextTutorLesson(user userState, factory tutorLessonFactory) (tutorLesson, error)
	saveAITutorLesson(lesson aiTutorLessonRecord) error
	getAITutorLesson(lessonID string) (aiTutorLessonRecord, bool, error)
	findApprovedAITutorLesson(language string, interfaceLanguage string, levelBand string, telegramID int64) (aiTutorLessonRecord, bool, error)
	aiTutorSessionCountForContext(telegramID int64, language string, interfaceLanguage string, levelBand string) (int, error)
	createAITutorSession(session aiTutorSessionRecord) error
	getAITutorSession(sessionID string) (aiTutorSessionRecord, bool, error)
	updateAITutorSessionStage(sessionID string, stage string, status string, completedAt string) error
	completedAITutorLessons(telegramID int64, limit int) ([]aiTutorCompletedLessonRecord, error)
	userCanRestartAITutorLesson(telegramID int64, lessonID string) (bool, error)
	saveAITutorAnswer(answer aiTutorAnswerRecord) error
	saveAITutorQualityCheck(check aiTutorQualityCheckRecord) error
	createAITutorWordReport(report aiTutorWordReportRecord) (aiTutorWordReportRecord, bool, error)
	getAITutorWordReport(reportID string) (aiTutorWordReportRecord, bool, error)
	countAITutorWordReports(telegramID int64, sessionID string, since time.Time) (int, int, error)
	countAITutorWordReportsSince(since time.Time) (int, error)
	setAITutorWordReportAdminMessage(reportID string, chatID string, messageID int64) error
	setAITutorWordReportFixPrompt(reportID string, chatID string, messageID int64) error
	resolveAITutorWordReport(reportID string, status aiTutorWordReportStatus, finalWord string, finalTranslation string) (aiTutorWordReportRecord, bool, error)
	findPendingAITutorWordReportByFixPrompt(chatID string, messageID int64) (aiTutorWordReportRecord, bool, error)
	updateAITutorLessonQuality(lessonID string, status string, postScore int) error
	scheduleAITutorReview(telegramID int64, lessonID string, dueAt string, intervalCode string) error
	clearLegacyTutorLessonBase() error
	referralInvitees(telegramID int64, limit int) ([]referralInviteeEntry, error)
	leaderboard(limit int) []leaderboardEntry
	languageLeaderboard(language string, limit int) []languageLeaderboardEntry
	dueReminderUsers(now time.Time, limit int) ([]reminderTarget, error)
	markReminderSent(telegramID int64, localDate string) error
	setReminderEnabled(telegramID int64, enabled bool) error
	setTimezone(telegramID int64, offsetMinutes int) error
}

func newStore(cfg config) (store, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.StorageDriver)) {
	case "", "sqlite":
		return newSQLiteStore(cfg.DatabasePath, cfg.DataPath, cfg.AITutorDatabasePath)
	case "json":
		return newJSONStore(cfg.DataPath)
	default:
		return nil, errors.New("unknown STORAGE_DRIVER: " + cfg.StorageDriver)
	}
}

func newJSONStore(path string) (*jsonStore, error) {
	if err := requireNonEmpty("DATA_PATH", path); err != nil {
		return nil, err
	}

	store := &jsonStore{
		path:                 path,
		users:                map[int64]userState{},
		tutorLessons:         map[string]tutorLesson{},
		tutorUserLessons:     map[int64]map[string]time.Time{},
		aiTutorLessons:       map[string]aiTutorLessonRecord{},
		aiTutorSessions:      map[string]aiTutorSessionRecord{},
		aiTutorAnswers:       map[string]aiTutorAnswerRecord{},
		aiTutorReviews:       map[string]aiTutorReviewRecord{},
		aiTutorQualityChecks: map[string]aiTutorQualityCheckRecord{},
		aiTutorWordReports:   map[string]aiTutorWordReportRecord{},
	}

	bytes, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return store, nil
	}
	if err := json.Unmarshal(bytes, &store.users); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *jsonStore) getOrCreateUser(telegramID int64, firstName string) (userState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	user, ok := s.users[telegramID]
	if !ok {
		user = newUserState(telegramID, now)
	}
	normalizeUser(&user, now)
	user.FirstName = firstName
	user.UpdatedAt = now
	s.users[telegramID] = user

	return user, s.saveLocked()
}

func (s *jsonStore) applyReferral(newUserID int64, inviterID int64) (referralApplication, error) {
	if inviterID == 0 || inviterID == newUserID {
		return referralApplication{}, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	newUser, ok := s.users[newUserID]
	if !ok {
		newUser = newUserState(newUserID, now)
	}
	normalizeUser(&newUser, now)
	if newUser.InvitedBy != 0 {
		return referralApplication{}, nil
	}

	inviter, ok := s.users[inviterID]
	if !ok {
		return referralApplication{}, nil
	}
	normalizeUser(&inviter, now)

	inviteeDays := referralInviteePremiumDays

	newUser.InvitedBy = inviterID
	newUser.Plan = "premium"
	newUser.PremiumUntil = premiumBase(newUser, now).Add(referralPremiumDuration(inviteeDays)).UTC()
	newUser.UpdatedAt = now

	inviter.ReferralCount++
	inviter.UpdatedAt = now

	s.users[newUserID] = newUser
	s.users[inviterID] = inviter
	levelReward := s.rewardReferralLevelLocked(newUserID, now)
	return referralApplication{
		Applied:            true,
		InviterID:          inviterID,
		InviteeID:          newUserID,
		InviteePremiumDays: inviteeDays,
		InviterPremiumDays: levelReward.InviterPremiumDays,
		NewUserUntil:       newUser.PremiumUntil,
		InviteeUntil:       newUser.PremiumUntil,
		InviterUntil:       levelReward.InviterUntil,
		ReferralCount:      inviter.ReferralCount,
	}, s.saveLocked()
}

func (s *jsonStore) rewardReferralLevel(telegramID int64) (referralLevelReward, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	reward := s.rewardReferralLevelLocked(telegramID, time.Now().UTC())
	if !reward.Applied {
		return reward, nil
	}
	return reward, s.saveLocked()
}

func (s *jsonStore) rewardReferralLevelLocked(telegramID int64, now time.Time) referralLevelReward {
	invitee, ok := s.users[telegramID]
	if !ok {
		return referralLevelReward{}
	}
	normalizeUser(&invitee, now)
	if invitee.InvitedBy == 0 || invitee.ReferralLevelRewarded || !referralLevelReached(invitee.XP) {
		s.users[telegramID] = invitee
		return referralLevelReward{}
	}
	inviter, ok := s.users[invitee.InvitedBy]
	if !ok {
		s.users[telegramID] = invitee
		return referralLevelReward{}
	}
	normalizeUser(&inviter, now)
	inviter.Plan = paidTierForExtension(inviter.Plan, "premium", now)
	inviter.PremiumUntil = premiumBase(inviter, now).Add(referralPremiumDuration(referralInviterLevelPremiumDays)).UTC()
	inviter.UpdatedAt = now
	invitee.ReferralLevelRewarded = true
	invitee.UpdatedAt = now

	s.users[telegramID] = invitee
	s.users[invitee.InvitedBy] = inviter
	return referralLevelReward{
		Applied:            true,
		InviteeID:          telegramID,
		InviterID:          invitee.InvitedBy,
		InviterPremiumDays: referralInviterLevelPremiumDays,
		InviterUntil:       inviter.PremiumUntil,
	}
}

func (s *jsonStore) creditReferralPurchase(telegramID int64, chargeID string, paidKopecks int64) (referralPurchaseReward, error) {
	if chargeID == "" || paidKopecks <= 0 {
		return referralPurchaseReward{}, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	buyer, ok := s.users[telegramID]
	if !ok || buyer.InvitedBy == 0 {
		return referralPurchaseReward{}, nil
	}
	for _, existingChargeID := range buyer.ReferralRewardChargeIDs {
		if existingChargeID == chargeID {
			return referralPurchaseReward{PaymentID: chargeID, BuyerID: telegramID}, nil
		}
	}

	direct, ok := s.users[buyer.InvitedBy]
	if !ok {
		return referralPurchaseReward{}, nil
	}
	if direct.TelegramID == 0 {
		direct.TelegramID = buyer.InvitedBy
	}
	directReward := referralDirectRewardKopecks(paidKopecks)
	indirectReward := int64(0)
	indirectID := direct.InvitedBy
	if indirectID == telegramID || indirectID == direct.TelegramID {
		indirectID = 0
	}

	if directReward > 0 {
		direct.ReferralBalanceKopecks += directReward
	}
	if indirectID != 0 {
		if indirect, ok := s.users[indirectID]; ok {
			indirectReward = referralIndirectRewardKopecks(paidKopecks)
			if indirectReward > 0 {
				indirect.ReferralBalanceKopecks += indirectReward
				indirect.UpdatedAt = time.Now().UTC()
				s.users[indirectID] = indirect
			}
		} else {
			indirectID = 0
		}
	}

	now := time.Now().UTC()
	direct.UpdatedAt = now
	buyer.ReferralRewardChargeIDs = append(buyer.ReferralRewardChargeIDs, chargeID)
	buyer.UpdatedAt = now
	s.users[telegramID] = buyer
	s.users[direct.TelegramID] = direct

	return referralPurchaseReward{
		Applied:               true,
		BuyerID:               telegramID,
		PaymentID:             chargeID,
		PaidKopecks:           paidKopecks,
		DirectInviterID:       direct.TelegramID,
		IndirectInviterID:     indirectID,
		DirectRewardKopecks:   directReward,
		IndirectRewardKopecks: indirectReward,
	}, s.saveLocked()
}

func (s *jsonStore) setMode(telegramID int64, mode string) error {
	return s.update(telegramID, func(user *userState) {
		user.Mode = mode
	})
}

func (s *jsonStore) setInterfaceLanguage(telegramID int64, language string) error {
	return s.update(telegramID, func(user *userState) {
		user.InterfaceLanguage = normalizeInterfaceLanguage(language)
		user.InterfaceSelected = true
	})
}

func (s *jsonStore) setLearningLanguage(telegramID int64, language string) error {
	return s.update(telegramID, func(user *userState) {
		user.LearningLanguage = normalizeLearningLanguage(language)
		user.LanguageSelected = true
		user.Mode = "idle"
	})
}

func (s *jsonStore) setUserLevel(telegramID int64, level string) error {
	return s.update(telegramID, func(user *userState) {
		user.Level = normalizeCEFRLevel(level)
	})
}

func (s *jsonStore) setLearningFocus(telegramID int64, focus string) error {
	return s.update(telegramID, func(user *userState) {
		user.LearningFocus = strings.TrimSpace(focus)
	})
}

func (s *jsonStore) setNavigationLayout(telegramID int64, layout navigationLayout) error {
	return s.update(telegramID, func(user *userState) {
		user.NavigationLayout = normalizeNavigationLayout(layout)
	})
}

func (s *jsonStore) addXP(telegramID int64, amount int) error {
	if amount <= 0 {
		return nil
	}
	return s.update(telegramID, func(user *userState) {
		user.XP += amount
	})
}

func (s *jsonStore) recordAITutorCompletion(telegramID int64, amount int) error {
	if amount <= 0 {
		return nil
	}
	return s.update(telegramID, func(user *userState) {
		user.XP += amount
		user.LessonCount++
		user.LessonsToday++
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *jsonStore) recordHabitLogin(telegramID int64) error {
	return s.update(telegramID, func(user *userState) {
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *jsonStore) claimDailyBonus(telegramID int64, date string, amount int) (bool, error) {
	date = strings.TrimSpace(date)
	if date == "" || amount <= 0 {
		return false, nil
	}
	now := time.Now().UTC()
	claimed := false
	err := s.update(telegramID, func(user *userState) {
		lastClaim := user.DailyBonusLastClaimedAt
		if lastClaim.IsZero() && len(user.DailyBonusClaims) > 0 {
			lastDate := user.DailyBonusClaims[len(user.DailyBonusClaims)-1]
			if parsed, err := time.Parse("2006-01-02", lastDate); err == nil {
				lastClaim = parsed
			}
		}
		if !lastClaim.IsZero() && now.Sub(lastClaim) < 24*time.Hour {
			return
		}
		for _, existing := range user.DailyBonusClaims {
			if existing == date {
				return
			}
		}
		user.DailyBonusClaims = append(user.DailyBonusClaims, date)
		if len(user.DailyBonusClaims) > 62 {
			user.DailyBonusClaims = user.DailyBonusClaims[len(user.DailyBonusClaims)-62:]
		}
		user.XP += amount
		user.DailyBonusLastClaimedAt = now
		dayDate := localDateForUser(*user, now)
		if date != "" {
			dayDate = date
		}
		day := recordHabitDay(user, now, true)
		day.Date = dayDate
		day.Claimed = true
		day.ClaimedAt = now
		user.HabitLog[dayDate] = day
		claimed = true
	})
	return claimed, err
}

func (s *jsonStore) saveLesson(telegramID int64, prompt string) error {
	return s.update(telegramID, func(user *userState) {
		user.Mode = "lesson"
		user.LastLessonPrompt = prompt
		user.LessonHistory = trimLessonHistory(append(user.LessonHistory, prompt))
		user.LessonCount++
		user.LessonsToday++
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *jsonStore) incrementPractice(telegramID int64) error {
	return s.update(telegramID, func(user *userState) {
		user.PracticeCount++
		user.PracticeToday++
		user.XP += 5
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *jsonStore) savePracticeHistory(telegramID int64, history []string) error {
	return s.update(telegramID, func(user *userState) {
		user.PracticeHistory = trimPracticeHistory(history, practiceMemoryLimit)
	})
}

func (s *jsonStore) incrementVoice(telegramID int64) error {
	return s.update(telegramID, func(user *userState) {
		user.VoiceCount++
		user.VoiceToday++
		user.XP += 3
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *jsonStore) addLearnedWord(telegramID int64, word vocabWord) (bool, int, error) {
	learned := false
	totalWords := 0
	err := s.update(telegramID, func(user *userState) {
		for _, existing := range user.LearnedWords {
			if normalizeLearningLanguage(existing.Language) == normalizeLearningLanguage(word.Language) &&
				(existing.ID == word.ID || legacyVocabID(existing.ID) == legacyVocabID(word.ID)) {
				totalWords = masteredWordCountForLanguage(*user, word.Language)
				return
			}
		}
		user.LearnedWords = append(user.LearnedWords, learnedWordEntry{
			ID:        word.ID,
			Language:  word.Language,
			Russian:   word.Russian,
			English:   word.English,
			Context:   word.Context,
			LearnedAt: time.Now().UTC(),
		})
		user.WordLessonCount++
		user.XP += 15
		learned = true
		totalWords = masteredWordCountForLanguage(*user, word.Language)
	})
	return learned, totalWords, err
}

func (s *jsonStore) markWordGameCorrect(telegramID int64, wordID string) error {
	return s.update(telegramID, func(user *userState) {
		user.WordGameCount++
		user.XP += 10
		incrementLearnedWordReview(user, wordID)
	})
}

func (s *jsonStore) markSpellingCorrect(telegramID int64, wordID string) error {
	return s.update(telegramID, func(user *userState) {
		incrementLearnedWordSpelling(user, wordID)
	})
}

func (s *jsonStore) addPhrasebookEntry(telegramID int64, entry phrasebookEntry) ([]phrasebookEntry, error) {
	var result []phrasebookEntry
	err := s.update(telegramID, func(user *userState) {
		now := time.Now().UTC()
		user.Phrasebook = append([]phrasebookEntry{entry}, user.Phrasebook...)
		user.Phrasebook = normalizePhrasebookEntries(user.Phrasebook, user.LearningLanguage, now)
		result = append([]phrasebookEntry(nil), user.Phrasebook...)
	})
	return result, err
}

func (s *jsonStore) removePhrasebookEntry(telegramID int64, id string) ([]phrasebookEntry, error) {
	var result []phrasebookEntry
	id = strings.TrimSpace(id)
	err := s.update(telegramID, func(user *userState) {
		next := user.Phrasebook[:0]
		for _, item := range user.Phrasebook {
			if item.ID != id {
				next = append(next, item)
			}
		}
		user.Phrasebook = normalizePhrasebookEntries(next, user.LearningLanguage, time.Now().UTC())
		result = append([]phrasebookEntry(nil), user.Phrasebook...)
	})
	return result, err
}

func (s *jsonStore) referralInvitees(telegramID int64, limit int) ([]referralInviteeEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 100
	}
	var items []referralInviteeEntry
	for _, user := range s.users {
		if user.InvitedBy != telegramID {
			continue
		}
		items = append(items, referralInviteeFromUser(user))
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].JoinedAt.After(items[j].JoinedAt)
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (s *jsonStore) extendPremium(telegramID int64, chargeID string, duration time.Duration, tier string) (time.Time, error) {
	var until time.Time
	err := s.update(telegramID, func(user *userState) {
		if chargeID != "" && user.LastPaymentChargeID == chargeID {
			until = user.PremiumUntil
			return
		}
		base := time.Now().UTC()
		if user.isPremium(base) {
			base = user.PremiumUntil
		}
		until = base.Add(duration)
		user.Plan = paidTierForExtension(user.Plan, tier, base)
		user.PremiumUntil = until.UTC()
		user.LastPaymentChargeID = chargeID
	})
	return until, err
}

// addMistakes appends new mistakes to the user's dictionary, avoiding duplicates
// (same Word+Correction pair). Keeps at most 200 entries total (oldest dropped).
func (s *jsonStore) addMistakes(telegramID int64, language string, entries []mistakeEntry) error {
	if len(entries) == 0 {
		return nil
	}
	language = normalizeLearningLanguage(language)
	return s.update(telegramID, func(user *userState) {
		existing := map[string]bool{}
		for _, m := range user.Mistakes {
			existing[normalizeLearningLanguage(m.Language)+"|"+m.Word+"|"+m.Correction] = true
		}
		for _, e := range entries {
			e.Language = language
			key := e.Language + "|" + e.Word + "|" + e.Correction
			if !existing[key] {
				user.Mistakes = append(user.Mistakes, e)
				existing[key] = true
			}
		}
		// cap at 200 entries, keep the newest
		if len(user.Mistakes) > 200 {
			user.Mistakes = user.Mistakes[len(user.Mistakes)-200:]
		}
	})
}

func (s *jsonStore) clearMistakes(telegramID int64, language string) error {
	language = normalizeLearningLanguage(language)
	return s.update(telegramID, func(user *userState) {
		filtered := user.Mistakes[:0]
		for _, mistake := range user.Mistakes {
			if normalizeLearningLanguage(mistake.Language) != language {
				filtered = append(filtered, mistake)
			}
		}
		user.Mistakes = filtered
	})
}

func (s *jsonStore) removeMistake(telegramID int64, language string, word string, correction string) error {
	language = normalizeLearningLanguage(language)
	return s.update(telegramID, func(user *userState) {
		filtered := user.Mistakes[:0]
		removed := false
		for _, mistake := range user.Mistakes {
			if !removed && normalizeLearningLanguage(mistake.Language) == language && mistake.Word == word && mistake.Correction == correction {
				removed = true
				continue
			}
			filtered = append(filtered, mistake)
		}
		user.Mistakes = filtered
	})
}

func (s *jsonStore) nextTutorLesson(user userState, factory tutorLessonFactory) (tutorLesson, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureTutorLessonMapsLocked()
	language, interfaceLanguage, level := tutorLessonStorageContext(user)
	seen := s.tutorUserLessons[user.TelegramID]
	keys := make([]string, 0, len(s.tutorLessons))
	for id, lesson := range s.tutorLessons {
		_, alreadySeen := seen[id]
		if tutorLessonMatchesContext(lesson, language, interfaceLanguage, level) && !alreadySeen {
			keys = append(keys, id)
		}
	}
	sort.SliceStable(keys, func(i, j int) bool {
		left := s.tutorLessons[keys[i]]
		right := s.tutorLessons[keys[j]]
		if left.LessonNumber != right.LessonNumber {
			return left.LessonNumber < right.LessonNumber
		}
		return left.ID < right.ID
	})
	if len(keys) > 0 {
		lesson := s.tutorLessons[keys[0]]
		s.assignTutorLessonLocked(user.TelegramID, lesson.ID)
		return tutorPersonalizeLessonForUser(lesson, user), nil
	}

	startSequence := s.tutorAssignedLessonCountLocked(user.TelegramID, language, interfaceLanguage, level)
	for attempt := 0; attempt < tutorLessonGenerationAttemptLimit(); attempt++ {
		lesson, err := factory(startSequence + attempt)
		if err != nil {
			return tutorLesson{}, err
		}
		if lesson.ID == "" {
			continue
		}
		if !tutorLessonMatchesContext(lesson, language, interfaceLanguage, level) {
			continue
		}
		if _, alreadySeen := seen[lesson.ID]; alreadySeen {
			continue
		}
		s.tutorLessons[lesson.ID] = lesson
		s.assignTutorLessonLocked(user.TelegramID, lesson.ID)
		return tutorPersonalizeLessonForUser(lesson, user), nil
	}
	return tutorLesson{}, errors.New("no unique tutor lesson available")
}

func (s *jsonStore) ensureTutorLessonMapsLocked() {
	if s.tutorLessons == nil {
		s.tutorLessons = map[string]tutorLesson{}
	}
	if s.tutorUserLessons == nil {
		s.tutorUserLessons = map[int64]map[string]time.Time{}
	}
}

func (s *jsonStore) assignTutorLessonLocked(telegramID int64, lessonID string) {
	if lessonID == "" {
		return
	}
	if s.tutorUserLessons[telegramID] == nil {
		s.tutorUserLessons[telegramID] = map[string]time.Time{}
	}
	s.tutorUserLessons[telegramID][lessonID] = time.Now().UTC()
}

func (s *jsonStore) tutorAssignedLessonCountLocked(telegramID int64, language string, interfaceLanguage string, level string) int {
	seen := s.tutorUserLessons[telegramID]
	if len(seen) == 0 {
		return 0
	}
	count := 0
	for lessonID := range seen {
		if tutorLessonMatchesContext(s.tutorLessons[lessonID], language, interfaceLanguage, level) {
			count++
		}
	}
	return count
}

func (s *jsonStore) saveAITutorLesson(lesson aiTutorLessonRecord) error {
	if strings.TrimSpace(lesson.ID) == "" {
		return errors.New("ai tutor lesson id is empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	now := formatDBTime(time.Now().UTC())
	if lesson.CreatedAt == "" {
		lesson.CreatedAt = now
	}
	lesson.UpdatedAt = now
	lesson.LearningLanguage = normalizeLearningLanguage(lesson.LearningLanguage)
	lesson.InterfaceLanguage = normalizeInterfaceLanguage(lesson.InterfaceLanguage)
	lesson.ExactLevel = normalizeCEFRLevel(lesson.ExactLevel)
	if strings.TrimSpace(lesson.LevelBand) == "" {
		lesson.LevelBand = aiTutorLevelBand(lesson.ExactLevel)
	} else {
		lesson.LevelBand = aiTutorLevelBand(lesson.LevelBand)
	}
	if lesson.Fingerprint == "" {
		lesson.Fingerprint = aiTutorFingerprint(lesson.Payload)
	}
	s.aiTutorLessons[lesson.ID] = lesson
	return nil
}

func (s *jsonStore) findApprovedAITutorLesson(language string, interfaceLanguage string, levelBand string, telegramID int64) (aiTutorLessonRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	language = normalizeLearningLanguage(language)
	interfaceLanguage = normalizeInterfaceLanguage(interfaceLanguage)
	levelBand = aiTutorLevelBand(levelBand)

	usedLessonIDs := map[string]bool{}
	for _, session := range s.aiTutorSessions {
		if session.TelegramID == telegramID {
			usedLessonIDs[session.LessonID] = true
		}
	}
	ids := make([]string, 0, len(s.aiTutorLessons))
	for id, lesson := range s.aiTutorLessons {
		if usedLessonIDs[id] {
			continue
		}
		if lesson.Status == aiTutorStatusApproved &&
			normalizeLearningLanguage(lesson.LearningLanguage) == language &&
			normalizeInterfaceLanguage(lesson.InterfaceLanguage) == interfaceLanguage &&
			aiTutorLevelBand(lesson.LevelBand) == levelBand {
			ids = append(ids, id)
		}
	}
	sort.SliceStable(ids, func(i, j int) bool {
		left := s.aiTutorLessons[ids[i]]
		right := s.aiTutorLessons[ids[j]]
		if left.PostScore != right.PostScore {
			return left.PostScore > right.PostScore
		}
		if left.PreflightScore != right.PreflightScore {
			return left.PreflightScore > right.PreflightScore
		}
		return left.CreatedAt < right.CreatedAt
	})
	if len(ids) == 0 {
		return aiTutorLessonRecord{}, false, nil
	}
	return s.aiTutorLessons[ids[0]], true, nil
}

func (s *jsonStore) getAITutorLesson(lessonID string) (aiTutorLessonRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	lesson, ok := s.aiTutorLessons[strings.TrimSpace(lessonID)]
	return lesson, ok, nil
}

func (s *jsonStore) aiTutorSessionCountForContext(telegramID int64, language string, interfaceLanguage string, levelBand string) (int, error) {
	if telegramID == 0 {
		return 0, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	language = normalizeLearningLanguage(language)
	interfaceLanguage = normalizeInterfaceLanguage(interfaceLanguage)
	levelBand = aiTutorLevelBand(levelBand)
	count := 0
	for _, session := range s.aiTutorSessions {
		if session.TelegramID != telegramID {
			continue
		}
		lesson, ok := s.aiTutorLessons[strings.TrimSpace(session.LessonID)]
		if !ok {
			continue
		}
		if normalizeLearningLanguage(lesson.LearningLanguage) == language &&
			normalizeInterfaceLanguage(lesson.InterfaceLanguage) == interfaceLanguage &&
			aiTutorLevelBand(lesson.LevelBand) == levelBand {
			count++
		}
	}
	return count, nil
}

func (s *jsonStore) createAITutorSession(session aiTutorSessionRecord) error {
	if strings.TrimSpace(session.ID) == "" {
		return errors.New("ai tutor session id is empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	now := formatDBTime(time.Now().UTC())
	if session.StartedAt == "" {
		session.StartedAt = now
	}
	if session.UpdatedAt == "" {
		session.UpdatedAt = now
	}
	if session.Status == "" {
		session.Status = aiTutorSessionActive
	}
	s.aiTutorSessions[session.ID] = session
	return nil
}

func (s *jsonStore) getAITutorSession(sessionID string) (aiTutorSessionRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	session, ok := s.aiTutorSessions[strings.TrimSpace(sessionID)]
	return session, ok, nil
}

func (s *jsonStore) completedAITutorLessons(telegramID int64, limit int) ([]aiTutorCompletedLessonRecord, error) {
	if telegramID == 0 {
		return nil, nil
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	items := make([]aiTutorCompletedLessonRecord, 0)
	for _, session := range s.aiTutorSessions {
		if session.TelegramID != telegramID || !aiTutorSessionIsCompleted(session) {
			continue
		}
		lesson, ok := s.aiTutorLessons[strings.TrimSpace(session.LessonID)]
		if !ok {
			continue
		}
		items = append(items, aiTutorCompletedLessonFromRecords(session, lesson))
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CompletedAt > items[j].CompletedAt
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (s *jsonStore) userCanRestartAITutorLesson(telegramID int64, lessonID string) (bool, error) {
	lessonID = strings.TrimSpace(lessonID)
	if telegramID == 0 || lessonID == "" {
		return false, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	for _, session := range s.aiTutorSessions {
		if session.TelegramID == telegramID && strings.TrimSpace(session.LessonID) == lessonID {
			return true, nil
		}
	}
	return false, nil
}

func (s *jsonStore) updateAITutorSessionStage(sessionID string, stage string, status string, completedAt string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	session, ok := s.aiTutorSessions[strings.TrimSpace(sessionID)]
	if !ok {
		return errors.New("ai tutor session not found")
	}
	session.CurrentStage = strings.TrimSpace(stage)
	if strings.TrimSpace(status) != "" {
		session.Status = strings.TrimSpace(status)
	}
	session.CompletedAt = strings.TrimSpace(completedAt)
	session.UpdatedAt = formatDBTime(time.Now().UTC())
	s.aiTutorSessions[session.ID] = session
	return nil
}

func (s *jsonStore) saveAITutorAnswer(answer aiTutorAnswerRecord) error {
	if strings.TrimSpace(answer.SessionID) == "" || strings.TrimSpace(answer.Stage) == "" {
		return errors.New("ai tutor answer session and stage are required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	if answer.CreatedAt == "" {
		answer.CreatedAt = formatDBTime(time.Now().UTC())
	}
	s.aiTutorAnswers[answer.SessionID+"|"+answer.Stage] = answer
	return nil
}

func (s *jsonStore) saveAITutorQualityCheck(check aiTutorQualityCheckRecord) error {
	if strings.TrimSpace(check.ID) == "" {
		return errors.New("ai tutor quality check id is empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	if check.CreatedAt == "" {
		check.CreatedAt = formatDBTime(time.Now().UTC())
	}
	s.aiTutorQualityChecks[check.ID] = check
	return nil
}

func (s *jsonStore) createAITutorWordReport(report aiTutorWordReportRecord) (aiTutorWordReportRecord, bool, error) {
	if strings.TrimSpace(report.ID) == "" {
		return aiTutorWordReportRecord{}, false, errors.New("ai tutor word report id is empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	report.TelegramID = report.TelegramID
	report.SessionID = strings.TrimSpace(report.SessionID)
	report.LessonID = strings.TrimSpace(report.LessonID)
	report.Stage = strings.TrimSpace(report.Stage)
	if report.Status == "" {
		report.Status = aiTutorWordReportPending
	}
	for _, existing := range s.aiTutorWordReports {
		if existing.TelegramID == report.TelegramID &&
			existing.SessionID == report.SessionID &&
			existing.Stage == report.Stage &&
			existing.Status == aiTutorWordReportPending {
			return existing, true, nil
		}
	}
	now := formatDBTime(time.Now().UTC())
	if report.CreatedAt == "" {
		report.CreatedAt = now
	}
	report.UpdatedAt = now
	s.aiTutorWordReports[report.ID] = report
	return report, false, nil
}

func (s *jsonStore) getAITutorWordReport(reportID string) (aiTutorWordReportRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	report, ok := s.aiTutorWordReports[strings.TrimSpace(reportID)]
	return report, ok, nil
}

func (s *jsonStore) countAITutorWordReports(telegramID int64, sessionID string, since time.Time) (int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	sessionID = strings.TrimSpace(sessionID)
	userCount := 0
	sessionCount := 0
	for _, report := range s.aiTutorWordReports {
		if !since.IsZero() && parseDBTime(report.CreatedAt).Before(since.UTC()) {
			continue
		}
		if report.TelegramID == telegramID {
			userCount++
		}
		if sessionID != "" && report.SessionID == sessionID {
			sessionCount++
		}
	}
	return userCount, sessionCount, nil
}

func (s *jsonStore) countAITutorWordReportsSince(since time.Time) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	count := 0
	for _, report := range s.aiTutorWordReports {
		if since.IsZero() || !parseDBTime(report.CreatedAt).Before(since.UTC()) {
			count++
		}
	}
	return count, nil
}

func (s *jsonStore) setAITutorWordReportAdminMessage(reportID string, chatID string, messageID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	report, ok := s.aiTutorWordReports[strings.TrimSpace(reportID)]
	if !ok {
		return errors.New("ai tutor word report not found")
	}
	report.AdminChatID = strings.TrimSpace(chatID)
	report.AdminMessageID = messageID
	report.UpdatedAt = formatDBTime(time.Now().UTC())
	s.aiTutorWordReports[report.ID] = report
	return nil
}

func (s *jsonStore) setAITutorWordReportFixPrompt(reportID string, chatID string, messageID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	report, ok := s.aiTutorWordReports[strings.TrimSpace(reportID)]
	if !ok {
		return errors.New("ai tutor word report not found")
	}
	report.FixPromptChatID = strings.TrimSpace(chatID)
	report.FixPromptMessageID = messageID
	report.UpdatedAt = formatDBTime(time.Now().UTC())
	s.aiTutorWordReports[report.ID] = report
	return nil
}

func (s *jsonStore) resolveAITutorWordReport(reportID string, status aiTutorWordReportStatus, finalWord string, finalTranslation string) (aiTutorWordReportRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	report, ok := s.aiTutorWordReports[strings.TrimSpace(reportID)]
	if !ok {
		return aiTutorWordReportRecord{}, false, errors.New("ai tutor word report not found")
	}
	if report.Status != aiTutorWordReportPending {
		return report, false, nil
	}
	now := formatDBTime(time.Now().UTC())
	report.Status = status
	report.FinalWord = strings.TrimSpace(finalWord)
	report.FinalTranslation = strings.TrimSpace(finalTranslation)
	report.ResolvedAt = now
	report.UpdatedAt = now
	s.aiTutorWordReports[report.ID] = report
	return report, true, nil
}

func (s *jsonStore) findPendingAITutorWordReportByFixPrompt(chatID string, messageID int64) (aiTutorWordReportRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	chatID = strings.TrimSpace(chatID)
	if chatID == "" || messageID == 0 {
		return aiTutorWordReportRecord{}, false, nil
	}
	for _, report := range s.aiTutorWordReports {
		if report.Status == aiTutorWordReportPending && report.FixPromptChatID == chatID && report.FixPromptMessageID == messageID {
			return report, true, nil
		}
	}
	return aiTutorWordReportRecord{}, false, nil
}

func (s *jsonStore) updateAITutorLessonQuality(lessonID string, status string, postScore int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	lesson, ok := s.aiTutorLessons[strings.TrimSpace(lessonID)]
	if !ok {
		return errors.New("ai tutor lesson not found")
	}
	lesson.Status = strings.TrimSpace(status)
	lesson.PostScore = postScore
	lesson.UpdatedAt = formatDBTime(time.Now().UTC())
	s.aiTutorLessons[lesson.ID] = lesson
	return nil
}

func (s *jsonStore) scheduleAITutorReview(telegramID int64, lessonID string, dueAt string, intervalCode string) error {
	if telegramID == 0 || strings.TrimSpace(lessonID) == "" || strings.TrimSpace(intervalCode) == "" || strings.TrimSpace(intervalCode) == "no_review" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureAITutorMapsLocked()
	now := formatDBTime(time.Now().UTC())
	id := itoa(int(telegramID)) + "|" + strings.TrimSpace(lessonID) + "|" + strings.TrimSpace(intervalCode)
	s.aiTutorReviews[id] = aiTutorReviewRecord{
		ID:           id,
		TelegramID:   telegramID,
		LessonID:     strings.TrimSpace(lessonID),
		DueAt:        strings.TrimSpace(dueAt),
		IntervalCode: strings.TrimSpace(intervalCode),
		Status:       "scheduled",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return nil
}

func (s *jsonStore) clearLegacyTutorLessonBase() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tutorLessons = map[string]tutorLesson{}
	s.tutorUserLessons = map[int64]map[string]time.Time{}
	return nil
}

func (s *jsonStore) ensureAITutorMapsLocked() {
	if s.aiTutorLessons == nil {
		s.aiTutorLessons = map[string]aiTutorLessonRecord{}
	}
	if s.aiTutorSessions == nil {
		s.aiTutorSessions = map[string]aiTutorSessionRecord{}
	}
	if s.aiTutorAnswers == nil {
		s.aiTutorAnswers = map[string]aiTutorAnswerRecord{}
	}
	if s.aiTutorReviews == nil {
		s.aiTutorReviews = map[string]aiTutorReviewRecord{}
	}
	if s.aiTutorQualityChecks == nil {
		s.aiTutorQualityChecks = map[string]aiTutorQualityCheckRecord{}
	}
	if s.aiTutorWordReports == nil {
		s.aiTutorWordReports = map[string]aiTutorWordReportRecord{}
	}
}

func (s *jsonStore) dueReminderUsers(now time.Time, limit int) ([]reminderTarget, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}
	var targets []reminderTarget
	for _, user := range s.users {
		normalizeUser(&user, now)
		if target, ok := reminderDue(user, now); ok {
			targets = append(targets, target)
			if len(targets) >= limit {
				break
			}
		}
	}
	return targets, nil
}

func (s *jsonStore) markReminderSent(telegramID int64, localDate string) error {
	return s.update(telegramID, func(user *userState) {
		user.LastReminderDate = localDate
	})
}

func (s *jsonStore) setReminderEnabled(telegramID int64, enabled bool) error {
	return s.update(telegramID, func(user *userState) {
		user.ReminderEnabled = enabled
	})
}

func (s *jsonStore) setTimezone(telegramID int64, offsetMinutes int) error {
	return s.update(telegramID, func(user *userState) {
		user.ReminderUTCOffset = offsetMinutes
		user.TimezoneSelected = true
	})
}

func (s *jsonStore) leaderboard(limit int) []leaderboardEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 10
	}
	entries := make([]leaderboardEntry, 0, len(s.users))
	for _, user := range s.users {
		normalizeUser(&user, time.Now().UTC())
		score := totalLeaderboardScore(user)
		level, title, _, _ := knowledgeLevel(score)
		name := user.FirstName
		if name == "" {
			name = "Ученик"
		}
		if score > 0 {
			entries = append(entries, leaderboardEntry{
				FirstName: name,
				Score:     score,
				XP:        user.XP,
				Words:     masteredWordCount(user),
				Mistakes:  len(user.Mistakes),
				Level:     level,
				Title:     title,
				Languages: leaderboardLanguages(user),
			})
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score == entries[j].Score {
			return entries[i].FirstName < entries[j].FirstName
		}
		return entries[i].Score > entries[j].Score
	})
	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries
}

func (s *jsonStore) languageLeaderboard(language string, limit int) []languageLeaderboardEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 10
	}
	language = normalizeLearningLanguage(language)
	entries := make([]languageLeaderboardEntry, 0, len(s.users))
	for _, user := range s.users {
		normalizeUser(&user, time.Now().UTC())
		entry := languageLeaderboardEntry{
			FirstName: user.FirstName,
			Words:     len(learnedWordsForLanguageCode(user, language)),
			Mistakes:  len(mistakesForLanguageCode(user, language)),
			Level:     user.Level,
		}
		if entry.FirstName == "" {
			entry.FirstName = "Ученик"
		}
		entry.Score = leaderboardActivityScore(entry.Words, entry.Mistakes)
		if entry.Score > 0 {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score == entries[j].Score {
			return entries[i].FirstName < entries[j].FirstName
		}
		return entries[i].Score > entries[j].Score
	})
	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries
}

func (s *jsonStore) update(telegramID int64, change func(*userState)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[telegramID]
	if !ok {
		user = newUserState(telegramID, time.Now().UTC())
	}
	normalizeUser(&user, time.Now().UTC())
	change(&user)
	user.UpdatedAt = time.Now().UTC()
	s.users[telegramID] = user
	s.rewardReferralLevelLocked(telegramID, time.Now().UTC())
	return s.saveLocked()
}

func (s *jsonStore) saveLocked() error {
	bytes, err := json.MarshalIndent(s.users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, bytes, 0600)
}

func newUserState(telegramID int64, now time.Time) userState {
	return userState{
		TelegramID:        telegramID,
		Plan:              "free",
		Mode:              "idle",
		InterfaceLanguage: "ru",
		InterfaceSelected: false,
		LearningLanguage:  "",
		LanguageSelected:  false,
		Level:             "A2",
		DailyDate:         todayUTC(now),
		ReminderEnabled:   true,
		ReminderUTCOffset: 180,
		ReminderHour:      19,
		CreatedAt:         now,
	}
}

func normalizeUser(user *userState, now time.Time) {
	if user.Plan == "" {
		user.Plan = "free"
	}
	if user.Mode == "" {
		user.Mode = "idle"
	}
	if user.InterfaceLanguage == "" {
		user.InterfaceLanguage = "ru"
	}
	if user.LearningLanguage == "" {
		user.LearningLanguage = "en"
	}
	for i := range user.LearnedWords {
		if user.LearnedWords[i].Language == "" {
			user.LearnedWords[i].Language = normalizeLearningLanguage(languageFromVocabID(user.LearnedWords[i].ID))
		}
		user.LearnedWords[i].Russian = cleanDictionaryDisplay(user.LearnedWords[i].Russian)
		user.LearnedWords[i].Context = normalizeVocabularyContext(user.LearnedWords[i].Context)
		if user.LearnedWords[i].ReviewCorrectCount < 0 {
			user.LearnedWords[i].ReviewCorrectCount = 0
		}
		if user.LearnedWords[i].SpellingCorrectCount < 0 {
			user.LearnedWords[i].SpellingCorrectCount = 0
		}
	}
	for i := range user.Mistakes {
		if user.Mistakes[i].Language == "" {
			user.Mistakes[i].Language = "en"
		}
	}
	user.Phrasebook = normalizePhrasebookEntries(user.Phrasebook, user.LearningLanguage, now)
	if user.Level == "" {
		user.Level = "A2"
	}
	if user.ReminderUTCOffset == 0 && !user.TimezoneSelected {
		user.ReminderUTCOffset = 180
	}
	if user.ReminderHour == 0 {
		user.ReminderHour = 19
	}
	if user.DailyDate != todayUTC(now) {
		user.DailyDate = todayUTC(now)
		user.LessonsToday = 0
		user.PracticeToday = 0
		user.VoiceToday = 0
	}
	user.HabitLog = normalizeHabitLog(user.HabitLog, now)
	user.NavigationLayout = normalizeNavigationLayout(user.NavigationLayout)
	if isPaidPlan(user.Plan) && (user.PremiumUntil.IsZero() || now.After(user.PremiumUntil)) {
		user.Plan = "free"
	}
}

func localDateForUser(user userState, now time.Time) string {
	offset := user.ReminderUTCOffset
	if offset == 0 && !user.TimezoneSelected {
		offset = 180
	}
	return now.UTC().Add(time.Duration(offset) * time.Minute).Format("2006-01-02")
}

func normalizeHabitLog(log map[string]habitDayEntry, now time.Time) map[string]habitDayEntry {
	if log == nil {
		return map[string]habitDayEntry{}
	}
	cutoff := now.UTC().AddDate(0, 0, -364).Format("2006-01-02")
	keys := make([]string, 0, len(log))
	for key := range log {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := map[string]habitDayEntry{}
	for _, key := range keys {
		if key < cutoff {
			continue
		}
		if _, err := time.Parse("2006-01-02", key); err != nil {
			continue
		}
		day := log[key]
		day.Date = key
		if day.UpdatedAt.IsZero() {
			day.UpdatedAt = now.UTC()
		}
		result[key] = day
	}
	if len(result) <= 365 {
		return result
	}
	keys = keys[:0]
	for key := range result {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for len(keys) > 365 {
		delete(result, keys[0])
		keys = keys[1:]
	}
	return result
}

func recordHabitDay(user *userState, now time.Time, login bool) habitDayEntry {
	if user.HabitLog == nil {
		user.HabitLog = map[string]habitDayEntry{}
	}
	date := localDateForUser(*user, now)
	day := user.HabitLog[date]
	day.Date = date
	day.Login = day.Login || login
	day.Lessons = user.LessonsToday
	day.Practice = user.PracticeToday
	day.Voice = user.VoiceToday
	day.Complete = day.Complete || (user.LessonsToday >= 1 && user.PracticeToday >= 1 && user.VoiceToday >= 1)
	for _, claimedDate := range user.DailyBonusClaims {
		if claimedDate == date {
			day.Claimed = true
			if day.ClaimedAt.IsZero() {
				day.ClaimedAt = user.DailyBonusLastClaimedAt
			}
			break
		}
	}
	day.UpdatedAt = now.UTC()
	user.HabitLog[date] = day
	user.HabitLog = normalizeHabitLog(user.HabitLog, now)
	return day
}

func referralInviteeFromUser(user userState) referralInviteeEntry {
	level, _, _, _ := knowledgeLevel(user.XP)
	joinedAt := user.CreatedAt
	if joinedAt.IsZero() {
		joinedAt = user.UpdatedAt
	}
	earned := "0 RUB"
	earnedUSDT := "0"
	if user.ReferralLevelRewarded {
		earned = "7 days Premium"
		earnedUSDT = "0"
	}
	return referralInviteeEntry{
		ID:            user.TelegramID,
		Name:          strings.TrimSpace(user.FirstName),
		Level:         user.Level,
		XPLevel:       level,
		XP:            user.XP,
		ReachedLevel3: referralLevelReached(user.XP),
		LevelRewarded: user.ReferralLevelRewarded,
		Earned:        earned,
		EarnedUSDT:    earnedUSDT,
		JoinedAt:      joinedAt,
	}
}

func normalizeNavigationLayout(layout navigationLayout) navigationLayout {
	return navigationLayout{
		FunctionRibbon: normalizeNavigationIDs(layout.FunctionRibbon, 32),
		MobilePinned:   normalizeNavigationIDs(layout.MobilePinned, 8),
		MobileMore:     normalizeNavigationIDs(layout.MobileMore, 32),
		MobileRail:     normalizeNavigationIDs(layout.MobileRail, 32),
	}
}

func normalizeNavigationIDs(values []string, limit int) []string {
	if limit <= 0 {
		limit = 32
	}
	allowed := map[string]bool{
		"home": true, "tutor": true, "lesson": true, "practice": true, "roleplay": true, "shadowing": true,
		"pronunciation": true, "words": true, "word-game": true, "spelling": true, "vocabulary": true,
		"phrasebook": true, "offline": true, "level": true, "progress": true, "awards": true,
		"leaderboard": true, "limits": true, "mistakes": true, "tools": true, "dashboard": true,
		"referral": true, "premium": true, "settings": true,
	}
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if !allowed[value] || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
		if len(result) >= limit {
			break
		}
	}
	return result
}

func normalizePhrasebookEntries(items []phrasebookEntry, fallbackLanguage string, now time.Time) []phrasebookEntry {
	seen := map[string]bool{}
	result := make([]phrasebookEntry, 0, len(items))
	for _, item := range items {
		item.Phrase = strings.TrimSpace(item.Phrase)
		if item.Phrase == "" {
			continue
		}
		key := strings.ToLower(item.Phrase)
		if seen[key] {
			continue
		}
		seen[key] = true
		item.ID = strings.TrimSpace(item.ID)
		if item.ID == "" {
			item.ID = "phrase-" + strconv.FormatInt(now.UnixNano()+int64(len(result)), 36)
		}
		item.Source = normalizePhrasebookSource(item.Source)
		if strings.TrimSpace(item.Language) == "" {
			item.Language = normalizeLearningLanguage(fallbackLanguage)
		}
		if item.CreatedAt.IsZero() {
			item.CreatedAt = now.UTC()
		}
		result = append(result, item)
		if len(result) >= 80 {
			break
		}
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result
}

func normalizePhrasebookSource(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "lesson", "practice", "roleplay", "mistake", "manual":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "manual"
	}
}

func incrementLearnedWordReview(user *userState, wordID string) {
	for i := range user.LearnedWords {
		if learnedWordIDMatches(user.LearnedWords[i], wordID) {
			user.LearnedWords[i].ReviewCorrectCount++
			return
		}
	}
}

func incrementLearnedWordSpelling(user *userState, wordID string) {
	for i := range user.LearnedWords {
		if learnedWordIDMatches(user.LearnedWords[i], wordID) {
			user.LearnedWords[i].SpellingCorrectCount++
			return
		}
	}
}

func learnedWordIDMatches(word learnedWordEntry, wordID string) bool {
	return word.ID == wordID || legacyVocabID(word.ID) == legacyVocabID(wordID)
}

func todayUTC(now time.Time) string {
	return now.UTC().Format("2006-01-02")
}

func (u userState) isPremium(now time.Time) bool {
	return isPaidPlan(u.Plan) && now.UTC().Before(u.PremiumUntil)
}

func (u userState) isPlatinum(now time.Time) bool {
	return u.Plan == "platinum" && now.UTC().Before(u.PremiumUntil)
}

func isPaidPlan(plan string) bool {
	return plan == "premium" || plan == "platinum"
}

func paidTierForExtension(current string, next string, base time.Time) string {
	next = strings.ToLower(strings.TrimSpace(next))
	if next != "platinum" {
		next = "premium"
	}
	if current == "platinum" && base.After(time.Now().UTC()) && next == "premium" {
		return "platinum"
	}
	return next
}
