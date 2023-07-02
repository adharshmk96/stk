package project

import (
	"path/filepath"
	"text/template"

	"github.com/adharshmk96/stk/pkg/project/tpl"
)

const gitIgnoreTemplate = `# If you prefer the allow list template instead of the deny list, see community template:
# https://github.com/github/gitignore/blob/main/community/Golang/Go.AllowList.gitignore
#
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, build with ` + "`go test -c`" + `
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# project generated files
*.db
`

func cmdFileTemplate(config *Config) []FileTemplate {
	return []FileTemplate{
		{
			Path:     filepath.Join(config.RootPath, "main.go"),
			Template: template.Must(template.New("main").Parse(string(tpl.MainTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.CmdPath, "root.go"),
			Template: template.Must(template.New("cmdroot").Parse(string(tpl.CmdRootTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.CmdPath, "version.go"),
			Template: template.Must(template.New("cmdversion").Parse(string(tpl.CmdVersionTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.CmdPath, "serve.go"),
			Template: template.Must(template.New("cmdserve").Parse(string(tpl.CmdServeTemplate()))),
		},
	}
}

func serverFileTemplate(config *Config) []FileTemplate {
	return []FileTemplate{
		{
			Path:     filepath.Join(config.RootPath, dirTree.ServerPath, "setup.go"),
			Template: template.Must(template.New("server_setup").Parse(string(tpl.ServerSetupTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.ServerPath, "routing.go"),
			Template: template.Must(template.New("server_routing").Parse(string(tpl.ServerRouterTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.ServerPath, "middleware.go"),
			Template: template.Must(template.New("server_middleware").Parse(string(tpl.ServerMiddlewareTemplate()))),
		},
	}
}

func handlerFileTemplate(config *Config) []FileTemplate {
	return []FileTemplate{
		{
			Path:     filepath.Join(config.RootPath, dirTree.HandlerPath, "handler.go"),
			Template: template.Must(template.New("handler").Parse(string(tpl.HandlerHandlersTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.HandlerPath, "ping.go"),
			Template: template.Must(template.New("handler_ping").Parse(string(tpl.HandlerPingTemplate()))),
		},
	}
}

func entityFileTemplate(config *Config) []FileTemplate {
	return []FileTemplate{
		// {
		// 	Path:     filepath.Join(config.RootPath, dirTree.EntityPath, "entity.go"),
		// 	Template: template.Must(template.New("entity").Parse(string(tpl.EntitiesHandlersTempalte()))),
		// },
		{
			Path:     filepath.Join(config.RootPath, dirTree.EntityPath, "handler.go"),
			Template: template.Must(template.New("entity_handler").Parse(string(tpl.EntityHandlersTempalte()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.EntityPath, "service.go"),
			Template: template.Must(template.New("entity_service").Parse(string(tpl.EntityServicesTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.EntityPath, "storage.go"),
			Template: template.Must(template.New("entity_storage").Parse(string(tpl.EntityStorageTemplate()))),
		},
	}
}

func serviceFileTemplate(config *Config) []FileTemplate {
	return []FileTemplate{
		{
			Path:     filepath.Join(config.RootPath, dirTree.ServicePath, "service.go"),
			Template: template.Must(template.New("service").Parse(string(tpl.ServiceServiceTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.ServicePath, "ping.go"),
			Template: template.Must(template.New("service_ping").Parse(string(tpl.ServicePingTemplate()))),
		},
	}
}

func storageFileTemplate(config *Config) []FileTemplate {
	return []FileTemplate{
		{
			Path:     filepath.Join(config.RootPath, dirTree.SqlitePath, "sqlite.go"),
			Template: template.Must(template.New("storage").Parse(string(tpl.StorageSqliteTemplate()))),
		},
		{
			Path:     filepath.Join(config.RootPath, dirTree.SqlitePath, "ping.go"),
			Template: template.Must(template.New("storage_ping").Parse(string(tpl.StorageSqlitePingTemplate()))),
		},
	}
}
