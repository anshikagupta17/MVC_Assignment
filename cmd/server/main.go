package main

import (
	"log"
	"net/http"

	"github.com/anshikagupta17/MVC_Assignment/core/controllers"
	db "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/middleware"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server is running"))
}

func main() {
	db.InitDB()

	http.HandleFunc("/", Home)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/profile", middleware.JWTMiddleware(controllers.Profile))
	http.HandleFunc("/village", middleware.JWTMiddleware(controllers.GetVillage))
	http.HandleFunc("/village/buildings", middleware.JWTMiddleware(controllers.GetVillageBuildings))
	http.Handle("/village/buildings/move", middleware.JWTMiddleware(http.HandlerFunc(controllers.MoveBuilding)))
	http.Handle("/village/collect", middleware.JWTMiddleware(http.HandlerFunc(controllers.ResourceCollection)))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
