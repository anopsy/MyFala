package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "surf_spots"
)

// TODO make function that checks connection do database
// to make getLocation and hetSurfable and populateConditions more readable
// TODO make function that connects to API to make windAtLocation and swellAtLocation more readable
type Location struct {
	Id   int
	Name string
	Lat  string
	Long string
}

func getLocation() []Location {

	listLocation := make([]Location, 0)
	//database query for location
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	rows, err := db.Query("SELECT id, name, lat, long FROM surf_spots")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var spot Location
		err = rows.Scan(&spot.Id, &spot.Name, &spot.Lat, &spot.Long)
		if err != nil {
			panic(err)
		}

		listLocation = append(listLocation, spot)

	}

	err = rows.Err()
	if err != nil {
		panic(err)

	}

	return listLocation

}

type Meta struct {
	Cost         int
	DailyQuota   int
	End          string
	Lat          float64
	Lng          float64
	Params       []string
	RequestCount int
	Start        string
}
type Wind struct {
	WindSpeed WindSpeed
	Time      string
}

type WindSpeed struct {
	Icon float64
	Noaa float64
	Sg   float64
}

type WindConditions struct {
	Hours []Wind
	Meta  Meta
}

func windAtLocation(x, y string) WindConditions {
	start := time.Now()
	startU := start.Unix()
	end := start.Add(time.Hour * 24)
	endU := end.Unix()

	params := url.Values{}
	params.Add("lat", x)
	params.Add("lng", y)
	params.Add("start", fmt.Sprintf("%d", startU))
	params.Add("end", fmt.Sprintf("%d", endU))
	params.Add("params", "windSpeed")
	url := "https://api.stormglass.io/v2/weather/point?"

	meClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds

	}
	req, err := http.NewRequest(http.MethodGet, url+params.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "4c4d2d92-3050-11ed-b970-0242ac130002-4c4d2e00-3050-11ed-b970-0242ac130002")

	res, getErr := meClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	windCond := WindConditions{}
	json.Unmarshal([]byte(body), &windCond)
	return windCond

}

type Swell struct {
	SwellHeight SwellHeight
	Time        string
}

type SwellHeight struct {
	Dwd   float64
	Icon  float64
	Meteo float64
	Noaa  float64
	Sg    float64
}

type Waves struct {
	Hours []Swell
	Meta  Meta
}

func swellAtLocation(x, y string) Waves {
	start := time.Now()
	startU := start.Unix()
	end := start.Add(time.Hour * 24)
	endU := end.Unix()

	params := url.Values{}
	params.Add("lat", x)
	params.Add("lng", y)
	params.Add("start", fmt.Sprintf("%d", startU))
	params.Add("end", fmt.Sprintf("%d", endU))
	params.Add("params", "swellHeight")
	url := "https://api.stormglass.io/v2/weather/point?"

	meClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds

	}
	req, err := http.NewRequest(http.MethodGet, url+params.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "f691cda0-015f-11ed-9a2a-0242ac130002-f691ce54-015f-11ed-9a2a-0242ac130002")

	res, getErr := meClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	waves := Waves{}
	json.Unmarshal([]byte(body), &waves)
	return waves
}

func populateConditions(list []Location) {
	for _, v := range list {
		listSwell := swellAtLocation(v.Lat, v.Long)
		listWind := windAtLocation(v.Lat, v.Long)
		listW := listWind.Hours
		listS := listSwell.Hours

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"dbname=%s sslmode=disable",
			host, port, user, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			panic(err)
		}

		for i, u := range listW {
			s := listS[i]
			if s.Time == u.Time {
				sqlStatement := `
    INSERT INTO surfspot_conditions (spot_id, name, time_stamp, swell, wind, surfable)
    VALUES ($1, $2, $3, $4, $5, $6)`
				isSurf := IsSurfable(s.SwellHeight.Icon, u.WindSpeed.Icon)
				_, err := db.Exec(sqlStatement, v.Id, v.Name, u.Time, s.SwellHeight.Icon, u.WindSpeed.Icon, isSurf)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func IsSurfable(s, w float64) bool {
	if s > 0.4 && w < 40.0 {
		return true
	}
	return false
}

type Surfable struct {
	Id       int
	Spot_id  int
	Name     string
	Time     time.Time
	Swell    float64
	Wind     float64
	Surfable bool
}

func getSurfable() []Surfable {
	listSurfable := make([]Surfable, 0)
	//database query for location
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	rows, err := db.Query("SELECT id, spot_id, name, time_stamp, swell, wind, surfable FROM surfspot_conditions where surfable = 't'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var entry Surfable
		err = rows.Scan(&entry.Id, &entry.Spot_id, &entry.Name, &entry.Time, &entry.Swell, &entry.Wind, &entry.Surfable)
		if err != nil {
			panic(err)
		}

		listSurfable = append(listSurfable, entry)

	}

	err = rows.Err()
	if err != nil {
		panic(err)

	}

	return listSurfable

}

type ClosestSpot struct {
	Id       int
	Name     string
	Lat      string
	Long     string
	Distance float64
	Time     time.Time
	Swell    float64
	Wind     float64
}

func calculateDistance(point1, point2 string, userLat, userLong float64) float64 {
	const PI float64 = 3.141592653589793
	p1, err := strconv.ParseFloat(point1, 64)
	if err != nil {
		panic(err)
	}
	p2, err := strconv.ParseFloat(point2, 64)
	if err != nil {
		panic(err)
	}

	radlat1 := float64(PI * p1 / 180)
	radlat2 := float64(PI * userLat / 180)

	theta := float64(p2 - userLong)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515 * 1.609344

	return dist
}

func listDistance(lat, long float64) []ClosestSpot {
	listSpots := getLocation()
	listSurf := getSurfable()
	cs := make([]ClosestSpot, 0)

	for _, w := range listSurf {
		for _, v := range listSpots {
			var spot ClosestSpot
			if v.Id == w.Spot_id {
				spot.Id = v.Id
				spot.Name = v.Name
				spot.Lat = v.Lat
				spot.Long = v.Long
				spot.Distance = calculateDistance(v.Lat, v.Long, lat, long)
				spot.Time = w.Time
				spot.Swell = w.Swell
				spot.Wind = w.Wind

				cs = append(cs, spot)
			}

		}
	}

	sort.Slice(cs, func(i, j int) bool { return cs[i].Distance < cs[j].Distance })
	return cs

}

func handlerFunc(w http.ResponseWriter, r *http.Request) {

	//listSpots := getLocation()
	//populateConditions(listSpots)
	myLat := 52.3802017
	myLong := 4.9121986
	fmt.Fprint(w, "<h1>Your closest surf spot is :</h1>")
	distance := listDistance(myLat, myLong)

	//fmt.Println(listSpots)
	//fmt.Println(listSurf)

	for _, v := range distance {

		fmt.Fprintf(w, "<h2>%v\n</h3>", v.Name)
		fmt.Fprintf(w, "<h3>It's %.2f km from you \n</h3>", v.Distance)
		fmt.Fprintf(w, "<h3>At %v  \n</h3>", v.Time)
		fmt.Fprintf(w, "<h3>The waves there are %v m\n</h3>", v.Swell)
		fmt.Fprintf(w, "<h3>The wind is %v\n</h3>", v.Wind)
		fmt.Println("     ")
		fmt.Println("     ")
		fmt.Println("     ")
	}

}

func main() {
	

	http.HandleFunc("/", handlerFunc)

	fmt.Println("Starting the server on:2137")
	http.ListenAndServe(":2137", nil)

}




