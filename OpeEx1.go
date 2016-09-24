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
	IssueRate	float64	`json:"issue_rate"`	// "1"
	IssueAmount	float64	`json:"issue_amount"`
	Issuer		string	`json:"issuer"`		// "FG"
	IssueYear	uint16	`json:"issue_year"`	// Fiscal Year
}

// Record of distribution
type Distribution struct {
	ProjectId	string	`json:"project_id"`	// {project_id} + "distribution"
	Currency	string	`json:"currency"`	// "JPY"
	IssueRate	float64	`json:"issue_rate"`	// "1"
	IssueAmount	float64	`json:"issue_amount"`
	Issuer		string	`json:"issuer"`		// "FG"
	IssueYear	uint16	`json:"issue_year"`	// Fiscal Year
	BKDept		string	`json:"bk_dept"`
	BKTeam		string	`json:"bk_team"`
	BKPerson	string	`json:"bk_person"`
	BKAmount	float64	`json:"bk_amount"`
	SCDept		string	`json:"sc_dept"`
	SCTeam		string	`json:"sc_team"`
	SCPerson	string	`json:"sc_person"`
	SCAmount	float64	`json:"sc_amount"`
	TBDept		string	`json:"tb_dept"`
	TBTeam		string	`json:"tb_team"`
	TBPerson	string	`json:"tb_person"`
	TBAmount	float64	`json:"tb_amount"`
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
	BKAmount	float64	`json:"bk_amount"`
	BKConfirmed	bool	`json:"bk_confirmed"`	// Yes: true, No: false	
	SCDept		string	`json:"sc_dept"`
	SCTeam		string	`json:"sc_team"`
	SCPerson	string	`json:"sc_person"`
	SCAmount	float64	`json:"sc_amount"`
	SCConfirmed	bool	`json:"sc_confirmed"`	// Yes: true, No: false	
	TBDept		string	`json:"tb_dept"`
	TBTeam		string	`json:"tb_team"`
	TBPerson	string	`json:"tb_person"`
	TBAmount	float64	`json:"tb_amount"`
	TBConfirmed	bool	`json:"tb_confirmed"`	// Yes: true, No: false
}

//
// Init
//
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Entering into Init()")

	var amount_record Amount

	// making a record
	amount_record = Amount {
		Entity:	"FG",
		Amount:	0,
	}
	bytes, err := json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error creating new record #####")
	}
	err = stub.PutState("FG", []byte(bytes))
	if err != nil {
		return nil, errors.New("##### OpeEx1: Unable to put the state #####")
	}

	amount_record = Amount {
		Entity:	"BK",
		Amount:	0,
	}
	bytes, err = json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error creating new record #####")
	}
	err = stub.PutState("BK", []byte(bytes))
	if err != nil {
		return nil, errors.New("##### OpeEx1: Unable to put the state #####")
	}

	amount_record = Amount {
		Entity:	"SC",
		Amount:	0,
	}
	bytes, err = json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error creating new record #####")
	}
	err = stub.PutState("SC", []byte(bytes))
	if err != nil {
		return nil, errors.New("##### OpeEx1:  Unable to put the state #####")
	}

	amount_record = Amount {
		Entity:	"TB",
		Amount:	0,
	}
	bytes, err = json.Marshal(amount_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error creating new record #####")
	}
	err = stub.PutState("TB", []byte(bytes))
	if err != nil {
		return nil, errors.New("##### OpeEx1: Unable to put the state #####")
	}

	// Nothing to do here, just return
	fmt.Println("Returning from Init()")
	return nil, nil
}

