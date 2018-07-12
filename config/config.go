// Package config reads from a specific local directory
// all files ending in '.properties' and creates a map
// containing the configuration.
//
// Example .properties file:
//
// # lines starting with hashtag are being ignored
// [General]
// filename = config.json
// configPath = ./test
//
package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Section map[string]string

type Configuration struct {
	C map[string]map[string]string
}

func GetConfiguration(path string) (*Configuration, error) {
	_, err := os.Stat(path)
	if err != nil {
		createConfigPath(path)
		return nil, fmt.Errorf("config file doesn't exist")
	}
	c := &Configuration{
		C: make(map[string]map[string]string, 0),
	}
	err = c.readConfigFiles(path)
	return c, err
}

func createConfigPath(path string) {
	os.Mkdir(path, 667)
}

func (c *Configuration) readConfigFiles(path string) error {
	fss, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", path, err)
	}
	fss = filterPropertiesFiles(fss)

	for _, fs := range fss {
		c.readFile(path + fs.Name())
	}
	return nil
}

func (c *Configuration) readFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", name, err)
	}
	fileScanner := bufio.NewScanner(f)
	var currentSectionName string
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.Trim(line, " ") == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			sectionName := getSectionName(line)
			currentSectionName = sectionName
			if _, ok := c.C[sectionName]; ok {
				continue
			}
			section := make(map[string]string, 0)
			c.C[sectionName] = section
			continue
		}
		kv := strings.Split(line, "=")
		c.setNewProperty(currentSectionName, strings.Trim(kv[0], " "), strings.Trim(kv[1], " "))
	}
	return nil
}

func (c *Configuration) setNewProperty(section, key, value string) {
	s, _ := c.C[section]
	s[key] = value
	c.C[section] = s
}

func getSectionName(line string) string {
	s := strings.Replace(line, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)
	return s
}

func filterPropertiesFiles(fss []os.FileInfo) []os.FileInfo {
	filtered := make([]os.FileInfo, 0)
	for _, fs := range fss {
		if strings.HasSuffix(fs.Name(), ".properties") {
			filtered = append(filtered, fs)
		}
	}
	return filtered
}
