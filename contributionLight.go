package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/amimof/huego"
)


func Hex2RGB(hex string) (color.RGBA, error) {
	var r, g, b uint8

	if len(hex) == 4 {
		fmt.Sscanf(hex, "#%1x%1x%1x", &r, &g, &b)
		r *= 17
		g *= 17
		b *= 17
	} else {
		fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}


type Response struct {
	Data struct {
		Viewer struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks []struct {
						ContributionDays []struct {
							Weekday int
							Date    string
							Color   string
						} `json:"contributionDays"`
					} `json:"weeks"`
					} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"viewer"`
	} `json:"data"`
}

func main(){
	yesterday := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	dateWithTimeSplit := strings.Split(yesterday, "T")[0]

	jsonData := map[string]string{
		"query": fmt.Sprintf(`
		{ 
			viewer { 
			  contributionsCollection(from: "%sT00:00:00" to: "%sT23:59:00") {
				contributionCalendar {
				  weeks {
					contributionDays {
					  weekday
					  date
					  color
					}
				  }
				}
			  }
			}
		  }
		`, 	dateWithTimeSplit, dateWithTimeSplit),
	}

	token := ""

	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonValue))
	request.Header.Set("Authorization", "bearer "+token)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	GitHubResponse := Response{}
	json.Unmarshal(body, &GitHubResponse)

	contributionDay := GitHubResponse.Data.Viewer.ContributionsCollection.ContributionCalendar.Weeks[0].ContributionDays[0]

	fmt.Println(contributionDay)

	bridge := huego.New("192.168.0.2", "aKJL1Pe9TNVGkdYy-I9twebO--Nvx9Q42HPVY-lH")
	l, err := bridge.GetLight(1)

	if err != nil {
		log.Fatal(err)
	}

	l.On()

	rgb, err := Hex2RGB(contributionDay.Color)

	fmt.Printf("%v\n", rgb)

	// q := color.RGBA{155, 233, 168, 150}	// #9be9a8
	// w := color.RGBA{64, 196, 99, 150} // #40c463
	// e := color.RGBA{48, 161, 78, 150} // #30a14e
	// r := color.RGBA{33, 110, 57, 250} // #216e39

	if err != nil {
		fmt.Println(err)
	}

	 l.Col(rgb)
	if err != nil {
		panic(err)
	}
  	fmt.Println(l.Name)
}