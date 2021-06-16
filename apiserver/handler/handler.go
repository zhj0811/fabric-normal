package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/zhj0811/fabric-normal/apiserver/common"
	"github.com/zhj0811/fabric-normal/apiserver/sdk"
	"github.com/zhj0811/fabric-normal/common/logging"
)

var logger = logging.NewLogger("debug", "handler")


type Data struct{
	Key string `json:"key"`
	Value string `json:"value"`
}

func Invoke(ctx *gin.Context) {
	req := &Data{}
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		logger.Errorf("Read request info failed %s", err.Error())
		Response(ctx, err, common.RequestFormatErr, nil)
		return
	}
	bytes, err := json.Marshal(req)
	if err !=nil {
		logger.Errorf("Marshal request info failed %s", err.Error())
		Response(ctx, err, common.MarshalJSONErr, nil)
		return
	}
	args := []string{"invoke", string(bytes)}
	res, err := sdk.Invoke(args)
	if err != nil {
		logger.Errorf("Chaincode invoke failed %s", err.Error())
		Response(ctx, err, common.InvokeErr, nil)
		return
	}
	logger.Infof("Upload data %s success", req.Key)
	Response(ctx, nil, common.Success, res.TransactionID)
	return
}

func Query(ctx *gin.Context) {
	tx := ctx.Query("key")
	args := []string{"query", string(tx)}
	res, err := sdk.Query(args)
	if err != nil {
		logger.Errorf("Chaincode query failed %s", err.Error())
		Response(ctx, err, common.QueryErr, nil)
		return
	}
	logger.Infof("Query res %+v", res)
	Response(ctx, nil, common.Success, res)
	return
}

