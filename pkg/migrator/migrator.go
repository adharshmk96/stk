package migrator

type MigratorConfig struct {
	RootDirectory string
	Database      string
	NumToMigrate  int
	DryRun        bool
}

func MigrateUp(config *MigratorConfig) (err error) {

	return nil

}
