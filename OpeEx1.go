package main

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
//	"encoding/json"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Project struct {
	ProjectId	string	`json:"project_id"`
	ProjectName	string	`json:"project_name"`
	BKamount	float64	`json:"bk_amount"`
	SCamount	float64	`json:"sc_amount"`
	TBamount	float64	`json:"tb_amount"`
	FGamount	float64	`json:"fg_amount"`
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var BK, SC, TB, FG float64
	var project Project
	var err error

	BK = 0
	SC = 0
	TB = 0
	FG = 0

	// Write the state (byte in string) to the ledger
	err = stub.PutState("BK", []byte(strconv.FormatFloat(BK, 'f', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("SC", []byte(strconv.FormatFloat(SC, 'f', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("TB", []byte(strconv.FormatFloat(TB, 'f', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("FG", []byte(strconv.FormatFloat(FG, 'f', -1, 64)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "issue" {
		// issue
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}

		// String to Float64
		var Current, Amount	float64
		var AmountStr		string
		var err error

		// String to Float64
		Amount, err = strconv.ParseFloat(args[0], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for Amount to be issued")
		}
		fmt.Printf("Invoke (issue): Amount = %f\n", Amount)

		AmountBytes, err := stub.GetState("FG")
		if err != nil {
			return nil, errors.New("Failed to get state for FG")
		}
		// String to Float64
		AmountStr = string(AmountBytes)
		Current, err = strconv.ParseFloat(AmountStr, 64)
		fmt.Printf("Invoke (issue): Current = %f, Issueing Amount = %f\n", Current, Amount)

		Current = Current + Amount
		err = stub.PutState("FG", []byte(strconv.FormatFloat(Current, 'f', -1, 64)))
		if err != nil {
			return nil, err
		}
		fmt.Printf("Invoke (issue): Current = %f\n", Current)

		if err != nil {
			return nil, err
		}
		return nil, nil

	} else if function == "distribute" {
		// distribute to entity
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2")
		}

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
		fmt.Printf("Invoke (distribute): Dest = %s, Amount = %f\n", Dest, Amount)

		// Get target amount
		AmountBytes, err := stub.GetState(Dest)
		if err != nil {
			return nil, errors.New("Failed to get state")
		}
		// String to Float64
		AmountStr = string(AmountBytes)
		Current, err = strconv.ParseFloat(AmountStr, 64)
		fmt.Printf("Invoke (distribute): Dest = %s, Current = %f, Adding Amount = %f\n", Dest, Current, Amount)

		// update target amount	
		Current = Current + Amount
		err = stub.PutState(Dest, []byte(strconv.FormatFloat(Current, 'f', -1, 64)))
		if err != nil {
			return nil, err
		}
		fmt.Printf("Invoke (distribute): Dest = %s, Current = %f\n", Dest, Current)

		// Get total amount
		AmountBytes, err = stub.GetState("FG")
		if err != nil {
			return nil, errors.New("Failed to get state for FG")
		}
		// String to Float64
		AmountStr = string(AmountBytes)
		Current, err = strconv.ParseFloat(AmountStr, 64)
		fmt.Printf("Invoke (distribute): Dest = FG, Current = %f, Adding Amount = %f\n", Current, Amount)

		// update target amount	
		Current = Current - Amount
		err = stub.PutState("FG", []byte(strconv.FormatFloat(Current, 'f', -1, 64)))
		if err != nil {
			return nil, err
		}
		fmt.Printf("Invoke (distribute): Dest = FG, Current = %f\n", Current)

		if err != nil {
			return nil, err
		}
		return nil, nil
	} else if function == "init" {
		var BK, SC, TB, FG float64
		var err error

		BK = 0
		SC = 0
		TB = 0
		FG = 0

		// Write the state (byte in string) to the ledger
		err = stub.PutState("BK", []byte(strconv.FormatFloat(BK, 'f', -1, 64)))
		if err != nil {
			return nil, err
		}

		err = stub.PutState("SC", []byte(strconv.FormatFloat(SC, 'f', -1, 64)))
		if err != nil {
			return nil, err
		}

		err = stub.PutState("TB", []byte(strconv.FormatFloat(TB, 'f', -1, 64)))
		if err != nil {
			return nil, err
		}

		err = stub.PutState("FG", []byte(strconv.FormatFloat(FG, 'f', -1, 64)))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return nil, errors.New("Function of that name doesn't exist.")
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

	// Get the state from the ledger
	AmountBytes, err := stub.GetState(Entity)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Entity + "\"}"
		return nil, errors.New(jsonResp)
	}
	// String to Float64
	AmountStr = string(AmountBytes)

	fmt.Printf("Query Response:%s\n", AmountStr)
	return AmountBytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
