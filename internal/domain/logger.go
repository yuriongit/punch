package domain

// LogLevel represents the log severity level for worker execution logs.
type LogLevel int

// LogLvl enums initialization
const (
	LogSucc LogLevel = iota + 1
	LogFata
	LogFini
	LogErro
)

// Load returns the string output of the enum
func (l LogLevel) Load() string {
	return [...]string{"LOG-SUCC", "LOG-FATA", "LOG-FINI", "LOG-ERRO"}[l-1]
}

// EnumIdx returns the index of the LogLvl enum.
func (l LogLevel) EnumIdx() int {
	return int(l)
}
