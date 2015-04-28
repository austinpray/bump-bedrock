package bedrock

import (
	"fmt"
	"github.com/jeffail/gabs"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var bedrockPath string
var composerJSONPath string
var changelogPath string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewBedrock(Path string) {
	bedrockPath = Path
	composerJSONPath = path.Join(bedrockPath, "composer.json")
	changelogPath = path.Join(bedrockPath, "CHANGELOG.md")
}

func GetComposerJson() *gabs.Container {
	dat, err := ioutil.ReadFile(composerJSONPath)
	check(err)
	jsonParsed, err := gabs.ParseJSON(dat)
	return jsonParsed
}

func UpdateComposerJSON(version string) {
	input, err := ioutil.ReadFile(composerJSONPath)
	check(err)

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "johnpbloch/wordpress") {
			lines[i] = fmt.Sprintf("    \"johnpbloch/wordpress\": \"%s\"", version)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(composerJSONPath, []byte(output), 0644)
	check(err)
}

func AddVersionNote(lines []string, i int, newWordPressVersion string) {
	lines = append(
		lines[:i],
		append(
			[]string{fmt.Sprintf("* Update to WordPress %s", newWordPressVersion)},
			lines[i:]...,
		)...,
	)
}

func AddTitle(lines []string, i int, nextBedrockVersion string) []string {
	date := time.Now().Format("2006-01-02")

	title := []string{fmt.Sprintf("### %s: %s", nextBedrockVersion, date), ""}

	return append(title, lines...)
}

func UpdateChangelog(version string) {
	input, err := ioutil.ReadFile(changelogPath)
	check(err)

	var currentBedrockVersion string
	var nextBedrockVersion string
	lines := strings.Split(string(input), "\n")

	// find current version
	r := regexp.MustCompile(`### (\d?\.){2}\d`)
	for _, line := range lines {
		if r.FindStringIndex(line) != nil {
			currentBedrockVersion = strings.Split(line, " ")[1]
			currentBedrockVersion = strings.Replace(currentBedrockVersion, ":", "", 1)
			break
		}
	}

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
	AddVersionNote(lines, i, version)
	lines = AddTitle(lines, i, nextBedrockVersion)

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(changelogPath, []byte(output), 0644)
	check(err)
}

func WordPressVersion() string {
	value := GetComposerJson().Path("require.johnpbloch/wordpress").Data().(string)
	return value
}

func UpdateWordPressVersion(version string) {
	UpdateComposerJSON(version)
	UpdateChangelog(version)
}
