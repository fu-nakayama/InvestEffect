package main

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"

//	"time"
//	"math"
//	"strings"
//	"encoding/json"
//	"sort"

)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var BK, SC, TB, Total float64
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	BK = 0
	SC = 0
	TB = 0
	Total = 0

	// String to Float64
	Total, err = strconv.ParseFloat(args[3], 64)
	if err != nil {
		return nil, errors.New("Expecting float value for Total")
	}
	fmt.Printf("Init: BK = %f, SC = %f, TB = %f, Total = %f\n", BK, SC, TB, Total)

	// Write the state (byte in string) to the ledger
	err = stub.PutState("BK", []byte(strconv.FormatFloat(BK, 'E', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("SC", []byte(strconv.FormatFloat(SC, 'E', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("TB", []byte(strconv.FormatFloat(TB, 'E', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("Total", []byte(strconv.FormatFloat(Total, 'E', -1, 64)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
//	if function == "distribute" {
//		// distribute to entity
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2")
		}
//	}
	
	var Dest		string
	var Current, Amount	float64
	var AmountStr		string
	var err error
	
	Dest = args[0]
	// String to Float64
	Amount, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, errors.New("Expecting float value for Amount to be moved")
	}
	fmt.Printf("Invoke: Dest = %s, Amount = %f\n", Dest, Amount)

	AmountBytes, err := stub.GetState(Dest)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	// String to Float64
	AmountStr = string(AmountBytes)
	Current, err = strconv.ParseFloat(AmountStr, 64)
	fmt.Printf("Invoke: Dest = %s, Current = %f, Addiing Amount = %f\n", Dest, Current, Amount)
	
	Current = Current + Amount
	err = stub.PutState(Dest, []byte(strconv.FormatFloat(Current, 'E', -1, 64)))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Invoke: Dest = %s, Current = %f\n", Dest, Current)

	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}

	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the entity to query")
	}

	Entity := args[0]
	var AmountStr	string
	var Amount	float64

	// Get the state from the ledger
	AmountBytes, err := stub.GetState(Entity)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Entity + "\"}"
		return nil, errors.New(jsonResp)
	}
	// String to Float64
	AmountStr = string(AmountBytes)
	Amount, err = strconv.ParseFloat(AmountStr, 64)

	fmt.Printf("Query Response:%s\n", AmountStr)
	return AmountBytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
