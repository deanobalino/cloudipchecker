/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the details of an ip address",
	Long:  `Get the details of an ip address`,
	Run: func(cmd *cobra.Command, args []string) {
		getIPDetails(args)
	},
}

var Source string
var IpAddress string
var SourceParam string

func getIPDetails(args []string) {
	//ipAddress := IpAddress
	if IpAddress != "" {
		if Source == "webjson" {
			SourceParam = "manual"
		} else {
			SourceParam = Source
		}
		fmt.Println("Performing ip lookup for: " + IpAddress + " using data source: " + Source)
		url := "https://cloudipchecker.azurewebsites.net/api/servicetags/" + SourceParam + "?ip=" + IpAddress
		resp, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {

			log.Fatal(err)
		}

		fmt.Println(string(body))
	} else {
		fmt.Println("ERROR: You need to provide an ip address with the --ip flag")
	}

}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//getCmd.PersistentFlags().String("format", "f", "Supports api or webjson (default)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().StringVarP(&Source, "source", "s", "webjson", "Choose which data source to use from 'api' or 'webjson'")
	getCmd.Flags().StringVarP(&IpAddress, "ip", "", "", "The IP Address to lookup, can be IPv4 or IPv6 (required)")
	rootCmd.MarkFlagRequired("ip")
}
