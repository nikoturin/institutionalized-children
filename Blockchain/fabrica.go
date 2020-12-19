package main

import(
	"encoding/json"
	"fmt"
	"strconv"
        "os"
        "strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	shell "github.com/ipfs/go-ipfs-api"
	sc "github.com/hyperledger/fabric/protos/peer"

	"bytes"
)

type SmartContract struct {

}

type Carro struct {
	Fabricante string `json:"fabricante"`
	Modelo string `json:"modelo"`
	Color string `json:"color"`
	Propietario string `json:"propietario"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "queryCarro" {
		return s.queryCarro(APIstub, args)
	} else if function == "crearCarro" {
		return s.crearCarro(APIstub, args)
	} else if function == "cambiarPropietario" {
		return s.cambiarPropietario(APIstub, args)
	} else if function == "inicializarLedger" {
		return s.inicializarLedger(APIstub)
	} else if function == "queryTodos" {
		return s.queryTodos(APIstub)
	} else if function == "set-to-ipfs"{
		return s.setToIpfs(APIstub, args)
	} else if function == "get-from-ipfs"{
		return s.getFromIpfs(APIstub,args)
	}

	return shim.Error("El metodo no existe")
}

// IPFS functions methods:

//func setToIpfs(stub shim.ChaincodeStubInterface, args []string) (string, error) {
func (s *SmartContract) setToIpfs(stub shim.ChaincodeStubInterface, args []string) sc.Response {
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


//func getFromIpfs(stub shim.ChaincodeStubInterface, args []string) (string, error) {
func (s *SmartContract) getFromIpfs(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
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


//// END

func (s *SmartContract) queryTodos(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "CAR0"
	endKey := "CAR99"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer

	buffer.WriteString("[")
	flagMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if flagMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{'Llave':'")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("', 'Valor':'")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("'}")
		flagMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	set_results, err := setToIpfs(buffer.String())
	get_hash, err:= getFromIpfs();
	fmt.Printf("query todos: %s", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) crearCarro(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	var carro = Carro{Fabricante: args[1], Modelo: args[2], Color: args[3], Propietario: args[4]}

	carroAsBytes, _ := json.Marshal(carro)

	APIstub.PutState(args[0], carroAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) cambiarPropietario(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	carroAsBytes, _ := APIstub.GetState(args[0])
	carro := Carro{}
	json.Unmarshal(carroAsBytes, &carro)

	carro.Propietario = args[1]

	carroAsBytes, _ = json.Marshal(carro)
	APIstub.PutState(args[0],carroAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryCarro(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Numero de argumentos incorrecto, se espera solo un valor")
	}
	carroAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(carroAsBytes)
}

func (s *SmartContract) inicializarLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	carros := []Carro{
		Carro{Fabricante:"BMW", Modelo:"S7", Color:"Rojo", Propietario:"Cristian"},
		Carro{Fabricante: "Chevrolet", Modelo: "Sail", Color:"Gris", Propietario:"Ramses"},
		Carro{Fabricante: "Tesla", Modelo:"T1", Color: "Rojo", Propietario: "Pedro" },
		Carro{Fabricante: "Apache", Modelo: "Triciclo", Color: "Rojo", Propietario: "Chabelo"},
	}
	i := 0
	for i < len(carros) {
		fmt.Println("i:", i)
		carroAsBytes, _ := json.Marshal(carros[i])
		APIstub.PutState("CAR"+strconv.Itoa(i), carroAsBytes)
		fmt.Println("Carro agregado:", carros[i])
		i = i + 1
	}
	return shim.Success(nil) 
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("Error al momento de crear el smart contract: %s", err)
	}
}
