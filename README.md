# gowac

gowac is a simple web server written in Go with some optimizations for deploying to Azure Web App for Containers

The easiest way to get started is to click the button below. This will create a Web App for Containers that deploys the bartr/gowac container from Docker Hub. You can accomplish the same thing via the portal or command line using the bartr/gowac container.

[![Deploy to Azure](https://azuredeploy.net/deploybutton.svg)](https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fraw.githubusercontent.com%2Fbartr%2Fgowac%2Fmaster%2Fazuredeploy.json)

The docker folder contains developer and release builds that build the gowac container. The developer build contains developer tools and is not suitable for production, but is a great way to explore Web App for Containers using Go.

For more information on Web App for containers, click [here](https://azure.microsoft.com/en-us/services/app-service/containers/)