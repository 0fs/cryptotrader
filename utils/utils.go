package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
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
