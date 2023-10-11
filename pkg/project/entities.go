package project

type ProjectConfig struct {
	PkgName  string
	AppName  string
	RootPath string
}

// TODO Compose this ?
type ModuleConfig struct {
	PkgName      string
	AppName      string
	RootPath     string
	ModName      string
	ExportedName string
}
