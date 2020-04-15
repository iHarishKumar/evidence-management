# evidence-management
 A DLT based evidence management system wherein the entire life-cycle of a specific crime/case file can be tracked with tamper-proof data. Here, any digital data that is specific to the case, will be hashed and pushed to the blockchain. Hence, reducing the risk of trust and mismanagement in the system. The parties involved are Police, Forensics, Court, Lawyers.

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

* Once the above configuration is done, it's time to setup the Node SDK.

  * Before we start the server, we need to write some config files so that the Node SDK will be able to interact with the Fabric without any problems.
  * Firstly, create a folder `artifacts/ca-config` .
  * Inside this folder, create `n` number of folders with org names. (`n` -> number of orgs)
  * Inside every folder, include `fabric-ca-server-config.yaml` file with the appropriate contents.
  * Also, create `n` number of files inside `artifacts` folder with `<orgname>.yaml` file name. (`n` -> number of orgs)

* Now get the containers up by running the following command.

  `./restart.sh`

* In the other command window, run the testAPI file which includes the steps for until invoking the chaincode.

  `./testAPI.sh`

# Version History

### 1.0 Fetaures

* Encrypt/Decrypt metadata.
* Hash file data.
* CounchDB add on.
* JWT Authentication.
* Metadata extraction using Exiftool-Vendored.



