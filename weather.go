package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var apiHTTP string

const head = `<!DOCTYPE html>
<html>
  	<head>
    	<meta charset="utf-8">
		<title>Weather</title>
	</head>
	<body>
	<form action="/weather" method="get">
		<label for="city">City: </label>
		<input type="text" id="city" name="city">
		<input type="submit" value="Search">
	</form>
`

const content = `<div class="name"><h1><b>{{.Name}}</b></h1></div>
	<div class="coord">
	Coordinates:</br>
	Longitude {{.Coord.Lon}}</br>
	Latitude  {{.Coord.Lat}}	
	</div>

	<div class="weather">
	Weather:</br>
	{{$weather := index .Weather 0}}
	ID {{$weather.ID}}</br>
	Main {{$weather.Main}}</br>
	Description {{$weather.Description}}</br>
	Icon {{$weather.Icon}}
	</div>

	<div class="base">Base: {{.Base}}</div>

	<div class="main">
	Main: </br>
	Temperature {{.Main.Temp}} °C</br>
	Feels like {{.Main.FeelsLike}} °C</br>
	Temperature min {{.Main.TempMin}} °C</br>
	Temperature max {{.Main.TempMax}} °C</br>
	Pressure {{.Main.Pressure}} mm Hg\n</br>
	Humidity {{.Main.Humidity}} %</br>
	</div>

	<div class="visibility">Visibility: {{.Visibility}} m</div>

	<div class="wind">
	Wind: </br>
	Speed {{.Wind.Speed}} m/s </br>
	Degree {{.Wind.Deg}} </br>
	</div>

	<div class="rain">Rain: {{.Rain.H}}</div>

	<div class="clouds">Clouds: {{.Clouds.All}}</div>

	<div class="dt">Dt: {{.Dt}}</div>

	<div class="sys">
	Sys:</br>
	Type {{.Sys.Type}}</br>
	ID {{.Sys.ID}}</br>
	Country {{.Sys.Country}}</br>
	Sunrise {{.Sys.Sunrise}}</br>
	Sunset {{.Sys.Sunset}}
	</div>

	<div class="timezone">Timezone: {{.Timezone}}</div>

	<div class="id">ID: {{.ID}}</div>
`
const foot = `</body>
</html>
`

type Coord struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float32 `json:"temp"`
	FeelsLike float32 `json:"feels_like"`
	TempMin   float32 `json:"temp_min"`
	TempMax   float32 `json:"temp_max"`
	Pressure  float32 `json:"pressure"`
	Humidity  float32 `json:"humidity"`
}

type Wind struct {
	Speed int `json:"speed"`
	Deg   int `json:"deg"`
}

type Rain struct {
	H float32 `json:"1h"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

type Response struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Rain       Rain      `json:"rain"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func makeHTTP(name string, key string) string {
	return fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", name, key)
}

func getData(httpStr string) []byte {
	response, err := http.Get(httpStr)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData
}

func informer(w http.ResponseWriter, req *http.Request) {
	var response Response
	var weatherData []byte
	var cityName string

	err := req.ParseForm() //parse args
	check(err)

	for key, values := range req.Form { // range over map
		for _, value := range values { // range over []string
			switch key {
			case "city":
				cityName = value
				break
			default:
				fmt.Fprintf(w, "invalid request")
			}
		}
	}

	fmt.Println(cityName)

	keyAPI := "6b4866b74a0e31e0bd0ccdc1db1de0dc"
	apiHTTP := makeHTTP(cityName, keyAPI)

	weatherData = getData(apiHTTP)

	json.Unmarshal([]byte(weatherData), &response)

	if cityName == "" {
		fmt.Fprintf(w, head)
		fmt.Fprintf(w, foot)
	} else {
		fmt.Fprintf(w, head)
		tmpl, _ := template.New("content").Parse(content)
		tmpl.Execute(w, response)
		fmt.Fprintf(w, foot)

		fmt.Println(response.Weather[0].Description)

	}

	fmt.Println(response.Name)
}

func main() {
	http.HandleFunc("/weather", informer) // http://localhost:8090/weather
	http.ListenAndServe(":8090", nil)

}
