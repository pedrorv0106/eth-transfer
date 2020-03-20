package main

import (
	"time"
	"github.com/eth-transfer/models"
	"github.com/eth-transfer/utils"
	"github.com/eth-transfer/app"
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"context"
	"math/big"
	_ "bytes"
	_ "fmt"
	"strings"
	"encoding/hex"
)
var (
	db  *gorm.DB
	blockchain models.Blockchain
	ethClient *ethclient.Client
	ethRPCClient *rpc.Client
	addresses [] models.Address
)
func main() {
	app.Init_Logger("../log/blockchain")
	app.LoadEnvVars()
	db = app.InitDB(app.GetEnv("DB_URL", ""))
	db.AutoMigrate(&models.Address{})
	db.AutoMigrate(&models.Blockchain{})
	db.AutoMigrate(&models.FeeTransaction{})
	db.First(&blockchain, "currency = ?", "eth")
	if blockchain.ID == 0 {
		blockchain = models.Blockchain{Currency: "eth", Height: utils.Stoi(app.GetEnv("BLOCK_NUMBER", "7549517"))}
		db.Create(&blockchain)
	}
	db.Find(&addresses)
	ethClient = utils.InitEthClient(app.GetEnv("INFURA_URL", ""))
	ethRPCClient = utils.InitEthRPCClient(context.Background(), app.GetEnv("INFURA_URL", ""))
	
	forever := make(chan bool)
	go func() {
		for {
			process_blockchain()
			time.Sleep(5 * time.Second)
		}
	}()

	logger.Infof(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func process_blockchain() {
	blocksLimit := 10
	block_number, err := utils.GetBlockNumber()
	latestBlockNumber := (int)(block_number.Int64());
	if err != nil {
		logger.Fatal(err)
		return
	}

	if blockchain.Height >= latestBlockNumber {
		logger.Info("Skip synchronization. No new blocks detected height: ", blockchain.Height, "latest_block: ", latestBlockNumber)
		return
	}
	fromBlock := blockchain.Height
	toBlock := utils.Ternary(latestBlockNumber > fromBlock + blocksLimit, fromBlock + blocksLimit, latestBlockNumber).(int)
	for blockID := fromBlock; blockID <= toBlock; blockID ++ {
		block, err := ethClient.BlockByNumber(context.Background(), big.NewInt(int64(blockID)))
		if (err != nil) {
			logger.Fatal(err)
		}
		build_deposits(block)
		
		logger.Info("Finished processing in block number ", blockID)
	}
	blockchain.Height = toBlock + 1
	db.Save(&blockchain)
}

func build_deposits(block *types.Block) ([]string){
	depositTxids := make([]string, 0)
	transactions := block.Transactions()
	for _, transaction := range transactions {
		txid := transaction.Hash().Hex()
		for _, address := range addresses {
			data := hex.EncodeToString(transaction.Data())
			if len(data) == 136 && data[0:8] == "a9059cbb" && strings.Contains(data, strings.Replace(strings.ToLower(address.Address), "0x", "", -1)) {
				depositTxids = append(depositTxids, txid)
				var feeTransactions [] models.FeeTransaction
				db.Where("txid = ?", txid).Find(&feeTransactions)
				if len(feeTransactions) > 0 {
					continue;
				}
				gas := int64(100000 * 3000000000)                // 100000(gas Limit) * 5gwei
				feeTxid, _ :=utils.SendETH(address.Address, big.NewInt(gas), app.GetEnv("COLD_WALLET_PRIVATE_KEY", ""))
				db.Create(&models.FeeTransaction{Txid: txid, FeeTxid: feeTxid, Amount: data[72:136], State: "pending", PrivKey: address.PrivKey})
			}
		}
	}
	return depositTxids
}
