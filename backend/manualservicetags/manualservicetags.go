package manualservicetags

import (
	"fmt"
	"io"
	"net/http"
	"log"
	"strings"
	"encoding/json"
	"os"
	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
)

var (
	jsonDownloadUrl string
)

type Values struct {
	ChangeNumber int64  `json:"changeNumber"`
	Cloud        string `json:"cloud"`
	Values       []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Properties struct {
			AddressPrefixes []string `json:"addressPrefixes"`
			ChangeNumber    int64    `json:"changeNumber"`
			NetworkFeatures []string `json:"networkFeatures"`
			Platform        string   `json:"platform"`
			Region          string   `json:"region"`
			RegionID        int64    `json:"regionId"`
			SystemService   string   `json:"systemService"`
		} `json:"properties"`
	} `json:"values"`
}
type ServiceTag struct {
	SystemService string `json:"Service"`
	Region string `json:"Region"`
	AddressPrefix string `json:"AddressPrefix"`
}

func GetServiceTags(w http.ResponseWriter, r *http.Request) {
	doc, err := goquery.NewDocument("https://www.microsoft.com/en-us/download/confirmation.aspx?id=56519")
	if err != nil {
		log.Fatal(err)
	} else {
	}
	//Get all Links in the HTML file
	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		//Find the link to the JSON download
		if strings.Contains(link, "ServiceTags_Public") {
			jsonDownloadUrl = link
		}
	})
	//Download the file and store it locally
	DownloadFile("webservicetags/service-tags-download.json", jsonDownloadUrl)
	if err != nil {
		panic(err)
	}

	//Read the JSON file as a variable
	serviceTagJson, err := ioutil.ReadFile("webservicetags/service-tags-download.json")
	if err != nil {
		fmt.Println(err)
	}
	//Convert the Bytes to a string
	strJson := string(serviceTagJson)

	//Marshal the JSON 
	rawData := json.RawMessage(strJson)
	bytes, err := rawData.MarshalJSON()
	if err != nil {
		panic(err)
	}

	//Create the struct of Values 'v'
	var v Values
	err = json.Unmarshal(bytes, &v)
	if err != nil {
		panic(err)
	}
	//Respond to the client with prettyJSON
	resJSON, _ := json.MarshalIndent(v, "", "    ")
	fmt.Fprintf(w, string(resJSON))
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}