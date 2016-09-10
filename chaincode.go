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
    "strings"
    "encoding/json"
    "fmt"
    "strconv"
    "math/rand"
    "time"
    "log"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)
var Dnprefix = "Dn:"
var Perprefix = "Per:"
var Reqprefix = "Req:"
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Donation struct {
    Id string `json:"id"`
    Who string `json:"who"`
    Rid string  `json:"rid"`
    Money int   `json:"money"`
}

type Request struct {
    Id string `json:"id"`
    Who string `json:"who"`
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

type AllRequest struct {
    AllRequests []Request `json:"allRequests"`
}



func main() {
    err := shim.Start(new(SimpleChaincode))
    
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}


func(t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    
    var allrsi []Request
    allRs := AllRequest{AllRequests: allrsi}
    allJson,_ := json.Marshal(&allRs)
    stub.PutState("allRequests", allJson)

    var names = [3]string{"Lucy", "Andy", "David"}
    var MyReqs, MyDons []string
    for _, v := range names {
        var person Person
        person = Person{Id: v, Name: v, MyRequests: MyReqs, MyDonations: MyDons}
        pb, err := json.Marshal(&person)
        if err != nil {
            return nil, errors.New("failed to init persons' instance")
        }
        preKey := Perprefix + person.Id
        stub.PutState(preKey, pb)
    } 
    log.Println("init function has done!")
    return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

     if function == "createDonation" {
         return t.createDonation(stub, args)
     }

     if function == "createRequest" {
         return t.createRequest(stub, args)
     }
     return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) createDonation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
     //args: ["jack", "requestid", money] 
     var from, toRid,dId string
     var money int
     var err error
   
     if len(args) != 3 {
         return nil, errors.New("My createDonation. Incorrect number of arguments. Expecting 3")
     }
     from = args[0]
     toRid = args[1]
     money, err = strconv.Atoi(args[2])
     if err != nil {
        return nil, errors.New("money cannot convert to number")
     }

     var donation Donation
    str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    bytes := []byte(str)
    result := []byte{}
    r:= rand.New(rand.NewSource(time.Now().UnixNano()))
   for i := 0; i < 6; i++ {
      result = append(result, bytes[r.Intn(len(bytes))])
   }
     donationId :=string(result)
     donation = Donation{Id: donationId, Rid: toRid, Who: from, Money: money}
     djson, err := json.Marshal(&donation)
     if err != nil {
        return nil, err
     }
     dId = Dnprefix+donation.Id
     stub.PutState(dId, djson)
     
     
     
     var person Person
     var myReqs, myDons []string
    //  update person data
     var perkey = Perprefix+ from
     personByte, err := stub.GetState(perkey)
     if personByte == nil {
        person = Person{Id: from, Name: from, MyRequests: myReqs, MyDonations: myDons}
        myDonations := person.MyDonations
         if myDonations == nil {
         myDonations = make([]string, 0)
             }
         myDonations = append(myDonations, donation.Id)
         person.MyDonations = myDonations
        perJson,err := json.Marshal(&person)
        if err !=nil{
            return nil, errors.New("failed to JSON person instance")    
            }
        stub.PutState(perkey,perJson)
        } else {
        err = json.Unmarshal(personByte, &person)
        if err !=nil{
            return nil, errors.New("failed to Unmarshal person instance")    
        }
        myDonations2 := person.MyDonations
        if myDonations2 == nil {
             myDonations2 = make([]string, 0)
        }
        myDonations2 = append(myDonations2, donation.Id)
        person.MyDonations = myDonations2
        perJson2,err := json.Marshal(&person)
        if err !=nil{
            return nil, errors.New("failed to JSON person instance")    
            }
        stub.PutState(perkey,perJson2)
     } 
     toReid := Reqprefix+ toRid 
     requestByte, err := stub.GetState(toReid)
     if err != nil {
           return nil, errors.New("request did not exist")
     }

     var request Request
     err = json.Unmarshal(requestByte, &request)
     if err != nil {
           return nil, errors.New("failed to Unmarshal request instance")
     }
     request.CurrentMoney += money
     donationList := request.DonationList 
     if donationList == nil {
          donationList = make([]string, 0)
     }
     donationList = append(donationList, donation.Id)
     request.DonationList = donationList
     requestJson,err := json.Marshal(&request)
        if err !=nil{
            return nil, errors.New("failed to JSON person instance")    
            }
    stub.PutState(toReid,requestJson)

    allRis, err := stub.GetState("allRequests")
    var allR AllRequest
    err = json.Unmarshal(allRis, &allR)
    if err != nil {
         return nil, errors.New("failed to Unmarshal AllRequest instance")    
    }
    reques := allR.AllRequests
    for _,v := range reques { 
        isEqu := strings.EqualFold(v.Id, request.Id)
        if isEqu == true {
            v.CurrentMoney += money
            dl2 := v.DonationList
            if dl2 == nil {
                dl2 = make([]string, 0)
            }
            dl2 = append(dl2, donationId)
            v.DonationList = dl2
            break
        }
    }
    allR.AllRequests = reques
    requesJson,err := json.Marshal(&allR)
    stub.PutState("allRequests", requesJson)
    return nil, nil     
}

