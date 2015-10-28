package main

import (
	"net/http"
	"github.com/httprouter"
	"github.com/minervad/controller"
	"gopkg.in/mgo.v2"
)

func main() {
	//instantiate location controller(mongo session) and httprouter
	hr := httprouter.New()
	loc := controllers.NewLocationController(getSession())

	// Get(read) a location resource
	hr.GET("/locations/:location_id", loc.GetLocation)

	// Create a new address
	hr.POST("/locations", loc.CreateLocation)

	// Update an existinf address
	hr.PUT("/locations/:location_id", loc.UpdateLocation)

	// Remove an existing address
	hr.DELETE("/locations/:location_id", loc.RemoveLocation)

	// start server
	http.ListenAndServe("localhost:8080", hr)
}

// create new mongo session
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://admin:admin@ds045464.mongolab.com:45464/go_273")
	//mongodb://minerva:amma123lab@ds045464.mongolab.com:41164/go_273

	if err != nil {
		panic(err)
	}
	
	s.SetMode(mgo.Monotonic, true)
	return s
}