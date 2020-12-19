package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	shell "github.com/ipfs/go-ipfs-api"
)
// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
return shim.Success(nil)
}

func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string
	var err error
	if fn == "set-to-ipfs" {
		result, err = setToIpfs(stub, args)
	} else if fn == "get-from-ipfs" {
		result, err = getFromIpfs(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(result))
}

func setToIpfs(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a value")
	}
	sh := shell.NewShell("ipfs_host:5001")
	// ipfs add
	cid, err := sh.Add(strings.NewReader(args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return "", err
	}
	// show hash
	fmt.Println("added ", cid)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return cid, nil
}

func getFromIpfs(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a hash")
	}
	sh := shell.NewShell("ipfs_host:5001")
	hashKey := args[0]
	// ipfs cat
	catResult, err := sh.Cat(hashKey)
	defer catResult.Close()
	if err != nil {
	fmt.Fprintf(os.Stderr, "error: %s", err)
		return "", err
	}
	// read string from buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(catResult)
	message := buf.String()
	return message, nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
	fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
