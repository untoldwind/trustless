package commands

var GlobalFlags = &struct {
	Debug      bool
	LogFormat  string
	LogFile    string
	ConfigFile string
}{}
