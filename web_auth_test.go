package main

import (
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestWebPasswordHashRoundTrip(t *testing.T) {
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if !verifyWebPassword("strong-password", hash) {
		t.Fatal("expected password to verify")
	}
	if verifyWebPassword("wrong-password", hash) {
		t.Fatal("expected wrong password to fail")
	}
}

func TestWebPasswordConfirmation(t *testing.T) {
	if err := validateWebPasswordConfirmation("strong-password", "strong-password"); err != nil {
		t.Fatal(err)
	}
	if err := validateWebPasswordConfirmation("strong-password", ""); err == nil {
		t.Fatal("expected empty confirmation to fail")
	}
	if err := validateWebPasswordConfirmation("strong-password", "another-password"); err == nil {
		t.Fatal("expected mismatched confirmation to fail")
	}
}

func TestWebAuthBrokerVerifiesCodeOnSite(t *testing.T) {
	broker := newWebAuthBroker()
	now := time.Now().UTC()
	request, err := broker.create(0, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := broker.beginTelegram(request.Token, telegramWebProfile{ID: 1001, FirstName: "Demo", Username: "Demor22"}); err != nil {
		t.Fatal(err)
	}
	completed, err := broker.verifySiteCode(request.Token, request.Code, now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if completed.Profile.ID != 1001 || completed.Profile.Username != "Demor22" || !completed.Completed {
		t.Fatalf("unexpected completed auth request: %#v", completed)
	}
	if _, _, err := broker.consume(request.Token, now.Add(time.Minute)); err == nil {
		t.Fatal("expected verified request to be consumed and removed")
	}
}

func TestSQLiteWebAccountCreateAndDuplicate(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	account, err := store.createWebAccount(-42, "Learner_1", hash)
	if err != nil {
		t.Fatal(err)
	}
	if account.UserID != -42 || account.Login != "learner_1" {
		t.Fatalf("unexpected account: %#v", account)
	}
	_, storedHash, ok, err := store.getWebAccountByLogin("LEARNER_1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || !verifyWebPassword("strong-password", storedHash) {
		t.Fatal("expected stored account and valid password hash")
	}
	if _, err := store.createWebAccount(-43, "learner_1", hash); !errors.Is(err, errWebLoginTaken) {
		t.Fatalf("expected duplicate login error, got %v", err)
	}
}

func TestSQLiteWebAccountUpdatePassword(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	oldHash, err := hashWebPassword("old-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-42, "learner_1", oldHash); err != nil {
		t.Fatal(err)
	}
	newHash, err := hashWebPassword("new-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.updateWebAccountPassword(-42, newHash); err != nil {
		t.Fatal(err)
	}
	_, storedHash, ok, err := store.getWebAccountByLogin("learner_1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected stored account")
	}
	if verifyWebPassword("old-password", storedHash) {
		t.Fatal("old password should no longer verify")
	}
	if !verifyWebPassword("new-password", storedHash) {
		t.Fatal("new password should verify")
	}
}

func TestSQLiteTelegramLinkMigratesWebDataWhenTelegramIsEmpty(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-42, "site_user", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-42, "site_user"); err != nil {
		t.Fatal(err)
	}
	if err := store.addXP(-42, 35); err != nil {
		t.Fatal(err)
	}

	account, user, err := store.authenticateTelegramWebAccount(-42, telegramWebProfile{ID: 1001, FirstName: "Demo", Username: "Demor22"})
	if err != nil {
		t.Fatal(err)
	}
	if account.UserID != 1001 || account.Login != "site_user" {
		t.Fatalf("unexpected linked account: %#v", account)
	}
	if user.TelegramID != 1001 || user.XP != 35 {
		t.Fatalf("expected web progress to move to telegram user, got %#v", user)
	}
	if _, ok, err := store.getUser(-42); err != nil || ok {
		t.Fatalf("expected old web user to be removed, ok=%v err=%v", ok, err)
	}
}

func TestSQLiteTelegramLinkMergesExistingTelegramProgress(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-42, "site_user", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-42, "site_user"); err != nil {
		t.Fatal(err)
	}
	if err := store.addXP(-42, 35); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(1001, "telegram_user"); err != nil {
		t.Fatal(err)
	}
	if err := store.addXP(1001, 80); err != nil {
		t.Fatal(err)
	}

	account, user, err := store.authenticateTelegramWebAccount(-42, telegramWebProfile{ID: 1001, FirstName: "Demo", Username: "Demor22"})
	if err != nil {
		t.Fatal(err)
	}
	if account.UserID != 1001 || account.Login != "site_user" {
		t.Fatalf("unexpected linked account: %#v", account)
	}
	if user.TelegramID != 1001 || user.XP != 115 {
		t.Fatalf("expected telegram and web progress to merge, got %#v", user)
	}
	if _, ok, err := store.getUser(-42); err != nil || ok {
		t.Fatalf("expected old web user to be removed, ok=%v err=%v", ok, err)
	}
}

