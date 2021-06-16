package sdk

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/viper"
)

const (
	channelID      = "mychannel"
	ccID = "mycc"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
	peer1          = "peer0.org1.example.com"
)

var (
	client *channel.Client
)

//InitSDKs 初始化sdk
func InitSDKs(configFilePath string) error {
	configOpt := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		return fmt.Errorf("fatal to create new sdk: %s", err.Error())
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("User1"), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err = channel.New(clientChannelContext)
	if err != nil {
		return fmt.Errorf("failed to create new channel client: %s", err)
	}

	viper.SetEnvPrefix("core")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigFile(configFilePath)

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error when initializing SDK config: %s", err)
	}
	fmt.Println("viper read in config success.")
	return nil
}

//Invoke invoke
func Invoke(in []string) (channel.Response, error) {
	return client.Execute(channel.Request{ChaincodeID: ccID, Fcn: in[0], Args:  [][]byte{[]byte(in[1])}},
		channel.WithRetry(retry.DefaultChannelOpts))
}

//Query 查询
func Query(in []string) ([]byte, error) {
	response, err := client.Query(channel.Request{ChaincodeID: ccID, Fcn: in[0],Args:  [][]byte{[]byte(in[1])}},
		channel.WithRetry(retry.DefaultChannelOpts),
	)
	if err != nil {
		return nil, err
	}
	return response.Payload, nil
}

//GetUserId 获取userid
func GetUserId() string {
	return viper.GetString("user.id")
}

//func GetBlockHeightByTxID(txId string) (uint64, error) {
//	block, err := hfcSdk.GetBlockByTxID(txId, "")
//	if err != nil {
//		return 0, err
//	}
//	return block.BlockNum, nil
//}
//
//func GetFilterTxByTxID(txId string) (*parseBlock.FilterTx, error) {
//	tx, err := hfcSdk.GetFilterTxByTxID(txId, "")
//	if err != nil {
//		return nil, err
//	}
//	return tx, nil
//}
