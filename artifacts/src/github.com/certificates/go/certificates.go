package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("certificates")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### certificates chaincode Init ###########")

	return shim.Success(nil)
}

// Invoke function that receives invokes and queries
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### certificates Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	// check if user is an university
	err := checkUserAttr(stub, "role", "university")

	if err != nil {
		// user doesn't have the expected role or attrs retrieval error
		return shim.Error(err.Error())
	}

	// check if user has email
	val, attrErr := getUserAttr(stub, "email")

	if attrErr != nil {
		// error getting user's email
		return shim.Error(attrErr.Error())
	}

	if function == "initLedger" {
		// init empty structures in ledger
		return t.initLedger(stub)
	}
	if function == "issueBadge" {
		// Add email to arguments at 1st position
		arguments := append([]string{val}, args...)
		// issue a Badge
		return t.issueBadge(stub, arguments)
	}
	if function == "issueCertificate" {
		// Add email to arguments at 1st position
		arguments := append([]string{val}, args...)
		// Issue a certificate
		return t.issueCertificate(stub, arguments)
	}
	if function == "getCertificate" {
		// Add email to arguments at 1st position
		// arguments := []string{val}

		// get a certificate
		return t.getCertificate(stub, args)
	}

	errorMsg := "Unknown action, check the first argument, must be one of 'setCertificate' or 'getCertificate'. But got:" + args[0]
	logger.Errorf(errorMsg)
	return shim.Error(errorMsg)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