func TestSQLiteTelegramLinkProgressMergeChoiceKeepsPrimarySettings(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-42, "site_user", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-42, "site_user"); err != nil {
		t.Fatal(err)
	}
	if err := store.addXP(-42, 35); err != nil {
		t.Fatal(err)
	}
	if err := store.setUserLevel(-42, "B2"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(1001, "telegram_user"); err != nil {
		t.Fatal(err)
	}
	if err := store.addXP(1001, 80); err != nil {
		t.Fatal(err)
	}
	if err := store.setUserLevel(1001, "A1"); err != nil {
		t.Fatal(err)
	}
	profile := telegramWebProfile{ID: 1001, FirstName: "Demo", Username: "Demor22"}
	challenge, err := store.telegramWebProgressMergeChallenge(-42, profile)
	if err != nil {
		t.Fatal(err)
	}
	if !challenge.Required || challenge.Web.XP != 35 || challenge.Telegram.XP != 80 {
		t.Fatalf("unexpected merge challenge: %#v", challenge)
	}

	_, user, err := store.authenticateTelegramWebAccountWithChoice(-42, profile, "web")
	if err != nil {
		t.Fatal(err)
	}
	if user.XP != 115 || user.Level != "B2" {
		t.Fatalf("expected merged XP with web settings kept, got %#v", user)
	}
}

