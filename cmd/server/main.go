package main

import (
	"log"
	"net/http"

	"github.com/anshikagupta17/MVC_Assignment/core/controllers"
	DB "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/middleware"
	"github.com/anshikagupta17/MVC_Assignment/db"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server is running"))
}

func Routes() {

	http.HandleFunc("/", Home)

	// Auth
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/profile", middleware.JWTMiddleware(controllers.Profile))

	// Village
	http.HandleFunc("/village", middleware.JWTMiddleware(controllers.GetVillage))
	http.HandleFunc("/village/state", middleware.JWTMiddleware(http.HandlerFunc(controllers.VillageState)))

	// Buildings
	http.HandleFunc("/village/buildings", middleware.JWTMiddleware(controllers.GetVillageBuildings))
	http.Handle("/village/buildings/move", middleware.JWTMiddleware(http.HandlerFunc(controllers.MoveBuilding)))
	http.Handle("/village/buildings/build", middleware.JWTMiddleware(http.HandlerFunc(controllers.AddBuilding)))
	http.Handle("/village/buildings/upgrade", middleware.JWTMiddleware(http.HandlerFunc(controllers.UpgradeBuilding)))

	// Resources
	http.Handle("/village/collect", middleware.JWTMiddleware(http.HandlerFunc(controllers.ResourceCollection)))

	// Troops
	http.Handle("/village/troops", middleware.JWTMiddleware(http.HandlerFunc(controllers.GetVillageTroops)))
	http.Handle("/village/troops/train", middleware.JWTMiddleware(http.HandlerFunc(controllers.TrainTroops)))
	http.Handle("/village/troops/upgrade", middleware.JWTMiddleware(http.HandlerFunc(controllers.UpgradeTroops)))

	// Battle
	http.Handle("/battle/matchmake", middleware.JWTMiddleware(http.HandlerFunc(controllers.FindOpponent)))
	http.Handle("/battle/attack", middleware.JWTMiddleware(http.HandlerFunc(controllers.AttackVillage)))
}

func main() {

	DB.InitDB()
	if err := db.SeedAll(DB.Conn); err != nil {
		log.Fatal(err)
	}

	Routes()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
