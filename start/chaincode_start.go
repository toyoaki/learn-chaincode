/*
Copyright IBM Corp 2016 All Rights Reserved.

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
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Printf("Oi Pati 1")
	fmt.Sprintf("Oi Pati 2")
	fmt.Println("Oi Pati 3")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "ShipperShip" {													//initialize the chaincode state, used as reset
		return t.ShipperShip(stub, "init", args)
	} else if function == "LogisticProviderShip" {													//initialize the chaincode state, used as reset
		return t.LogisticProviderShip(stub, "init", args)	
	}

	fmt.Println("invoke did not find func: " + function)					//error
	return nil, errors.New("Received unknown function invocation: " + function)
}

// Salvar dados do embarcador
func (t *SimpleChaincode) ShipperShip(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// salvar dados físicos
	return t.SendEvent(stub, "ShipperShip", args);
}

func (t *SimpleChaincode) LogisticProviderShip(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// salvar dados físicos
	return t.SendEvent(stub, "LogisticProviderShip", args);
}

// Simulate event sending
func (t *SimpleChaincode) SendEvent(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return t.Quote(stub, function, args);
}

func (t *SimpleChaincode) Quote(stub shim.ChaincodeStubInterface, source string, args []string) ([]byte, error) {
	return nil, errors.New("Quote: " + source)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("hi there " + function)						//error
		return nil, nil;
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}