func TestSQLiteTelegramLinkReplacesTemporaryTelegramWebAccount(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-42, "site_user", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-42, "site_user"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.addLearnedWord(-42, vocabWord{ID: "en:apple", Language: "en", Russian: "яблоко", English: "apple"}); err != nil {
		t.Fatal(err)
	}
	markTestWordMastered(t, store, -42, "en:apple")
	if _, err := store.createWebAccount(1001, "tg_1001", telegramDisabledPasswordHash(1001)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(1001, "tg_1001"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.addLearnedWord(1001, vocabWord{ID: "en:house", Language: "en", Russian: "дом", English: "house"}); err != nil {
		t.Fatal(err)
	}
	markTestWordMastered(t, store, 1001, "en:house")

	account, user, err := store.authenticateTelegramWebAccount(-42, telegramWebProfile{ID: 1001, FirstName: "Demo", Username: "Demor22"})
	if err != nil {
		t.Fatal(err)
	}
	if account.UserID != 1001 || account.Login != "site_user" || !account.PasswordSet {
		t.Fatalf("expected temporary telegram account to be replaced by web login, got %#v", account)
	}
	if user.TelegramID != 1001 || user.FirstName != "Demor22" || len(user.LearnedWords) != 2 {
		t.Fatalf("expected merged telegram user with username display, got %#v", user)
	}
	if _, ok, err := store.getUser(-42); err != nil || ok {
		t.Fatalf("expected old web user to be removed, ok=%v err=%v", ok, err)
	}
	overall := store.leaderboard(10)
	if len(overall) != 1 || overall[0].FirstName != "Demor22" || overall[0].Score != 30 {
		t.Fatalf("expected linked telegram username in unified leaderboard, got %#v", overall)
	}
}

func TestSQLiteTelegramLinkRejectsReplacingLinkedTelegram(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(1001, "linked_user", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(1001, "linked_user"); err != nil {
		t.Fatal(err)
	}

	if _, _, err := store.authenticateTelegramWebAccount(1001, telegramWebProfile{ID: 2002, FirstName: "Other"}); !errors.Is(err, errWebAccountTGLocked) {
		t.Fatalf("expected locked telegram error, got %v", err)
	}
	account, ok, err := store.getWebAccountByUserID(1001)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || account.Login != "linked_user" {
		t.Fatalf("expected original account to stay linked, ok=%v account=%#v", ok, account)
	}
	if _, ok, err := store.getWebAccountByUserID(2002); err != nil || ok {
		t.Fatalf("expected no account for replacement telegram, ok=%v err=%v", ok, err)
	}
}

func TestSQLiteTelegramLinkRejectsUsedTelegramForAnotherWebLogin(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-42, "site_user", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-42, "site_user"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(1001, "telegram_user", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(1001, "telegram_user"); err != nil {
		t.Fatal(err)
	}

	if _, _, err := store.authenticateTelegramWebAccount(-42, telegramWebProfile{ID: 1001, FirstName: "Telegram"}); !errors.Is(err, errWebTelegramUsed) {
		t.Fatalf("expected used telegram error, got %v", err)
	}
	siteAccount, ok, err := store.getWebAccountByUserID(-42)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || siteAccount.Login != "site_user" {
		t.Fatalf("expected web-only account to stay untouched, ok=%v account=%#v", ok, siteAccount)
	}
	telegramAccount, ok, err := store.getWebAccountByUserID(1001)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || telegramAccount.Login != "telegram_user" {
		t.Fatalf("expected telegram account to stay untouched, ok=%v account=%#v", ok, telegramAccount)
	}
}

func TestSQLiteReferralCodeGeneratedAndLookup(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()

	first, err := store.getOrCreateUser(-101, "first")
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.getOrCreateUser(-102, "second")
	if err != nil {
		t.Fatal(err)
	}
	if len(first.ReferralCode) != webReferralCodeLength || len(second.ReferralCode) != webReferralCodeLength {
		t.Fatalf("expected 8-char codes, got %q and %q", first.ReferralCode, second.ReferralCode)
	}
	if first.ReferralCode == second.ReferralCode {
		t.Fatalf("expected unique codes, both were %q", first.ReferralCode)
	}
	found, ok, err := store.getUserByReferralCode(first.ReferralCode)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || found.TelegramID != first.TelegramID {
		t.Fatalf("unexpected referral lookup: ok=%v user=%#v", ok, found)
	}
}

func TestSQLiteApplyWebReferralRewardsOnce(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()

	inviter, err := store.getOrCreateUser(-201, "inviter")
	if err != nil {
		t.Fatal(err)
	}
	result, err := store.applyWebReferral(-202, inviter.ReferralCode)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Applied {
		t.Fatal("expected referral to apply")
	}
	newUser, err := store.mustGetUser(-202)
	if err != nil {
		t.Fatal(err)
	}
	refreshedInviter, err := store.mustGetUser(-201)
	if err != nil {
		t.Fatal(err)
	}
	if newUser.InvitedBy != inviter.TelegramID {
		t.Fatalf("expected invited_by %d, got %d", inviter.TelegramID, newUser.InvitedBy)
	}
	if !newUser.isPremium(result.NewUserUntil.Add(-1)) {
		t.Fatalf("expected new user premium until %s", result.NewUserUntil)
	}
	if refreshedInviter.ReferralCount != 1 {
		t.Fatalf("unexpected referral count: %#v", refreshedInviter)
	}
	if refreshedInviter.isPremium(time.Now().UTC()) || result.InviterPremiumDays != 0 {
		t.Fatalf("inviter should not receive premium before invitee reaches level 3: result=%#v user=%#v", result, refreshedInviter)
	}
	if err := store.addXP(newUser.TelegramID, 280); err != nil {
		t.Fatal(err)
	}
	rewardedInviter, err := store.mustGetUser(inviter.TelegramID)
	if err != nil {
		t.Fatal(err)
	}
	rewardedInvitee, err := store.mustGetUser(newUser.TelegramID)
	if err != nil {
		t.Fatal(err)
	}
	if !rewardedInvitee.ReferralLevelRewarded {
		t.Fatal("expected invitee referral level reward marker")
	}
	if !rewardedInviter.isPremium(time.Now().UTC()) {
		t.Fatalf("expected inviter premium after invitee reaches level 3: %#v", rewardedInviter)
	}
	firstUntil := rewardedInviter.PremiumUntil
	if err := store.addXP(newUser.TelegramID, 1000); err != nil {
		t.Fatal(err)
	}
	rewardedInviterAgain, err := store.mustGetUser(inviter.TelegramID)
	if err != nil {
		t.Fatal(err)
	}
	if !rewardedInviterAgain.PremiumUntil.Equal(firstUntil) {
		t.Fatalf("expected level reward only once, premium changed from %s to %s", firstUntil, rewardedInviterAgain.PremiumUntil)
	}

	second, err := store.applyWebReferral(-202, inviter.ReferralCode)
	if err != nil {
		t.Fatal(err)
	}
	if second.Applied {
		t.Fatal("expected second referral application to be ignored")
	}
	again, err := store.mustGetUser(-201)
	if err != nil {
		t.Fatal(err)
	}
	if again.ReferralCount != 1 {
		t.Fatalf("expected referral count to stay 1, got %d", again.ReferralCount)
	}
}

func TestSQLiteReferralInviteeBonusAndLevelReward(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()

	inviter, err := store.getOrCreateUser(-301, "inviter")
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 11; i++ {
		inviteeID := -400 - int64(i)
		result, err := store.applyReferral(inviteeID, inviter.TelegramID)
		if err != nil {
			t.Fatal(err)
		}
		if !result.Applied || result.InviterPremiumDays != 0 || result.InviteePremiumDays != referralInviteePremiumDays {
			t.Fatalf("invite %d result = %#v, want no inviter premium before level 3 and invitee days %d", i, result, referralInviteePremiumDays)
		}
		invitee, err := store.mustGetUser(inviteeID)
		if err != nil {
			t.Fatal(err)
		}
		if invitee.InvitedBy != inviter.TelegramID || !invitee.isPremium(time.Now()) {
			t.Fatalf("invitee %d did not receive referral premium: %#v", i, invitee)
		}
		if i == 1 {
			if err := store.addXP(inviteeID, 280); err != nil {
				t.Fatal(err)
			}
		}
	}
	refreshedInviter, err := store.mustGetUser(inviter.TelegramID)
	if err != nil {
		t.Fatal(err)
	}
	if refreshedInviter.ReferralCount != 11 {
		t.Fatalf("referral count = %d, want 11", refreshedInviter.ReferralCount)
	}
	if !refreshedInviter.isPremium(time.Now().UTC()) {
		t.Fatalf("expected inviter premium after first invitee reaches level 3: %#v", refreshedInviter)
	}
}

func TestSQLiteReferralPurchaseCreditsTwoLevelsOnce(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()

	grand, err := store.getOrCreateUser(-501, "grand")
	if err != nil {
		t.Fatal(err)
	}
	direct, err := store.getOrCreateUser(-502, "direct")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.applyReferral(direct.TelegramID, grand.TelegramID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.applyReferral(-503, direct.TelegramID); err != nil {
		t.Fatal(err)
	}

	reward, err := store.creditReferralPurchase(-503, "payment-1", rubToKopecks(1000))
	if err != nil {
		t.Fatal(err)
	}
	if !reward.Applied || reward.DirectRewardKopecks != rubToKopecks(200) || reward.IndirectRewardKopecks != rubToKopecks(50) {
		t.Fatalf("unexpected reward: %#v", reward)
	}
	refreshedDirect, err := store.mustGetUser(direct.TelegramID)
	if err != nil {
		t.Fatal(err)
	}
	refreshedGrand, err := store.mustGetUser(grand.TelegramID)
	if err != nil {
		t.Fatal(err)
	}
	if refreshedDirect.ReferralBalanceKopecks != rubToKopecks(200) || refreshedGrand.ReferralBalanceKopecks != rubToKopecks(50) {
		t.Fatalf("unexpected balances: direct=%d grand=%d", refreshedDirect.ReferralBalanceKopecks, refreshedGrand.ReferralBalanceKopecks)
	}

	duplicate, err := store.creditReferralPurchase(-503, "payment-1", rubToKopecks(1000))
	if err != nil {
		t.Fatal(err)
	}
	if duplicate.Applied {
		t.Fatal("expected duplicate payment reward to be ignored")
	}
	refreshedDirect, _ = store.mustGetUser(direct.TelegramID)
	refreshedGrand, _ = store.mustGetUser(grand.TelegramID)
	if refreshedDirect.ReferralBalanceKopecks != rubToKopecks(200) || refreshedGrand.ReferralBalanceKopecks != rubToKopecks(50) {
		t.Fatalf("duplicate changed balances: direct=%d grand=%d", refreshedDirect.ReferralBalanceKopecks, refreshedGrand.ReferralBalanceKopecks)
	}
}

func TestSQLiteLeaderboardIsUnifiedAcrossWebAndTelegram(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-1, "web_one", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-2, "web_two", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-1, "web_one"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-2, "web_two"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(1001, "telegram_user"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(1002, "xp_only_telegram"); err != nil {
		t.Fatal(err)
	}
	if err := store.addXP(-1, 40); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.addLearnedWord(-2, vocabWord{ID: "en:apple", Language: "en", Russian: "яблоко", English: "apple"}); err != nil {
		t.Fatal(err)
	}
	markTestWordMastered(t, store, -2, "en:apple")
	if _, _, err := store.addLearnedWord(1001, vocabWord{ID: "en:house", Language: "en", Russian: "дом", English: "house"}); err != nil {
		t.Fatal(err)
	}
	markTestWordMastered(t, store, 1001, "en:house")
	if err := store.addXP(1002, 1000); err != nil {
		t.Fatal(err)
	}

	overall := store.leaderboard(10)
	if len(overall) != 2 {
		t.Fatalf("expected unified leaderboard with 2 language-active entries, got %d: %#v", len(overall), overall)
	}
	names := map[string]int{}
	for _, entry := range overall {
		names[entry.FirstName] = entry.Score
	}
	if names["web_two"] != 15 || names["telegram_user"] != 15 {
		t.Fatalf("unexpected unified leaders: %#v", overall)
	}
	if _, ok := names["web_one"]; ok {
		t.Fatalf("xp-only web account should not enter leaderboard: %#v", overall)
	}
	if _, ok := names["xp_only_telegram"]; ok {
		t.Fatalf("xp-only telegram account should not enter leaderboard: %#v", overall)
	}

	byLanguage := store.languageLeaderboard("en", 10)
	if len(byLanguage) != 2 {
		t.Fatalf("expected unified language leaderboard, got %d: %#v", len(byLanguage), byLanguage)
	}
}

func markTestWordMastered(t *testing.T, store *sqliteStore, userID int64, wordID string) {
	t.Helper()
	for i := 0; i < learnedWordMasteryThreshold; i++ {
		if err := store.markWordGameCorrect(userID, wordID); err != nil {
			t.Fatalf("markWordGameCorrect(%s): %v", wordID, err)
		}
	}
}
