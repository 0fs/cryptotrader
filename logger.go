package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"runtime"
)

var logger zerolog.Logger

type LogHook struct{}

func initLogger() {
	logger = log.Hook(LogHook{})
	logger.Info().Int("CPU", runtime.NumCPU()).Msg("")
}

func (h LogHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
}
