package project

type Config struct {
	PkgName  string
	AppName  string
	RootPath string
	IsGit    bool
}

type ModuleConfig struct {
	ModName      string
	ExportedName string
}
