# m4 - docker

## Structure

There are two sets of docker build files in this solution

* debug - builds an image suitable for developers to connect to the running container via SSH and debug issues while running in the App Services environment. It includes the full Go environment plus additional developer tools. It is not suitable for production.
  * note that App Services for Containers provides a reverse proxy for the ssh connection and it is only accessible via the Azure Portal using Azure subscription credentials.
* release - release builds a very small container with the bare minimum to run the app and is suitable for production environments

## App Services for containers

Both docker builds are optimized for running in App Services for Containers. The containers will also work in aks, k8s or other orchestrators, but they depend on having a reverse proxy in front of them for SSL offloading, firewall, etc.
