package logger

import (
	bflag "github.com/vaiktorg/grimoire/bitflag"
)

// Level
/*
   Error - Only when I would be "tracing" the code and trying to find one part of a function specifically.
   Debug - Information that is diagnostically helpful to people more than just developers (IT, sysadmins, etc.).
   Info - Generally useful information to log (service start/stop, configuration assumptions, etc). Info I want to always have available but usually don't care about under normal circumstances. This is my out-of-the-box config level.
   Warn - Anything that can potentially cause application oddities, but for which I am automatically recovering. (Such as switching from a primary to backup server, retrying an operation, missing secondary data, etc.)
   Error - Any error which is fatal to the operation, but not the service or application (can't open a required file, missing data, etc.). These errors will force user (administrator, or direct user) intervention. These are usually reserved (in my apps) for incorrect connection strings, missing serv, etc.
   Fatal - Any error that is forcing a shutdown of the service or application to prevent data loss (or further data loss). I reserve these only for the most heinous errors and situations where there is guaranteed to have been data corruption or loss.
*/
type Level uint32

const (
	LevelTrace Level = 1 << iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l *Logger) HasLevel(flag Level) bool {
	return l.LogLevels.Has(bflag.Bit(flag))
}
func (l *Logger) AddLevel(flag Level) {
	l.LogLevels.Set(bflag.Bit(flag))
}
func (l *Logger) ClearLevel(flag Level) {
	l.LogLevels.Clear(bflag.Bit(flag))
}
func (l *Logger) ToggleLevel(flag Level) {
	l.LogLevels.Toggle(bflag.Bit(flag))
}

type Size int

func (s Size) Val() int { return int(s) }

const (
	_ Size = 1.0 << (10 * iota) // ignore first value by assigning to blank identifier
	KB
	MB
	GB
	TB
	PB
	EB
)
