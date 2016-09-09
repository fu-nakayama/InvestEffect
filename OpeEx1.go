package main

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Current amount
type Amount struct {
	Entity		string	`json:"entity"`		// "FG" | BK" | "SC" | "TB"
	Amount		float64	`json:"amount"`
}

// Record of issue
type Issue struct {
	ProjectId	string	`json:"project_id"`	// {project_id} + "issue"
	Currency	string	`json:"currency"`	// "JPY"
	MinAmount	float64	`json:"min_amount"`
	IssueAmount	float64	`json:"issue_amount"`
	Issuer		string	`json:"issuer"`		// "FG"
	IssueYear	uint16	`json:"issue_year"`	// Fiscal Year
	Confirmed	bool	`json:"confirmed"`	// Yes: true, No: false
}

// Record of distribution
type Distribution struct {
	ProjectId	string	`json:"project_id"`	// {project_id} + "issue"
	Currency	string	`json:"currency"`	// "JPY"
	MinAmount	float64	`json:"min_amount"`	// "0.01"
	IssueAmount	float64	`json:"issue_amount"`
	Issuer		float64	`json:"issuer"`
	IssueYear	uint16	`json:"issue_year"`
	BKDept		string	`json:"bk_dept"`
	BKTeam		string	`json:"bk_team"`
	BKPerson	string	`json:"bk_person"`
	BKConfirmed	bool	`json:"bk_confirmed"`	// Yes: true, No: false	
	SCDept		string	`json:"sc_dept"`
	SCTeam		string	`json:"sc_team"`
	SCPerson	string	`json:"sc_person"`
	SCConfirmed	bool	`json:"sc_confirmed"`	// Yes: true, No: false	
	TBDept		string	`json:"tb_dept"`
	TBTeam		string	`json:"tb_team"`
	TBPerson	string	`json:"tb_person"`
	TBConfirmed	bool	`json:"tb_confirmed"`	// Yes: true, No: false	
}

type Project struct {
	ProjectId	string	`json:"project_id"`
	BKamount	float64	`json:"bk_amount"`
	SCamount	float64	`json:"sc_amount"`
	TBamount	float64	`json:"tb_amount"`
	FGamount	float64	`json:"fg_amount"`
}

//
// Init
//
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var amount_record Amount

	// making a record
	amount_record = Amount {
		Entity:	"FG",
		Amount:	0,
	}
	bytes, err := json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("Error creating new record")
	}
	err = stub.PutState(project_id, []byte(bytes))
	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	amount_record = Amount {
		Entity:	"BK",
		Amount:	0,
	}
	bytes, err := json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("Error creating new record")
	}
	err = stub.PutState(project_id, []byte(bytes))
	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	amount_record = Amount {
		Entity:	"SC",
		Amount:	0,
	}
	bytes, err := json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("Error creating new record")
	}
	err = stub.PutState(project_id, []byte(bytes))
	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	amount_record = Amount {
		Entity:	"TB",
		Amount:	0,
	}
	bytes, err := json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("Error creating new record")
	}
	err = stub.PutState(project_id, []byte(bytes))
	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	// Nothing to do here, just return
	return nil, nil
}

//
// Invoke
//
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "issue" {
		// issue (ProjectId, BKamount, SCamount, TBamount, FGamount)
		if len(args) != 5 {
			return nil, errors.New("Incorrect number of arguments. Expecting 5")
		}

		// String to Float64
		var bk_amount, sc_amount, tb_amount, fg_amount	float64
		var project_id									string
		var err											error

		// Set Arguments to local variables
		project_id = args[0]

		bk_amount, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for bk_amount to be issued")
		}
		sc_amount, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for sc_amount to be issued")
		}
		tb_amount, err = strconv.ParseFloat(args[3], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for tb_amount to be issued")
		}
		fg_amount, err = strconv.ParseFloat(args[4], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for fg_amount to be issued")
		}
		fmt.Printf("Invoke (issue): project_id = %f\n", project_id)
		fmt.Printf("Invoke (issue): bk_amount = %f\n", bk_amount)
		fmt.Printf("Invoke (issue): sc_amount = %f\n", sc_amount)
		fmt.Printf("Invoke (issue): tb_amount = %f\n", tb_amount)
		fmt.Printf("Invoke (issue): fg_amount = %f\n", fg_amount)

		// making a record
		var project_record Project
		project_record = Project {
			ProjectId:	args[0],
			BKamount:	bk_amount,
			SCamount:	sc_amount,
			TBamount:	tb_amount,
			FGamount:	fg_amount,
		}

		bytes, err := json.Marshal(project_record)
		if err != nil {
			return nil, errors.New("Error creating new record")
		}

		err = stub.PutState(project_id, []byte(bytes))
		if err != nil {
			return nil, errors.New("Unable to put the state")
		}

		return nil, nil
	}

	// Error
	fmt.Println("Invoke did not find function: " + function)
	return nil, errors.New("Received unknown function for Invoke")
}

//
// Query callback representing the query of a chaincode
//
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "get_invest_summary" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed");
			return nil, errors.New("Query: Incorrect number of arguments passed")
		}

		project_id := args[0]
		return t.get_invest_summary(stub, project_id)
	}

	// Error
	fmt.Println("Query did not find function: " + function)
	return nil, errors.New("Received unknown function for Query")
}

//
// get_invest_summary
//
func (t *SimpleChaincode) get_invest_summary(stub *shim.ChaincodeStub, project_id string) ([]byte, error) {
	var err				error
	var project_record	Project

	// Get the state from the ledger
	project_summary_asbytes, err := stub.GetState(project_id)
	if err != nil {
		return nil, errors.New("Error: Failed to get state for project_id: " + project_id)
	}

	if err = json.Unmarshal(project_summary_asbytes, &project_record) ; err != nil {return nil, errors.New("Error unmarshalling data "+string(project_summary_asbytes))}
	fmt.Printf("Invoke (issue): project_id = %f\n", project_id)
	fmt.Printf("Invoke (issue): bk_amount = %f\n", project_record.BKamount)
	fmt.Printf("Invoke (issue): sc_amount = %f\n", project_record.SCamount)
	fmt.Printf("Invoke (issue): tb_amount = %f\n", project_record.TBamount)
	fmt.Printf("Invoke (issue): fg_amount = %f\n", project_record.FGamount)

	bytes, err := json.Marshal(project_record)
	if err != nil {
		return nil, errors.New("Error creating returning record")
	}
	return []byte(bytes), nil
}

//
// Main
//
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting OpeEx1 chaincode: %s", err)
	}
}
