package config

type PhpConfig map[string]*PhpInfo

type PhpInfo struct {
	Version        string
	ExactVersion   string
	InstallPath    string
	IsCLI          bool
	Global         bool
	Extensions     []string
	NeedsRebuild   bool
	FpmPidLocation string
	FpmSocket      string
	PhpIniLocation string
	FpmConfig      string
	PoolConfig     string
	PeclPath       string
}
