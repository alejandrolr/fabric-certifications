package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) getCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var key string // key to retrieve info
	var err error

	if len(args) != 1 {
		logger.Infof("key = %s, len = %d\n", args, len(args))
		return shim.Error("Incorrect number of arguments. Expecting key to query")
	}

	key = args[0]

	// Get the state from the ledger
	KeyValBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	if KeyValBytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + key + "\",\"Value\":\"" + string(KeyValBytes) + "\"}"
	logger.Infof("Query Response:%s\n", jsonResp)
	return shim.Success(KeyValBytes)
}
