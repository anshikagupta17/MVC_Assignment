package models

type VillageTroop struct {
	TroopID   int64  `json:"troop_id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	Quantity  int    `json:"quantity"`
	Damage    int    `json:"damage"`
	MaxHealth int    `json:"max_health"`
}

type UpgradeTroopRequest struct {
	TroopID int `json:"troop_id"`
}
