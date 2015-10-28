package controllers

import (
    "io/ioutil"
	"github.com/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
    "strconv"
)


// LocationController represents the controller for operating on the InputAddress resource
type LocationController struct {
		session *mgo.Session
	}

//Specs about IO structs


type InputAddress struct {
		Name   string        `json:"name"`
		Address string 		`json:"address"`
		City string			`json:"city"`
		State string		`json:"state"`
		Zip string			`json:"zip"`
	}



type OutputAddress struct {

		Id     bson.ObjectId `json:"_id" bson:"_id,omitempty"`
		Name   string        `json:"name"`
		Address string 		`json:"address"`
		City string			`json:"city" `
		State string		`json:"state"`
		Zip string			`json:"zip"`

		Coordinate struct{
			Lat string 		`json:"lat"`
			Lang string 	`json:"lang"`
		}
	}

//struct for google response

type GoogleResponse struct {
	Results []GoogleResult
}

type GoogleResult struct {

	Address      string               `json:"formatted_address"`
	AddressParts []GoogleAddressPart `json:"address_components"`
	Geometry     Geometry
	Types        []string
}

type GoogleAddressPart struct {

	Name      string `json:"long_name"`
	ShortName string `json:"short_name"`
	Types     []string
}

type Geometry struct {

	Bounds   Bounds
	Location Point
	Type     string
	Viewport Bounds
}

type Bounds struct {
	NorthEast, SouthWest Point
}

type Point struct {
	Lat float64
	Lng float64
}



// reference to a LocationController for mongo session
func NewLocationController(s *mgo.Session) *LocationController {
	return &LocationController{s}
}

//The func for google's response
func getGoogLocation(address string) OutputAddress{
	client := &http.Client{}

	reqURL := "http://maps.google.com/maps/api/geocode/json?address="
	reqURL += url.QueryEscape(address)
	reqURL += "&sensor=false";
	fmt.Println("URL formed: "+ reqURL)
	req, err := http.NewRequest("GET", reqURL , nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error sending req to google: ", err);	
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response: ", err);	
	}

	var res GoogleResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("error unmashalling response: ", err);	
	}

	var ret OutputAddress
	ret.Coordinate.Lat = strconv.FormatFloat(res.Results[0].Geometry.Location.Lat,'f',7,64)
	ret.Coordinate.Lang = strconv.FormatFloat(res.Results[0].Geometry.Location.Lng,'f',7,64)

	return ret;
}

// GetLocation retrieves one location resource
func (uc LocationController) GetLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("location_id")
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

    
    oid := bson.ObjectIdHex(id)
	var o OutputAddress
	if err := uc.session.DB("go_273").C("Locations").FindId(oid).One(&o); err != nil {
        w.WriteHeader(404)
        return
    }
	// Marshal into JSON structure
	uj, _ := json.Marshal(o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}



//create a new Location resource
func (uc LocationController) CreateLocation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var u InputAddress
	var oA OutputAddress

	json.NewDecoder(r.Body).Decode(&u)	
	googResCoor := getGoogLocation(u.Address + "+" + u.City + "+" + u.State + "+" + u.Zip);
    fmt.Println("resp is: ", googResCoor.Coordinate.Lat, googResCoor.Coordinate.Lang);
	oA.Id = bson.NewObjectId()
	oA.Name = u.Name
	oA.Address = u.Address
	oA.City= u.City
	oA.State= u.State
	oA.Zip = u.Zip
	oA.Coordinate.Lat = googResCoor.Coordinate.Lat
	oA.Coordinate.Lang = googResCoor.Coordinate.Lang

	// Write the user to mongo
	uc.session.DB("go_273").C("Locations").Insert(oA)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(oA)
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// remove an existing location resource
func (uc LocationController) RemoveLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("location_id")
	
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB("go_273").C("Locations").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}

//UpdateLocation updates an existing location resource
func (uc LocationController) UpdateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var i InputAddress
	var o OutputAddress

	id := p.ByName("location_id")
	// fmt.Println(id)
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }
    oid := bson.ObjectIdHex(id)
	
	if err := uc.session.DB("go_273").C("Locations").FindId(oid).One(&o); err != nil {
        w.WriteHeader(404)
        return
    }	

	json.NewDecoder(r.Body).Decode(&i)	
    //Trying to get the lat lang!!!--------------------
	googResCoor := getGoogLocation(i.Address + "+" + i.City + "+" + i.State + "+" + i.Zip);
    fmt.Println("resp is: ", googResCoor.Coordinate.Lat, googResCoor.Coordinate.Lang);

	
	o.Address = i.Address
	o.City = i.City
	o.State = i.State
	o.Zip = i.Zip
	o.Coordinate.Lat = googResCoor.Coordinate.Lat
	o.Coordinate.Lang = googResCoor.Coordinate.Lang

	// Write the user to mongo
	c := uc.session.DB("go_273").C("Locations")
	
	id2 := bson.M{"_id": oid}
	err := c.Update(id2, o)
	if err != nil {
		panic(err)
	}
	
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(o)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}