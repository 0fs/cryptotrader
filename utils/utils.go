package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"strconv"
)

func Stf(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	return f
}

func Fts(f float64) string {
	return fmt.Sprintf("%.8f", f)
}

func StochRsi(period int, rsi []float64) ([]float64, []float64, []float64) {
	stochRsi := make([]float64, len(rsi))
	smoothK := make([]float64, len(rsi))
	smoothD := make([]float64, len(rsi))

	var min, max float64
	for i := 0; i < len(rsi); i++ {
		max = -1.0
		min = 101.0
		d := IntMax(i-period, 0)
		for j := i; j > d; j-- {
			min = math.Min(min, rsi[j])
			max = math.Max(max, rsi[j])
		}

		stochRsi[i] = (rsi[i] - min) / (max - min) * 100

		// SmoothK
		d = IntMax(i-3, 0)
		s := 0.0
		for j := i; j > d; j-- {
			s += stochRsi[j]
		}
		smoothK[i] = s / 3.0

		// SmoothD
		s = 0.0
		for j := i; j > d; j-- {
			s += smoothK[j]
		}
		smoothD[i] = s / 3.0
	}

	return stochRsi, smoothK, smoothD
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
