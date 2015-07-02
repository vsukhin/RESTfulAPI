package models

// Структура для организации хранения цены услуги
type ApiRange struct {
	Begin int `json:"begin" db:"begin"` // Нижняя граница диапазона
	End   int `json:"end" db:"end"`     // Верхняя граница диапазона
}

type ApiSupplierPrice struct {
	Supplier_ID       int64 `db:"supplier_id"`       // Идентификатор поставщика
	Customer_Table_ID int64 `db:"customer_table_id"` // Идентификатор таблицы
}

type ApiSMSHLRPrice struct {
	Supplier_ID        int64    `json:"supplierId" db:"supplier_id"`              // Идентификатор поставщика
	Mobile_Operator_ID int      `json:"mobileOperatorId" db:"mobile_operator_id"` // Идентификатор мобильного оператора
	AmountRange        ApiRange `json:"range" db:"range"`                         // Диапазон количества sms или hlr запросов
	Price              float64  `json:"price" db:"price"`                         // Стоимость sms или hlr запроса для оператора
}

type ApiRecognizePrice struct {
	Supplier_ID   int64   `json:"supplierId" db:"supplier_id"`        // Идентификатор поставщика
	Product_ID    int     `json:"recognizeProductId" db:"product_id"` // Идентификатор позиции
	Price         float64 `json:"price" db:"price"`                   // Стоимость позиции
	Increase      bool    `json:"orderIncrease" db:"increase"`        // Использование наценки
	PriceIncrease float64 `json:"priceIncrease" db:"price_increase"`  // Наценка
}

type ApiVerifyPrice struct {
	Supplier_ID int64   `json:"supplierId" db:"supplier_id"`           // Идентификатор поставщика
	Product_ID  int     `json:"verificationProductId" db:"product_id"` // Идентификатор позиции
	Price       float64 `json:"price" db:"price"`                      // Стоимость позиции
}

// Конструктор создания объекта цены услуги в api
func NewApiRange(begin int, end int) *ApiRange {
	return &ApiRange{
		Begin: begin,
		End:   end,
	}
}

func NewApiSupplierPrice(supplier_id int64, customer_table_id int64) *ApiSupplierPrice {
	return &ApiSupplierPrice{
		Supplier_ID:       supplier_id,
		Customer_Table_ID: customer_table_id,
	}
}

func NewApiSMSHLRPrice(supplier_id int64, mobile_operator_id int, amountrange ApiRange, price float64) *ApiSMSHLRPrice {
	return &ApiSMSHLRPrice{
		Supplier_ID:        supplier_id,
		Mobile_Operator_ID: mobile_operator_id,
		AmountRange:        amountrange,
		Price:              price,
	}
}

func NewApiRеcognizePrice(supplier_id int64, product_id int, price float64, increase bool, priceincrease float64) *ApiRecognizePrice {
	return &ApiRecognizePrice{
		Supplier_ID:   supplier_id,
		Product_ID:    product_id,
		Price:         price,
		Increase:      increase,
		PriceIncrease: priceincrease,
	}
}

func NewApiVerifyPrice(supplier_id int64, product_id int, price float64) *ApiVerifyPrice {
	return &ApiVerifyPrice{
		Supplier_ID: supplier_id,
		Product_ID:  product_id,
		Price:       price,
	}
}
