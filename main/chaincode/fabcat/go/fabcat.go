package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Cat struct {
	Name  string
	Type  string
	Color string
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCat" {
		return s.queryCat(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createCat" {
		return s.createCat(APIstub, args)
	} else if function == "queryAllCats" {
		return s.queryAllCats(APIstub)
	} else if function == "changeCatName" {
		return s.changeCatName(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCat(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	catAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(catAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	cats := []Cat{
		Cat{Name: "Da Hua", Type: "American Shorthair", Color: "White"},
		Cat{Name: "Xiao Huang", Type: "Ragdoll", Color: "Yellow"},
		Cat{Name: "Xiao Mi", Type: "Leopard Cat", Color: "Black"},
		Cat{Name: "Lai Fu", Type: "Garfield", Color: "Gray"},
		Cat{Name: "Mi Fan", Type: "Persian Cat", Color: "White,Yellow"},
	}

	i := 0
	for i < len(cats) {
		fmt.Println("i is ", i)
		catAsBytes, _ := json.Marshal(cats[i])
		APIstub.PutState("CAT"+strconv.Itoa(i), catAsBytes)
		fmt.Println("Added", cats[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createCat(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var cat = Cat{Name: args[1], Type: args[2], Color: args[3]}

	catAsBytes, _ := json.Marshal(cat)
	APIstub.PutState(args[0], catAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllCats(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CAT0"
	endKey := "CAT999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[\n")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",\n")
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

	fmt.Printf("- queryAllCats:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeCatName(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	catAsBytes, _ := APIstub.GetState(args[0])
	cat := Cat{}

	json.Unmarshal(catAsBytes, &cat)
	cat.Name = args[1]

	catAsBytes, _ = json.Marshal(cat)
	APIstub.PutState(args[0], catAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