//
// Invoke
//
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Entering into Invoke: " + function)

	if function == "issue" {			// issue //
		// (ProjectId, Issueamount)
		if len(args) != 2 {
			return nil, errors.New("##### OpeEx1:  Incorrect number of arguments. Expecting 2 arguments for issue #####")
		}

		// String to Float64
		var issue_amount	float64
		var project_id		string
		var err			error

		// Set Arguments to local variables
		project_id = args[0]

		issue_amount, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return nil, errors.New("##### OpeEx1: Expecting float value for issue_amount to be issued #####")
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
			IssueRate:	1,
			IssueAmount:	issue_amount,
			Issuer:		"FG",
			IssueYear:	year,
		}
		bytes, err := json.Marshal(issue_record)
		if err != nil {
			return nil, errors.New("##### OpeEx1: Error creating new Issue record #####")
		}
		issue_key := project_id + "issue"
		err = stub.PutState(issue_key, []byte(bytes))
		if err != nil {
			return nil, errors.New("##### OpeEx1: Unable to put the state for Issue #####")
		}

		fmt.Println("Returning from Invoke: " + function)
		return nil, nil
	} else if function == "project" {		// project //
		// (ProjectId, ProjectName, InvestType, InvestAmount,
		//  AMCPercent, GCCPercent, GMCPercent, RBBCPercent, CICPercent,
		//  BKDept, BKTeam, BKPerson, BKAmount,
		//  SCDept, SCTeam, SCPerson, SCAmount,
		//  TBDept, TBTeam, TBPerson, TBAmount)
		if len(args) != 21 {
			return nil, errors.New("##### OpeEx1: Incorrect number of arguments. Expecting 21 arguments for project #####")
		}

		// String to Float64
		var project_id, project_name, invest_type				string
		var invest_amount							float64
		var amc_percent, gcc_percent, gmc_percent, rbbc_percent, cic_percent	float64
		var bk_dept, bk_team, bk_person						string
		var sc_dept, sc_team, sc_person						string
		var tb_dept, tb_team, tb_person						string
		var bk_amount, sc_amount, tb_amount					float64
		var bk_confirmed, sc_confirmed, tb_confirmed				bool
		var err			error

		// Set Arguments to local variables
		project_id =	args[0]
		project_name = 	args[1]
		invest_type = 	args[2]
		invest_amount, err = strconv.ParseFloat(args[3], 64)
		if err != nil {
			invest_amount = 0
		}
		amc_percent, err = strconv.ParseFloat(args[4], 64)
		if err != nil {
			amc_percent = 0
		}
		gcc_percent, err = strconv.ParseFloat(args[5], 64)
		if err != nil {
			gcc_percent = 0
		}
		gmc_percent, err = strconv.ParseFloat(args[6], 64)
		if err != nil {
			gmc_percent = 0
		}
		rbbc_percent, err = strconv.ParseFloat(args[7], 64)
		if err != nil {
			rbbc_percent = 0
		}
		cic_percent, err = strconv.ParseFloat(args[8], 64)
		if err != nil {
			cic_percent = 0
		}
		bk_dept = 	args[9]
		bk_team = 	args[10]
		bk_person = 	args[11]
		bk_amount, err = strconv.ParseFloat(args[12], 64)
		if err != nil {
			bk_amount = 0
			bk_confirmed = true
		}		
		sc_dept = 	args[13]
		sc_team = 	args[14]
		sc_person = 	args[15]
		sc_amount, err = strconv.ParseFloat(args[16], 64)
		if err != nil {
			sc_amount = 0
			sc_confirmed = true
		}		
		tb_dept = 	args[17]
		tb_team = 	args[18]
		tb_person = 	args[19]
		tb_amount, err = strconv.ParseFloat(args[20], 64)
		if err != nil {
			tb_amount = 0
			tb_confirmed = true
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
			BKAmount:	bk_amount,
			BKConfirmed:	bk_confirmed,	
			SCDept:		sc_dept,
			SCTeam:		sc_team,
			SCPerson:	sc_person,
			SCAmount:	sc_amount,
			SCConfirmed:	sc_confirmed,	
			TBDept:		tb_dept,
			TBTeam:		tb_team,
			TBPerson:	tb_person,
			TBAmount:	tb_amount,
			TBConfirmed:	tb_confirmed,
		}
		bytes, err := json.Marshal(project_record)
		if err != nil {
			return nil, errors.New("##### OpeEx1: Error on creating new Project record #####")
		}
		project_key := project_id + "project"
		err = stub.PutState(project_key, []byte(bytes))
		if err != nil {
			return nil, errors.New("##### OpeEx1: Unable to put the state for Project #####")
		}

		fmt.Println("Returning from Invoke: " + function)
		return nil, nil
	} else if function == "receivable" {		// receivable //
		// (ProjectId, AMCPercent, AMCAmount,
		//  GCCPercent, GCCAmount, GMCPercent, GMCAmount,
		//  RBBCPercent, RBBCAmount, CICPercent, CICAmount)
		if len(args) != 11 {
			return nil, errors.New("##### OpeEx1: Incorrect number of arguments. Expecting 11 arguments for receivable #####")
		}

		// String to Float64
		var project_id												string
		var amc_percent, amc_amount, gcc_percent, gcc_amount		float64
		var gmc_percent, gmc_amount, rbbc_percent, rbbc_amount		float64
		var cic_percent, cic_amount									float64
		var err														error

		// Set Arguments to local variables
		project_id =	args[0]
		amc_percent, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			amc_percent = 0
		}
		amc_amount, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			amc_amount = 0
		}
		gcc_percent, err = strconv.ParseFloat(args[3], 64)
		if err != nil {
			gcc_percent = 0
		}
		gcc_amount, err = strconv.ParseFloat(args[4], 64)
		if err != nil {
			gcc_amount = 0
		}
		gmc_percent, err = strconv.ParseFloat(args[5], 64)
		if err != nil {
			gmc_percent = 0
		}
		gmc_amount, err = strconv.ParseFloat(args[6], 64)
		if err != nil {
			gmc_amount = 0
		}
		rbbc_percent, err = strconv.ParseFloat(args[7], 64)
		if err != nil {
			rbbc_percent = 0
		}
		rbbc_amount, err = strconv.ParseFloat(args[8], 64)
		if err != nil {
			rbbc_amount = 0
		}
		cic_percent, err = strconv.ParseFloat(args[9], 64)
		if err != nil {
			cic_percent = 0
		}
		cic_amount, err = strconv.ParseFloat(args[10], 64)
		if err != nil {
			cic_amount = 0
		}
		
		// making a Receivable record
		var receivable_record Receivable
		receivable_record = Receivable {
			ProjectId:	project_id,
			Currency:	"JPY",
			AMCPercent:	amc_percent,
			AMCAmount:	amc_amount,
			GCCPercent:	gcc_percent,
			GCCAmount:	gcc_amount,
			GMCPercent:	gmc_percent,
			GMCAmount:	gmc_amount,
			RBBCPercent:	rbbc_percent,
			RBBCAmount:	rbbc_amount,
			CICPercent:	cic_percent,
			CICAmount:	cic_amount,
		}
		bytes, err := json.Marshal(receivable_record)
		if err != nil {
			return nil, errors.New("##### OpeEx1: Error on creating new Receivable record #####")
		}
		receivable_key := project_id + "receivable"
		err = stub.PutState(receivable_key, []byte(bytes))
		if err != nil {
			return nil, errors.New("##### OpeEx1: Unable to put the state for Receivable #####")
		}

		fmt.Println("Returning from Invoke: " + function)
		return nil, nil
	} else if function == "distribution" {		// distribution //
		// (ProjectId, IssueAmount,
		//  BKDept, BKTeam, BKPerson, BKAmount,
		//  SCDept, SCTeam, SCPerson, SCAmount,
		//  TBDept, TBTeam, TBPerson, TBAmount)
		if len(args) != 14 {
			return nil, errors.New("##### OpeEx1: Incorrect number of arguments. Expecting 14 arguments for distribution #####")
		}

		// String to Float64
		var project_id						string
		var issue_amount, bk_amount, sc_amount, tb_amount	float64
		var bk_dept, bk_team, bk_person				string
		var sc_dept, sc_team, sc_person				string
		var tb_dept, tb_team, tb_person				string
		var err															error

		// Set Arguments to local variables
		project_id =	args[0]
		issue_amount, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			issue_amount = 0
		}
		bk_dept =		args[2]
		bk_team =		args[3]
		bk_person =		args[4]
		bk_amount, err = strconv.ParseFloat(args[5], 64)
		if err != nil {
			bk_amount = 0
		}
		sc_dept =		args[6]
		sc_team =		args[7]
		sc_person =		args[8]
		sc_amount, err = strconv.ParseFloat(args[9], 64)
		if err != nil {
			sc_amount = 0
		}
		tb_dept =		args[10]
		tb_team =		args[11]
		tb_person =		args[12]
		tb_amount, err = strconv.ParseFloat(args[13], 64)
		if err != nil {
			tb_amount = 0
		}
		
		// making a Distribution record
		var distribution_record Distribution
		distribution_record = Distribution {
			ProjectId:	project_id,
			Currency:	"JPY",
			IssueRate:	1,
			IssueAmount:	issue_amount,
			Issuer:	"FG",
			BKDept:	bk_dept,
			BKTeam:	bk_team,
			BKPerson:	bk_person,
			BKAmount:	bk_amount,
			SCDept:	sc_dept,
			SCTeam:	sc_team,
			SCPerson:	sc_person,
			SCAmount:	sc_amount,
			TBDept:	tb_dept,
			TBTeam:	tb_team,
			TBPerson:	tb_person,
			TBAmount:	tb_amount,
		}
		bytes, err := json.Marshal(distribution_record)
		if err != nil {
			return nil, errors.New("##### OpeEx1: Error on creating new Distribution record #####")
		}
		distribution_key := project_id + "distribution"
		err = stub.PutState(distribution_key, []byte(bytes))
		if err != nil {
			return nil, errors.New("##### OpeEx1: Unable to put the state for Distribution #####")
		}

		fmt.Println("Returning from Invoke: " + function)
		return nil, nil
	} else if function == "confirm" {		// project //
		// (ProjectId, Entity)
		fmt.Println("Entering into confirm")
		if len(args) != 2 {
			return nil, errors.New("##### OpeEx1: Incorrect number of arguments. Expecting 2 arguments for confirm #####")
		}

		// Get the state from the ledger
		var project_record Project
		project_id := args[0]
		project_key := args[0] + "project"

		fmt.Println("Calling GetState in confirm")
		project_asbytes, err := stub.GetState(project_key)
		if err != nil {
			return nil, errors.New("##### OpeEx1: Failed to get state for project_id: " + project_id + " #####")
		}

		fmt.Println("Calling Unmarshal in confirm")
		err = json.Unmarshal(project_asbytes, &project_record)
		if err != nil {
			return nil, errors.New("##### OpeEx1: Error unmarshalling data " + string(project_asbytes) + " #####")
		}
		fmt.Printf("Invoke (confirm): project_id = %s\n",	project_id)
		fmt.Printf("Invoke (confirm): project_name = %s\n",	project_record.ProjectName)
		fmt.Printf("Invoke (confirm): confirmed = %t\n",	project_record.Confirmed)
		fmt.Printf("Invoke (confirm): invest_type = %s\n",	project_record.InvestType)
		fmt.Printf("Invoke (confirm): invest_amount = %f\n",	project_record.InvestAmount)
		fmt.Printf("Invoke (confirm): amc_percent = %f\n",	project_record.AMCPercent)
		fmt.Printf("Invoke (confirm): gcc_percent = %f\n",	project_record.GCCPercent)
		fmt.Printf("Invoke (confirm): gmc_percent = %f\n",	project_record.GMCPercent)
		fmt.Printf("Invoke (confirm): rbbc_percent = %f\n",	project_record.RBBCPercent)
		fmt.Printf("Invoke (confirm): cic_percent = %f\n",	project_record.CICPercent)
		fmt.Printf("Invoke (confirm): bk_dept = %s\n",		project_record.BKDept)
		fmt.Printf("Invoke (confirm): bk_team = %s\n",		project_record.BKTeam)
		fmt.Printf("Invoke (confirm): bk_person = %s\n",	project_record.BKPerson)
		fmt.Printf("Invoke (confirm): bk_amount = %f\n",	project_record.BKAmount)
		fmt.Printf("Invoke (confirm): bk_confirmed = %t\n",	project_record.BKConfirmed)
		fmt.Printf("Invoke (confirm): sc_dept = %s\n",		project_record.SCDept)
		fmt.Printf("Invoke (confirm): sc_team = %s\n",		project_record.SCTeam)
		fmt.Printf("Invoke (confirm): sc_person = %s\n",	project_record.SCPerson)
		fmt.Printf("Invoke (confirm): sc_amount = %f\n",	project_record.SCAmount)
		fmt.Printf("Invoke (confirm): sc_confirmed = %t\n",	project_record.SCConfirmed)
		fmt.Printf("Invoke (confirm): tb_dept = %s\n",		project_record.TBDept)
		fmt.Printf("Invoke (confirm): tb_team = %s\n",		project_record.TBTeam)
		fmt.Printf("Invoke (confirm): tb_person = %s\n",	project_record.TBPerson)
		fmt.Printf("Invoke (confirm): tb_amount = %f\n",	project_record.TBAmount)
		fmt.Printf("Invoke (confirm): tb_confirmed = %t\n",	project_record.TBConfirmed)

		if args[1] == "BK" {
			project_record.BKConfirmed = true
			fmt.Printf("Invoke (confirm): project_id: %s (BK) has been confirmed\n", project_id)
		} else if args[1] == "SC" {
			project_record.SCConfirmed = true
			fmt.Printf("Invoke (confirm): project_id: %s (SC) has been confirmed\n", project_id)
		} else if args[1] == "TB" {
			project_record.TBConfirmed = true
			fmt.Printf("Invoke (confirm): project_id: %s (TB) has been confirmed\n", project_id)
		} else {
			return nil, errors.New("##### OpeEx1: Expecting entity name to be confirmed #####")
		}
		if project_record.BKConfirmed == true && 
		   project_record.SCConfirmed == true &&
		   project_record.TBConfirmed == true {
		   	project_record.Confirmed = true
			fmt.Printf("Invoke (confirm): project_id: %s has been confirmed\n", project_id)
		}

		fmt.Println("Calling Marshal in confirm")
		bytes, err := json.Marshal(project_record)
		if err != nil {
			return nil, errors.New("##### OpeEx1: Error creating new Project record #####")
		}
		fmt.Println("Calling PutState in confirm")
		err = stub.PutState(project_key, []byte(bytes))
		if err != nil {
			return nil, errors.New("##### OpeEx1: Unable to put the state for Project #####")
		}

		fmt.Println("Returning from Invoke: " + function)
		return nil, nil
	}

	// Error
	fmt.Println("Invoke did not find function: " + function)
	return nil, errors.New("##### OpeEx1: Received unknown function for Invoke #####")
}

