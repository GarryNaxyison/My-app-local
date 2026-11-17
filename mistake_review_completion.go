package main

import (
	"log"
	"time"
)

const mistakeReviewXP = 8

func (s *sqliteStore) completeMistakeReview(userID int64, mode string, mistake mistakeEntry) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`UPDATE users SET xp=xp+?, mode='idle', updated_at=?
		WHERE telegram_id=? AND mode=? AND learning_language=?`,
		mistakeReviewXP, formatDBTime(time.Now().UTC()), userID, mode, normalizeLearningLanguage(mistake.Language))
	if err != nil {
		return false, err
	}
	changed, err := result.RowsAffected()
	if err != nil || changed == 0 {
		return false, err
	}
	result, err = tx.Exec(`DELETE FROM mistakes WHERE telegram_id=? AND language=? AND word=? AND correction=? AND added_at=?`,
		userID, normalizeLearningLanguage(mistake.Language), mistake.Word, mistake.Correction, formatDBTime(mistake.AddedAt))
	if err != nil {
		return false, err
	}
	changed, err = result.RowsAffected()
	if err != nil || changed == 0 {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	if _, err := s.rewardReferralLevel(userID); err != nil {
		log.Printf("mistake review referral reward pending for user %d: %v", userID, err)
	}
	return true, nil
}

func (s *jsonStore) completeMistakeReview(userID int64, mode string, mistake mistakeEntry) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	previous, exists := s.users[userID]
	if !exists || previous.Mode != mode || normalizeLearningLanguage(previous.LearningLanguage) != normalizeLearningLanguage(mistake.Language) {
		return false, nil
	}
	index := mistakeLearningReviewIndex(previous.Mistakes, mistakeLearningReviewID(mistake))
	if index < 0 {
		return false, nil
	}
	user := previous
	user.Mistakes = append(append([]mistakeEntry(nil), previous.Mistakes[:index]...), previous.Mistakes[index+1:]...)
	user.XP += mistakeReviewXP
	user.Mode = "idle"
	user.UpdatedAt = time.Now().UTC()
	s.users[userID] = user
	inviter, inviterExists := s.users[user.InvitedBy]
	s.rewardReferralLevelLocked(userID, user.UpdatedAt)
	if err := s.saveLocked(); err != nil {
		if inviterExists {
			s.users[user.InvitedBy] = inviter
		}
		s.users[userID] = previous
		return false, err
	}
	return true, nil
}
