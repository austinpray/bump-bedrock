package main

import (
	"encoding/json"
	"fmt"
	"github.com/austinpray/bump-bedrock/bedrock"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"os"
)

type Commit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Tag struct {
	Name       string `json:"name"`
	ZipballUrl string `json:"zipball_url"`
	TarballUrl string `json:"tarball_url"`
	Commit     Commit `json:"commit"`
}

type Tags []Tag

var APITagsUrl string = "https://api.github.com/repos/johnpbloch/wordpress/tags"

func GetWordPressTags() Tags {
	response, err := http.Get(APITagsUrl)
	res := Tags{}
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		json.Unmarshal([]byte(contents), &res)
	}
	return res
}

func Bump(c *cli.Context) {
	dir := c.Args().First()
	if dir == "" {
		fmt.Println("You need to specify a bedrock dir")
		os.Exit(1)
	}

	bedrock.NewBedrock(dir)

	bedrock.UpdateWordPressVersion(GetWordPressTags()[0].Name)
}

func GetVersion(c *cli.Context) {
	fmt.Println(GetWordPressTags()[0].Name)
}

func main() {
	app := cli.NewApp()
	app.Usage = "bump that bedrock version son"

	app.Commands = []cli.Command{
		{
			Name:    "getversion",
			Aliases: []string{"getv"},
			Usage:   "Get the most recent WordPress version",
			Action:  GetVersion,
		},
		{
			Name:   "bump",
			Usage:  "Execute a bump. Update Changelog, Composer.json",
			Action: Bump,
		},
	}

	app.Run(os.Args)
}