//
// Query callback representing the query of a chaincode
//
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Entering into Query: " + function)

	if function == "get_current_amount" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed");
			return nil, errors.New("##### OpeEx1: Query: Incorrect number of arguments passed #####")
		}

		entity := args[0]
		fmt.Println("Executing Query: " + function)
		return t.get_current_amount(stub, entity)
	} else if function == "get_project" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed");
			return nil, errors.New("##### OpeEx1: Query: Incorrect number of arguments passed #####")
		}

		project_id := args[0]
		fmt.Println("Executing Query: " + function)
		return t.get_project(stub, project_id)
	} else if function == "get_issue" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed");
			return nil, errors.New("##### OpeEx1: Query: Incorrect number of arguments passed #####")
		}

		project_id := args[0]
		fmt.Println("Executing Query: " + function)
		return t.get_issue(stub, project_id)
	} else if function == "get_distribution" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed");
			return nil, errors.New("##### OpeEx1: Query: Incorrect number of arguments passed #####")
		}

		project_id := args[0]
		fmt.Println("Executing Query: " + function)
		return t.get_distribution(stub, project_id)
	} else if function == "get_receivable" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed");
			return nil, errors.New("##### OpeEx1: Query: Incorrect number of arguments passed #####")
		}

		project_id := args[0]
		fmt.Println("Executing Query: " + function)
		return t.get_receivable(stub, project_id)
	}	

	// Error
	fmt.Println("Query did not find function: " + function)
	return nil, errors.New("##### OpeEx1: Received unknown function for Query #####")
}

