package models

// Структура для хранения ftp доступа
type ApiInputFtp struct {
	Ready             bool   `json:"ready" db:"ready"`               // Можно отображать или сообщать
	Customer_Table_ID int64  `json:"tableId" db:"customer_table_id"` // Идентификатор таблицы связанной с данным ftp
	Host              string `json:"host" db:"host"`                 // Имя хоста или ip
	Port              int    `json:"port" db:"port"`                 // Порт
	Path              string `json:"path" db:"path"`                 // Путь от корня
	Login             string `json:"login" db:"login"`               // Логин
	Password          string `json:"password" db:"password"`         // Пароль
}

type DtoInputFtp struct {
	Order_ID          int64  `db:"order_id"`          // Идентификатор заказа
	Ready             bool   `db:"ready"`             // Можно отображать или сообщать
	Customer_Table_ID int64  `db:"customer_table_id"` // Идентификатор таблицы связанной с данным ftp
	Host              string `db:"host"`              // Имя хоста или ip
	Port              int    `db:"port"`              // Порт
	Path              string `db:"path"`              // Путь от корня
	Login             string `db:"login"`             // Логин
	Password          string `db:"password"`          // Пароль
}

// Конструктор создания объекта ftp доступа в api
func NewApiInputFtp(ready bool, customer_table_id int64, host string, port int, path string, login string, password string) *ApiInputFtp {
	return &ApiInputFtp{
		Ready:             ready,
		Customer_Table_ID: customer_table_id,
		Host:              host,
		Port:              port,
		Path:              path,
		Login:             login,
		Password:          password,
	}
}

// Конструктор создания объекта ftp доступа в бд
func NewDtoInputFtp(order_id int64, ready bool, customer_table_id int64, host string, port int, path string, login string, password string) *DtoInputFtp {
	return &DtoInputFtp{
		Order_ID:          order_id,
		Ready:             ready,
		Customer_Table_ID: customer_table_id,
		Host:              host,
		Port:              port,
		Path:              path,
		Login:             login,
		Password:          password,
	}
}
