package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *SimpleChaincode) initLedger(stub shim.ChaincodeStubInterface) pb.Response {

	// Get issuer-list state from the ledger
	KeyValBytes, err := stub.GetState(ISSUER_LIST)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for issuer-list \"}"
		return shim.Error(jsonResp)
	}

	if KeyValBytes != nil {
		// if issuer-list exists do nothing
		jsonResp := "{\"Error\":\"issuer-list is already created, aborting!\"}"
		logger.Infof(jsonResp)
		return shim.Error(jsonResp)
	}

	// issuer-list doesn't exist, creating
	logger.Infof("Empty issuer-list, creating...")
	issuerList := IssuerList{}

	// Marshal issuerList
	out, marshalErr := json.Marshal(issuerList)

	if marshalErr != nil {
		// error marshaling list
		return shim.Error(marshalErr.Error())
	}

	// Write the state into the ledger
	err = stub.PutState(ISSUER_LIST, []byte(out))
	if err != nil {
		return shim.Error(err.Error())
	}

	returnMessage := "Successfully created Issuer List"
	logger.Infof(returnMessage)
	return shim.Success([]byte(returnMessage))

}

// Marshal an elem and upload to ledger using the provided key string
func marshalAndPutState(stub shim.ChaincodeStubInterface, elem interface{}, key string) error {

	// Marshal
	out, marshalErr := json.Marshal(elem)
	if marshalErr != nil {
		// error marshaling
		return errors.New(marshalErr.Error())
	}

	// Write the state into the ledger (KEY, unique)
	err := stub.PutState(key, []byte(out))
	if err != nil {
		return errors.New(err.Error())
	}

	logger.Infof("%T created in ledger with key: %s!", elem, key)

	return nil
}

// Function that retrieves an element from ledger using the provided key
//
// It retrieves a map[string] element, that can be converted to the desired struct
// using the function FillStruct
func getStructFromLedger(stub shim.ChaincodeStubInterface, key string) (map[string]interface{}, error) {

	var element map[string]interface{}

	// Get state from the ledger
	elementBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "Failed to get state for " + key
		logger.Errorf(jsonResp)
		return nil, errors.New(jsonResp)
	}

	if elementBytes == nil {
		// element doesn't exists (return empty map[])
		jsonResp := key + " doesn't exist"
		logger.Warningf(jsonResp)
		return make(map[string]interface{}), nil
	}

	json.Unmarshal(elementBytes, &element)

	return element, nil
}

// Function that convert map[string] to struct
func FillStruct(m map[string]interface{}, s interface{}) {
	j, _ := json.Marshal(m)
	json.Unmarshal(j, s)
}
