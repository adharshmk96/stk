package project

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
)

const cmdDir = "cmd"
const serverDir = "server"
const pkgDir = "pkg"
const entityDir = "entity"
const httpDir = "http"
const handlerDir = "handler"
const transportDir = "transport"
const validatorDir = "validator"
const serviceDir = "service"
const storageDir = "storage"
const infraDir = "infra"
const configDir = "config"
const svrerrDir = "svrerr"
const sqliteDir = "sqlite"

var cmdPath = filepath.Join(cmdDir)
var serverPath = filepath.Join(serverDir)
var entityPath = filepath.Join(pkgDir, entityDir)
var handlerPath = filepath.Join(pkgDir, httpDir, handlerDir)
var transportPath = filepath.Join(pkgDir, httpDir, transportDir)
var validatorPath = filepath.Join(pkgDir, httpDir, validatorDir)
var servicePath = filepath.Join(pkgDir, serviceDir)
var storagePath = filepath.Join(pkgDir, storageDir)
var configPath = filepath.Join(pkgDir, infraDir, configDir)
var svrerrPath = filepath.Join(pkgDir, infraDir, svrerrDir)
var sqlitePath = filepath.Join(storagePath, sqliteDir)

var dirTree = DirTree{
	CmdPath:       cmdPath,
	ServerPath:    serverPath,
	EntityPath:    entityPath,
	HandlerPath:   handlerPath,
	TransportPath: transportPath,
	ValidatorPath: validatorPath,
	ServicePath:   servicePath,
	StoragePath:   storagePath,
	ConfigPath:    configPath,
	SvrerrPath:    svrerrPath,
	SqlitePath:    sqlitePath,
}

var dirNames = DirNames{
	CmdDir:        cmdDir,
	ServerDir:     serverDir,
	EntityDir:     entityDir,
	PkgDir:        pkgDir,
	HandlerDir:    handlerDir,
	TransportDir:  transportDir,
	ValidatorDir:  validatorDir,
	ServiceDir:    serviceDir,
	StorageDir:    storageDir,
	InfraDir:      infraDir,
	ConfigDir:     configDir,
	SvrerrDir:     svrerrDir,
	SqliteRepoDir: sqliteDir,
}

var dirList []string

func init() {
	paths := reflect.ValueOf(dirTree)
	for i := 0; i < paths.NumField(); i++ {
		dirList = append(dirList, paths.Field(i).String())
	}
}

func CreateProjectStructure(dirs []string) {
	// create dirs
	for _, dir := range dirs {
		log.Println("Creating dir: ", dir)
		os.MkdirAll(filepath.Join(dir), 0755)
	}
}
