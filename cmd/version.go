package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var appRepo = "pattonjp/localcluster"

// BuildInfo versioning info from build
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

// Print prints version info to stdout
func (bi *BuildInfo) Print() {
	fmt.Println("localcluster version: ", bi.Version, bi.Date, bi.Commit)
}

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "displays the current version",
		Run: func(cmd *cobra.Command, args []string) {
			currBuild.Print()
		},
	}
	cmd.AddCommand(updateCmd())
	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			err := doSelfUpdate()
			if err != nil {
				fmt.Printf("Could not update version \n %v+", err)
			}
		},
	}
	return cmd
}

func doSelfUpdate() error {
	v, err := semver.Parse(currBuild.Version)
	if currBuild.Version != "dev" && err != nil {
		return err
	}

	updater, err := getUpdater()
	if err != nil {
		fmt.Println("failed", err)
		return err
	}

	latest, err := updater.UpdateSelf(v, appRepo)
	if err != nil {
		fmt.Println("Binary update failed:", err)
		return nil
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		fmt.Println("Current binary is the latest version", currBuild.Version)
	} else {
		fmt.Println("Successfully updated to version", latest.Version)
		fmt.Println("Release note:\n", latest.ReleaseNotes)
	}
	return nil
}

func getUpdater() (*selfupdate.Updater, error) {
	config := selfupdate.Config{}
	updater, err := selfupdate.NewUpdater(config)
	if err != nil {
		fmt.Println("could not create the updater")
		return nil, err
	}

	return updater, nil
}

func updateCheck() bool {
	updater, err := getUpdater()
	if err != nil {
		return false
	}

	latest, found, err := updater.DetectLatest(appRepo)
	if err != nil || !found {
		fmt.Println("Error occurred while detecting version:", err)
		return false
	}

	v, _ := semver.Parse(currBuild.Version)

	if !found || latest.Version.LTE(v) {
		fmt.Println("Current version is the latest")
		return false
	}

	msg := "\nA new version (%s) is available. \n\n%sDo you want to update now from %s? (y/n): "
	fmt.Printf(msg, latest.Version, latest.ReleaseNotes, currBuild.Version)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || (input == "n\n") {
		fmt.Println("ok... maybe later")
		return false
	}

	err = doSelfUpdate()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
