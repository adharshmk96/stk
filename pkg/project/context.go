package project

import (
	"os"
	"path"

	"github.com/adharshmk96/stk/consts"
	"github.com/adharshmk96/stk/pkg/commands"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
)

const (
	DEFAULT_WORKDIR = "./"
	CONFIG_FILENAME = ".stk.yaml"
)

type Context struct {
	IsGitRepo  bool
	IsGoModule bool

	PackageName string
	AppName     string

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

func NewContext(args []string) *Context {

	goCmd := commands.NewGoCmd()
	gitCmd := commands.NewGitCmd()

	workDir := viper.GetString("project.workdir")

	os.Chdir(workDir)

	packageName := GetPackageName(args)
	appName := GetAppNameFromPkgName(packageName)

	isGitRepo := gitCmd.IsRepo()
	isGoMod := goCmd.IsMod()

	ctx := &Context{
		PackageName: packageName,
		AppName:     appName,

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
	viper.Set("description", "This project is generated using stk.")
	viper.Set("author", "STK")

	// Migrator configs
	viper.Set(consts.CONFIG_MIGRATOR_WORKDIR, "./stk-migrations")
	viper.Set(consts.CONFIG_MIGRATOR_DB_TYPE, "sqlite3")
	viper.Set(consts.CONFIG_MIGRATOR_DB_FILEPATH, "stk.db")

	// Create the config file
	configPath := path.Join(ctx.WorkDir, CONFIG_FILENAME)
	err := viper.WriteConfigAs(configPath)
	if err != nil {
		return err
	}

	return nil

}
