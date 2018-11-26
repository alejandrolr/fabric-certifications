package main

import (
	"reflect"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *SimpleChaincode) issueBadge(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// check function args. In description, email is omitted in description
	// because it is obtained from the certificate.
	if len(args) != 8 {
		return shim.Error(`Incorrect number of arguments. Expecting 7:\n
		1) Issuer Name, 2) Issuer URL, 3) Badge Name (id), 4) Badge Description\n
		5) Badge Criteria, 6) Job Description for Signature, 7) Name for Signature`)
	}

	logger.Infof("Action: issue Badge")
	issuerEmail, issuerName, issuerUrl := args[0] /*Issuer ID*/, args[1], args[2]
	badgeName, badgeDesc, criteria := args[3], args[4], args[5]
	badgeJobDesc, badgeSigName := args[6], args[7]

	var issuer Issuer
	var issuerList IssuerList
	issuerExists := false
	issuerIndexInIssuerList := 0

	var badge Badge
	var badgeFromLedger Badge

	// 1. Get issuerList from ledger
	// -----------------------------
	issuerListMap, err := getStructFromLedger(stub, ISSUER_LIST)
	if err != nil {
		// error retrieving issuer-list
		return shim.Error(err.Error())
	}
	// Convert map[string] to IssuerList struct
	FillStruct(issuerListMap, &issuerList)

	// 2. Check if issuer exist in ledger
	// ----------------------------------
	for i, issuerSum := range issuerList.IssuerSummary {
		if strings.Compare(issuerSum.Email, issuerEmail) == 0 {
			logger.Infof("Issuer already exists, skipping issuer creation...")
			issuerExists = true
			issuerIndexInIssuerList = i
		}
	}

	// 3. If issuer doesn't exist in ledger, create it and append
	//    its summary to issuerList
	// ----------------------------------------------------------
	if !issuerExists {
		logger.Infof("Issuer doesn't exist, creating...")
		// Create Issuer struct
		issuer = createIssuer(issuerEmail, issuerUrl, issuerEmail, issuerName)

		// Create IssuerSummary struct
		issuerSummary := createIssuerSummary(issuer)

		err = marshalAndPutState(stub, issuer, issuerEmail)
		if err != nil {
			// error marshaling or putting state into ledger
			return shim.Error(err.Error())
		}

		// Append issuerSummary to issuerList
		issuerList.IssuerSummary = append(issuerList.IssuerSummary, issuerSummary)
		issuerIndexInIssuerList = len(issuerList.IssuerSummary) - 1
	} else {
		// 4. Issuer exists, get issuer from ledger
		// ----------------------------------------
		issuerMap, err := getStructFromLedger(stub, issuerEmail)
		if err != nil {
			// error retrieving issuer
			return shim.Error(err.Error())
		}
		// Convert map[string] to struct
		FillStruct(issuerMap, &issuer)

		// check if user is not empty
		if reflect.DeepEqual(issuer, Issuer{}) {
			// issuer doesn't exist in ledger, aborting
			return shim.Error("Issuer is empty in ledger, aborting!")
		}

		logger.Infof("Issuer %s found in ledger!", issuerEmail)
	}

	// 5. Create a Badge (it includes issuer) and write it to the ledger
	// -----------------------------------------------------------------
	// Badge ID will be the badge name without spaces and in lowercase
	badgeID := BADGE_PREFIX + strings.ToLower(strings.Replace(badgeName, " ", "", -1))

	// get Badge from ledger
	badgeFromLedgerMap, err := getStructFromLedger(stub, badgeID)
	if err != nil {
		// error retrieving badge
		return shim.Error(err.Error())
	}
	// Convert map[string] to Badge struct
	FillStruct(badgeFromLedgerMap, &badgeFromLedger)

	// Check if badge exists (not empty)
	if reflect.DeepEqual(badgeFromLedger, Badge{}) {
		logger.Infof("Badge doesn't exist, creating...")
		// creating elements from parameters
		badgeSignatureLines := createSignatureLines(badgeJobDesc, badgeSigName)
		badgeCriteria := createCriteria(criteria)
		// Create Badge
		badge = createBadge(badgeID, badgeName, badgeDesc, issuer, badgeCriteria, badgeSignatureLines)
	} else {
		// if badge exist, abort
		logger.Errorf("Badge exists, aborting...")
		return shim.Error("Badge already exists, aborting!")
	}

	// Write the badge into the ledger (KEY: id, unique)
	err = marshalAndPutState(stub, badge, badge.Id)
	if err != nil {
		// error marshaling or putting state into ledger
		return shim.Error(err.Error())
	}

	// 6. Append badgeID to badgeIDs in IssuerSummary
	// ----------------------------------------------
	issuerList.IssuerSummary[issuerIndexInIssuerList].BagdeIDs =
		append(issuerList.IssuerSummary[issuerIndexInIssuerList].BagdeIDs, badgeID)

	// Write the state into the ledger (KEY: issuer-list, unique)
	err = marshalAndPutState(stub, issuerList, ISSUER_LIST)
	if err != nil {
		// error marshaling or putting state into ledger
		return shim.Error(err.Error())
	}

	returnMessage := "Successfully updated blockchain: CREATED Badge and UPDATED issuerList"
	return shim.Success([]byte(returnMessage))
}
