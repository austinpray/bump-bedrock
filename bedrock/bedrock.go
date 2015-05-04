package bedrock

import (
	"fmt"
	"github.com/jeffail/gabs"
	"github.com/mcuadros/go-version"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type BedrockRepo interface {
	UpdateWordPressVersion(version string) string
}

type BedrockRepoInstance struct {
	bedrockPath, composerJSONPath, changelogPath string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewBedrock(Path string) BedrockRepoInstance {
	return BedrockRepoInstance{
		bedrockPath:      Path,
		composerJSONPath: path.Join(Path, "composer.json"),
		changelogPath:    path.Join(Path, "CHANGELOG.md"),
	}
}

func (b BedrockRepoInstance) GetComposerJson() *gabs.Container {
	dat, err := ioutil.ReadFile(b.composerJSONPath)
	check(err)
	jsonParsed, err := gabs.ParseJSON(dat)
	return jsonParsed
}

func (b BedrockRepoInstance) UpdateComposerJSON(version string) {
	cmd := exec.Command("composer.phar", "require", "johnpbloch/wordpress", version, "--no-update", "--no-progress")
	cmd.Dir = b.bedrockPath
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (b BedrockRepoInstance) AddVersionNote(lines []string, i int, newWordPressVersion string) []string {
	lines = append(
		lines[:i],
		append(
			[]string{fmt.Sprintf("* Update to WordPress %s", newWordPressVersion)},
			lines[i:]...,
		)...,
	)
	return lines
}

func (b BedrockRepoInstance) AddTitle(lines []string, i int, nextBedrockVersion string) []string {
	date := time.Now().Format("2006-01-02")

	title := []string{fmt.Sprintf("### %s: %s", nextBedrockVersion, date), ""}

	return append(title, lines...)
}

func (b BedrockRepoInstance) GetCurrentBedrockVersion(lines []string) string {
	var currentBedrockVersion string

	r := regexp.MustCompile(`### (\d?\.){2}\d`)
	for _, line := range lines {
		if r.FindStringIndex(line) != nil {
			currentBedrockVersion = strings.Split(line, " ")[1]
			currentBedrockVersion = strings.Replace(currentBedrockVersion, ":", "", 1)
			break
		}
	}
	return currentBedrockVersion
}

func (b BedrockRepoInstance) UpdateChangelog(version string) {
	input, err := ioutil.ReadFile(b.changelogPath)
	check(err)

	var currentBedrockVersion string
	var nextBedrockVersion string
	lines := strings.Split(string(input), "\n")

	// find current version
	currentBedrockVersion = b.GetCurrentBedrockVersion(lines)

	// calculate bumped version
	nextBedrockVersionArray := strings.Split(currentBedrockVersion, ".")
	minorVersion, _ := strconv.Atoi(nextBedrockVersionArray[2])
	minorVersion++
	nextBedrockVersionArray[2] = strconv.Itoa(minorVersion)
	nextBedrockVersion = strings.Join(nextBedrockVersionArray, ".")

	// check for HEAD
	i := 0
	if strings.Contains(lines[0], "### HEAD") {
		lines = append(lines[:i], lines[i+2:]...)
	} else {
		lines = append([]string{""}, lines...)
	}
	lines = b.AddVersionNote(lines, i, version)
	lines = b.AddTitle(lines, i, nextBedrockVersion)

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(b.changelogPath, []byte(output), 0644)
	check(err)
}

func (b BedrockRepoInstance) WordPressVersion() string {
	value := b.GetComposerJson().Path("require.johnpbloch/wordpress").Data().(string)
	return value
}

func (b BedrockRepoInstance) UpdateWordPressVersion(v string) string {
	if version.Compare(b.WordPressVersion(), v, "<") {
		b.UpdateComposerJSON(v)
		b.UpdateChangelog(v)
		return "updated successfully"
	} else {
		return "nothing to update"
	}
}