//
// get_issue
//
func (t *SimpleChaincode) get_issue(stub *shim.ChaincodeStub, project_id string) ([]byte, error) {
	var err			error
	var issue_record	Issue

	// Get the state from the ledger
	issue_key := project_id + "issue"
	issue_asbytes, err := stub.GetState(issue_key)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Failed to get state for project_id: " + project_id + " #####")
	}
	err = json.Unmarshal(issue_asbytes, &issue_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error unmarshalling data " + string(issue_asbytes) + " #####")
	}
	fmt.Printf("Query (get_issue): project_id = %s\n",	project_id)
	fmt.Printf("Query (get_issue): currency = %s\n",	issue_record.Currency)
	fmt.Printf("Query (get_issue): issue_rate = %f\n",	issue_record.IssueRate)
	fmt.Printf("Query (get_issue): issue_amount = %f\n",	issue_record.IssueAmount)
	fmt.Printf("Query (get_issue): issuer = %s\n",		issue_record.Issuer)
	fmt.Printf("Query (get_issue): issue_year = %d\n",	issue_record.IssueYear)

	bytes, err := json.Marshal(issue_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error creating returning record #####")
	}
	return []byte(bytes), nil
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
		return nil, errors.New("##### OpeEx1: Failed to get state for project_id: " + project_id + " #####")
	}
	err = json.Unmarshal(project_asbytes, &project_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error unmarshalling data " + string(project_asbytes) + " #####")
	}
	if project_record.project_id == nil {
		return nil, errors.New("##### OpeEx1: project_id: " + project_id + " was not found #####")
	}
	fmt.Printf("Query (get_project): project_id = %s\n",	project_id)
	fmt.Printf("Query (get_project): project_name = %s\n",	project_record.ProjectName)
	fmt.Printf("Query (get_project): confirmed = %t\n",	project_record.Confirmed)
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
	fmt.Printf("Query (get_project): bk_amount = %f\n",	project_record.BKAmount)
	fmt.Printf("Query (get_project): bk_confirmed = %t\n",	project_record.BKConfirmed)
	fmt.Printf("Query (get_project): sc_dept = %s\n",	project_record.SCDept)
	fmt.Printf("Query (get_project): sc_team = %s\n",	project_record.SCTeam)
	fmt.Printf("Query (get_project): sc_person = %s\n",	project_record.SCPerson)
	fmt.Printf("Query (get_project): sc_amount = %f\n",	project_record.SCAmount)
	fmt.Printf("Query (get_project): sc_confirmed = %t\n",	project_record.SCConfirmed)
	fmt.Printf("Query (get_project): tb_dept = %s\n",	project_record.TBDept)
	fmt.Printf("Query (get_project): tb_team = %s\n",	project_record.TBTeam)
	fmt.Printf("Query (get_project): tb_person = %s\n",	project_record.TBPerson)
	fmt.Printf("Query (get_project): tb_amount = %f\n",	project_record.TBAmount)
	fmt.Printf("Query (get_project): tb_confirmed = %t\n",	project_record.TBConfirmed)

	bytes, err := json.Marshal(project_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error creating returning record #####")
	}
	return []byte(bytes), nil
}

