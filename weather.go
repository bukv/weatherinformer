package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var clear map[string]func() //create a map for storing clear screen funcs

//Main  in response
type Main struct {
	Temp     float32 `json:"temp"`
	Pressure float32 `json:"pressure"`
	Humidity float32 `json:"humidity"`
}

//Response data
type Response struct {
	Main Main   `json:"main"`
	Name string `json:"name"`
}

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux clear terminal screen
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows clear cmd screen
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//CallClear func for clear screen in cmd or terminal
func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
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

func playSound() error {
	f, err := os.Open("C:\\Users\\user\\Desktop\\1.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func showWeather(link string) {
	var response Response
	var weatherData []byte
	var temperature, pressure, humidity float32
	for {
		weatherData = getData(link)
		json.Unmarshal([]byte(weatherData), &response)

		if temperature != response.Main.Temp || pressure != response.Main.Pressure || humidity != response.Main.Humidity {

			CallClear() //Clear screen

			fmt.Printf("***** %s *****\n", response.Name)
			fmt.Printf("Тemperature = %v ͦ C\n", response.Main.Temp)
			fmt.Printf("Pressure    = %v mm Hg\n", response.Main.Pressure)
			fmt.Printf("Humidity    = %v %s\n", response.Main.Humidity, "%")
			fmt.Println("")
			fmt.Println("Updated: ", time.Now())
			fmt.Println("")
			fmt.Print("[press Ctrl+C to exit]\n")
			fmt.Println("")

			temperature = float32(response.Main.Temp)
			pressure = float32(response.Main.Pressure)
			humidity = float32(response.Main.Humidity)

			if err := playSound(); err != nil {
				log.Fatal(err)
			}

			//uncomment next line if you want see all API data
			//fmt.Println(string(weatherData))
		}
		//waiting 5 min
		time.Sleep(300 * time.Second)
	}
}

func main() {
	var apiHTTP string
	var cityName string
	var keyAPI string

	fmt.Println("Enter city name:")
	fmt.Scan(&cityName)
	fmt.Println("")

	keyAPI = "6b4866b74a0e31e0bd0ccdc1db1de0dc"
	apiHTTP = makeHTTP(cityName, keyAPI)

	showWeather(apiHTTP)
}
