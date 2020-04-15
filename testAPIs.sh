#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi

starttime=$(date +%s)

# Print the usage message
function printHelp () {
  echo "Usage: "
  echo "  ./testAPIs.sh -l golang|node"
  echo "    -l <language> - chaincode language (defaults to \"golang\")"
}
# Language defaults to "golang"
LANGUAGE="golang"

# Parse commandline args
while getopts "h?l:" opt; do
  case "$opt" in
    h|\?)
      printHelp
      exit 0
    ;;
    l)  LANGUAGE=$OPTARG
    ;;
  esac
done

##set chaincode path
function setChaincodePath(){
	LANGUAGE=`echo "$LANGUAGE" | tr '[:upper:]' '[:lower:]'`
	case "$LANGUAGE" in
		"golang")
		CC_SRC_PATH="chaincode"
		;;
		"node")
		CC_SRC_PATH="$PWD/chaincode"
		;;
		*) printf "\n ------ Language $LANGUAGE is not supported yet ------\n"$
		exit 1
	esac
}

setChaincodePath

#---------Enroll Users----------------
#---------Start-----------------------
echo "POST request Enroll on Police  ..."
echo
POLICE_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=Police')

echo $POLICE_TOKEN
POLICE_TOKEN=$(echo $POLICE_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "Police token is $POLICE_TOKEN"
echo
echo "POST request Enroll on Forensics ..."
echo
FORENSICS_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=Forensics')
echo $FORENSICS_TOKEN
FORENSICS_TOKEN=$(echo $FORENSICS_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "Forensics token is $FORENSICS_TOKEN"
echo
echo
echo "POST request Enroll on Lawyers ..."
echo
LAWYERS_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=Lawyers')
echo $LAWYERS_TOKEN
LAWYERS_TOKEN=$(echo $LAWYERS_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "Lawyers token is $LAWYERS_TOKEN"
echo
echo
echo "POST request Enroll on Court ..."
echo
COURT_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=Court')
echo $COURT_TOKEN
COURT_TOKEN=$(echo $COURT_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "Court token is $COURT_TOKEN"
echo
echo
#---------Enroll Users----------------
#---------End-------------------------

#---------Create Channel--------------
#---------Start-----------------------

echo "POST request Create channel  ..."
echo
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization: Bearer $POLICE_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"channelName":"mychannel",
	"channelConfigPath":"../channel-artifacts/mychannel.tx"
}'
echo
echo
#---------Create Channel--------------
#---------End-------------------------

#---------Join Channel----------------
#---------Start-----------------------

sleep 5
echo "POST request Join channel on Police"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $POLICE_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.police.example.com","peer1.police.example.com"]
}'
echo
echo

echo "POST request Join channel on Forensics"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $FORENSICS_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.forensics.example.com","peer1.forensics.example.com"]
}'
echo
echo

echo "POST request Join channel on Lawyers"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $LAWYERS_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.lawyers.example.com","peer1.lawyers.example.com"]
}'
echo
echo

echo "POST request Join channel on Court"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $COURT_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.court.example.com","peer1.court.example.com"]
}'
echo
echo

#---------Join Channel----------------
#---------End-------------------------

#---------Install Channel-------------
#---------Start-----------------------

echo "POST Install chaincode on Police"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $POLICE_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.police.example.com\",\"peer1.police.example.com\"],
	\"chaincodeName\":\"mycc\",
	\"chaincodePath\":\"$CC_SRC_PATH\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"v1\"
}"
echo
echo

echo "POST Install chaincode on Forensics"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $FORENSICS_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.forensics.example.com\",\"peer1.forensics.example.com\"],
	\"chaincodeName\":\"mycc\",
	\"chaincodePath\":\"$CC_SRC_PATH\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"v1\"
}"
echo
echo

echo "POST Install chaincode on Lawyers"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $LAWYERS_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.lawyers.example.com\",\"peer1.lawyers.example.com\"],
	\"chaincodeName\":\"mycc\",
	\"chaincodePath\":\"$CC_SRC_PATH\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"v1\"
}"
echo
echo

echo "POST Install chaincode on Court"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $COURT_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.court.example.com\",\"peer1.court.example.com\"],
	\"chaincodeName\":\"mycc\",
	\"chaincodePath\":\"$CC_SRC_PATH\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"v1\"
}"
echo
echo

#---------Install Channel-------------
#---------End-------------------------

#---------Instantiate Channel---------
#---------Start-----------------------

echo "POST instantiate chaincode on POLICE"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes \
  -H "authorization: Bearer $POLICE_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["peer0.police.example.com"],
	"chaincodeName":"mycc",
	"chaincodeVersion":"v1",
	"chaincodeType": "$LANGUAGE",
	"args":["a","100","b","200"]
}'

echo
echo


echo "POST invoke chaincode on peers of Forensics, Police, Court and Lawyers"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc?key=abc \
  -H "authorization: Bearer $POLICE_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.police.example.com"],
	"fcn":"invoke",
  "operation":"createCase",
	"args": ["createCase","Case-Name", "Case - Description", "CREATED", "/Users/harishgunjalli/Desktop/2018-19/Screenshot 2019-08-27 at 1.06.57 PM.png"]
}'
echo
echo


echo "Total execution time : $(($(date +%s)-starttime)) secs ..."
