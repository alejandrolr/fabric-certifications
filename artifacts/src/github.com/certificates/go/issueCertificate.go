package main

import (
	"reflect"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *SimpleChaincode) issueCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Infof("Action: issue Certificate")
	var returnMessage string

	// check function args. In description, email is omitted in description
	// because it is obtained from the certificate.
	if len(args) != 7 {
		return shim.Error(`Incorrect number of arguments. Expecting 6:\n
		1) issuedOn, 2) Recipient Email, 3) recipient Name, 4) Recipient Public Key, 
		5) Certificate location, 6) Badge ID (name of the badge without spaces in lowercase)\n`)
	}

	// Parameters
	issuerEmail, issuedOn := args[0], args[1]
	recipientEmail := args[2]
	recipientName, recipientPubKey := args[3], args[4]
	location := args[5]
	badgeKey := args[6]

	var issuerList IssuerList
	issuerExists := false
	issuerIndexInIssuerList := 0

	var badgeFromLedger Badge

	var cert Certificate
	var certFromLedger Certificate

	// 1. Check if badge exists in ledger and if it is owned by the certificate issuer
	// -------------------------------------------------------------------------------
	// get Badge from ledger
	badgeFromLedgerMap, err := getStructFromLedger(stub, BADGE_PREFIX+badgeKey)
	if err != nil {
		// error retrieving badge
		return shim.Error(err.Error())
	}
	// Convert map[string] to Badge struct
	FillStruct(badgeFromLedgerMap, &badgeFromLedger)

	// Check if badge exists (not empty)
	if reflect.DeepEqual(badgeFromLedger, Badge{}) {
		logger.Errorf("Provided Badge doesn't exist, aborting...")
		return shim.Error("Badge doesn't exist, aborting")
	}

	// Check if badge is owned by issuer
	if strings.Compare(badgeFromLedger.Issuer.Id, issuerEmail) != 0 {
		return shim.Error("Badge is not owned by " + issuerEmail)
	}

	// 2. Create certificate if it doesn't exist
	// -----------------------------------------
	// create Certificate ID (cert:recipientEmail-badgeKey)
	certID := CERT_PREFIX + recipientEmail + "-" + badgeKey

	// get certificate from ledger
	certFromLedgerMap, err := getStructFromLedger(stub, certID)
	if err != nil {
		// error retrieving cert
		return shim.Error(err.Error())
	}
	// Convert map[string] to Certificate struct
	FillStruct(certFromLedgerMap, &certFromLedger)

	// Check if certificate exists (not empty)
	if reflect.DeepEqual(certFromLedger, Certificate{}) {
		logger.Infof("Certificate doesn't exist, creating...")

		// creating elements from parameters
		rec := createRecipient(recipientEmail)
		recProf := createRecipientProfile(recipientPubKey, recipientName)
		ver := createVerification(location)

		// Create Certificate
		cert = createCertificate(certID, issuedOn, rec, recProf, ver, badgeFromLedger)
	} else {
		// if certificate exist, abort
		logger.Errorf("Certificate exists, aborting...")
		return shim.Error("Certificate already exists, aborting!")
	}

	// Write the cert into the ledger (KEY: certId, unique)
	err = marshalAndPutState(stub, cert, cert.Id)
	if err != nil {
		// error marshaling or putting state into ledger
		return shim.Error(err.Error())
	}

	// 3. Append certID to certIDs in IssuerSummary and update issuer-list
	// -------------------------------------------------------------------
	// Get issuerList from ledger
	issuerListMap, err := getStructFromLedger(stub, ISSUER_LIST)
	if err != nil {
		// error retrieving issuer-list
		return shim.Error(err.Error())
	}
	// Convert map[string] to IssuerList struct
	FillStruct(issuerListMap, &issuerList)

	// Find issuer position in issuerList
	for i, issuerSum := range issuerList.IssuerSummary {
		if strings.Compare(issuerSum.Email, issuerEmail) == 0 {
			logger.Infof("Found Issuer: " + issuerEmail)
			issuerExists = true
			issuerIndexInIssuerList = i
		}
	}
	if !issuerExists {
		logger.Errorf("Issuer doesn't exist")
		return shim.Error("Issuer doesn't exist")
	}

	// Append certID to certIDs slice
	issuerList.IssuerSummary[issuerIndexInIssuerList].CertIDs =
		append(issuerList.IssuerSummary[issuerIndexInIssuerList].CertIDs, certID)

	// Write the state into the ledger (KEY: issuer-list, unique)
	err = marshalAndPutState(stub, issuerList, ISSUER_LIST)
	if err != nil {
		// error marshaling or putting state into ledger
		return shim.Error(err.Error())
	}

	returnMessage = "Successfully updated blockchain: CREATED Certificate and UPDATED issuerList"
	return shim.Success([]byte(returnMessage))
}
