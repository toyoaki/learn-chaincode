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
	"math/rand"	
	"strconv"
	"encoding/json"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"	
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Order struct {	
	OrderId                             string      `json:"order_id"`
    ClientId                            string      `json:"client_id"`
    LogisticProviderId                  string      `json:"logistic_provider_id"`
    InvoiceValue                        int         `json:"invoice_value"`
    OriginZipCode                       string      `json:"origin_zip_code"`
    DestinationZipCode                  string      `json:"destination_zip_code"`    
    ClientWeight                        int         `json:"client_weight"`
    ClientWidth                         int         `json:"client_width"`
    ClientLength                        int         `json:"client_length"`
    ClientHeight                        int         `json:"client_height"`
    ClientFinalShippingCost             int         `json:"client_final_shipping_cost"`
    LogisticProviderWeight              int         `json:"logistic_provider_weight"`
    LogisticProviderWidth               int         `json:"logistic_provider_width"`
    LogisticProviderLength              int         `json:"logistic_provider_length"`
    LogisticProviderHeight              int         `json:"logistic_provider_height"`
    LogisticProviderFinalShippingCost   int         `json:"logistic_provider_final_shipping_cost"`
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
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "ShipperShip" {
		return t.ShipperShip(stub, args)
	} else if function == "LogisticProviderShip" {
		return t.LogisticProviderShip(stub, args)	
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Salvar dados do embarcador
func (t *SimpleChaincode) ShipperShip(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {	
	if len(args) != 9 {
		return nil, errors.New("[ShipperShip] Incorrect number of arguments. Expecting 9")
	}

	var order Order
	var err error
	var bytes []byte

	order.OrderId 			      = "order-" + strconv.Itoa(rand.Intn(100000000))
    order.ClientId 			      = args[0]
    order.LogisticProviderId      = args[1]
    order.InvoiceValue, err       = strconv.Atoi(string(args[2]))
    order.OriginZipCode           = args[3]
    order.DestinationZipCode      = args[4]
    order.ClientWeight, err       = strconv.Atoi(string(args[5]))
    order.ClientWidth, err        = strconv.Atoi(string(args[6]))
    order.ClientLength, err       = strconv.Atoi(string(args[7]))
    order.ClientHeight, err       = strconv.Atoi(string(args[8]))

	bytes, err = json.Marshal(order)
	
	if err != nil { 		
		return nil, errors.New("[ShipperShip] Error marshalling order")
	}

	err = stub.PutState("order", bytes)

	if err != nil { 
		return nil, errors.New("[ShipperShip] Unable to put the state") 
	}
    
	return t.SendEvent(stub, "ShipperShip", args);
}

func (t *SimpleChaincode) LogisticProviderShip(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {	
	return t.SendEvent(stub, "LogisticProviderShip", args);
}

// Simulate event sending
func (t *SimpleChaincode) SendEvent(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return t.Quote(stub, function, args);
}

func (t *SimpleChaincode) Quote(stub shim.ChaincodeStubInterface, source string, args []string) ([]byte, error) {
	return nil, errors.New("Quote: " + source + " --- " + strings.Join(args," "))
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "findAll" {
		resultsIterator, err := stub.RangeQueryState("order-0", "order-999999999")

		if err != nil {
			return nil, errors.New("[Query] Unknown error")
		}

		defer resultsIterator.Close()

		result := "{orders: ["
		
		for resultsIterator.HasNext() {
			queryKeyAsStr, queryValAsBytes, err := resultsIterator.Next()

			fmt.Println(queryKeyAsStr)

			if err != nil {
				return nil, errors.New("[Query] Unknown error")
			}

			result += string(queryValAsBytes) + ","
		}

		if len(result) == 1 {
			result = "]}"
		} else {
			result = result[:len(result)-1] + "]}"
		}

		return []byte(result), nil
	}

	return nil, errors.New("[Query] Received unknown function query: " + function)
}
