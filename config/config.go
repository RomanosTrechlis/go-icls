// Package config reads from a specific local directory
// all files ending in '.properties' and creates a map
// containing the configuration.
package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Section map[string]string

// Configuration contains a map of sections with
// the corresponding key-value pairs
//
// Example of .properties file:
//
// # lines starting with hashtag are being ignored
// [General]
// filename = test.file
// path = newpath
// key=value.new
//
// [Work]
// hours = 8
// salary=1000
// key=value
//
type Configuration struct {
	// C contains a map with keys the section and
	// values the key-value pairs per section
	C map[string]map[string]string
}

// GetConfigurationFromDir returns the Configuration struct
// with all key-value pairs found inside the .properties
// files of a directory, by section.
func GetConfigurationFromDir(path string) (*Configuration, error) {
	_, err := os.Stat(path)
	if err != nil {
		createConfigPath(path)
		return nil, fmt.Errorf("config file doesn't exist")
	}

	c := &Configuration{
		C: make(map[string]map[string]string, 0),
	}

	fss, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", path, err)
	}
	fss = filterPropertiesFiles(fss...)

	for _, fs := range fss {
		c.readFile(filepath.Join(path, fs.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", path+fs.Name(), err)
		}
	}
	return c, err
}

// GetConfigurationFromSingleFile returns the Configuration struct with
// all key-value pairs found inside the specific .properties file, by section.
func GetConfigurationFromSingleFile(filename string) (*Configuration, error) {
	c := &Configuration{
		C: make(map[string]map[string]string, 0),
	}
	fs, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read fileinfo for %s: %v", filename, err)
	}
	fss := filterPropertiesFiles(fs)
	for _, _ = range fss {
		err = c.readFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", filename, err)
		}
	}
	return c, nil
}

func createConfigPath(path string) {
	os.Mkdir(path, 667)
}

func (c *Configuration) readFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", name, err)
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	var currentSectionName string
	for fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			currentSectionName = c.processSection(line)
			continue
		}
		kv := strings.Split(line, "=")
		c.setNewProperty(currentSectionName, strings.Trim(kv[0], " "), strings.Trim(kv[1], " "))
	}
	return nil
}

func (c *Configuration) processSection(line string) string {
	sectionName := getSectionName(line)
	if _, ok := c.C[sectionName]; ok {
		return sectionName
	}
	section := make(map[string]string, 0)
	c.C[sectionName] = section
	return sectionName
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

func filterPropertiesFiles(fss ...os.FileInfo) []os.FileInfo {
	filtered := make([]os.FileInfo, 0)
	for _, fs := range fss {
		if strings.HasSuffix(fs.Name(), ".properties") {
			filtered = append(filtered, fs)
		}
	}
	return filtered
}
