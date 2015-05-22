/* Config package provides methods and data structures to work with system configuration */

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	gelf "github.com/probkiizokna/go-gelf"
	yaml "gopkg.in/yaml.v2"
)

const (
	CONFIG_NAME        = "conf/application.yml"
	CONFIG_APPLICATION = "getloyalty"
)

type ModeName string
type LevelName string

type MysqlConfiguration struct {
	Driver   string // Драйвер
	Host     string // Хост базы данных
	Port     int16  // Порт подключения по протоколу tcp/ip
	Type     string // Тип подключения к базе данных socket | tcp
	Socket   string // Путь к socket файлу
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных
}

var (
	Configuration struct {
		WorkingDirectory string `yaml:"WorkingDirectory"` // Рабочая директория сервера, сразу после запуска приложение меняет текущую директорию
		TempDirectory    string `yaml:"TempDirectory"`    // Путь к временной директории сервера
		Server           struct {
			Host               string        `yaml:"Host"` // IP адрес или имя хоста на котором поднимается сервер, можно указывать 0.0.0.0 для всех ip адресов
			Port               uint32        `yaml:"Port"` // tcp/ip порт занимаемый сервером
			Address            string        // Консолидированный адрес Host:Port
			PublicAddress      string        `yaml:"PublicAddress"`      // Публичный адрес на котором сервер доступен извне
			Socket             string        `yaml:"Socket"`             // Unix socket на котором поднимается сервер, только для unix-like операционных систем Linux, Unix, Mac
			Mode               string        `yaml:"Mode"`               // Режим работы
			ReadTimeout        time.Duration `yaml:"ReadTimeout"`        // Время в милисекундах ожидания запроса
			WriteTimeout       time.Duration `yaml:"WriteTimeout"`       // Время в милисекундах ожидания выдачи ответа
			MaxHeaderBytes     int           `yaml:"MaxHeaderBytes"`     // Максимальный размер заголовка http запроса в байтах
			KeepAlive          int           `yaml:"KeepAlive"`          // Режим работы соединения Keep Alive
			DocumentRoot       string        `yaml:"DocumentRoot"`       // Корень http сервера
			SessionTimeout     time.Duration `yaml:"SessionTimeout"`     // Время в милисекундах ожидания завершения неактивной сессии
			DefaultLanguage    string        `yaml:"DefaultLanguage"`    // Язык по умолчанию в формате ISO 639-2
			AvailableLanguages []string      `yaml:"AvailableLanguages"` // Список языков в формате ISO 639-2
			FileStorage        string        `yaml:"FileStorage"`        // Путь к хранилищу файлов
			TemplateStorage    string        `yaml:"TemplateStorage"`    // Путь к хранилищу шаблонов
			FileTimeout        time.Duration `yaml:"FileTimeout"`        // Время в милисекундах истечения жизни файла
			ResourceStorage    string        `yaml:"ResourceStorage"`    // Путь к хранилищу строк локализациим
			TableTimeout       time.Duration `yaml:"TableTimeout"`       // Время в милисекундах истечения жизни временной таблицы
			MessageTimeout     time.Duration `yaml:"MessageTimeout"`     // Время в милисекундах возможности изменения сообщения
		} `yaml:"Server"`

		Logger struct { // Система логирования
			Mode    []string               `yaml:"Mode"`   // Режим логирования, перечисляются включенные режимы логирования
			Levels  map[ModeName]LevelName `yaml:"Levels"` // Уровень логирования для каждого режима логирования
			Format  string                 `yaml:"Format"` // Формат строки лога
			File    string                 `yaml:"File"`   // Режим вывода в файл, путь и имя файла лога
			Graylog struct {               // Настройки подключения к graylog серверу
				Host        string               `yaml:"Host"`        // IP адрес или имя хоста Graylog сервера
				Port        uint32               `yaml:"Port"`        // Порт на котором находится Graylog сервер
				Proto       string               `yaml:"Proto"`       // Протокол передачи данных, возможные значения: tcp, udp. По умолчанию: udp
				Source      string               `yaml:"Source"`      // Наименование источника логов
				ChunkSize   uint32               `yaml:"ChunkSize"`   // Максимальный размер отправляемого пакета
				Compression gelf.CompressionType `yaml:"Compression"` // Сжатие передаваемых пакетов данных
			} `yaml:"Graylog"`
		} `yaml:"Logs"`

		MySql []struct { // Реквизиты подключения к базе данных
			Driver   string `yaml:"Driver"`   // Название драйвера
			Host     string `yaml:"Host"`     // Ip адрес или имя хоста базы данных
			Port     int16  `yaml:"Port"`     // Порт подключения для режима tcp/ip
			Type     string `yaml:"Type"`     // Тип или режим подключения к базе данных. Возможные значения: socket, tcp
			Socker   string `yaml:"Socker"`   // Путь и имя сокета базы данных
			Name     string `yaml:"Name"`     // Название базы данных
			Login    string `yaml:"Login"`    // Логин подключения к базе данных
			Password string `yaml:"Password"` // Пароль подключения к базе данных
			Charset  string `yaml:"Charset"`  // Кодировка данных
		} `yaml:"MySQL"`

		Mail struct { //Реквизиты для почтовых отправлений
			Host     string `yaml:"Host"`     // Ip адрес или имя хоста почтового сервера
			Port     int    `yaml:"Port"`     // Порт почтового сервера
			Sender   string `yaml:"Sender"`   // Email рассылки
			Login    string `yaml:"Login"`    // Логин подключения к почтовому серверу
			Password string `yaml:"Password"` // Пароль подключения к почтовому серверу
		} `yaml:"Mail"`
	}
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
