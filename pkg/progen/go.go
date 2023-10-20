package progen

func initGoMod() error {
	return runCommand("go", "mod", "init")
}

func goModTidy() error {
	return runCommand("go", "mod", "tidy")
}
