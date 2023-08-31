package updater

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/fatih/color"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

var (
	layout = "2006-01-02 15:04:05"
)

type VersionManager interface {
	UpdateToLatest() error
	CheckForUpdateGuarded()
	Print()
}

func New(version, commit, date, repo string) VersionManager {
	return &buildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
		AppRepo: repo,
	}
}

// buildInfo versioning info from build
type buildInfo struct {
	Version string
	Commit  string
	Date    string
	AppRepo string
}

func (bi *buildInfo) name() string {
	return filepath.Base(os.Args[0])
}

func (bi *buildInfo) print(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	c := color.New(color.FgHiYellow)
	c.Printf(format, args...)
}

func (bi *buildInfo) printErr(msg string, err error) {
	c := color.New(color.FgHiRed)
	c.Println(msg, err)
}

// Print prints version info to stdout
func (bi *buildInfo) Print() {
	bi.print("%s \n version: %s \n built:   %s \n commit:  %s \n",
		bi.name(),
		bi.Semver().String(),
		bi.Date,
		bi.Commit,
	)
}

func (bi *buildInfo) Semver() semver.Version {
	v, err := semver.Parse(bi.Version)
	if err != nil {

		v = semver.Version{}
	}
	return v
}

func (bi *buildInfo) UpdateToLatest() error {
	v := bi.Semver()
	updater, err := bi.getUpdater()
	if err != nil {
		return err
	}

	latest, err := updater.UpdateSelf(v, bi.AppRepo)
	if err != nil {
		bi.printErr("Binary update failed: ", err)
		return nil
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		bi.print("Current binary is the latest version %s \n", bi.Version)
	} else {
		bi.print("Successfully updated to version %s \n", latest.Version)
		bi.print("Release note:\n%s\n", latest.ReleaseNotes)
	}
	return nil
}

func (bi *buildInfo) getUpdater() (*selfupdate.Updater, error) {
	config := selfupdate.Config{}
	updater, err := selfupdate.NewUpdater(config)
	if err != nil {
		bi.printErr("could not create the updater", err)
		return nil, err
	}

	return updater, nil
}

func (bi *buildInfo) updateCheck() bool {
	updater, err := bi.getUpdater()
	if err != nil {
		return false
	}

	latest, found, err := updater.DetectLatest(bi.AppRepo)
	if err != nil || !found {
		bi.printErr("Error occurred while detecting version:", err)
		return false
	}

	v, _ := semver.Parse(bi.Version)

	if !found || latest.Version.LTE(v) {
		return false
	}

	msg := "\nA new version (%s) is available. \n\n%sDo you want to update now from %s? (y/n): "
	bi.print(msg, latest.Version, latest.ReleaseNotes, bi.Version)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || (input == "n\n") {
		bi.print("ok... maybe later")
		return false
	}

	err = bi.UpdateToLatest()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (bi *buildInfo) getTempFilePath() string {
	return filepath.Join(os.TempDir(), bi.name()+"-lastcheck")

}

func (bi *buildInfo) getLastReadDate() time.Time {
	lastDate := time.Now()
	found := false
	if buf, err := os.ReadFile(bi.getTempFilePath()); err == nil {
		fStr := string(buf)
		fStr = strings.ReplaceAll(fStr, "\n", "")
		fStr = strings.ReplaceAll(fStr, "\r", "")
		if pdate, err := time.Parse(layout, fStr); err == nil {
			lastDate = pdate
		} else {
			bi.printErr("parsing date", err)
		}
		found = true
	}
	if !found {
		bi.writeLastReadDate(lastDate)
	}
	return lastDate
}

func (bi *buildInfo) writeLastReadDate(d time.Time) bool {
	err := os.WriteFile(bi.getTempFilePath(), []byte(d.Format(layout)), 0755)
	return err == nil
}

// CheckForUpdateGuarded use a temp file to avoid
// rate limiting from the github api
func (bi *buildInfo) CheckForUpdateGuarded() {
	if bi.Version == "" {
		return
	}

	lastDate := bi.getLastReadDate()
	nextCheck := lastDate.Add(24 * time.Hour)
	now := time.Now()
	if nextCheck.After(now) {
		return
	}
	bi.updateCheck()
	bi.writeLastReadDate(now)
}
