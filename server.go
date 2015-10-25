package main

import (
	// Standard library packages
	"net/http"

	// Third party packages
	"github.com/httprouter"
	"github.com/minervad/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Get a LocationController instance
	uc := controllers.NewLocationController(getSession())

	// Get a location resource
	r.GET("/locations/:location_id", uc.GetLocation)

	// Create a new address
	r.POST("/locations", uc.CreateLocation)

	// Update an address
	r.PUT("/locations/:location_id", uc.UpdateLocation)

	// Remove an existing address
	r.DELETE("/locations/:location_id", uc.RemoveLocation)

	// Fire up the server
	http.ListenAndServe("localhost:8080", r)
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://admin:admin@ds045464.mongolab.com:45464/go_273")
	//mongodb://minerva:amma123lab@ds045464.mongolab.com:41164/go_273

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	
	s.SetMode(mgo.Monotonic, true)
	// Deliver session
	return s
}