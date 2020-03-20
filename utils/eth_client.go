package utils

import (
	"log"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
	"context"
	"math/big"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"
	"encoding/hex"
)
var (
	rpcClient * rpc.Client
	ethClient *ethclient.Client
)
func InitEthClient(INFURA_URL string) *ethclient.Client {
	client, err := ethclient.Dial(INFURA_URL)
	if err != nil {
		log.Fatal(err)
	}
	ethClient = client
	return client
}

func InitEthRPCClient(ctx context.Context, rawUrl string) *rpc.Client {
	c, err := rpc.DialContext(ctx, rawUrl)
	if err != nil {
		log.Fatal(err)
	}
	rpcClient = c
	return rpcClient
}

func GetBlockNumber() (*big.Int, error) {
	var result hexutil.Big
	err := rpcClient.CallContext(context.Background(), &result, "eth_blockNumber")
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}

func SendETH(to string, amount *big.Int, privKey string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privKey)
    if err != nil {
        return "", err
    }

    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        return "", errors.New("error casting public key to ECDSA")
    }

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        return "", err
    }

    value := amount // in wei (1 eth)
    gasLimit := uint64(21000)                // in units
    gasPrice, err := ethClient.SuggestGasPrice(context.Background())
    if err != nil {
        return "", err
    }

    toAddress := common.HexToAddress(to)
    var data []byte
    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

    chainID, err := ethClient.NetworkID(context.Background())
    if err != nil {
        return "", err
    }

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        return "", err
    }

    err = ethClient.SendTransaction(context.Background(), signedTx)
    if err != nil {
        return "", err
	}
	return signedTx.Hash().Hex(), nil
}

func SendToken(to string, amount string, tokenAddress string, privKey string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privKey)
    if err != nil {
        return "", err
    }

    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        return "", errors.New("error casting public key to ECDSA")
    }

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        return "", err
    }

    value := big.NewInt(0) // 
    gasLimit := uint64(100000)                // in units
    gasPrice, err := ethClient.SuggestGasPrice(context.Background())
    if err != nil {
        return "", err
    }

    toAddress := common.HexToAddress(tokenAddress)
	
	hexData := fmt.Sprintf("%s%064s%s", "a9059cbb", strings.Replace(to, "0x", "", -1), amount)
	data, _ := hex.DecodeString(hexData)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

    chainID, err := ethClient.NetworkID(context.Background())
    if err != nil {
        return "", err
    }

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        return "", err
    }

    err = ethClient.SendTransaction(context.Background(), signedTx)
    if err != nil {
        return "", err
	}
	return signedTx.Hash().Hex(), nil
}
