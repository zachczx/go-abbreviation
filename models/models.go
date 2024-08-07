package models

type Abv struct {
	Short     string `json:"short"`
	Long      string `json:"long"`
	Initial   string
	Metaphone string
	// Remarks string `json:"remarks"`
}

type List struct {
	List []Abv `json:"abbreviations"`
}

type AbvDb struct {
	Short     string
	Long      string
	Initial   string
	Metaphone string
}

type AbvDbList struct {
	AbvDbList []AbvDb
}
