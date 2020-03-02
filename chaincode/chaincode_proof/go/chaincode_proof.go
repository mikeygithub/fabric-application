package main

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/hyperledger/fabric-chaincode-go/shim"
//	"strconv"
//	sc "github.com/hyperledger/fabric-protos-go/peer"
//	"time"
//)
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)
//凭证信息结构体
type Proof struct {
	//时间
	Time string `json:time`
	//凭证文件路近
	FilePath string `json:filepath`
	//文件hash
	HashCode string `json:hashcode`
	//持有者
	Owner string `json:owner`
	//是否过期
	Overdue bool `json:overdue`
}

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
//type Car struct {
//	Make   string `json:"make"`
//	Model  string `json:"model"`
//	Colour string `json:"colour"`
//	Owner  string `json:"owner"`
//}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryProof" {
		return s.queryProof(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createProof" {
		return s.createProof(APIstub, args)
	} else if function == "queryAllProof" {
		return s.queryAllProof(APIstub)
	//} else if function == "changeProofOwner" {
		//return s.changeProofOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryProof(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(carAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	proof := []Proof{
		Proof{Time:time.Local.String(),FilePath:"/home/mikey/fabric1",HashCode:"2e1ecb697ab70115c7d5113af2779d1ba05bf800f72ec5c2566a14ea50b59723",Owner:"Mikey",Overdue:false},
		Proof{Time:time.Local.String(),FilePath:"/home/mikey/fabric2",HashCode:"3e1ecb697ab70115c7d5113af2779d1ba05bf800f72ec5c2566a14ea50b59723",Owner:"Leo",Overdue:false},
		Proof{Time:time.Local.String(),FilePath:"/home/mikey/fabric3",HashCode:"4e1ecb697ab70115c7d5113af2779d1ba05bf800f72ec5c2566a14ea50b59723",Owner:"Don",Overdue:false},
		Proof{Time:time.Local.String(),FilePath:"/home/mikey/fabric4",HashCode:"5e1ecb697ab70115c7d5113af2779d1ba05bf800f72ec5c2566a14ea50b59723",Owner:"Mike",Overdue:false},
		Proof{Time:time.Local.String(),FilePath:"/home/mikey/fabric5",HashCode:"6e1ecb697ab70115c7d5113af2779d1ba05bf800f72ec5c2566a14ea50b59723",Owner:"Lucy",Overdue:false},
	}

	i := 0
	for i < len(proof) {
		fmt.Println("i is ", i)
		carAsBytes, _ := json.Marshal(proof[i])
		APIstub.PutState("Proof"+strconv.Itoa(i), carAsBytes)
		fmt.Println("Added", proof[i])
		i = i + 1
	}

	return shim.Success(nil)
}
//创建凭证
func (s *SmartContract) createProof(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var proof = Proof{ Time: args[1], FilePath: args[2], HashCode: args[3], Owner: args[4], Overdue: false}
	proofAsBytes, _ := json.Marshal(proof)
	APIstub.PutState(args[0], proofAsBytes)

	return shim.Success(nil)
}
//查询凭证
func (s *SmartContract) queryAllProof(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Proof0"
	endKey := "Proof999"

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
		buffer.WriteString("{\"Key\":")
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

	fmt.Printf("- queryAllProofs:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//func (s *SmartContract) changeProofOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
//
//	if len(args) != 2 {
//		return shim.Error("Incorrect number of arguments. Expecting 2")
//	}
//
//	carAsBytes, _ := APIstub.GetState(args[0])
//	car := Car{}
//
//	json.Unmarshal(carAsBytes, &car)
//	car.Owner = args[1]
//
//	carAsBytes, _ = json.Marshal(car)
//	APIstub.PutState(args[0], carAsBytes)
//
//	return shim.Success(nil)
//}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
