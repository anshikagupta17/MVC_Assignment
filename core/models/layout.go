package models

type BuildingData struct {
	TypeId int `json:"type_id"`
	Level  int `json:"level"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type DefenseLayout struct {
	Buildings []BuildingData `json:"buildings"`
}

type ResourcesLayout struct {
	Buildings []BuildingData `json:"buildings"`
}
