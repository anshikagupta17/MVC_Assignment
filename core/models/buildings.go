package models

type MoveBuildingRequest struct {
	BuildingInstanceId int64 `json:"building_instance_id"`
	X                  int   `json:"x"`
	Y                  int   `json:"y"`
}

type AddBuildingRequest struct {
	BuildingID int64 `json:"building_id"`
	X          int   `json:"x"`
	Y          int   `json:"y"`
}

type UpgradeBuildingRequest struct {
	BuildingInstanceID int64 `json:"building_instance_id"`
}

type ShopBuilding struct {
	BuildingID   int64  `json:"building_id"`
	Name         string `json:"name"`
	Cost         int    `json:"cost"`
	CostType     string `json:"cost_type"`
	SizeX        int    `json:"size_x"`
	SizeY        int    `json:"size_y"`
	CurrentCount int    `json:"current_count"`
	MaxQuantity  int    `json:"max_quantity"`
}
