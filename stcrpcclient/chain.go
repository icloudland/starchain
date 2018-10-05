// Copyright (c) 2014-2017 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package stcrpcclient

import (
	"encoding/json"

	"github.com/icloudland/starchain/stcjson"
	"github.com/icloudland/starchain/net/rpchttp"
	"fmt"
)

// FutureGetBlockCountResult is a future promise to deliver the result of a
// GetBlockCountAsync RPC invocation (or an applicable error).
type FutureGetBlockCountResult chan *response

// Receive waits for the response promised by the future and returns the number
// of blocks in the longest block chain.
func (r FutureGetBlockCountResult) Receive() (int64, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal the result as an int64.
	var count int64
	err = json.Unmarshal(res, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetBlockCountAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetBlockCount for the blocking version and more details.
func (c *Client) GetBlockCountAsync() FutureGetBlockCountResult {
	cmd := stcjson.NewGetBlockCountCmd()
	return c.sendCmd(cmd)
}

// GetBlockCount returns the number of blocks in the longest block chain.
func (c *Client) GetBlockCount() (int64, error) {
	return c.GetBlockCountAsync().Receive()
}

type FutureGetBlockResult chan *response

func (r FutureGetBlockResult) Receive() (*rpchttp.BlockInfo, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(res[:]))
	var blockInfo rpchttp.BlockInfo
	err = json.Unmarshal(res, &blockInfo)
	if err != nil {
		return nil, err
	}
	return &blockInfo, nil

}

func (c *Client) GetBlockAsync(blockHeight int64) FutureGetBlockResult {
	cmd := stcjson.NewGetBlockCmd(blockHeight)
	return c.sendCmd(cmd)
}

func (c *Client) GetBlock(blockHeight int64) (*rpchttp.BlockInfo, error) {
	return c.GetBlockAsync(blockHeight).Receive()
}

type FutureGetRawTransactionResult chan *response

func (r FutureGetRawTransactionResult) Receive() (*rpchttp.Transactions, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(res[:]))
	var tx rpchttp.Transactions
	err = json.Unmarshal(res, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (c *Client) GetRawTransactionAsync(txHash string) FutureGetRawTransactionResult {

	cmd := stcjson.NewGetRawTransactionCmd(txHash)
	return c.sendCmd(cmd)
}

func (c *Client) GetRawTransaction(txHash string) (*rpchttp.Transactions, error) {
	return c.GetRawTransactionAsync(txHash).Receive()
}

type FutureSendRawTransactionResult chan *response

func (r FutureSendRawTransactionResult) Receive() (string, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return "", err
	}
	fmt.Println(string(res[:]))
	var tx string
	err = json.Unmarshal(res, &tx)
	if err != nil {
		return "", err
	}
	return tx, nil
}

func (c *Client) SendRawTransactionAsync(rawTx string) FutureSendRawTransactionResult {

	cmd := stcjson.NewSendRawTransactionCmd(rawTx)
	return c.sendCmd(cmd)
}

func (c *Client) SendRawTransaction(rawTx string) (string, error) {
	return c.SendRawTransactionAsync(rawTx).Receive()
}
