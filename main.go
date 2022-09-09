package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "surf_spots"
)

// TODO function that takes data from API and saves it to a data base
// database of lat/long for dutch coast
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

		fmt.Println("Successfully connected!")

		for i, u := range listW {
			s := listS[i]
			if s.Time == u.Time {
				sqlStatement := `
    INSERT INTO conditions (spot_id, time_stamp, swell, wind)
    VALUES ($1, $2, $3, $4)`
				_, err := db.Exec(sqlStatement, v.Id, u.Time, s.SwellHeight.Icon, u.WindSpeed.Icon)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
func main() {

	listSpots := getLocation()
	fmt.Println(listSpots)
	populateConditions(listSpots)

}
