module github.com/deanobalino/cloudipchecker/backend

go 1.15

require (
	github.com/Azure/azure-sdk-for-go v51.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/deanobalino/cloud_ip_checker/apiservicetags v0.0.0-00010101000000-000000000000
	github.com/deanobalino/cloud_ip_checker/manualservicetags v0.0.0-00010101000000-000000000000
	github.com/deanobalino/cloud_ip_checker/webservicetags v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0 // indirect
)

replace github.com/deanobalino/cloudipchecker/backend/apiservicetags => ./apiservicetags

replace github.com/deanobalino/cloudipchecker/backend/webservicetags => ./webservicetags
