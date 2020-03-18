package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

type buyer_seller struct {
}

type Case_Asset struct {
	Case_Id            string `json:"Case_Id"`
	Case_Name          string `json:"Case_Name"`
	Case_Description   string `json:"Case_Description"`
	Status             string `json:"Status"`
	FIR_Image_URL      string `json:"FIR_Image_URL"`
	FIR_HASH           string `json:"FIR_HASH"`
	FIR_METATDATA_HASH string `json:"FIR_METADATA_HASH"`
	Date_Of_Creation   string `json:"Date_Of_Creation"`
}

type Document_Asset struct {
	Document_Id            string `json:"Document_Id"`
	Document_Description   string `json:"Document_Description"`
	Document_URL           string `json:"Document_URL"`
	Document_HASH          string `json:"Document_HASH"`
	Document_METADATA_HASH string `json:"Document_METADATA_HASH"`
	Case_Id                string `json:"Case_Id"`
}

type FIR_Asset struct {
	FIR_ID            string `json:"FIR_ID"`
	FIR_Description   string `json:"FIR_Description"`
	FIR_URL           string `json:"FIR_URL"`
	FIR_DOC_HASH      string `json:"FIR_DOC_HASH"`
	FIR_METADATA_HASH string `json:"FIR_METADATA_HASH"`
	Case_Id           string `json:"Case_Id"`
}

type CounterNO struct {
	Counter int `json:"counter"`
}

type DocumentNO struct {
	Counter int `json:"counter"`
}

