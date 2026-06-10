package models

import (
	"encoding/json"
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
	VillageId     int64
	BuildingId    int64
	Level         int64
	Quantity      int64
	UpgradeEndsAt int64
	X             int64
	Y             int64
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
