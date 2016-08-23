## CSM-MySQL
This is sample implementation of CSM for mysql. This holds the extensions for CSM, UCP's sdl and instance definition for deploying a CSM-MySQL service, and Dockerfile for creating CSM-MySQL container image.

use following command to start the CSM-MySQL container and run it on local (running on the same machine) docker host
```
make all
```

Here are the details of various scripts/extensions you will find in this directory:

### SIDECAR_HOME
This is where all the extensions are located for workspace, connection and setup. All these extensions are shell scripts, you can look at how they implement the contract for Sidecar. They accept necessary command line parameters and write appropriate output the provided file.

### deployment
In this directory you will find two files
1. sdl.json - this the service defintion for deploying CSM-MySQL on UCP, it has two components, first one is mysql and second one for csm-mysql
2. instance.json - this the instance definition for deploying the sdl.json to UCP.

### scripts
This contains helper scripts for running various make commands

### Dockerfile
This is the docker file which is used to create CSM-MySQL container image. This
1. Takes catalog-service-manager:latest as the base image
2. Installs mysql-client on the image so that shell scripts can use that client to talk to remote mysql (running on seperate container)
2. Copies all the extensions from SIDECAR_HOME into the container at appropriate location
3. Sets the necessary environment variables
4. Uses Sidecar's entrypoint, so that when container starts, CSM service is started automatically.
