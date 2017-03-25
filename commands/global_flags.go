package commands

// GlobalFlags holds the values of all global commandline flags
var GlobalFlags = &struct {
	Debug      bool
	LogFormat  string
	LogFile    string
	ConfigFile string
}{}
