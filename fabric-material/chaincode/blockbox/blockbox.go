/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Vehicle structure, with 4 properties.  Structure tags are used by encoding/json library
type Vehicle struct {
	ID         string `json:"id" bson:"_id" `
	Department string `json:"department" bson:"department"`
	Vin        string `json:"vin" bson:"vin"`
	Picture    string `json:"picture" bson:"picture"`
	StaticInfo     StaticInfo      `json:"staticInfo" bson:"staticInfo" `
    TrackingInfo   TrackingInfo    `json:"trackingInfo"  bson:"trackingInfo" `
    SystemWarnings []SystemWarning `json:"warnings" bson:"warnings"`
    SystemErrors   []SystemWarning `json:"errors" bson:"errors"`

}
type SystemWarning struct {
	Cause    string    `json:"cause" bson:"cause"`
	Date     time.Time `json:"date" bson:"date"`
	Critical bool      `json:"critical" bson:"critical"`
}

type TrackingInfo struct {
	MILTime               []IntTracking    `json:"milTime" bson:"milTime"`
	BatteryChargeLevel    []FloatTracking  `json:"batteryChargeLevel" bson:"batteryChargeLevel"`
	BatteryChargingStatus []StringTracking `json:"batteryChargingStatus" bson:"batteryChargingStatus"`
	Charges               []IntTracking    `json:"charges" bson:"charges"`
	CoolantTemp           []FloatTracking  `json:"coolantTemp" bson:"coolantTemp"`
	EngineLoad            []FloatTracking  `json:"engineLoad" bson:"engineLoad"`
	EngineRuntime         []IntTracking    `json:"engineRuntime" bson:"engineRuntime"`
	FuelLevel             []FloatTracking  `json:"fuelLevel" bson:"fuelLevel"`
	KilowattKM            []IntTracking    `json:"kilowattKM" bson:"kilowattKM"`
	Location              []GeoTracking    `json:"location" bson:"location"`
	LockStatus            []BoolTracking   `json:"lockStatus" bson:"lockStatus"`
	OutsideTemp           []FloatTracking  `json:"outsideTemp" bson:"outsideTemp"`
	Rpm                   []IntTracking    `json:"rpm" bson:"rpm"`
	Speed                 []IntTracking    `json:"speed" bson:"speed"`
	ThrottlePos           []FloatTracking  `json:"throttlePos" bson:"throttlePos"`
}
type IntTracking struct {
	TimeStamp time.Time `json:"t" bson:"t"`
	Value     int       `json:"v" bson:"v"`
}

type FloatTracking struct {
	TimeStamp time.Time `json:"t" bson:"t"`
	Value     float32   `json:"v" bson:"v"`
}

type StringTracking struct {
	TimeStamp time.Time `json:"t" bson:"t"`
	Value     string    `json:"v" bson:"v"`
}

type GeoTracking struct {
	TimeStamp time.Time `json:"t" bson:"t"`
	Value     Location  `json:"v" bson:"v"`
}

type BoolTracking struct {
	TimeStamp time.Time `json:"t" bson:"t"`
	Value     bool      `json:"v" bson:"v"`
}

