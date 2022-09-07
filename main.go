package main

import (
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
	dbname = "conditions"
)

//TODO function that takes data from API and saves it to a data base
//database of lat/long for dutch coast

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

func windAtLocation(x, y float64) WindConditions {
	start := time.Now()
	startU := start.Unix()
	end := start.Add(time.Hour * 24)
	endU := end.Unix()

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", 52.4907))
	params.Add("lng", fmt.Sprintf("%f", 4.6023))
	params.Add("start", fmt.Sprintf("%d", startU))
	params.Add("end", fmt.Sprintf("%d", endU))
	params.Add("params", "windSpeed")
	url := "https://api.stormglass.io/v2/weather/point?"

	meClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds

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

func swellAtLocation(x, y float64) Waves {
	start := time.Now()
	startU := start.Unix()
	end := start.Add(time.Hour * 24)
	endU := end.Unix()

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", x))
	params.Add("lng", fmt.Sprintf("%f", y))
	params.Add("start", fmt.Sprintf("%d", startU))
	params.Add("end", fmt.Sprintf("%d", endU))
	params.Add("params", "swellHeight")
	url := "https://api.stormglass.io/v2/weather/point?"

	meClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds

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
func main() {
	/*
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

	   	sqlStatement := `
	       INSERT INTO conditions (spot_id, time_stamp, swell, wind)
	       VALUES ($1, $2, $3, $4)` //TODO ogarniecia

	   	_, err = db.Exec(sqlStatement, 1, Time.Wind, Swell, Windspeed)
	   	if err != nil {
	   		panic(err)
	   	}
	*/
	lat := 52.4907
	long := 4.6023
	listSwell := swellAtLocation(lat, long)
	fmt.Println(listSwell)
	listWind := windAtLocation(lat, long)
	fmt.Println(listWind)

}
