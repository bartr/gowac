# gowac - docker

## Structure

There are two sets of docker build files in this solution

* debug - builds an image suitable for developers to connect to the running container via SSH and debug issues while running in the Web Apps for Containers environment. It includes the full Go environment plus additional developer tools. It is not suitable for production.
  * note that Web Apps for Containers provides a reverse proxy for the ssh connection and it is only accessible via the Azure Portal using Azure subscription credentials.
* release - release builds a very small container with the bare minimum to run the app and is suitable for production environments

## Web Apps for Containers

Both docker builds are optimized for running in Web Apps for Containers. The containers will also work in aks, k8s or other orchestrators, but they depend on having a reverse proxy for SSL offloading, firewall, etc.
