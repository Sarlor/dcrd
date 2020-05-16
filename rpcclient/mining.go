// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpcclient

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil/v3"
	chainjson "github.com/decred/dcrd/rpc/jsonrpc/types/v2"
)

// FutureGenerateResult is a future promise to deliver the result of a
// GenerateAsync RPC invocation (or an applicable error).
type FutureGenerateResult cmdRes

// Receive waits for the response promised by the future and returns a list of
// block hashes generated by the call.
func (r *FutureGenerateResult) Receive() ([]*chainhash.Hash, error) {
	res, err := receiveFuture(r.ctx, r.c)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a list of strings.
	var result []string
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	}

	// Convert each block hash to a chainhash.Hash and store a pointer to
	// each.
	convertedResult := make([]*chainhash.Hash, len(result))
	for i, hashString := range result {
		convertedResult[i], err = chainhash.NewHashFromStr(hashString)
		if err != nil {
			return nil, err
		}
	}

	return convertedResult, nil
}

// GenerateAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See Generate for the blocking version and more details.
func (c *Client) GenerateAsync(ctx context.Context, numBlocks uint32) *FutureGenerateResult {
	cmd := chainjson.NewGenerateCmd(numBlocks)
	return (*FutureGenerateResult)(c.sendCmd(ctx, cmd))
}

// Generate generates numBlocks blocks and returns their hashes.
func (c *Client) Generate(ctx context.Context, numBlocks uint32) ([]*chainhash.Hash, error) {
	return c.GenerateAsync(ctx, numBlocks).Receive()
}

// FutureGetGenerateResult is a future promise to deliver the result of a
// GetGenerateAsync RPC invocation (or an applicable error).
type FutureGetGenerateResult cmdRes

// Receive waits for the response promised by the future and returns true if the
// server is set to mine, otherwise false.
func (r *FutureGetGenerateResult) Receive() (bool, error) {
	res, err := receiveFuture(r.ctx, r.c)
	if err != nil {
		return false, err
	}

	// Unmarshal result as a boolean.
	var result bool
	err = json.Unmarshal(res, &result)
	if err != nil {
		return false, err
	}

	return result, nil
}

// GetGenerateAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetGenerate for the blocking version and more details.
func (c *Client) GetGenerateAsync(ctx context.Context) *FutureGetGenerateResult {
	cmd := chainjson.NewGetGenerateCmd()
	return (*FutureGetGenerateResult)(c.sendCmd(ctx, cmd))
}

// GetGenerate returns true if the server is set to mine, otherwise false.
func (c *Client) GetGenerate(ctx context.Context) (bool, error) {
	return c.GetGenerateAsync(ctx).Receive()
}

// FutureSetGenerateResult is a future promise to deliver the result of a
// SetGenerateAsync RPC invocation (or an applicable error).
type FutureSetGenerateResult cmdRes

// Receive waits for the response promised by the future and returns an error if
// any occurred when setting the server to generate coins (mine) or not.
func (r *FutureSetGenerateResult) Receive() error {
	_, err := receiveFuture(r.ctx, r.c)
	return err
}

// SetGenerateAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See SetGenerate for the blocking version and more details.
func (c *Client) SetGenerateAsync(ctx context.Context, enable bool, numCPUs int) *FutureSetGenerateResult {
	cmd := chainjson.NewSetGenerateCmd(enable, &numCPUs)
	return (*FutureSetGenerateResult)(c.sendCmd(ctx, cmd))
}

// SetGenerate sets the server to generate coins (mine) or not.
func (c *Client) SetGenerate(ctx context.Context, enable bool, numCPUs int) error {
	return c.SetGenerateAsync(ctx, enable, numCPUs).Receive()
}

// FutureGetHashesPerSecResult is a future promise to deliver the result of a
// GetHashesPerSecAsync RPC invocation (or an applicable error).
type FutureGetHashesPerSecResult cmdRes

