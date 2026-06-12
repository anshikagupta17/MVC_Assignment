package models

import (
	"time"
)

type Mine struct {
	ID             int64
	BuildingId     int64
	Level          int
	TypeId         int64
	ProductionRate int
	LastCollected  time.Time
	UpgradeEndsAt  *time.Time
}

const (
	TownHall      = 1
	Cannon        = 2
	ArcherTower   = 3
	Wall          = 4
	Mortar        = 5
	GoldMine      = 6
	ElixirMine    = 7
	GoldStorage   = 8
	ElixirStorage = 9
	ArmyCamp      = 10
)
