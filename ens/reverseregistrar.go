// Copyright 2017 Orinoco Payments
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ens

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	etherutils "github.com/orinocopay/go-etherutils"
	"github.com/orinocopay/go-etherutils/ens/reverseregistrarcontract"
)

// ReverseRegistrar obtains the reverse registrar contract for a chain
func ReverseRegistrar(client *ethclient.Client) (registrar *reverseregistrarcontract.ReverseRegistrarContract, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	// Obtain a registry contract
	registry, err := RegistryContract(client)
	if err != nil {
		return
	}

	// Obtain the registry address from the registrar
	nameHash, err := NameHash("addr.reverse")
	registrarAddress, err := registry.Owner(nil, nameHash)
	if err != nil {
		return
	}
	if registrarAddress == UnknownAddress {
		err = errors.New("no registrar for that network")
	}

	registrar, err = reverseregistrarcontract.NewReverseRegistrarContract(registrarAddress, client)
	return
}

// CreateReverseRegistrarSession creates a session suitable for multiple calls
func CreateReverseRegistrarSession(chainID *big.Int, wallet *accounts.Wallet, account *accounts.Account, passphrase string, contract *reverseregistrarcontract.ReverseRegistrarContract, gasLimit *big.Int, gasPrice *big.Int) *reverseregistrarcontract.ReverseRegistrarContractSession {
	// Create a signer
	signer := etherutils.AccountSigner(chainID, wallet, account, passphrase)

	// Return our session
	session := &reverseregistrarcontract.ReverseRegistrarContractSession{
		Contract: contract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     account.Address,
			Signer:   signer,
			GasPrice: gasPrice,
			GasLimit: gasLimit,
		},
	}

	return session
}

// SetName sets the name for the sending address
func SetName(session *reverseregistrarcontract.ReverseRegistrarContractSession, name string) (tx *types.Transaction, err error) {
	tx, err = session.SetName(name)
	return
}
