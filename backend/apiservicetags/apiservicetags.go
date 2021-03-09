package apiservicetags

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

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

type allServiceTagDetails []serviceTagDetail

var serviceTagDetails = allServiceTagDetails{}


var (
	
	location = "westeurope"
	//ipAddress = "1.1.1.1"
	//ipAddress = "13.77.53.49"
)



// Initialize Azure Auth
func AzureAuth() autorest.Authorizer {
	auth, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Println("There was an error authenticating with the Azure CLI credentials")
		panic(err)
	}
	return auth
}

//GetServiceTags from Azure API
func GetServiceTagsFromAPI(subID string, location string, ipAddress string) {
		// establish a service tags client
	
	serviceTagClient := network.NewServiceTagsClient(subID)
	//check if auth is detected
	if AzureAuth() == nil {
		log.Panicln("No Azure CLI auth detected")
	}
	//Auth the client
	serviceTagClient.Authorizer = AzureAuth()
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("Authorized Client with Azure using CLI credentials")
	fmt.Println("-------------------------------------------------------------")

	//Call the Azure API to get the tags
	listTags, err := serviceTagClient.List(context.Background(), location)
	if err != nil {
		log.Println(err)
	} else {
		//We have the service ta data sucessfully
		fmt.Println("-------------------------------------------------------------")
		fmt.Println("Retrieved Service Tag Information from the Azure API")
		fmt.Println("-------------------------------------------------------------")
		//place the values into a variable
		values := *listTags.Values
		//Loop through the values, adding them to a map.
		for _, value := range values {
			prefixes := value.Properties.AddressPrefixes
			for _, prefix := range *prefixes {
				_, cidrA, _ := net.ParseCIDR(prefix)
				ipB := net.ParseIP(ipAddress)
				if cidrA.Contains(ipB) { //&& *value.Properties.Region != "" {
					var newServiceTag serviceTagDetail
					newServiceTag.Region = *value.Properties.Region
					newServiceTag.Service = *value.Properties.SystemService
					newServiceTag.AddressPrefix = prefix
					serviceTagDetails = append(serviceTagDetails, newServiceTag)
				} 
			}
		}
	}
}

func GetServiceTags(w http.ResponseWriter, r *http.Request) {
	ipAddress := r.URL.Query().Get("ip")
	if ipAddress != "" {
		subID := os.Getenv("SUBSCRIPTION_ID")
		fmt.Println(subID)
		GetServiceTagsFromAPI(subID, location, ipAddress)
		if len(serviceTagDetails) > 0 {
			//json.NewEncoder(w).Encode(serviceTagDetails)
			var response correctRes 
			response.Status = 200
			response.Values = serviceTagDetails
			resJSON, _ := json.MarshalIndent(response, "", "    ")
			serviceTagDetails = nil
			fmt.Fprintf(w, string(resJSON))
		} else {
			ReturnError(404, "IP Address not found in Azure Service Tag data returned from the Service Tag API. Try '/api/servicetags/manual?ip=' which may have more data", w, r)
		}
	} else {
		ReturnError(400, "Please pass in an 'ip' query parameter to search for an IP",w ,r)
	}
}

func ReturnError(status int, message string, w http.ResponseWriter, r *http.Request) {
	var errorResponse errorRes
	errorResponse.Status = status
	errorResponse.Error = message
	resJSON, _ := json.MarshalIndent(errorResponse, "", "    ")
	fmt.Fprintf(w, string(resJSON))
} 