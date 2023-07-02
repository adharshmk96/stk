package project

import "text/template"

type Config struct {
	PkgName  string
	AppName  string
	RootPath string
	DirTree
	DirNames
}

type FileTemplate struct {
	Path     string
	Template *template.Template
}

type DirTree struct {
	CmdPath       string
	ServerPath    string
	EntityPath    string
	HandlerPath   string
	TransportPath string
	ValidatorPath string
	ServicePath   string
	StoragePath   string
	SqlitePath    string
	ConfigPath    string
	SvrerrPath    string
}

type DirNames struct {
	CmdDir        string
	ServerDir     string
	EntityDir     string
	PkgDir        string
	HandlerDir    string
	TransportDir  string
	ValidatorDir  string
	ServiceDir    string
	StorageDir    string
	InfraDir      string
	ConfigDir     string
	SvrerrDir     string
	SqliteRepoDir string
}
