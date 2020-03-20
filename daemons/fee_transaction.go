package main

import (
	"time"
	"github.com/eth-transfer/models"
	"github.com/eth-transfer/utils"
	"github.com/eth-transfer/app"
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"context"
	_ "math/big"
	_ "bytes"
	_ "fmt"
	_ "strings"
	_ "encoding/hex"
)
var (
	db  *gorm.DB
	ethClient *ethclient.Client
	ethRPCClient *rpc.Client
	addresses [] models.Address
	tokenAddress string
	tokenDecimal int
	coldWalletAddress string
)
func main() {
	app.Init_Logger("../log/fee_transaction")
	app.LoadEnvVars()
	db = app.InitDB(app.GetEnv("DB_URL", ""))
	ethClient = utils.InitEthClient(app.GetEnv("INFURA_URL", ""))
	ethRPCClient = utils.InitEthRPCClient(context.Background(), app.GetEnv("INFURA_URL", ""))
	
	forever := make(chan bool)
	go func() {
		for {
			var feeTransactions [] models.FeeTransaction
			db.Where("state = ?", "pending").Find(&feeTransactions)
			for _, feeTransaction := range feeTransactions {
				process(feeTransaction)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	logger.Infof(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	
}

func process(feeTransaction models.FeeTransaction) {
	txReceipt, err:= ethClient.TransactionReceipt(context.Background(), common.HexToHash(feeTransaction.FeeTxid))
	if (err != nil) {
		return
	}
	if txReceipt.Status == 0 {
		return
	}
	txid, err := utils.SendToken(app.GetEnv("COLD_WALLET_ADDRESS", ""), feeTransaction.Amount, app.GetEnv("TOKEN_ADDRESS", ""), feeTransaction.PrivKey)
	if (err != nil) {
		return;
	}
	logger.Info(txid)
	feeTransaction.State = "confirmed"
	db.Save(&feeTransaction)
}

