package client

// Структура для обогащения возрастом
type Age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   *int   `json:"age"`
}

// Структура для обогащения полом
type Sex struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      *string `json:"gender"`
	Probability float32 `json:"probability"`
}

// Структура для обогащения национальностью
type Nation struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []country `json:"country"`
}

type country struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}
