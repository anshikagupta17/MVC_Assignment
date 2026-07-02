package models

type AttackTroop struct {
	TroopID  int64 `json:"troop_id"`
	Quantity int   `json:"quantity"`
}

type AttackRequest struct {
	DefenderID int64         `json:"defender_id"`
	Troops     []AttackTroop `json:"troops"`
}

type MatchmakingResult struct {
	VillageID     int64             `json:"village_id"`
	TownhallLevel int               `json:"townhall_level"`
	Trophies      int               `json:"trophies"`
	Gold          int               `json:"gold"`
	Elixir        int               `json:"elixir"`
	Buildings     []VillageBuilding `json:"buildings"`
}

type BattleResult struct {
	Stars                 int  `json:"stars"`
	DestructionPercentage int  `json:"destruction"`
	TownhallDestroyed     bool `json:"townhall_destroyed"`

	LootGold   int `json:"loot_gold"`
	LootElixir int `json:"loot_elixir"`

	TrophyChange int `json:"trophy_change"`
}

type OpponentResponse struct {
	VillageID     int64             `json:"village_id"`
	TownhallLevel int               `json:"townhall_level"`
	Trophies      int               `json:"trophies"`
	Gold          int               `json:"gold"`
	Elixir        int               `json:"elixir"`
	Buildings     []VillageBuilding `json:"buildings"`
}
