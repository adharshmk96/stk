package progen

type Config struct {
	PkgName      string
	AppName      string
	RootPath     string
	ModName      string
	ExportedName string

	IsGoModule bool
	IsGitRepo  bool
}
