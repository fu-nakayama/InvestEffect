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
	var BKamount, SCamount, TBamount float64
	var err error

	BKamount = 0
	SCamount = 0
	TBamount = 0

	fmt.Printf("Init: BKamount = %f, SCamount = %f, TBamount = %f\n", BKamount, SCamount, TBamount)

	// Write the state to the ledger
	err = stub.PutState("BK", []byte(strconv.FormatFloat(BKamount, 'E', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("SC", []byte(strconv.FormatFloat(SCamount, 'E', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("TB", []byte(strconv.FormatFloat(TBamount, 'E', -1, 64)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error
	var BKamount, SCamount, TBamount float64

	BKamount, err = strconv.FormatFloat(stub.GetState("BK"), 'E', -1, 64)
	SCamount, err = strconv.FormatFloat(stub.GetState("SC"), 'E', -1, 64)
	TBamount, err = strconv.FormatFloat(stub.GetState("TB"), 'E', -1, 64)
	fmt.Printf("Invoke: BKamount = %f, SCamount = %f, TBamount = %f\n", BKamount, SCamount, TBamount)

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}

	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	Entity = args[0]

	// Get the state from the ledger
	Statebytes, err := stub.GetState(Entity)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Entity + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Statebytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + Entity + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Entity\":\"" + Entity + "\",\"Amount\":\"" + string(Statebytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Statebytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
