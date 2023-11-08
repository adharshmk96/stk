package sqlmigrator

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/adharshmk96/stk/consts"
	"github.com/adharshmk96/stk/pkg/utils"
	"github.com/spf13/viper"
)

const (
	DEFAULT_LOG_FILE = ".commit-status"
)

type MigrationType string

const (
	MigrationUp   MigrationType = "up"
	MigrationDown MigrationType = "down"
)

type Database string

const (
	PostgresDB Database = "postgres"
	MySQLDB    Database = "mysql"
	SQLiteDB   Database = "sqlite"
)

type Migrations []*MigrationFileEntry

type MigrationFileEntry struct {
	Number       int
	Name         string
	Committed    bool
	UpFilePath   string
	DownFilePath string
}

func ParseMigrationEntry(migrationEntry string) (*MigrationFileEntry, error) {
	parts := strings.Split(migrationEntry, "_")
	partLength := len(parts)

	if partLength == 0 {
		return nil, ErrInvalidMigration
	}

	commit_status := parts[partLength-1]
	if commit_status != "up" && commit_status != "down" {
		return nil, ErrInvalidMigration
	}

	name := strings.Join(parts[1:partLength-1], "_")

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, ErrInvalidMigration
	}

	rawMigration := &MigrationFileEntry{
		Name:      name,
		Number:    number,
		Committed: commit_status == "up",
	}

	return rawMigration, nil
}

func (r *MigrationFileEntry) String() string {
	m_String := fmt.Sprintf("%d", r.Number)
	if r.Name != "" {
		m_String = m_String + "_" + r.Name
	}
	return m_String
}

func (r *MigrationFileEntry) EntryString() string {
	entryString := fmt.Sprintf("%d", r.Number)
	if r.Name != "" {
		entryString += "_" + r.Name
	}
	if r.Committed {
		entryString += "_up"
	} else {
		entryString += "_down"
	}
	return entryString
}

func (r *MigrationFileEntry) FileNames(extention string) (string, string) {
	fileName := fmt.Sprintf("%d", r.Number)
	if r.Name != "" {
		fileName += "_" + r.Name
	}
	upFileName := fileName + "_up." + extention
	downFileName := fileName + "_down." + extention
	return upFileName, downFileName
}

func (r *MigrationFileEntry) LoadFileContent() (string, string) {

	upContent, err := os.ReadFile(r.UpFilePath)
	if err != nil {
		return "", ""
	}

	downContent, err := os.ReadFile(r.DownFilePath)
	if err != nil {
		return "", ""
	}

	return string(upContent), string(downContent)
}

type Context struct {
	WorkDir    string
	LogFile    string
	Database   Database
	DryRun     bool
	Migrations Migrations
}

func DefaultContextConfig() (string, Database, string) {
	rootDirectory := viper.GetString(consts.CONFIG_MIGRATOR_WORKDIR)
	dbChoice := viper.GetString(consts.CONFIG_MIGRATOR_DB_TYPE)
	logFile := utils.GetFirst(viper.GetString(consts.CONFIG_MIGRATOR_LOGFILE), DEFAULT_LOG_FILE)

	dbType := SelectDatabase(dbChoice)
	subDir := SelectSubDirectory(dbType)

	workDir := path.Join(rootDirectory, subDir)

	return workDir, dbType, logFile
}

func NewContext(workDir string, dbType Database, logFile string, dry bool) *Context {

	ctx := &Context{
		WorkDir:  workDir,
		Database: dbType,
		LogFile:  logFile,
		DryRun:   dry,
	}

	err := InitializeMigrationsFolder(ctx)
	if err != nil {
		return nil
	}

	return ctx
}

func (ctx *Context) LoadMigrationEntries() error {
	migrations := []*MigrationFileEntry{}
	entires, err := ReadLines(path.Join(ctx.WorkDir, ctx.LogFile))
	if err != nil {
		return err
	}

	for _, entry := range entires {
		migration, err := ParseMigrationEntry(entry)
		if err != nil {
			return err
		}

		upFileName, downFileName := migration.FileNames(SelectExtention(ctx.Database))
		migration.UpFilePath = path.Join(ctx.WorkDir, upFileName)
		migration.DownFilePath = path.Join(ctx.WorkDir, downFileName)

		migrations = append(migrations, migration)
	}

	ctx.Migrations = migrations
	return nil
}

func (ctx *Context) WriteMigrationEntries() error {
	filePath := path.Join(ctx.WorkDir, ctx.LogFile)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	for _, migration := range ctx.Migrations {
		_, err := file.WriteString(migration.EntryString() + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

type MigrationDBEntry struct {
	Number    int
	Name      string
	Direction string
	Created   time.Time
}
