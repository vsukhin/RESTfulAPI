/* Config package provides methods and data structures to work with system configuration */

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	yaml "gopkg.in/yaml.v2"
)

const (
	CONFIG_NAME        = "conf/application.yml"
	CONFIG_APPLICATION = "getloyalty"
)

func InitConfig() (err error) {
	// Поиск файла конфигурации в текущей директории и по пути выше до корня проекта
	configPath, err := seekConfigFile(CONFIG_NAME, CONFIG_APPLICATION)
	if err != nil {
		logger.Fatalf("Can't find application config file: %v", err)
		return err
	} else {
		if configPath == "" {
			logger.Fatalf("Can't find application config file %s inside of the project %s", CONFIG_NAME, CONFIG_APPLICATION)
			return errors.New("Config path error")
		}
	}

	configData, err := ioutil.ReadFile(filepath.Join(configPath, CONFIG_NAME))
	if err != nil {
		logger.Fatalf("Can't read from configuration file: %v", err)
		return err
	}
	if err = yaml.Unmarshal(configData, &Configuration); err != nil {
		logger.Fatalf("Can't unmarshal data from yaml to configuration structure: %v", err)
		return err
	} else {
		Configuration.Server.Address = fmt.Sprintf("%s:%d", Configuration.Server.Host, Configuration.Server.Port)

		// Подготовка значения переменных для дальнейшего использования (удаление / в конце строки если есть)
		var rex *regexp.Regexp
		rex, _ = regexp.Compile(`(/+)$`)
		Configuration.Server.PublicAddress = rex.ReplaceAllString(Configuration.Server.PublicAddress, "")
	}

	return nil
}

func seekConfigFile(configName string, rootDir string) (configPath string, err error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err = os.Stat(filepath.Join(currentPath, configName)); !os.IsNotExist(err) {
			break
		} else {
			absPath, err := filepath.Abs(currentPath)
			if err != nil {
				return "", err
			}
			if filepath.Base(absPath) == rootDir {
				return "", nil
			}
			if absPath == "/" {
				return "", nil
			}
			currentPath = filepath.Join(currentPath, "../")
		}
	}

	return currentPath, nil
}
