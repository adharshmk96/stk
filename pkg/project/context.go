package project

import (
	"github.com/adharshmk96/stk/pkg/git"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
)

type Context struct {
	PackageName string
	AppName     string
	Modules     []string

	IsGitRepo  bool
	IsGoModule bool

	WorkDir string
}

type TemplateConfig struct {
	PkgName      string
	AppName      string
	ModName      string
	ExportedName string
}

func setDefaults() {
	viper.SetDefault("project.workdir", ".")
	viper.SetDefault("project.modules", []string{"ping"})
}

func NewContext(args []string) *Context {

	setDefaults()

	packageName := GetPackageName(args)
	appName := GetAppNameFromPkgName(packageName)

	isGitRepo := git.IsRepo()
	isGoMod := IsGoModule()

	// todo, get this from config with defaults
	workDir := viper.GetString("project.workdir")
	modules := viper.GetStringSlice("project.modules")

	return &Context{
		PackageName: packageName,
		AppName:     appName,
		Modules:     modules,

		IsGoModule: isGoMod,
		IsGitRepo:  isGitRepo,

		WorkDir: workDir,
	}
}

func GetTemplateConfig(ctx *Context, module string) *TemplateConfig {

	moduleName := strcase.ToLowerCamel(module)
	exportedName := strcase.ToCamel(moduleName)

	return &TemplateConfig{
		PkgName:      ctx.PackageName,
		AppName:      ctx.AppName,
		ModName:      moduleName,
		ExportedName: exportedName,
	}
}
