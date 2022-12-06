package task

func ApkInstall(pkgs ...string) []string {
	// TODO the --no-cache parameter and the image names should be ApkInstallOpts
	return append([]string{"apk", "add", "--no-cache"}, pkgs...)
}

func GoCompileStatic() []string {
	return []string{"go", "build"}
}

func RunWithDapr(appId string, command string, options ...string) []string {
	return append([]string{"dapr", "run", "-a", appId, "-H", "3500", "-p", "3000", "--", command}, options...)
}
