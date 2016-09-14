package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
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
}

// Record of distribution
type Distribution struct {
	ProjectId	string	`json:"project_id"`	// {project_id} + "dist"
	Currency	string	`json:"currency"`	// "JPY"
	MinAmount	float64	`json:"min_amount"`	// "0.01"
	IssueAmount	float64	`json:"issue_amount"`
	Issuer		float64	`json:"issuer"`		// "FG"
	IssueYear	uint16	`json:"issue_year"`	// Fiscal Year
	BKDept		string	`json:"bk_dept"`
	BKTeam		string	`json:"bk_team"`
	BKPerson	string	`json:"bk_person"`
	BKamount	float64	`json:"bk_amount"`
	SCDept		string	`json:"sc_dept"`
	SCTeam		string	`json:"sc_team"`
	SCPerson	string	`json:"sc_person"`
	SCamount	float64	`json:"sc_amount"`
	TBDept		string	`json:"tb_dept"`
	TBTeam		string	`json:"tb_team"`
	TBPerson	string	`json:"tb_person"`
	TBamount	float64	`json:"tb_amount"`
}

// Record of receivable
type Receivable struct {
	ProjectId	string	`json:"project_id"`	// {project_id} + "receivable"
	Currency	string	`json:"currency"`	// "JPY"
	AMCPercent	float64	`json:"amc_percent"`
	AMCAmount	float64	`json:"amc_amount"`
	GCCPercent	float64	`json:"gcc_percent"`
	GCCAmount	float64	`json:"gcc_amount"`
	GMCPercent	float64	`json:"gmc_percent"`
	GMCAmount	float64	`json:"gmc_amount"`
	RBBCPercent	float64	`json:"rbbc_percent"`
	RBBCAmount	float64	`json:"rbbc_amount"`
	CICPercent	float64	`json:"cic_percent"`
	CICAmount	float64	`json:"cic_amount"`
}

