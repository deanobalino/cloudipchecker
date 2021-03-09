package webservicetags

import (
	"fmt"
	"io"
	"net/http"
	"log"
	"strings"
	"encoding/json"
	"os"
	"io/ioutil"
	"net"

	"github.com/PuerkitoBio/goquery"
)

var (
	jsonDownloadUrl string
)
type allServiceTagDetails []serviceTagDetail

var serviceTagDetails = allServiceTagDetails{}

type serviceTagDetail struct {
	Region        string `json:"Region"`
	Service       string `json:"Service"`
	AddressPrefix string `json:"AddressPrefix"`
}

type correctRes struct {
	Status int `json:Status"`
	Values interface {} //`json:"Values"`
}

type errorRes struct {
	Status int `json:Status"`
	Error string `json:Error"`
}

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
	ipAddress := r.URL.Query().Get("ip")
	if ipAddress != "" {
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
	v := Values{}
	err = json.Unmarshal(bytes, &v)
	if err != nil {
		panic(err)
	}
	CheckServiceTags(v, ipAddress, w , r)
	if len(serviceTagDetails) > 0 {
		//json.NewEncoder(w).Encode(serviceTagDetails)
		var response correctRes 
		response.Status = 200
		response.Values = serviceTagDetails
		resJSON, _ := json.MarshalIndent(response, "", "    ")
		serviceTagDetails = nil
		fmt.Fprintf(w, string(resJSON))
	} else {
		ReturnError(404, "IP Address not found in the Azure Service Tag JSON downloaded from Azure website", w, r)
	}
	//Respond to the client with prettyJSON
	//resJSON, _ := json.MarshalIndent(v, "", "    ")
	//fmt.Fprintf(w, string(resJSON))
	} else {
		ReturnError(400, "Please pass in an 'ip' query parameter to search for an IP",w ,r)
	}
	
}

func CheckServiceTags(v Values, ipAddress string, w http.ResponseWriter, r *http.Request) {
	for _, value := range v.Values {
		prefixes := value.Properties.AddressPrefixes
		for _, prefix := range prefixes {
			_, cidrA, _ := net.ParseCIDR(prefix)
			ipB := net.ParseIP(ipAddress)
			if cidrA.Contains(ipB) { 
				var newServiceTag serviceTagDetail
				newServiceTag.Region = value.Properties.Region
				newServiceTag.Service = value.Properties.SystemService
				newServiceTag.AddressPrefix = prefix
				serviceTagDetails = append(serviceTagDetails, newServiceTag)
			} 
		}
	}

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

func ReturnError(status int, message string, w http.ResponseWriter, r *http.Request) {
	var errorResponse errorRes
	errorResponse.Status = status
	errorResponse.Error = message
	resJSON, _ := json.MarshalIndent(errorResponse, "", "    ")
	fmt.Fprintf(w, string(resJSON))
} 