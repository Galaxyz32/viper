package viper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/magiconair/properties"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func findCWD() (string, error) {
	serverFile, err := filepath.Abs(os.Args[0])

	if err != nil {
		return "", fmt.Errorf("Can't get absolute path for executable: %v", err)
	}

	path := filepath.Dir(serverFile)
	realFile, err := filepath.EvalSymlinks(serverFile)

	if err != nil {
		if _, err = os.Stat(serverFile + ".exe"); err == nil {
			realFile = filepath.Clean(serverFile + ".exe")
		}
	}

	if err == nil && realFile != serverFile {
		path = filepath.Dir(realFile)
	}

	return path, nil
}

func unmarshallConfigReader(in io.Reader, c map[string]interface{}, configType string) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

	switch strings.ToLower(configType) {
	case "yaml", "yml":
		if err := yaml.Unmarshal(buf.Bytes(), &c); err != nil {
			return ConfigParseError{err}
		}

	case "json":
		if err := json.Unmarshal(buf.Bytes(), &c); err != nil {
			return ConfigParseError{err}
		}

	case "toml":
		if _, err := toml.Decode(buf.String(), &c); err != nil {
			return ConfigParseError{err}
		}

	case "properties", "props", "prop":
		var p *properties.Properties
		var err error
		if p, err = properties.Load(buf.Bytes(), properties.UTF8); err != nil {
			return ConfigParseError{err}
		}
		for _, key := range p.Keys() {
			value, _ := p.Get(key)
			c[key] = value
		}
	}

	insensitiviseMap(c)
	return nil
}