type Project struct {
	ProjectId	string	`json:"project_id"`	// {project_id} + "project"
	ProjectName	string	`json:"project_name"`
	InvestType	string	`json:"invest_type"`
	InvestAmount	float64	`json:"invest_amount"`
	Confirmed	bool	`json:"confirmed"`	// Yes: true, No: false	
	AMCPercent	float64	`json:"amc_percent"`
	GCCPercent	float64	`json:"gcc_percent"`
	GMCPercent	float64	`json:"gmc_percent"`
	RBBCPercent	float64	`json:"rbbc_percent"`
	CICPercent	float64	`json:"cic_percent"`
	BKDept		string	`json:"bk_dept"`
	BKTeam		string	`json:"bk_team"`
	BKPerson	string	`json:"bk_person"`
	BKamount	float64	`json:"bk_amount"`
	BKConfirmed	bool	`json:"bk_confirmed"`	// Yes: true, No: false	
	SCDept		string	`json:"sc_dept"`
	SCTeam		string	`json:"sc_team"`
	SCPerson	string	`json:"sc_person"`
	SCamount	float64	`json:"sc_amount"`
	SCConfirmed	bool	`json:"sc_confirmed"`	// Yes: true, No: false	
	TBDept		string	`json:"tb_dept"`
	TBTeam		string	`json:"tb_team"`
	TBPerson	string	`json:"tb_person"`
	TBamount	float64	`json:"tb_amount"`
	TBConfirmed	bool	`json:"tb_confirmed"`	// Yes: true, No: false
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
	err = stub.PutState("FG", []byte(bytes))
	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	amount_record = Amount {
		Entity:	"BK",
		Amount:	0,
	}
	bytes, err = json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("Error creating new record")
	}
	err = stub.PutState("BK", []byte(bytes))
	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	amount_record = Amount {
		Entity:	"SC",
		Amount:	0,
	}
	bytes, err = json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("Error creating new record")
	}
	err = stub.PutState("SC", []byte(bytes))
	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	amount_record = Amount {
		Entity:	"TB",
		Amount:	0,
	}
	bytes, err = json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("Error creating new record")
	}
	err = stub.PutState("TB", []byte(bytes))
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

	if function == "issue" {			// issue //
		// issue (ProjectId, Issueamount)
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments for issue.")
		}

		// String to Float64
		var issue_amount	float64
		var project_id		string
		var err			error

		// Set Arguments to local variables
		project_id = args[0]

		issue_amount, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for issue_amount to be issued")
		}
		fmt.Printf("Invoke (issue): project_id = %s\n", project_id)
		fmt.Printf("Invoke (issue): issue_amount = %f\n", issue_amount)

		// Get current date and time
		t := time.Now()

		// making a Issue record
		var year	uint16
		var month 	uint8
		year =		uint16(t.Year())
		month =		uint8(t.Month())
		if month < 4 {
			year = year + 1
		}
		
		var issue_record Issue
		issue_record = Issue {
			ProjectId:	project_id,
			Currency:	"JPY",
			MinAmount:	0,
			IssueAmount:	issue_amount,
			Issuer:		"FG",
			IssueYear:	year,
		}
		bytes, err := json.Marshal(issue_record)
		if err != nil {
			return nil, errors.New("Error creating new Issue record")
		}
		issue_key := project_id + "issue"
		err = stub.PutState(issue_key, []byte(bytes))
		if err != nil {
			return nil, errors.New("Unable to put the state for Issue")
		}
		return nil, nil
	} else if function == "project" {		// project //
		// issue (ProjectId, ProjectName, InvestType, InvestAmount,
		//        AMCPercent, GCCPercent, GMCPercent, RBBCPercent, CICPercent,
		//        BKDept, BKTeam, BKPerson, BKamount,
		//        SCDept, SCTeam, SCPerson, SCamount,
		//        TBDept, TBTeam, TBPerson, TBamount)
		if len(args) != 21 {
			return nil, errors.New("Incorrect number of arguments. Expecting 21 arguments for project.")
		}

		// String to Float64
		var project_id, project_name, invest_type				string
		var invest_amount							float64
		var amc_percent, gcc_percent, gmc_percent, rbbc_percent, cic_percent	float64
		var bk_dept, bk_team, bk_person						string
		var sc_dept, sc_team, sc_person						string
		var tb_dept, tb_team, tb_person						string
		var bk_amount, sc_amount, tb_amount					float64
		var err			error

		// Set Arguments to local variables
		project_id =	args[0]
		project_name = 	args[1]
		invest_type = 	args[2]
		invest_amount, err = strconv.ParseFloat(args[3], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for invest_amount to be issued")
		}
		amc_percent, err = strconv.ParseFloat(args[4], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for amc_percent to be issued")
		}
		gcc_percent, err = strconv.ParseFloat(args[5], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for gcc_percent to be issued")
		}
		gmc_percent, err = strconv.ParseFloat(args[6], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for gmc_percent to be issued")
		}
		rbbc_percent, err = strconv.ParseFloat(args[7], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for rbbc_percent to be issued")
		}
		cic_percent, err = strconv.ParseFloat(args[8], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for cic_percent to be issued")
		}
		bk_dept = 	args[9]
		bk_team = 	args[10]
		bk_person = 	args[11]
		bk_amount, err = strconv.ParseFloat(args[12], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for bk_amount to be issued")
		}		
		sc_dept = 	args[13]
		sc_team = 	args[14]
		sc_person = 	args[15]
		sc_amount, err = strconv.ParseFloat(args[16], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for sc_amount to be issued")
		}		
		tb_dept = 	args[17]
		tb_team = 	args[18]
		tb_person = 	args[19]
		tb_amount, err = strconv.ParseFloat(args[20], 64)
		if err != nil {
			return nil, errors.New("Expecting float value for tb_amount to be issued")
		}		
		
		// making a Project record
		var project_record Project
		project_record = Project {
			ProjectId:	project_id,
			ProjectName:	project_name,
			InvestType:	invest_type,
			InvestAmount:	invest_amount,
			Confirmed:	false,
			AMCPercent:	amc_percent,
			GCCPercent:	gcc_percent,
			GMCPercent:	gmc_percent,
			RBBCPercent:	rbbc_percent,
			CICPercent:	cic_percent,
			BKDept:		bk_dept,
			BKTeam:		bk_team,
			BKPerson:	bk_person,
			BKamount:	bk_amount,
			BKConfirmed:	false,	
			SCDept:		sc_dept,
			SCTeam:		sc_team,
			SCPerson:	sc_person,
			SCamount:	sc_amount,
			SCConfirmed:	false,	
			TBDept:		tb_dept,
			TBTeam:		tb_team,
			TBPerson:	tb_person,
			TBamount:	tb_amount,
			TBConfirmed:	false,
		}
		bytes, err := json.Marshal(project_record)
		if err != nil {
			return nil, errors.New("Error creating new Project record")
		}
		project_key := project_id + "project"
		err = stub.PutState(project_key, []byte(bytes))
		if err != nil {
			return nil, errors.New("Unable to put the state for Project")
		}

		return nil, nil
	} else if function == "confirm" {		// project //
		// issue (ProjectId, Entity)
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments for confirm.")
		}

		// Get the state from the ledger
		var project_record Project
		project_id := args[0]
		project_asbytes, err := stub.GetState(project_id)
		if err != nil {
			return nil, errors.New("Error: Failed to get state for project_id: " + project_id)
		}
		if err = json.Unmarshal(project_asbytes, &project_record) ; err != nil {
			return nil, errors.New("Error unmarshalling data "+string(project_asbytes))
		}

		if args[1] == "BK" {
			project_record.BKConfirmed = true
		} else if args[1] == "SC" {
			project_record.SCConfirmed = true
		} else if args[1] == "TB" {
			project_record.TBConfirmed = true
		} else {
			return nil, errors.New("Expecting entity name to be confirmed")
		}
		if project_record.BKConfirmed == true && 
		   project_record.SCConfirmed == true &&
		   project_record.TBConfirmed == true {
		   	project_record.Confirmed = true
		}

		bytes, err := json.Marshal(project_record)
		if err != nil {
			return nil, errors.New("Error creating new Project record")
		}
		confirm_key := project_id + "confirm"
		err = stub.PutState(confirm_key, []byte(bytes))
		if err != nil {
			return nil, errors.New("Unable to put the state for Project")
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

	if function == "get_project" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed");
			return nil, errors.New("Query: Incorrect number of arguments passed")
		}

		project_id := args[0]
		return t.get_project(stub, project_id)
	} else if function == "get_receivable" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed");
			return nil, errors.New("Query: Incorrect number of arguments passed")
		}

		project_id := args[0]
		return t.get_receivable(stub, project_id)
	}	

	// Error
	fmt.Println("Query did not find function: " + function)
	return nil, errors.New("Received unknown function for Query")
}