//
// get_distribution
//
func (t *SimpleChaincode) get_distribution(stub *shim.ChaincodeStub, project_id string) ([]byte, error) {
	var err						error
	var distribution_record		Distribution

	// Get the state from the ledger
	distribution_key := project_id + "distribution"
	distribution_asbytes, err := stub.GetState(distribution_key)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Failed to get state for project_id: " + project_id + " #####")
	}
	err = json.Unmarshal(distribution_asbytes, &distribution_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error unmarshalling data " + string(distribution_asbytes) + " #####")
	}
	fmt.Printf("Query (get_distribution): project_id = %s\n",	project_id)
	fmt.Printf("Query (get_distribution): currency = %s\n",		distribution_record.Currency)
	fmt.Printf("Query (get_distribution): issue_rate = %f\n",	distribution_record.IssueRate)
	fmt.Printf("Query (get_distribution): issue_amount = %f\n",	distribution_record.IssueAmount)
	fmt.Printf("Query (get_distribution): issuer = %s\n",		distribution_record.Issuer)
	fmt.Printf("Query (get_distribution): issue_year = %d\n",	distribution_record.IssueYear)
	fmt.Printf("Query (get_distribution): bk_dept = %s\n",		distribution_record.BKDept)
	fmt.Printf("Query (get_distribution): bk_team = %s\n",		distribution_record.BKTeam)
	fmt.Printf("Query (get_distribution): bk_person = %s\n",	distribution_record.BKPerson)
	fmt.Printf("Query (get_distribution): bk_amount = %f\n",	distribution_record.BKAmount)
	fmt.Printf("Query (get_distribution): sc_dept = %s\n",		distribution_record.SCDept)
	fmt.Printf("Query (get_distribution): sc_team = %s\n",		distribution_record.SCTeam)
	fmt.Printf("Query (get_distribution): sc_person = %s\n",	distribution_record.SCPerson)
	fmt.Printf("Query (get_distribution): sc_amount = %f\n",	distribution_record.SCAmount)
	fmt.Printf("Query (get_distribution): tb_dept = %s\n",		distribution_record.TBDept)
	fmt.Printf("Query (get_distribution): tb_team = %s\n",		distribution_record.TBTeam)
	fmt.Printf("Query (get_distribution): tb_person = %s\n",	distribution_record.TBPerson)
	fmt.Printf("Query (get_distribution): tb_amount = %f\n",	distribution_record.TBAmount)

	bytes, err := json.Marshal(distribution_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error creating returning record #####")
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
		return nil, errors.New("##### OpeEx1: Failed to get state for project_id: " + project_id + " #####")
	}
	err = json.Unmarshal(receivable_asbytes, &receivable_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error unmarshalling data " + string(receivable_asbytes) + " #####")
	}
	fmt.Printf("Query (get_receivable): project_id = %s\n",		project_id)
	fmt.Printf("Query (get_receivable): currency = %s\n",		receivable_record.Currency)
	fmt.Printf("Query (get_receivable): amc_percent = %f\n",	receivable_record.AMCPercent)
	fmt.Printf("Query (get_receivable): amc_amount = %f\n",		receivable_record.AMCAmount)
	fmt.Printf("Query (get_receivable): gcc_percent = %f\n",	receivable_record.GCCPercent)
	fmt.Printf("Query (get_receivable): gcc_amount = %f\n",		receivable_record.GCCAmount)
	fmt.Printf("Query (get_receivable): gmc_percent = %f\n",	receivable_record.GMCPercent)
	fmt.Printf("Query (get_receivable): gmc_amount = %f\n",		receivable_record.GMCAmount)
	fmt.Printf("Query (get_receivable): rbbc_percent = %f\n",	receivable_record.RBBCPercent)
	fmt.Printf("Query (get_receivable): rbbc_amount = %f\n",	receivable_record.RBBCAmount)
	fmt.Printf("Query (get_receivable): cic_percent = %f\n",	receivable_record.CICPercent)
	fmt.Printf("Query (get_receivable): cic_amount = %f\n",		receivable_record.CICAmount)

	bytes, err := json.Marshal(receivable_record)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Error creating returning record #####")
	}
	return []byte(bytes), nil
}

//
// get_current_amount
//
func (t *SimpleChaincode) get_current_amount(stub *shim.ChaincodeStub, entity string) ([]byte, error) {
	var err		error
	var AmountStr	string
	var Amount	float64

	// Get the state from the ledger
	AmountBytes, err := stub.GetState(entity)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Failed to get state for entity: " + entity + " #####")
	}

	// Bytes to String
	AmountStr = string(AmountBytes)

	// String to Float64
	Amount, err = strconv.ParseFloat(AmountStr, 64)
	if err != nil {
		return nil, errors.New("##### OpeEx1: Expecting float value for Amount to be issued #####")
	}
	fmt.Printf("Query (get_receivable): entity = %s\n",	entity)
	fmt.Printf("Query (get_receivable): amount = %f\n",	Amount)

	return []byte(AmountStr), nil
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