type Location struct {
	GeoJSONType string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type StaticInfo struct {
	Brand string `json:"brand" bson:"brnad"`
	Consumption float32 `json:"consumption" bson:"consumption"`
	Displacement int `json:"displacement" bson:"displacement"`
	Engine string `json:"engine" bson:"engine"`
	Make time.Time `json:"make" bson:"make"`
	Model string `json:"model" bson:"model"`
	Weight int `json:"weight" bson:"weight"`
}


/*
 * The Init method is called when the Smart Contract "fabVehicle" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabVehicle"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryVehicle" {
		return s.queryVehicle(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAllVehicles" {
		return s.queryAllVehicles(APIstub)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryVehicle(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	VehicleAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(VehicleAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	v1, err := generateTestVehicle("rand1")
	if err != nil {
		shim.Error(err.Error())
	}

	v2, err := generateTestVehicle("rand2")
	if err != nil {
		shim.Error(err.Error())
	}

	Vehicles := []Vehicle{
		v1, v2,
	}

	i := 0
	for i < len(Vehicles) {
		fmt.Println("i is ", i)
		VehicleAsBytes, _ := json.Marshal(Vehicles[i])
		APIstub.PutState(Vehicles[i].ID, VehicleAsBytes)
		fmt.Println("Added", Vehicles[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/* func (s *SmartContract) createVehicle(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var Vehicle = Vehicle{Make: args[1], Model: args[2], Colour: args[3], Owner: args[4]}

	VehicleAsBytes, _ := json.Marshal(Vehicle)
	APIstub.PutState(args[0], VehicleAsBytes)

	return shim.Success(nil)
} */

func (s *SmartContract) queryAllVehicles(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "001"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Id\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllVehicles:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/* func (s *SmartContract) changeVehicleOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	VehicleAsBytes, _ := APIstub.GetState(args[0])
	Vehicle := Vehicle{}

	json.Unmarshal(VehicleAsBytes, &Vehicle)
	Vehicle.Owner = args[1]

	VehicleAsBytes, _ = json.Marshal(Vehicle)
	APIstub.PutState(args[0], VehicleAsBytes)

	return shim.Success(nil)
} */

func generateTestVehicle(s string) (Vehicle, error) {

	h := fnv.New64()
	_, err := h.Write([]byte(s))
	if err != nil {
		return Vehicle{}, err
	}
	randSeed := h.Sum64()
	rand.Seed(int64(randSeed))


	makeDate := time.Unix(rand.Int63n(1567515734), 0)
	updateDate := time.Unix(rand.Int63n(1567515734), 0)

	carN1 := Vehicle{
		ID:             "001",
		Department:     "ABC",
		Vin:            "VW-01-NM",
		Picture:        "Picture1",
	}

	newInfo := TrackingInfo{
		MILTime:               []IntTracking{IntTracking{TimeStamp: updateDate, Value: 300}},
		BatteryChargeLevel:    []FloatTracking{FloatTracking{TimeStamp: updateDate, Value: 40}},
		BatteryChargingStatus: []StringTracking{StringTracking{TimeStamp: updateDate, Value: "charging"}},
		Charges:               []IntTracking{IntTracking{TimeStamp: updateDate, Value: 30}},
		CoolantTemp:           []FloatTracking{FloatTracking{TimeStamp: updateDate, Value: 60.5}},
		EngineLoad:            []FloatTracking{FloatTracking{TimeStamp: updateDate, Value: 400}},
		EngineRuntime:         []IntTracking{IntTracking{TimeStamp: updateDate, Value: 259}},
		FuelLevel:             []FloatTracking{FloatTracking{TimeStamp: updateDate, Value: 80.6}},
		KilowattKM:            []IntTracking{IntTracking{TimeStamp: updateDate, Value: 300}},
		Location: []GeoTracking{GeoTracking{TimeStamp: updateDate, Value: Location{
			GeoJSONType: "location",
			Coordinates: []float64{30, 30},
		}}},
		LockStatus:  []BoolTracking{BoolTracking{TimeStamp: updateDate}},
		OutsideTemp: []FloatTracking{FloatTracking{TimeStamp: updateDate, Value: 40}},
		Rpm:         []IntTracking{IntTracking{TimeStamp: updateDate, Value: 3000}},
		Speed:       []IntTracking{IntTracking{TimeStamp: updateDate, Value: 80}},
		ThrottlePos: []FloatTracking{FloatTracking{TimeStamp: updateDate, Value: 30}},
	}



	carN1.TrackingInfo = newInfo

	carN1.StaticInfo = StaticInfo{
		Brand:        "VW",
		Consumption:  30,
		Displacement: 40,
		Engine:       "e",
		Make:         makeDate,
		Model:        "Polo",
		Weight:       2000,
	}

	sw0 := SystemWarning{Cause: "Oil Pressure", Date: updateDate, Critical: true}
	sw1 := SystemWarning{Cause: "Low Fuel Indicator", Date: updateDate, Critical: false}
	sw2 := SystemWarning{Cause: "ABS Break", Date: updateDate, Critical: false}
	sw3 := SystemWarning{Cause: "High RPM", Date: updateDate, Critical: true}
	sw4 := SystemWarning{Cause: "Oil Pressure", Date: updateDate, Critical: true}

	sysWarnings0 := []SystemWarning{sw0, sw3, sw4}
	sysWarnings1 := []SystemWarning{sw1, sw2, sw3}

	carN1.SystemWarnings = sysWarnings0
	carN1.SystemErrors = sysWarnings1

	return carN1, nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
