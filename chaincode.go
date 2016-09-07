
/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
    "errors"
    "encoding/json"
    "fmt"
    "log"
    "strconv"
    
     "time"
     "math/rand"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
var idPrefix = "Dn:"
var perPrefix = "P:"
type Donation struct {
    Id string `json:"id"`
    Who string `json:"who"`
    Rid string  `json:"rid"`
    Money int   `json:"money"`
}

type Request struct {
    Id string `json:"id"`
    Name string  `json:"name"`
    Description string `json:"description"`
    ExpectedMoney int `json:"expectedMoney"`
    CurrentMoney int  `json:"currentMoney"`
    DonationList []string `json:"donationList"`
}


type Person struct {
    Id string `json:"id"`
    Name string `json:"name"`
    MyRequests []string `json:"myRequests"`
    MyDonations []string `json:"myDonations"`
}



func main() {
    err := shim.Start(new(SimpleChaincode))
    
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}
func (t *SimpleChaincode) createDonation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
     var err error
     var donation Donation
     
     r := rand.New(rand.NewSource(time.Now().UnixNano()))
     ID := strconv.Itoa(r.Intn(10))
      fmt.Println("ID")
     donation = Donation{Id: ID,Who: args[0],Rid: args[1],Money: 10000}
     donationBytes,err :=json.Marshal(&donation)
     if err !=nil {
           fmt.Println("error creating donation" + donation.Id)
           return nil, errors.New("Error creating donation " + donation.Id)
     }
    err = stub.PutState(idPrefix+donation.Id,donationBytes)
    fmt.Println("Donation Done")
    return nil, nil
     
}

func(t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

    var request Request
    var donationLts []string
    request = Request{Id: "rid", Name: "Donation Go", Description: "Wanna to go to University", ExpectedMoney: 10000, CurrentMoney: 0, DonationList: donationLts}
    rjson, err := json.Marshal(&request)
    if err != nil {
        return nil, err
    }
    stub.PutState("requestid", rjson)
    log.Println("init function has done!")
    return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
   fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createDonation" {													//initialize the chaincode state, used as reset
		return t.createDonation(stub,args)
	}
	  fmt.Println("invoke did not find func: " + function)
     return nil, errors.New("Received unknown function invocation")
}


func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)
    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    log.Println("Get into read function")
 
    var key, jsonResp string
    var err error

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }
    if valAsbytes == nil {
        return []byte("cannot find the key's value of the chaincode"), nil
    }

    return valAsbytes, nil
}