//
// get_project
//
func (t *SimpleChaincode) get_project(stub *shim.ChaincodeStub, project_id string) ([]byte, error) {
	var err			error
	var project_record	Project

	// Get the state from the ledger
	project_key := project_id + "project"
	project_asbytes, err := stub.GetState(project_key)
	if err != nil {
		return nil, errors.New("Error: Failed to get state for project_id: " + project_id)
	}

	if err = json.Unmarshal(project_asbytes, &project_record) ; err != nil {
		return nil, errors.New("Error unmarshalling data "+string(project_asbytes))
	}
	fmt.Printf("Query (get_project): project_id = %s\n",	project_id)
	fmt.Printf("Query (get_project): project_name = %s\n",	project_record.ProjectName)
	fmt.Printf("Query (get_project): invest_type = %s\n",	project_record.InvestType)
	fmt.Printf("Query (get_project): invest_amount = %f\n",	project_record.InvestAmount)
	fmt.Printf("Query (get_project): amc_percent = %f\n",	project_record.AMCPercent)
	fmt.Printf("Query (get_project): gcc_percent = %f\n",	project_record.GCCPercent)
	fmt.Printf("Query (get_project): gmc_percent = %f\n",	project_record.GMCPercent)
	fmt.Printf("Query (get_project): rbbc_percent = %f\n",	project_record.RBBCPercent)
	fmt.Printf("Query (get_project): cic_percent = %f\n",	project_record.CICPercent)
	fmt.Printf("Query (get_project): bk_dept = %s\n",	project_record.BKDept)
	fmt.Printf("Query (get_project): bk_team = %s\n",	project_record.BKTeam)
	fmt.Printf("Query (get_project): bk_person = %s\n",	project_record.BKPerson)
	fmt.Printf("Query (get_project): bk_amount = %f\n",	project_record.BKamount)
	fmt.Printf("Query (get_project): sc_dept = %s\n",	project_record.SCDept)
	fmt.Printf("Query (get_project): sc_team = %s\n",	project_record.SCTeam)
	fmt.Printf("Query (get_project): sc_person = %s\n",	project_record.SCPerson)
	fmt.Printf("Query (get_project): sc_amount = %f\n",	project_record.SCamount)

	bytes, err := json.Marshal(project_record)
	if err != nil {
		return nil, errors.New("Error creating returning record")
	}
	return []byte(bytes), nil
}

//
// get_receivable
//
func (t *SimpleChaincode) get_receivable(stub *shim.ChaincodeStub, project_id string) ([]byte, error) {
	var err			error
	var receivable_record	Receivable

	// Get the state from the ledger
	receivable_key := project_id + "receivable"
	receivable_asbytes, err := stub.GetState(receivable_key)
	if err != nil {
		return nil, errors.New("Error: Failed to get state for project_id: " + project_id)
	}

	if err = json.Unmarshal(receivable_asbytes, &receivable_record) ; err != nil {
		return nil, errors.New("Error unmarshalling data "+string(receivable_asbytes))
	}
	fmt.Printf("Query (get_receivable): project_id = %s\n",		project_id)
	fmt.Printf("Query (get_receivable): currency = %s\n",		project_record.Currency)
	fmt.Printf("Query (get_receivable): amc_percent = %f\n",	project_record.AMCPercent)
	fmt.Printf("Query (get_receivable): amc_amount = %f\n",		project_record.AMCAmount)
	fmt.Printf("Query (get_receivable): gcc_percent = %f\n",	project_record.GCCPercent)
	fmt.Printf("Query (get_receivable): gcc_amount = %f\n",		project_record.GCCCamount)
	fmt.Printf("Query (get_receivable): gmc_percent = %f\n",	project_record.GMCPercent)
	fmt.Printf("Query (get_receivable): gmc_amount = %f\n",		project_record.GMCAmount)
	fmt.Printf("Query (get_receivable): rbbc_percent = %f\n",	project_record.RBBCPercent)
	fmt.Printf("Query (get_receivable): rbbc_amount = %f\n",	project_record.RBBCAmount)
	fmt.Printf("Query (get_receivable): cic_percent = %f\n",	project_record.CICPercent)
	fmt.Printf("Query (get_receivable): cic_amount = %f\n",		project_record.CICAmount)

	bytes, err := json.Marshal(receivable_record)
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
