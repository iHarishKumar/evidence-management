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

type evidence_management struct {
}

type Case_Asset struct {
	Case_Id               string `json:"Case_Id"`
	Case_Name             string `json:"Case_Name"`
	Case_Description      string `json:"Case_Description"`
	Status                string `json:"Status"`
	DBCOLLECTION          string `json:"DBCOLLECTION"`
	FIR_ENCRYPT           string `json:"FIR_ENCRYPT"`
	FIR_METATDATA_ENCRYPT string `json:"FIR_METADATA_ENCRYPT"`
	Date_Of_Creation      string `json:"Date_Of_Creation"`
	FIR_HASH              string `json:"FIR_HASH"`
	FIR_METADATA_HASH     string `json:"FIR_METADATA_HASH"`
	FIR_DATA_HASH         string `json:"FIR_DATA_HASH"`
}

//Artificat Asset
type Document_Asset struct {
	Document_Id               string `json:"Document_Id"`
	Document_Description      string `json:"Document_Description"`
	DBCOLLECTION              string `json:"DBCOLLECTION"`
	Document_ENCRYPT          string `json:"Document_HASH"`
	Document_METADATA_ENCRYPT string `json:"Document_METADATA_HASH"`
	Case_Id                   string `json:"Case_Id"`
	Document_HASH             string `json:"Document_HASH"`
	Document_METADATA_HASH    string `json:"Document_METADATA_HASH"`
	Document_Data_HASH        string `json:"Document_Data_HASH"`
}

type FIR_Asset struct {
	FIR_ID            string `json:"FIR_ID"`
	FIR_Description   string `json:"FIR_Description"`
	DBCOLLECTION      string `json:"DBCOLLECTION"`
	FIR_DOC_HASH      string `json:"FIR_DOC_HASH"`
	FIR_METADATA_HASH string `json:"FIR_METADATA_HASH"`
	Case_Id           string `json:"Case_Id"`
}

type ACCUSED struct {
	Name             string `json: "Name"`
	ID               string `json: "ID"`
	Photo            string `json: "Photo"`
	Accused_For      string `json: "Accused_For"`
	Term             string `json: "Term"`
	Case_Id          string `json: "Case_Id"`
	Serving_Sentence bool   `json: "Serving_Sentence"`
	On_Bail          bool   `json: "On_Bail"`
	Start_Date       string `json: "Start_Date"`
}

type SUSPECT struct {
	Name        string `json: "Name"`
	ID          string `json: "ID"`
	Reason      string `json: "Reason"`
	Case_Id     string `json: "Case_Id"`
	Photo       string `json: "Photo"`
	Description string `json: "Description"`
	Notes       string `json: "Notes"`
}

type VICTIM struct {
	Name        string `json: "Name"`
	ID          string `json: "ID"`
	Case_Id     string `json: "Case_Id"`
	Photo       string `json: "Photo"`
	IsAlive     bool   `json: "IsAlive"`
	Description string `json: "Description"`
	Report      string `json: "Report"`
}

