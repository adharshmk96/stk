package project

type Config struct {
	PkgName  string
	AppName  string
	RootPath string
}

type ModuleConfig struct {
	ModName      string
	ExportedName string
}
