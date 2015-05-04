package bedrock

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewBedrock(t *testing.T) {
	b := NewBedrock("/yee/lmao/test")
	assert.Equal(t, "/yee/lmao/test", b.bedrockPath)
	assert.Equal(t, "/yee/lmao/test/composer.json", b.composerJSONPath)
	assert.Equal(t, "/yee/lmao/test/CHANGELOG.md", b.changelogPath)
}

func getTmpDir(test string) string {
	return fmt.Sprintf(
		"/tmp/bump-bedrock-test/%s-%s-%s",
		strconv.FormatInt(time.Now().UnixNano(), 10),
		strconv.Itoa(os.Getpid()),
		test,
	)
}

func makeTmpDir(test string) string {
	dir := getTmpDir(test)
	cpCmd := exec.Command("mkdir", "-p", dir)
	cpCmd.Run()
	return dir
}

func TestGetComposerJSON(t *testing.T) {
	tmpRepo := makeTmpDir("getComposerJSON")
	t.Log(tmpRepo)
	srcFolder := "./fixtures/"
	destFolder := tmpRepo
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd.Run()
	if err != nil {
		panic(err)
	}
	b := NewBedrock(tmpRepo).GetComposerJson()
	assert.Equal(t, "roots/bedrock", b.Path("name").Data().(string))
	assert.Equal(t, "4.2.1", b.Path("require.johnpbloch/wordpress").Data().(string))
}

func TestUpdateComposerJSON(t *testing.T) {
	tmpRepo := makeTmpDir("updateComposerJSON")
	t.Log(tmpRepo)
	srcFolder := "./fixtures/"
	destFolder := tmpRepo
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd.Run()
	if err != nil {
		panic(err)
	}
	b := NewBedrock(tmpRepo)
	assert.Equal(t, "4.2.1", b.GetComposerJson().Path("require.johnpbloch/wordpress").Data().(string))
	b.UpdateComposerJSON("4.0.4")
	assert.Equal(t, "4.0.4", b.GetComposerJson().Path("require.johnpbloch/wordpress").Data().(string))
}

func TestWordPressVersion(t *testing.T) {
	tmpRepo := makeTmpDir("TestWordPressVersion")
	t.Log(tmpRepo)
	srcFolder := "./fixtures/"
	destFolder := tmpRepo
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd.Run()
	if err != nil {
		panic(err)
	}
	b := NewBedrock(tmpRepo)
	assert.Equal(t, "4.2.1", b.WordPressVersion())
}

func TestUpdateWordPressVersion(t *testing.T) {
	tmpRepo := makeTmpDir("updateWordPressVersion")
	t.Log(tmpRepo)
	srcFolder := "./fixtures/"
	destFolder := tmpRepo
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd.Run()
	if err != nil {
		panic(err)
	}
	b := NewBedrock(tmpRepo)
	assert.Equal(t, "nothing to update", b.UpdateWordPressVersion("4.2.1"))
	assert.Equal(t, "updated successfully", b.UpdateWordPressVersion("100.2.1"))
}

func TestUpdateChangelog(t *testing.T) {
	tmpRepo := makeTmpDir("updateChangelog")
	t.Log(tmpRepo)
	srcFolder := "./fixtures/"
	destFolder := tmpRepo
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd.Run()
	if err != nil {
		panic(err)
	}
	b := NewBedrock(tmpRepo)

	b.UpdateChangelog("100.100.100")

	input, err := ioutil.ReadFile(b.changelogPath)
	check(err)

	lines := strings.Split(string(input), "\n")

	assert.Equal(t, "1.3.7", b.GetCurrentBedrockVersion(lines))

	assert.Equal(t, "* Update to WordPress 100.100.100", lines[2])

	b.changelogPath = tmpRepo + "/CHANGELOG-head.md"

	b.UpdateChangelog("100.100.100")

	assert.Equal(t, "1.3.7", b.GetCurrentBedrockVersion(lines))

	assert.Equal(t, "* Update to WordPress 100.100.100", lines[2])

	t.Log(b.changelogPath)
}

func TestAddVersionNote(t *testing.T) {
	lines := []string{
		"before",
		"after",
	}

	b := NewBedrock("yee")

	out := b.AddVersionNote(lines, 0, "100.100.100")

	t.Log(out)

	assert.Equal(t, "* Update to WordPress 100.100.100", out[0])
}
