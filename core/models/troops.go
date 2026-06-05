package models

type TroopData struct {
	TypeId   int `json:"type_id"`
	Level    int `json:"level"`
	Quantity int `json:"quantity"`
}

type Troop struct {
	Troops []TroopData `json:"troops"`
}