func (t *SimpleChaincode) createRequest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
     //args: [jack, projectName, description, expectedMoney]
     if len(args) != 4 {
          return nil, errors.New("Incorrect number of arguments. Expecting 4")
     }
     var name, projectName, description string
     var expectedMoney int
     var err error
     name = args[0]
     projectName = args[1]
     description = args[2]
     expectedMoney, err = strconv.Atoi(args[3])
     if err != nil {
        return nil, errors.New("money cannot convert to number")
     }
     
     var request Request
     var dl []string
    str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    bytes := []byte(str)
    result := []byte{}
    r:= rand.New(rand.NewSource(time.Now().UnixNano()))
   for i := 0; i < 6; i++ {
      result = append(result, bytes[r.Intn(len(bytes))])
   }
     requestId :=string(result)
     request = Request{Id: requestId, Who: name, Name: projectName, Description: description, ExpectedMoney: expectedMoney, CurrentMoney: 0, DonationList: dl}
     rj, err := json.Marshal(&request)
     if err != nil {
            return nil, errors.New("failed to Marshal request instance")    
     }
     var rkey string
     rkey = Reqprefix + request.Id
     stub.PutState(rkey, rj)


     perkey := Perprefix + request.Who
     personByte, err := stub.GetState(perkey)
     if err !=nil{
         return nil, errors.New("failed to get person instance")    
     }
     var person Person
     var myReqs, myDons []string
     if personByte == nil {
         person = Person{Id: name, Name: name, MyRequests: myReqs, MyDonations: myDons}
     } else {
        err := json.Unmarshal(personByte, &person)
        if err !=nil{
            return nil, errors.New("failed to Unmarshal person instance")    
        }
     }
     myRes := person.MyRequests
     if myRes == nil {
        myRes = make([]string, 0)
     }
     myRes = append(myRes, request.Id)
     person.MyRequests = myRes
     pj,_ := json.Marshal(person)
     pkey := Perprefix + person.Id
     stub.PutState(pkey, pj)

     allJson, _ := stub.GetState("allRequests")
     var allrs3 AllRequest
     err = json.Unmarshal(allJson, &allrs3)
     if err != nil {
         return nil, errors.New("failed to Unmarshal AllRequest instance")    
     }
     allRs2 := allrs3.AllRequests
     if allRs2 == nil {
         allRs2 = []Request{}
     }
     allRs2 = append(allRs2, request)
     allrs3.AllRequests = allRs2
     allJson2,_ := json.Marshal(&allrs3)
     stub.PutState("allRequests", allJson2)
     
     return nil, nil
}




func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    log.Println("query is running " + function)
    log.Println(function)
    log.Println(args[0])
    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }

    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}

func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
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