type CounterNO struct {
	Counter int `json:"counter"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(evidence_management))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}

}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *evidence_management) Init(APIstub shim.ChaincodeStubInterface) pb.Response {

	// Initializing Case Counter
	counter := [...]string{"CaseCounter", "DocumentCounter", "FIRCounter", "PDCounter", "AccusedCounter", "SuspectCounter", "VictimCounter"}
	for _, element := range counter {
		CounterBytes, _ := APIstub.GetState(element)
		if CounterBytes == nil {
			var count = CounterNO{Counter: 0}
			counterBytes, _ := json.Marshal(count)
			err := APIstub.PutState(element+"NO", counterBytes)
			if err != nil {
				return shim.Error(fmt.Sprintf("Failed to initiate %s counter", element))
			}
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
func (t *evidence_management) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
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
	} else if action == "putPrivateData" {
		return t.putPrivateData(stub, args)
	} else if action == "getPrivateData" {
		return t.getPrivateData(stub, args)
	} else if action == "addAccused" {
		return t.addAccused(stub, args)
	} else if action == "addSuspect" {
		return t.addSuspect(stub, args)
	} else if action == "addVictim" {
		return t.addVictim(stub, args)
	}

	fmt.Println("invoke did not find func: " + action) //error

	return shim.Error("Received unknown function")
}

// ===== Example: Ad hoc rich query ========================================================
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *evidence_management) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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
func (t *evidence_management) GetTxTimestampChannel(APIstub shim.ChaincodeStubInterface) (string, error) {
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
func (t *evidence_management) queryAsset(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
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

// Example for putting private data
func (t *evidence_management) putPrivateData(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, Required 2")
	}
	PrDCounter := getCounter(APIstub, "PrDCounterNo")
	PrDCounter++

	PrDCounterBytes, err := json.Marshal(PrDCounter)

	if err != nil {
		return shim.Error("Failed to Marsha")
	}

	err1 := APIstub.PutPrivateData("collectionMedium", "abc"+strconv.Itoa(PrDCounter), PrDCounterBytes)

	// Increment Private Data counter
	incrementCounter(APIstub, "PrDCounterNo")

	if err1 != nil {
		return shim.Error("Something went wrong")
	}
	return shim.Success(nil)
}

func (t *evidence_management) addAccused(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments, Required 7")
	}

	//To check each argument is not null
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(string(i+1) + "st argument must be a non-empty string")
		}
	}

	//err1 := APIstub.PutPrivateData("collectionMedium", "abc"+strconv.Itoa(PrDCounter), PrDCounterBytes)
	CaseBytes, _ := APIstub.GetState(args[5])

	CaseAsset := Case_Asset{}

	json.Unmarshal(CaseBytes, &CaseAsset)

	if CaseBytes == nil {
		return shim.Error("Invalid Case Id")
	}
	AccusedCounter := getCounter(APIstub, "AccusedCounterNO")
	AccusedCounter++

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := t.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		return shim.Error("Returning error in Transaction TimeStamp")
	}
	// Checks are done and the counter is also ready. Create ACCUSED asset.
	AccusedAsset := ACCUSED{ID: CaseAsset.Case_Id + "Accused" + strconv.Itoa(AccusedCounter), Case_Id: CaseAsset.Case_Id, Name: args[2], Photo: args[3], Accused_For: args[4], Term: args[6], Serving_Sentence: true, On_Bail: false, Start_Date: txTimeAsPtr}

	AccusedAssetAsBytes, errMarshal := json.Marshal(AccusedAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Case: %s", errMarshal))
	}

	errPut := APIstub.PutPrivateData("collectionAccused", AccusedAsset.ID, AccusedAssetAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create Case Asset: %s", AccusedAsset.ID))
	}

	//TO Increment the Case Counter
	incrementCounter(APIstub, "AccusedCounterNO")

	fmt.Println("Success in creating Case Asset %v", AccusedAsset)
	return shim.Success(nil)

}
func (t *evidence_management) addSuspect(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments, Required 8")
	}

	//To check each argument is not null
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(string(i+1) + "st argument must be a non-empty string")
		}
	}

	//err1 := APIstub.PutPrivateData("collectionMedium", "abc"+strconv.Itoa(PrDCounter), PrDCounterBytes)
	CaseBytes, _ := APIstub.GetState(args[4])

	CaseAsset := Case_Asset{}

	json.Unmarshal(CaseBytes, &CaseAsset)

	if CaseBytes == nil {
		return shim.Error("Invalid Case Id")
	}
	SuspectCounter := getCounter(APIstub, "SuspectCounterNO")
	SuspectCounter++

	// Checks are done and the counter is also ready. Create Suspect asset.
	SuspectAsset := SUSPECT{ID: CaseAsset.Case_Id + "Suspect" + strconv.Itoa(SuspectCounter), Case_Id: CaseAsset.Case_Id, Name: args[2], Reason: args[3], Photo: args[5], Description: args[6], Notes: args[7]}

	SuspectAssetAsBytes, errMarshal := json.Marshal(SuspectAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Case: %s", errMarshal))
	}

	errPut := APIstub.PutPrivateData("collectionSuspect", SuspectAsset.ID, SuspectAssetAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create Case Asset: %s", SuspectAsset.ID))
	}

	//TO Increment the Case Counter
	incrementCounter(APIstub, "SuspectCounterNO")

	fmt.Println("Success in creating Case Asset %v", SuspectAsset)
	return shim.Success(nil)

}

func (t *evidence_management) addVictim(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments, Required 2")
	}

	//To check each argument is not null
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(string(i+1) + "st argument must be a non-empty string")
		}
	}

	//err1 := APIstub.PutPrivateData("collectionMedium", "abc"+strconv.Itoa(PrDCounter), PrDCounterBytes)
	CaseBytes, _ := APIstub.GetState(args[3])

	CaseAsset := Case_Asset{}

	json.Unmarshal(CaseBytes, &CaseAsset)

	if CaseBytes == nil {
		return shim.Error("Invalid Case Id")
	}
	SuspectCounter := getCounter(APIstub, "VictimCounterNO")
	SuspectCounter++

	val, errVal := strconv.ParseBool(args[5])
	if errVal != nil {
		return shim.Error("Error converting string to bool")
	}
	// Checks are done and the counter is also ready. Create Suspect asset.
	VictimAsset := VICTIM{ID: CaseAsset.Case_Id + "Victim" + strconv.Itoa(SuspectCounter), Case_Id: CaseAsset.Case_Id, Name: args[2], Photo: args[4], IsAlive: val, Description: args[6], Report: args[7]}

	VictimAssetAsBytes, errMarshal := json.Marshal(VictimAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Case: %s", errMarshal))
	}

	errPut := APIstub.PutPrivateData("collectionVictim", VictimAsset.ID, VictimAssetAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create Case Asset: %s", VictimAsset.ID))
	}

	//TO Increment the Case Counter
	incrementCounter(APIstub, "VictimCounterNO")

	fmt.Println("Success in creating Case Asset %v", VictimAsset)
	return shim.Success(nil)

}

func (t *evidence_management) getPrivateData(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, Required 2")
	}

	val, err := APIstub.GetPrivateData(args[2], args[2])

	if err != nil {
		return shim.Error("Error fetching private data from collection `collectionMedium`")
	}

	return shim.Success(val)

}

// update Case Attributes
func (t *evidence_management) updateCase(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments, Required 8")
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

	if caseAsset.Status == "CLOSED" {
		return shim.Error("Cannot update CLOSED case")
	}

	//Reopen case

	caseAsset.Case_Name = args[2]
	caseAsset.Case_Description = args[3]
	caseAsset.Status = args[4]
	caseAsset.DBCOLLECTION = args[5]
	caseAsset.FIR_ENCRYPT = args[6]
	caseAsset.FIR_METATDATA_ENCRYPT = args[7]
	caseAsset.FIR_HASH = args[8]
	caseAsset.FIR_METADATA_HASH = args[9]

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
func (t *evidence_management) createCase(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	//To check number of arguments are 7
	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments, Required 10 arguments")
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
	var status = args[3]
	if status != "CREATED" {
		status = "CREATED"
	}
	var comAsset = Case_Asset{Case_Id: "Case" + strconv.Itoa(caseCounter), Case_Name: args[1], Case_Description: args[2], Status: status, DBCOLLECTION: args[4], FIR_ENCRYPT: args[5], FIR_METATDATA_ENCRYPT: args[6], Date_Of_Creation: txTimeAsPtr, FIR_HASH: args[7], FIR_METADATA_HASH: args[8], FIR_DATA_HASH: args[9]}

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
func (t *evidence_management) createFIR(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	firCounter := getCounter(APIstub, "FIRCounterNO")
	firCounter++

	firAsset := FIR_Asset{FIR_ID: "FIR" + strconv.Itoa(firCounter), FIR_Description: args[2], DBCOLLECTION: args[3], FIR_DOC_HASH: args[4], FIR_METADATA_HASH: args[5], Case_Id: caseAsset.Case_Id}

	firAssetAsBytes, errMarshal := json.Marshal(firAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in FIR: %s", errMarshal))
	}
	errPut := APIstub.PutState(firAsset.FIR_ID, firAssetAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create FIR Asset: %s", errPut))
	}

	// Update counter state
	incrementCounter(APIstub, "FIRCounterNO")
	fmt.Println("Success in creating FIR Asset %v", firAsset)

	return shim.Success(nil)
}

func (t *evidence_management) createDoc(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// To check number of arguments
	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments, Required 9 arguments")
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
	docCounter := getCounter(APIstub, "DocumentCounterNO")
	docCounter++

	docAsset := Document_Asset{Document_Id: "Document" + strconv.Itoa(docCounter), Document_Description: args[2], DBCOLLECTION: args[3], Document_ENCRYPT: args[4], Document_METADATA_ENCRYPT: args[5], Case_Id: caseAsset.Case_Id, Document_HASH: args[6], Document_METADATA_HASH: args[7], Document_Data_HASH: args[8]}

	docAssetAsBytes, errMarshal := json.Marshal(docAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Document: %s", errMarshal))
	}
	errPut := APIstub.PutState(docAsset.Document_Id, docAssetAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create Document Asset: %s", errPut))
	}

	// Update counter state
	incrementCounter(APIstub, "DocumentCounterNO")
	fmt.Println("Success in creating Document Asset %v", docAsset)

	return shim.Success(nil)
}

// update Case Status
func (t *evidence_management) updateCaseStatus(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

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

	if caseAsset.Status == "CLOSED" {
		return shim.Error("Case is already CLOSED")
	}

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
func (t *evidence_management) queryAllAsset(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

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
func (t *evidence_management) getHistoryForRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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
