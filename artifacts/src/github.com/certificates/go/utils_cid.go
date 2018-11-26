/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Function that checks user's attributes.
// Arguments:
// - stub (shim.ChaincodeStubInterface)
// - attribute to search
// - desired role as a string
//
// Return:
// - nil error if user has the attribute-value
// - Error string if not
func checkUserAttr(stub shim.ChaincodeStubInterface, attr, desiredAttr string) error {
	// Check user's attr
	val, ok, err := cid.GetAttributeValue(stub, attr)
	var errMessage string

	if err != nil {
		// There was an error trying to retrieve the attribute
		errMessage = "Error trying to retrieve the " + attr + " attribute"
		return errors.New(errMessage)
	}
	if !ok {
		// The client identity does not possess the attribute
		errMessage = "User do not have " + attr + " attribute"
		return errors.New(errMessage)
	}

	if strings.Compare(val, desiredAttr) != 0 {
		// user is not an university
		errMessage = "User " + attr + " is not equal to " + desiredAttr
		return errors.New(errMessage)
	}

	// logger.Infof("val = %s, ok = %s, err = %s\n", val, ok, err)
	return nil
}

// Function that get user's attributes.
// Arguments:
// - stub (shim.ChaincodeStubInterface)
// - attribute to search
//
// Return (val, err)
// - (attributeValue, nil) if OK
// - ("", Error) if Error
func getUserAttr(stub shim.ChaincodeStubInterface, attr string) (string, error) {
	// Check user's attr
	val, ok, err := cid.GetAttributeValue(stub, attr)
	var errMessage string

	if err != nil {
		// There was an error trying to retrieve the attribute
		errMessage = "Error trying to retrieve the '" + attr + "' attribute"
		return "", errors.New(errMessage)
	}
	if !ok {
		// The client identity does not possess the attribute
		errMessage = "User does not have '" + attr + "' attribute"
		return "", errors.New(errMessage)
	}

	if len(val) <= 0 {
		// attribute is empty
		errMessage = "User '" + attr + "' attribute is empty"
		return "", errors.New(errMessage)
	}

	return val, nil
}
