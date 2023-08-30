package utils

func install(dep Dep) {
	if dep.Version == "" {
		dep.Version = "latest"
	}
	execSteam(asdf, "plugin", "add", dep.Name)
	execSteam(asdf, "install", dep.Name, dep.Version)

}
