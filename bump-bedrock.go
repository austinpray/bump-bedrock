package main

import (
	"encoding/json"
	"fmt"
	"github.com/austinpray/bump-bedrock/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/austinpray/bump-bedrock/bedrock"
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

func GetWordPressTags(url string) Tags {
	response, err := http.Get(url)
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

func Bump(b bedrock.BedrockRepo, newWordPressVersion string) {
	fmt.Println(b.UpdateWordPressVersion(newWordPressVersion))
}

func GetVersion(tags Tags) {
	fmt.Println(tags[0].Name)
}

func main() {

	var APITagsUrl string = "https://api.github.com/repos/johnpbloch/wordpress/tags"

	app := cli.NewApp()
	app.Usage = "bump that bedrock version son"

	app.Commands = []cli.Command{
		{
			Name:    "getversion",
			Aliases: []string{"getv"},
			Usage:   "Get the most recent WordPress version",
			Action: func(c *cli.Context) {
				GetVersion(GetWordPressTags(APITagsUrl))
			},
		},
		{
			Name:  "bump",
			Usage: "Execute a bump. Update Changelog, Composer.json",
			Action: func(c *cli.Context) {
				Bump(
					bedrock.NewBedrock(c.Args().First()),
					GetWordPressTags(APITagsUrl)[0].Name,
				)
			},
		},
	}

	app.Run(os.Args)
}