type FIRNo struct {
	Counter int `json:"counter"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(buyer_seller))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}

}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *buyer_seller) Init(APIstub shim.ChaincodeStubInterface) pb.Response {

	// Initializing Case Counter
	CaseCounterBytes, _ := APIstub.GetState("CaseCounterNO")
	if CaseCounterBytes == nil {
		var CaseCounter = CounterNO{Counter: 0}
		CaseCounterBytes, _ := json.Marshal(CaseCounter)
		err := APIstub.PutState("CaseCounterNO", CaseCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to Initiate Case Counter"))
		}
	}

	DocumentCounterBytes, _ := APIstub.GetState("DocumentCounterNo")
	if DocumentCounterBytes == nil {
		var DocumentCounter = DocumentNO{Counter: 0}
		DocumentCounterBytes, _ := json.Marshal(DocumentCounter)
		err := APIstub.PutState("DocumentCounterNo", DocumentCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to Initiate Document Counter"))
		}
	}

	FIRCounterBytes, _ := APIstub.GetState("FIRCounterNo")
	if FIRCounterBytes == nil {
		var FIRCounter = FIRNo{Counter: 0}
		FIRCounterBytes, _ := json.Marshal(FIRCounter)
		err := APIstub.PutState("FIRCounterNo", FIRCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to Initiate FIR Counter"))
		}
	}

	return shim.Success(nil)

}

///Start of Private Function

//getCounter to the latest value of the counter based on the Asset Type provided as input parameter
func getCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	fmt.Sprintf("Counter Current Value %d of Asset Type %s", counterAsset.Counter, AssetType)

	return counterAsset.Counter
}

//incrementCounter to the increase value of the counter based on the Asset Type provided as input parameter by 1
func incrementCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.Counter++
	counterAsBytes, _ = json.Marshal(counterAsset)

	err := APIstub.PutState(AssetType, counterAsBytes)
	if err != nil {

		fmt.Sprintf("Failed to Increment Counter")

	}
	return counterAsset.Counter
}

///End of Private Function

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *buyer_seller) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("function is ==> :" + function)
	action := args[0]
	fmt.Println(" action is ==> :" + action)
	fmt.Println(args)

	if action == "queryAsset" {
		return t.queryAsset(stub, args)
	} else if action == "queryAllAsset" {
		return t.queryAllAsset(stub, args)
	} else if action == "getHistoryForRecord" {
		return t.getHistoryForRecord(stub, args)
	} else if action == "createCase" {
		return t.createCase(stub, args)
	} else if action == "updateCase" {
		return t.updateCase(stub, args)
	} else if action == "updateCaseStatus" {
		return t.updateCaseStatus(stub, args)
	} else if action == "createFIR" {
		return t.createFIR(stub, args)
	} else if action == "createDoc" {
		return t.createDoc(stub, args)
	}

	fmt.Println("invoke did not find func: " + action) //error

	return shim.Error("Received unknown function")
}

// ===== Example: Ad hoc rich query ========================================================
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *buyer_seller) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	queryString := args[1]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	//fmt.Println("GetQueryResultForQueryString() : getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	fmt.Println(resultsIterator)
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		queryResponseStr := string(queryResponse.Value)
		fmt.Println(queryResponseStr)
		buffer.WriteString(queryResponseStr)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	//fmt.Println("GetQueryResultForQueryString(): getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

//E-com

// GetTxTimestampChannel Function gets the Transaction time when the chain code was executed it remains same on all the peers where chaincode executes
func (t *buyer_seller) GetTxTimestampChannel(APIstub shim.ChaincodeStubInterface) (string, error) {
	txTimeAsPtr, err := APIstub.GetTxTimestamp()
	if err != nil {
		fmt.Printf("Returning error in TimeStamp \n")
		return "Error", err
	}
	fmt.Printf("\t returned value from APIstub: %v\n", txTimeAsPtr)
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()

	return timeStr, nil
}

// queryAsset Function gets the assets based on Id provided as input
func (t *buyer_seller) queryAsset(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, Required 2")
	}

	fmt.Println("In Query Asset")

	AssetAsBytes, _ := APIstub.GetState(args[1])

	if AssetAsBytes == nil {
		return shim.Error("Could not locate Asset")

	}

	return shim.Success(AssetAsBytes)
}

// update Case Attributes
func (t *buyer_seller) updateCase(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments, Required 6")
	}

	if len(args[1]) == 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	caseBytes, _ := APIstub.GetState(args[1])

	if caseBytes == nil {
		return shim.Error("Cannot Find Case Asset ")
	}

	caseAsset := Case_Asset{}

	json.Unmarshal(caseBytes, &caseAsset)

	caseAsset.Case_Name = args[2]
	caseAsset.Case_Description = args[3]
	caseAsset.Status = args[4]
	caseAsset.FIR_Image_URL = args[5]
	caseAsset.FIR_HASH = args[6]
	caseAsset.FIR_METATDATA_HASH = args[7]

	comAssetAsBytes, errMarshal := json.Marshal(caseAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error: %s", errMarshal))
	}

	errPut := APIstub.PutState(caseAsset.Case_Id, comAssetAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to Update Case: %s", caseAsset.Case_Id))
	}

	fmt.Println("Success in updating Case Asset %v ", caseAsset)

	return shim.Success(nil)
}

// create Case
func (t *buyer_seller) createCase(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	//To check number of arguments are 7
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments, Required 6 arguments")
	}

	//To check each argument is not null
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(string(i+1) + "st argument must be a non-empty string")
		}
	}

	caseCounter := getCounter(APIstub, "CaseCounterNO")
	caseCounter++

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := t.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		return shim.Error("Returning error in Transaction TimeStamp")
	}

	var comAsset = Case_Asset{Case_Id: "Case" + strconv.Itoa(caseCounter), Case_Name: args[1], Case_Description: args[2], Status: args[3], FIR_Image_URL: args[4], FIR_HASH: args[5], FIR_METATDATA_HASH: args[6], Date_Of_Creation: txTimeAsPtr}

	comAssetAsBytes, errMarshal := json.Marshal(comAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Case: %s", errMarshal))
	}

	errPut := APIstub.PutState(comAsset.Case_Id, comAssetAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create Case Asset: %s", comAsset.Case_Id))
	}

	//TO Increment the Case Counter
	incrementCounter(APIstub, "CaseCounterNO")

	fmt.Println("Success in creating Case Asset %v", comAsset)

	return shim.Success(nil)

}

// Create FIR Asset
func (t *buyer_seller) createFIR(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// To check number of arguments
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments, Required 5 arguments")
	}
	// To check each argument is not null
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(string(i+1) + "st argument mush be a non-empty string")
		}

	}

	caseBytes, _ := APIstub.GetState(args[1])

	if caseBytes == nil {
		return shim.Error("Cannot Find Case Asset ")
	}
	caseAsset := Case_Asset{}

	json.Unmarshal(caseBytes, &caseAsset)
	firCounter := getCounter(APIstub, "FIRCounterNo")
	firCounter++

	firAsset := FIR_Asset{FIR_ID: "FIR" + strconv.Itoa(firCounter), FIR_Description: args[2], FIR_URL: args[3], FIR_DOC_HASH: args[4], FIR_METADATA_HASH: args[5], Case_Id: caseAsset.Case_Id}

	firAssetAsBytes, errMarshal := json.Marshal(firAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in FIR: %s", errMarshal))
	}
	errPut := APIstub.PutState(firAsset.FIR_ID, firAssetAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create FIR Asset: %s", errPut))
	}

	// Update counter state
	incrementCounter(APIstub, "FIRCounterNo")
	fmt.Println("Success in creating FIR Asset %v", firAsset)

	return shim.Success(nil)
}

func (t *buyer_seller) createDoc(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// To check number of arguments
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments, Required 5 arguments")
	}
	// To check each argument is not null
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(string(i+1) + "st argument mush be a non-empty string")
		}

	}

	caseBytes, _ := APIstub.GetState(args[1])

	if caseBytes == nil {
		return shim.Error("Cannot find Case Asset")
	}

	caseAsset := Case_Asset{}

	json.Unmarshal(caseBytes, &caseAsset)
	docCounter := getCounter(APIstub, "DocumentCounterNo")
	docCounter++

	docAsset := Document_Asset{Document_Id: "Document" + strconv.Itoa(docCounter), Document_Description: args[2], Document_URL: args[3], Document_HASH: args[4], Document_METADATA_HASH: args[5], Case_Id: caseAsset.Case_Id}

	docAssetAsBytes, errMarshal := json.Marshal(docAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Document: %s", errMarshal))
	}
	errPut := APIstub.PutState(docAsset.Document_Id, docAssetAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create Document Asset: %s", errPut))
	}

	// Update counter state
	incrementCounter(APIstub, "DocumentCounterNo")
	fmt.Println("Success in creating Document Asset %v", docAsset)

	return shim.Success(nil)
}

// update Case Status
func (t *buyer_seller) updateCaseStatus(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, Required 3")
	}

	if len(args[1]) == 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[2]) == 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	caseBytes, _ := APIstub.GetState(args[1])

	if caseBytes == nil {
		return shim.Error("Cannot Find Case Asset ")
	}

	caseAsset := Case_Asset{}

	json.Unmarshal(caseBytes, &caseAsset)

	caseAsset.Status = args[2]

	updateCaseBytes, errMarshal := json.Marshal(caseAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error: %s", errMarshal))
	}

	errPutCase := APIstub.PutState(caseAsset.Case_Id, updateCaseBytes)

	if errPutCase != nil {
		return shim.Error(fmt.Sprintf("Failed to Update Case: %s", caseAsset.Case_Id))
	}

	return shim.Success(nil)
}

// query all assets
func (t *buyer_seller) queryAllAsset(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	startKey := ""

	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {

		return shim.Error(err.Error())

	}

	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		// respValue := string(queryResponse.Value)
		if err != nil {

			return shim.Error(err.Error())

		}

		// Add a comma before array members, suppress it for the first array member

		if bArrayMemberAlreadyWritten == true {

			buffer.WriteString(",")

		}

		buffer.WriteString("{\"Key\":")

		buffer.WriteString("\"")

		buffer.WriteString(queryResponse.Key)

		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")

		// Record is a JSON object, so we write as-is

		buffer.WriteString(string(queryResponse.Value))

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true

	}

	buffer.WriteString("]")

	fmt.Printf("- queryAllAssets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

// get History For Record
func (t *buyer_seller) getHistoryForRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	recordKey := args[1]

	fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

	resultsIterator, err := stub.GetHistoryForKey(recordKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the key/value pair
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON vehiclePart)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
