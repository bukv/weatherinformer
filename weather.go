package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

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

type weatherInfo struct {
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
	IconLink   string
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func makeHTTP(name string, key string) string {
	return fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", name, key)
}

func iconHTTP(iconCode string) string {
	return fmt.Sprintf("http://openweathermap.org/img/wn/%s@2x.png", iconCode)
}

func getData(httpStr string) []byte {
	response, err := http.Get(httpStr)
	check(err)

	responseData, err := ioutil.ReadAll(response.Body)
	check(err)

	return responseData
}

func parseCity(req *http.Request) string {
	var city string

	err := req.ParseForm() //parse args
	check(err)

	for key, values := range req.Form { // range over map
		for _, value := range values { // range over []string
			switch key {
			case "city":
				city = value
				break
			default:
				fmt.Println("invalid request")
			}
		}
	}

	if city == "" {
		city = "london"
	}

	return city
}

func informer(w http.ResponseWriter, req *http.Request) {
	var data weatherInfo
	var weatherData []byte
	var cityName string
	var apiHTTP string
	var keyAPI string

	cityName = parseCity(req)

	keyAPI = "" //enter your API key here
	apiHTTP = makeHTTP(cityName, keyAPI)

	weatherData = getData(apiHTTP)

	json.Unmarshal([]byte(weatherData), &data)

	data.IconLink = iconHTTP(data.Weather[0].Icon)

	tmpl := template.Must(template.ParseFiles("html/index.html"))

	tmpl.Execute(w, data)
}

func main() {
	req := mux.NewRouter()
	req.HandleFunc("/", informer)

	req.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/", http.FileServer(http.Dir("styles/"))))

	http.Handle("/", req) // http://localhost:8090/
	http.ListenAndServe(":8090", nil)
}
