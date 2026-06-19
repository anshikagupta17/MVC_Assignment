package models

import (
	"encoding/json"
	"time"
)

type Village struct {
	ID            int64           `json:"id"`
	UserId        int64           `json:"user_id"`
	TownhallLevel int             `json:"townhall_level"`
	Gold          int             `json:"gold"`
	Elixir        int             `json:"elixir"`
	HousingSpace  int             `json:"housing_space"`
	Trophies      int             `json:"trophies"`
	Layout        json.RawMessage `json:"layout"`
}

type VillageBuilding struct {
	ID            int64      `json:"id"`
	VillageId     int64      `json:"village_id"`
	BuildingId    int64      `json:"building_id"`
	Level         int64      `json:"level"`
	UpgradeEndsAt *time.Time `json:"upgrade_ends_at"`
	X             int64      `json:"x"`
	Y             int64      `json:"y"`
	LastUpdate    *time.Time `json:"last_update"`
	SizeX         int        `json:"size_x"`
	SizeY         int        `json:"size_y"`
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
