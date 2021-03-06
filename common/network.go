package common

import (
	math2 "github.com/idena-network/idena-go/common/math"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

const (
	ShortSessionFlips      = 6
	LongSessionTesters     = 10
	ShortSessionExtraFlips = 2

	WordDictionarySize = 3300
	WordPairsPerFlip   = 3

	MaxFlipSize    = 1024 * 600
	MaxProfileSize = 1024 * 1024
)

const (
	MinTotalScore       = 0.75
	MinShortScore       = 0.6
	MinLongScore        = 0.75
	MinHumanTotalScore  = 0.92
	MinFlipsForVerified = 13
	MinFlipsForHuman    = 24
	StakeToBalanceCoef  = 0.75
	LastScoresCount     = 10
)

func ShortSessionExtraFlipsCount() uint {
	return ShortSessionExtraFlips
}

func ShortSessionFlipsCount() uint {
	return ShortSessionFlips
}

func LongSessionFlipsCount(networkSize int) uint {
	_, flipsPerIdentity := NetworkParams(networkSize)
	totalFlips := uint64(flipsPerIdentity * networkSize)
	return uint(math2.Max(5, math2.Min(totalFlips, uint64(flipsPerIdentity*LongSessionTesters))))
}

func NetworkParams(networkSize int) (epochDuration int, flips int) {
	if networkSize == 0 {
		return 1, 0
	}
	epochDuration = int(math.Round(math.Pow(float64(networkSize), 0.33)))
	flips = 3
	return
}

func NormalizedEpochDuration(validationTime time.Time, networkSize int) time.Duration {
	const day = 24 * time.Hour
	baseEpochDays, _ := NetworkParams(networkSize)
	if baseEpochDays < 21 {
		return day * time.Duration(baseEpochDays)
	}
	if validationTime.Weekday() != time.Saturday {
		return day * 20
	}
	return day * time.Duration(math2.MinInt(28, int(math.Round(math.Pow(float64(networkSize), 0.33)/21)*21)))
}

func GodAddressInvitesCount(networkSize int) uint16 {
	return uint16(math2.MinInt(500, math2.MaxInt(50, networkSize/3)))
}

func EncodeScore(shortPoints float32, shortFlipsCount uint32) byte {
	p := byte(shortPoints*2) << 4
	q := byte(shortFlipsCount) & 0x0f

	return p | q
}

func DecodeScore(score byte) (shortPoints float32, shortFlipsCount uint32) {
	p := score >> 4
	q := score & 0x0f

	return float32(p) / 2, uint32(q)
}

func CalculateIdentityScores(scores []byte, totalShortPoints float32, totalShortFlipsCount uint32) (totalPoints float32, totalFlips uint32) {
	var sumPoints float32
	var sumFlips uint32

	for _, item := range scores {
		a, b := DecodeScore(item)
		sumPoints += a
		sumFlips += b
	}

	// consider old points
	multiplier := decimal.NewFromInt32(1).Sub(decimal.NewFromFloat32(0.1).Mul(decimal.NewFromInt32(int32(len(scores)))))
	addPoints, _ := decimal.NewFromFloat32(totalShortPoints).Mul(multiplier).Float64()
	addFlips := decimal.NewFromInt32(int32(totalShortFlipsCount)).Mul(multiplier).Round(0).IntPart()

	sumPoints += float32(addPoints)
	sumFlips += uint32(addFlips)

	return sumPoints, sumFlips
}
