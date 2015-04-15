package cliUtils

import (
	"os"
	"os/user"
	"testing"
)

func TestFindConfigFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error reading current working directory")
	}
	cases := []struct {
		configFileName string
		startingPoint  string
		expectedPath   string
	}{
		{
			".config",
			"/testdata/configSearchTree/configTreeDeeper/configTreeDeepest",
			"/testdata/configSearchTree/configTreeDeeper/",
		},
		{
			".config",
			"/testdata/configSearchTree/configTreeDeeper",
			"/testdata/configSearchTree/configTreeDeeper/",
		},
	}
	for _, c := range cases {
		c.startingPoint = cwd + c.startingPoint
		c.expectedPath = cwd + c.expectedPath + c.configFileName
		retPath := findConfigFile(c.startingPoint)
		if retPath != c.expectedPath {
			t.Errorf("Find config file error! Sent: %s, Expected: %s, Got: %s", c.startingPoint, c.expectedPath, retPath)
		}
	}
}

func TestGetPath(t *testing.T) {
	basePath := "./testdata/"

	usr, _ := user.Current()
	homeDir := usr.HomeDir

	cases := []struct {
		testPath        string
		expectedAbsPath string
	}{
		{
			"/etc/ssl/certs/tls-client.crt",
			"/etc/ssl/certs/tls-client.crt",
		},
		{
			"~/.config",
			homeDir + "/.config",
		},
	}
	for _, c := range cases {
		retPath := getAbsPath(c.testPath, basePath)
		if retPath != c.expectedAbsPath {
			t.Errorf("Path resolution error! Sent: %s, Expected: %s, Got: %s", c.testPath, c.expectedAbsPath, retPath)
		}
	}
}
