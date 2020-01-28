# evidence-management
 A Hyperledger project on Evidence Management.

# Steps for setting up the network

* First things first, need to write the configtx.yaml and crypto-config.yaml files for the network you want to setup.

* Make sure you to include `bin` folder of `fabric-samples` at the same level as the root of this project.

* Once these files are written, now its time to update the template for the docker containers.

  * Make the following changes to the `docker-compose-e2e-template.yaml` file.
    * Include the number of CA's you want with the required configuration.
    * Find for `CA_PRIVATE_KEY` key, and then replace with `OrgName_CA_PRIVATE_KEY` (First letter caps).
    * Also, the name of the CA would be of the form `ca.orgname.example.com` which needs to be changed accordingly.
    * Once the CA names are done, need to do the same for the peers.
    * Along with these, need to edit the file `base/docker-compose.yaml` file accordingly.

* Now, create the folders `crypto-config` and `channel-artifacts` inside the root of this project

* Once all the necessary changes are done, run the following command for generating the crypto material.

  `bash byfn1.sh generate <orglist>`

  (To run this command, we need bash with 4.0+)

* Now get the containers up by running the following command.

  `./runApp.sh <ProjectName>`



