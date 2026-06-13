package models

import (
	"encoding/json"
	"time"
)

type Village struct {
	ID            int64
	UserId        int64
	TownhallLevel int
	Gold          int
	Elixir        int
	HousingSpace  int
	Trophies      int
	Layout        json.RawMessage
}

type VillageBuilding struct {
	ID            int64
	VillageId     int64
	BuildingId    int64
	Level         int64
	Quantity      int64
	UpgradeEndsAt *time.Time
	X             int64
	Y             int64
	LastUpdate    *time.Time
}

type BuildingMetadata struct {
	ID             int64
	SizeX          int64
	SizeY          int64
	Name           string
	UpgradeCost    int64
	CostType       string
	UpgradeTimeSec int64
}

type VillageResponse struct {
	ID            int64             `json:"id"`
	Gold          int               `json:"gold"`
	Elixir        int               `json:"elixir"`
	TownhallLevel int               `json:"townhall_level"`
	Layout        json.RawMessage   `json:"layout"`
	Buildings     []VillageBuilding `json:"buildings"`
}
