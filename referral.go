package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	referralInviteePremiumDays      = 7
	referralInviterLevelPremiumDays = 7
	referralRewardXPLevel           = 3
	referralDirectRewardPercent     = 20
	referralIndirectRewardPercent   = 5
	referralWithdrawalMinKopecks    = int64(1000 * 100)
	defaultUSDTExchangeRateRubles   = 72
)

type referralApplication struct {
	Applied            bool
	InviterID          int64
	InviteeID          int64
	InviteePremiumDays int
	InviterPremiumDays int
	NewUserUntil       time.Time
	InviteeUntil       time.Time
	InviterUntil       time.Time
	ReferralCode       string
	ReferralCount      int
}

type referralPurchaseReward struct {
	Applied               bool
	BuyerID               int64
	PaymentID             string
	PaidKopecks           int64
	DirectInviterID       int64
	IndirectInviterID     int64
	DirectRewardKopecks   int64
	IndirectRewardKopecks int64
}

type referralLevelReward struct {
	Applied            bool
	InviteeID          int64
	InviterID          int64
	InviterPremiumDays int
	InviterUntil       time.Time
}

func referralLevelReached(xp int) bool {
	level, _, _, _ := knowledgeLevel(xp)
	return level >= referralRewardXPLevel
}

func referralPremiumDuration(days int) time.Duration {
	if days <= 0 {
		return 0
	}
	return time.Duration(days) * 24 * time.Hour
}

func referralDirectRewardKopecks(paidKopecks int64) int64 {
	return referralPercentRewardKopecks(paidKopecks, referralDirectRewardPercent)
}

func referralIndirectRewardKopecks(paidKopecks int64) int64 {
	return referralPercentRewardKopecks(paidKopecks, referralIndirectRewardPercent)
}

func referralPercentRewardKopecks(paidKopecks int64, percent int) int64 {
	if paidKopecks <= 0 || percent <= 0 {
		return 0
	}
	return paidKopecks * int64(percent) / 100
}

func rubToKopecks(rubles int) int64 {
	if rubles <= 0 {
		return 0
	}
	return int64(rubles) * 100
}

func formatRubKopecks(kopecks int64) string {
	if kopecks < 0 {
		kopecks = 0
	}
	if kopecks%100 == 0 {
		return fmt.Sprintf("%d ₽", kopecks/100)
	}
	return fmt.Sprintf("%.2f ₽", float64(kopecks)/100)
}

func formatUSDTFromKopecks(kopecks int64) string {
	return formatUSDTFromKopecksAtRate(kopecks, float64(defaultUSDTExchangeRateRubles))
}

func formatUSDTFromKopecksAtRate(kopecks int64, usdtRubRate float64) string {
	if kopecks < 0 {
		kopecks = 0
	}
	if usdtRubRate <= 0 {
		usdtRubRate = float64(defaultUSDTExchangeRateRubles)
	}
	return fmt.Sprintf("%.2f USDT", float64(kopecks)/100/usdtRubRate)
}

func formatUSDTRubRate(usdtRubRate float64) string {
	return formatUSDTRubRateValue(usdtRubRate) + " ₽"
}

func formatUSDTRubRateValue(usdtRubRate float64) string {
	if usdtRubRate <= 0 {
		usdtRubRate = float64(defaultUSDTExchangeRateRubles)
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", usdtRubRate), "0"), ".")
}
