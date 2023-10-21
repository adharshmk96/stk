package project

import (
	"os"
	"path"

	"github.com/adharshmk96/stk/pkg/commands"
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

	GitCmd commands.GitCmd
	GoCmd  commands.GoCmd
}

type TemplateConfig struct {
	PkgName      string
	AppName      string
	ModName      string
	ExportedName string
}

func setDefaults() {
	viper.SetDefault("project.modules", []string{"ping"})
}

func NewContext(args []string) *Context {

	setDefaults()

	goCmd := commands.NewGoCmd()
	gitCmd := commands.NewGitCmd()

	workDir := viper.GetString("project.workdir")
	modules := viper.GetStringSlice("project.modules")

	os.Chdir(workDir)

	packageName := GetPackageName(args)
	appName := GetAppNameFromPkgName(packageName)

	isGitRepo := gitCmd.IsRepo()
	isGoMod := goCmd.IsMod()

	ctx := &Context{
		PackageName: packageName,
		AppName:     appName,
		Modules:     modules,

		IsGoModule: isGoMod,
		IsGitRepo:  isGitRepo,

		WorkDir: workDir,

		GitCmd: gitCmd,
		GoCmd:  goCmd,
	}

	return ctx
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

func (ctx *Context) WriteDefaultConfig() error {

	// project configs
	viper.Set("name", ctx.AppName)
	viper.Set("version", "v0.0.1")
	viper.Set("description", "This project is generated using stk.")
	viper.Set("author", "")

	// module configs
	viper.Set("project.modules", ctx.Modules)

	// Migrator configs
	viper.Set("migrator.workdir", "./stk-migrations")
	viper.Set("migrator.database", "sqlite3")
	viper.Set("migrator.sqlite.filepath", "stk.db")

	// Create the config file
	configPath := path.Join(ctx.WorkDir, ".stk.yaml")
	err := viper.WriteConfigAs(configPath)
	if err != nil {
		return err
	}

	return nil

}
