package config

import (
	"time"

	gelf "github.com/probkiizokna/go-gelf"
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

var Configuration struct {
	WorkingDirectory string        `yaml:"WorkingDirectory"` // Рабочая директория сервера, сразу после запуска приложение меняет текущую директорию
	TempDirectory    string        `yaml:"TempDirectory"`    // Путь к временной директории сервера
	FileStorage      string        `yaml:"FileStorage"`      // Путь к хранилищу файлов
	TemplateStorage  string        `yaml:"TemplateStorage"`  // Путь к хранилищу шаблонов
	ResourceStorage  string        `yaml:"ResourceStorage"`  // Путь к хранилищу строк локализациим
	TableTimeout     time.Duration `yaml:"TableTimeout"`     // Время в милисекундах истечения жизни временной таблицы
	MessageTimeout   time.Duration `yaml:"MessageTimeout"`   // Время в милисекундах возможности изменения сообщения
	FileTimeout      time.Duration `yaml:"FileTimeout"`      // Время в милисекундах истечения жизни файла
	SystemAccount    int64         `yaml:"SystemAccount"`    // Идентификатор объединения системы для финансового учета
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
		Socket   string `yaml:"Socket"`   // Путь и имя сокета базы данных
		Name     string `yaml:"Name"`     // Название базы данных
		Login    string `yaml:"Login"`    // Логин подключения к базе данных
		Password string `yaml:"Password"` // Пароль подключения к базе данных
		Charset  string `yaml:"Charset"`  // Кодировка данных
	} `yaml:"MySQL"`

	Cassandra struct { // Подключение к базе данных Apache Cassandra
		NumConns    int           `yaml:"NumConns"`    // Количество подключений
		Timeout     time.Duration `yaml:"Timeout"`     // Таймаут
		NumRetries  int           `yaml:"NumRetries"`  // Количество попыток повтора запросов при ошибке
		Keyspace    string        `yaml:"Keyspace"`    // Пространство ключей по умолчанию
		Consistency string        `yaml:"Consistency"` // Консистенция запросов
		Servers     []string      `yaml:"Servers"`     // Сервер или сервера при работе в кластере
		Login       string        `yaml:"Login"`       // Логин подключения к серверам
		Password    string        `yaml:"Password"`    // Пароль подключения к серверам
	} `yaml:"Cassandra"`

	Mail struct { //Реквизиты для почтовых отправлений
		Host     string `yaml:"Host"`     // Ip адрес или имя хоста почтового сервера
		Port     int    `yaml:"Port"`     // Порт почтового сервера
		Sender   string `yaml:"Sender"`   // Email рассылки
		Receiver string `yaml:"Receiver"` // Email поддержки
		Login    string `yaml:"Login"`    // Логин подключения к почтовому серверу
		Password string `yaml:"Password"` // Пароль подключения к почтовому серверу
	} `yaml:"Mail"`

	PerformerCommunication struct { // Сервер реализации взаимодействия с конечными поставщиками услуг предоставляющими своё API
		Host           string        `yaml:"Host"`           // Используется клиентом. Публичный адрес для подключения клиентов к серверу по TCP/IP протоколу
		Port           int16         `yaml:"Port"`           // Используется клиентом. Публичный порт для подключения клиентов к серверу по TCP/IP протоколу
		ReconnectDelay time.Duration `yaml:"ReconnectDelay"` // Ожидание перед повторными попытками переподключения клиента к серверу
		ServerHost     string        `yaml:"ServerHost"`     // Используется сервером. Ip адрес или имя хоста на котором работает сервер по TCP/IP протоколу
		ServerPort     int16         `yaml:"ServerPort"`     // Используется сервером. Номер порта на котором работает сервер по TCP/IP протоколу
		Socket         string        `yaml:"Socket"`         // Путь и имя сокета для связи с сервером по unix:socket протоколу
		Mode           string        `yaml:"Mode"`           // Режим связи с сервером (tcp или socket)
		Keys           struct {      // Ключи авторизации и защиты данных соединения
			Ca     []string `yaml:"CA"` // Корневые сертификаты
			Server struct { // Серверные ключи, указываются только в конфигурации на сервере
				Private string `yaml:"Private"` // Приватный ключ сервера
				Public  string `yaml:"Public"`  // Публичный ключ сервера
			} `yaml:"Server"`
			Client struct { // Клиентские ключи, указываются только в конфигурации на клиенте
				Private string `yaml:"Private"` // Приватный ключ клиента
				Public  string `yaml:"Public"`  // Публичный ключ клиента
			} `yaml:"Client"`
		} `yaml:"Keys"`
	} `yaml:"Performer Communication"`
}
