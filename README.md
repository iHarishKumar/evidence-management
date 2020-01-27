# evidence-management
 A Hyperledger project on Evidence Management.

# Follow the steps given below for setting up the network with 'N' number of organisations

* First things first, need to write the configtx.yaml and crypto-config.yaml files for the network you want to setup.
* Once these files are written, now its time to update the template for the docker containers.
  * Make the following changes to the docker-compose-e2e-template.yaml file.
    * Include the number of CA's you want with the required configuration.
    * Find for *CA_PRIVATE_KEY* key, and then replace with *OrgName_CA_PRIVATE_KEY* (First letter caps).
    * Also, the name of the ca would be of the form *ca.orgname.example.com* 

