package dummy

import (
	"math/rand"
	"reconciler/pkg/util"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const idLength = 10

// Generate random ID consist of uppercase alphanumeric
func GenerateRandomID() string {
	result := make([]byte, idLength)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func GenerateRandomAmount(minAmount float64, maxAmount float64) float64 {
	randomNum := minAmount + rand.Float64()*(maxAmount-minAmount)
	roundToTwoDecimals := util.RoundToTwoDecimals(randomNum)

	return roundToTwoDecimals
}

// Generate random type with 70% DEBIT(D), 30% CREDIT(C):
func GenerateRandomTrxType() string {
	r := rand.Float64() // Generates a float64 between 0.0 and 1.0
	if r < 0.7 {
		return util.TRX_TYPE_DEBIT
	}
	return util.TRX_TYPE_CREDIT
}

func GenerateRandomTime(start time.Time, end time.Time) time.Time {
	duration := end.Sub(start)
	randomDuration := time.Duration(rand.Int63n(int64(duration)))
	randomTime := start.Add(randomDuration)

	return randomTime
}
