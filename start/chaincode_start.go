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
	"time"
	"strconv"
	"encoding/json"	
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

func main() {
	fmt.Println("[IP] Start Contract")

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
	fmt.Println("[IP] invoke is running " + function)

	
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "ShipperShip" {
		orderId, err := t.ShipperShip(stub, args)
		return []byte(orderId), err
	} else if function == "LogisticProviderShip" {
		orderId, err := t.LogisticProviderShip(stub, args)
		return []byte(orderId), err
	}

	fmt.Println("[IP] invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Salvar dados do embarcador
func (t *SimpleChaincode) ShipperShip(stub shim.ChaincodeStubInterface, args []string) (string, error) {	
	fmt.Println("[IP][ShipperShip] save")

	if len(args) != 9 {
		return "", errors.New("[IP][ShipperShip] Incorrect number of arguments. Expecting 9")
	}

	var order Order
	var err error
	var bytes []byte

	order.OrderId 			      = "order-" + strconv.FormatInt(time.Now().Unix(), 10)
    order.ClientId 			      = args[0]
    order.LogisticProviderId      = args[1]
    order.InvoiceValue, err       = strconv.Atoi(string(args[2]))
    order.OriginZipCode           = args[3]
    order.DestinationZipCode      = args[4]
    order.ClientWeight, err       = strconv.Atoi(string(args[5]))
    order.ClientWidth, err        = strconv.Atoi(string(args[6]))
    order.ClientLength, err       = strconv.Atoi(string(args[7]))
    order.ClientHeight, err       = strconv.Atoi(string(args[8]))

    fmt.Println("[IP][ShipperShip] Args parsed for OrderId: " + order.OrderId)

	bytes, err = json.Marshal(order)

	if err != nil { 		
		return "", errors.New("[IP][ShipperShip] Error marshalling order")
	}

	fmt.Println("[IP][ShipperShip] Order marshalled for OrderId" + order.OrderId)

	err = stub.PutState(order.OrderId, bytes)

	fmt.Println("[IP][ShipperShip] PutState for OrderId: " + order.OrderId)

	if err != nil { 
		return "", errors.New("[ShipperShip] Unable to put the state") 
	}
    
	return t.SendEvent(stub, "ShipperShip", order.OrderId);
}

// Salvar dados do embarcador
func (t *SimpleChaincode) LogisticProviderShip(stub shim.ChaincodeStubInterface, args []string) (string, error) {	
	fmt.Println("[IP][LogisticProviderShip] save")

	if len(args) != 5 {
		return "", errors.New("[IP][LogisticProviderShip] Incorrect number of arguments. Expecting 5")
	}

	order, err := t.FindOrderById(stub, args[0])

    order.LogisticProviderWeight, err       = strconv.Atoi(string(args[1]))
    order.LogisticProviderWidth,  err       = strconv.Atoi(string(args[2]))
    order.LogisticProviderLength, err       = strconv.Atoi(string(args[3]))
    order.LogisticProviderHeight, err       = strconv.Atoi(string(args[4]))

	bytes, err := json.Marshal(order)

	if err != nil { 		
		return "", errors.New("[IP][LogisticProviderShip] Error marshalling order")
	}

	fmt.Println("[IP][LogisticProviderShip] Order marshalled for OrderId" + order.OrderId)

	err = stub.PutState(order.OrderId, bytes)

	fmt.Println("[IP][LogisticProviderShip] PutState for OrderId: " + order.OrderId)

	if err != nil { 
		return "", errors.New("[IP][LogisticProviderShip] Unable to put the state") 
	}
    
	return t.SendEvent(stub, "LogisticProviderShip", order.OrderId);
}

// Simulate event sending
func (t *SimpleChaincode) SendEvent(stub shim.ChaincodeStubInterface, source string, orderId string) (string, error) {
	fmt.Println("[IP][SendEvent] Send event with source " + source + " for OrderId: " + orderId)
	
	if source == "ShipperShip" {
		return t.QuoteForShipper(stub, orderId)
	} else if source == "LogisticProviderShip" {
		return t.QuoteForLogisticProvider(stub, orderId)	
	}

	return orderId, errors.New("[IP][SendEvent] Unknown source: " + source) 
}

func (t *SimpleChaincode) QuoteForShipper(stub shim.ChaincodeStubInterface, orderId string) (string, error) {
	fmt.Println("[IP][Quote] for orderId: " + orderId)

	order, err := t.FindOrderById(stub, orderId)

	order.ClientFinalShippingCost = order.InvoiceValue * (order.ClientWeight + order.ClientWidth + order.ClientLength + order.ClientHeight)

	fmt.Println("[IP][Quote] calculated final shipping cost: " + strconv.Itoa(order.ClientFinalShippingCost) + " for orderId: " + orderId)

	orderAsBytes, err := json.Marshal(order)	
	
	err = stub.PutState(order.OrderId, orderAsBytes)

	fmt.Println("[IP][Quote] saved order with orderId: " + orderId)

	return orderId, err
}

func (t *SimpleChaincode) QuoteForLogisticProvider(stub shim.ChaincodeStubInterface, orderId string) (string, error) {
	fmt.Println("[IP][QuoteForLogisticProvider] for orderId: " + orderId)

	order, err := t.FindOrderById(stub, orderId)

	order.LogisticProviderFinalShippingCost = order.InvoiceValue * (order.LogisticProviderWeight + order.LogisticProviderWidth + order.LogisticProviderLength + order.LogisticProviderHeight)

	fmt.Println("[IP][QuoteForLogisticProvider] calculated final shipping cost: " + strconv.Itoa(order.LogisticProviderFinalShippingCost) + " for orderId: " + orderId)

	orderAsBytes, err := json.Marshal(order)	
	
	err = stub.PutState(order.OrderId, orderAsBytes)

	fmt.Println("[IP][QuoteForLogisticProvider] saved order with orderId: " + orderId)

	return orderId, err
}

func (t *SimpleChaincode) FindOrderById(stub shim.ChaincodeStubInterface, orderId string) (Order, error){
	fmt.Println("[IP][FindOrderById] orderId: " + orderId)

	var order Order
	var valAsBytes []byte
	valAsBytes, err := stub.GetState(orderId)

	if err != nil {		
		return order, errors.New("[Quote] Unknown error")
	}
	
	json.Unmarshal(valAsBytes, &order)

	fmt.Println("[IP][FindOrderById] Order with orderId: " + orderId + " found")

	return order, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("[IP][Query] Incorrect number of arguments. Expecting 4")
	}

	if function == "findByClientIdAndLogisticProviderId" {
		resultsIterator, err := stub.RangeQueryState("order-0", "order-99999999999999")

		if err != nil {
			return nil, errors.New("[IP][Query] Unknown error")
		}

		clientId := args[0]
		logisticProviderId := args[1]
		pendingOrder, err := strconv.ParseBool(args[2])
		findAll, err := strconv.ParseBool(args[3])

		hasResult := false

		defer resultsIterator.Close()

		result := "{\"orders\": ["
		
		for resultsIterator.HasNext() {
			queryKeyAsStr, queryValAsBytes, err := resultsIterator.Next()

			fmt.Println("[IP][Query] hack: " + queryKeyAsStr)

			if err != nil {
				return nil, errors.New("[IP][Query] Unknown error")
			}

			var order Order
			json.Unmarshal(queryValAsBytes, &order)

			clientIdOk := clientId == "-1" || order.ClientId == clientId			 				
			logisticProviderIdOk := logisticProviderId == "-1" || order.LogisticProviderId == logisticProviderId
			
			var findPendingOk bool

			if pendingOrder {
				findPendingOk = (order.LogisticProviderFinalShippingCost == 0)
			} else {
				findPendingOk = (order.LogisticProviderFinalShippingCost != 0)
			}

			fmt.Println("[IP][Query] ClientId: " + clientId + " | LogisticProviderId: " + logisticProviderId + " | PendingOrder: " + strconv.FormatBool(pendingOrder) + " | ClientIdOk: " + strconv.FormatBool(clientIdOk) + " | LogisticProviderIdOk: " + strconv.FormatBool(logisticProviderIdOk) + " | FindPendingOk: " + strconv.FormatBool(findPendingOk) + " | FindAll: " + strconv.FormatBool(findAll))

			if(findAll || (clientIdOk && logisticProviderIdOk && findPendingOk)) {
				result += string(queryValAsBytes) + ","
				hasResult = true
			}
		}

		if hasResult {
			result = result[:len(result)-1] + "]}"
		} else {
			result = result + "]}"
		}

		return []byte(result), nil
	}

	return nil, errors.New("[IP][Query] Received unknown function query: " + function)
}
