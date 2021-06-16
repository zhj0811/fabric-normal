package handler

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric/common/flogging"
)

var logger = flogging.MustGetLogger("handler")

func QueryByKey(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Infof("Get value by key %s", args[0])
	return stub.GetState(args[0])
}

func SaveData(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Infof("Enter .....%s", function)
	err := stub.PutState(stub.GetTxID(), []byte(args[0]))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
