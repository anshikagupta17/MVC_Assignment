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