// Receive waits for the response promised by the future and returns a recent
// hashes per second performance measurement while generating coins (mining).
// Zero is returned if the server is not mining.
func (r *FutureGetHashesPerSecResult) Receive() (int64, error) {
	res, err := receiveFuture(r.ctx, r.c)
	if err != nil {
		return -1, err
	}

	// Unmarshal result as an int64.
	var result int64
	err = json.Unmarshal(res, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// GetHashesPerSecAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetHashesPerSec for the blocking version and more details.
func (c *Client) GetHashesPerSecAsync(ctx context.Context) *FutureGetHashesPerSecResult {
	cmd := chainjson.NewGetHashesPerSecCmd()
	return (*FutureGetHashesPerSecResult)(c.sendCmd(ctx, cmd))
}

// GetHashesPerSec returns a recent hashes per second performance measurement
// while generating coins (mining).  Zero is returned if the server is not
// mining.
func (c *Client) GetHashesPerSec(ctx context.Context) (int64, error) {
	return c.GetHashesPerSecAsync(ctx).Receive()
}

// FutureGetMiningInfoResult is a future promise to deliver the result of a
// GetMiningInfoAsync RPC invocation (or an applicable error).
type FutureGetMiningInfoResult cmdRes

// Receive waits for the response promised by the future and returns the mining
// information.
func (r *FutureGetMiningInfoResult) Receive() (*chainjson.GetMiningInfoResult, error) {
	res, err := receiveFuture(r.ctx, r.c)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a getmininginfo result object.
	var infoResult chainjson.GetMiningInfoResult
	err = json.Unmarshal(res, &infoResult)
	if err != nil {
		return nil, err
	}

	return &infoResult, nil
}

// GetMiningInfoAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetMiningInfo for the blocking version and more details.
func (c *Client) GetMiningInfoAsync(ctx context.Context) *FutureGetMiningInfoResult {
	cmd := chainjson.NewGetMiningInfoCmd()
	return (*FutureGetMiningInfoResult)(c.sendCmd(ctx, cmd))
}

// GetMiningInfo returns mining information.
func (c *Client) GetMiningInfo(ctx context.Context) (*chainjson.GetMiningInfoResult, error) {
	return c.GetMiningInfoAsync(ctx).Receive()
}

// FutureGetNetworkHashPS is a future promise to deliver the result of a
// GetNetworkHashPSAsync RPC invocation (or an applicable error).
type FutureGetNetworkHashPS cmdRes

// Receive waits for the response promised by the future and returns the
// estimated network hashes per second for the block heights provided by the
// parameters.
func (r *FutureGetNetworkHashPS) Receive() (int64, error) {
	res, err := receiveFuture(r.ctx, r.c)
	if err != nil {
		return -1, err
	}

	// Unmarshal result as an int64.
	var result int64
	err = json.Unmarshal(res, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// GetNetworkHashPSAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetNetworkHashPS for the blocking version and more details.
func (c *Client) GetNetworkHashPSAsync(ctx context.Context) *FutureGetNetworkHashPS {
	cmd := chainjson.NewGetNetworkHashPSCmd(nil, nil)
	return (*FutureGetNetworkHashPS)(c.sendCmd(ctx, cmd))
}

// GetNetworkHashPS returns the estimated network hashes per second using the
// default number of blocks and the most recent block height.
//
// See GetNetworkHashPS2 to override the number of blocks to use and
// GetNetworkHashPS3 to override the height at which to calculate the estimate.
func (c *Client) GetNetworkHashPS(ctx context.Context) (int64, error) {
	return c.GetNetworkHashPSAsync(ctx).Receive()
}

// GetNetworkHashPS2Async returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetNetworkHashPS2 for the blocking version and more details.
func (c *Client) GetNetworkHashPS2Async(ctx context.Context, blocks int) *FutureGetNetworkHashPS {
	cmd := chainjson.NewGetNetworkHashPSCmd(&blocks, nil)
	return (*FutureGetNetworkHashPS)(c.sendCmd(ctx, cmd))
}

// GetNetworkHashPS2 returns the estimated network hashes per second for the
// specified previous number of blocks working backwards from the most recent
// block height.  The blocks parameter can also be -1 in which case the number
// of blocks since the last difficulty change will be used.
//
// See GetNetworkHashPS to use defaults and GetNetworkHashPS3 to override the
// height at which to calculate the estimate.
func (c *Client) GetNetworkHashPS2(ctx context.Context, blocks int) (int64, error) {
	return c.GetNetworkHashPS2Async(ctx, blocks).Receive()
}

// GetNetworkHashPS3Async returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetNetworkHashPS3 for the blocking version and more details.
func (c *Client) GetNetworkHashPS3Async(ctx context.Context, blocks, height int) *FutureGetNetworkHashPS {
	cmd := chainjson.NewGetNetworkHashPSCmd(&blocks, &height)
	return (*FutureGetNetworkHashPS)(c.sendCmd(ctx, cmd))
}

// GetNetworkHashPS3 returns the estimated network hashes per second for the
// specified previous number of blocks working backwards from the specified
// block height.  The blocks parameter can also be -1 in which case the number
// of blocks since the last difficulty change will be used.
//
// See GetNetworkHashPS and GetNetworkHashPS2 to use defaults.
func (c *Client) GetNetworkHashPS3(ctx context.Context, blocks, height int) (int64, error) {
	return c.GetNetworkHashPS3Async(ctx, blocks, height).Receive()
}

// FutureGetWork is a future promise to deliver the result of a
// GetWorkAsync RPC invocation (or an applicable error).
type FutureGetWork cmdRes

// Receive waits for the response promised by the future and returns the hash
// data to work on.
func (r *FutureGetWork) Receive() (*chainjson.GetWorkResult, error) {
	res, err := receiveFuture(r.ctx, r.c)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a getwork result object.
	var result chainjson.GetWorkResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWorkAsync returns an instance of a type that can be used to get the result
// of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetWork for the blocking version and more details.
func (c *Client) GetWorkAsync(ctx context.Context) *FutureGetWork {
	cmd := chainjson.NewGetWorkCmd(nil)
	return (*FutureGetWork)(c.sendCmd(ctx, cmd))
}

// GetWork returns hash data to work on.
//
// See GetWorkSubmit to submit the found solution.
func (c *Client) GetWork(ctx context.Context) (*chainjson.GetWorkResult, error) {
	return c.GetWorkAsync(ctx).Receive()
}

// FutureGetWorkSubmit is a future promise to deliver the result of a
// GetWorkSubmitAsync RPC invocation (or an applicable error).
type FutureGetWorkSubmit cmdRes

// Receive waits for the response promised by the future and returns whether
// or not the submitted block header was accepted.
func (r *FutureGetWorkSubmit) Receive() (bool, error) {
	res, err := receiveFuture(r.ctx, r.c)
	if err != nil {
		return false, err
	}

	// Unmarshal result as a boolean.
	var accepted bool
	err = json.Unmarshal(res, &accepted)
	if err != nil {
		return false, err
	}

	return accepted, nil
}

// GetWorkSubmitAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetWorkSubmit for the blocking version and more details.
func (c *Client) GetWorkSubmitAsync(ctx context.Context, data string) *FutureGetWorkSubmit {
	cmd := chainjson.NewGetWorkCmd(&data)
	return (*FutureGetWorkSubmit)(c.sendCmd(ctx, cmd))
}

// GetWorkSubmit submits a block header which is a solution to previously
// requested data and returns whether or not the solution was accepted.
//
// See GetWork to request data to work on.
func (c *Client) GetWorkSubmit(ctx context.Context, data string) (bool, error) {
	return c.GetWorkSubmitAsync(ctx, data).Receive()
}

// FutureSubmitBlockResult is a future promise to deliver the result of a
// SubmitBlockAsync RPC invocation (or an applicable error).
type FutureSubmitBlockResult cmdRes

// Receive waits for the response promised by the future and returns an error if
// any occurred when submitting the block.
func (r *FutureSubmitBlockResult) Receive() error {
	res, err := receiveFuture(r.ctx, r.c)
	if err != nil {
		return err
	}

	if string(res) != "null" {
		var result string
		err = json.Unmarshal(res, &result)
		if err != nil {
			return err
		}

		return errors.New(result)
	}

	return nil

}

// SubmitBlockAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See SubmitBlock for the blocking version and more details.
func (c *Client) SubmitBlockAsync(ctx context.Context, block *dcrutil.Block, options *chainjson.SubmitBlockOptions) *FutureSubmitBlockResult {
	blockHex := ""
	if block != nil {
		blockBytes, err := block.Bytes()
		if err != nil {
			return (*FutureSubmitBlockResult)(newFutureError(ctx, err))
		}

		blockHex = hex.EncodeToString(blockBytes)
	}

	cmd := chainjson.NewSubmitBlockCmd(blockHex, options)
	return (*FutureSubmitBlockResult)(c.sendCmd(ctx, cmd))
}

// SubmitBlock attempts to submit a new block into the Decred network.
func (c *Client) SubmitBlock(ctx context.Context, block *dcrutil.Block, options *chainjson.SubmitBlockOptions) error {
	return c.SubmitBlockAsync(ctx, block, options).Receive()
}

// FutureRegenTemplateResult is a future promise to deliver the result of a
// RegenTemplate RPC invocation (or an applicable error).
type FutureRegenTemplateResult cmdRes

// Receive waits for the response and returns an error if any has occurred.
func (r *FutureRegenTemplateResult) Receive() error {
	_, err := receiveFuture(r.ctx, r.c)
	return err
}

// RegenTemplateAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See RegenTemplate for the blocking version and more details.
func (c *Client) RegenTemplateAsync(ctx context.Context) *FutureRegenTemplateResult {
	cmd := chainjson.NewRegenTemplateCmd()
	return (*FutureRegenTemplateResult)(c.sendCmd(ctx, cmd))
}

// RegenTemplate asks the node to regenerate its current block template. Note
// that template generation is currently asynchronous, therefore no guarantees
// are made for when or whether a new template will actually be available.
func (c *Client) RegenTemplate(ctx context.Context) error {
	return c.RegenTemplateAsync(ctx).Receive()
}
