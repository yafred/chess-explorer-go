package version

import "fmt"

// BuildTime ... When this tool was built
var BuildTime string

// GitRepository ... Where the source code is
var GitRepository string

// GitTag ... The git tag
var GitTag string

// GitCommit ... The exact level of code
var GitCommit string

// DisplayVersion ... Displays build information
func DisplayVersion() {
	fmt.Println("Build information:")
	if BuildTime == "" {
		fmt.Println("There is no build information ... this must be a local build.")
	} else {
		fmt.Println("Version: ", GitTag)
		fmt.Println("Source code:", GitRepository)
		fmt.Println("Built on:", BuildTime)
		fmt.Println("Git commit:", GitCommit)
	}
}
