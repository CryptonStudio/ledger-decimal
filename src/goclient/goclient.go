/*******************************************************************************
*   (c) 2018 ZondaX GmbH
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
********************************************************************************/

// A simple command line tool that outputs json messages representing transactions
// Usage: samples [0-3] [binary|text]
// Note: Use build_samples.sh script to update correctly update dependencies

package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/tendermint/go-crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/zondax/ledger-goclient"
)

func PrintSampleFunc(message bank.SendMsg, output string) {

	res := message.GetSignBytes()

	if output == "binary" {
		fmt.Print(res)
	} else if output == "text" {
		fmt.Print(string(res))
	}
}

func ParseArgs(numberOfSamples int) (int, string, int) {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 2 {
		fmt.Println("Not enought arguments")
		fmt.Printf("USAGE: samples SampleNumber[0-%d] OutputType[\"binary|text\"]\n", numberOfSamples-1)

		return 0, "", -1
	}
	sampleIndex, err := strconv.Atoi(argsWithoutProg[0])
	if err != nil {
		fmt.Println("Could not parse argument " + err.Error())
		return 0, "", -1
	}
	if sampleIndex < 0 && sampleIndex >= numberOfSamples {
		fmt.Printf("Number must be betweem 0-%d\n", numberOfSamples-1)
		return 0, "", -1
	}
	sampleOutput := argsWithoutProg[1]
	if sampleOutput != "binary" && sampleOutput != "text" {
		fmt.Println("Wrong OutputType, only binary|text allowed")
		return 0, "", -1
	}
	return sampleIndex, sampleOutput, 0
}

func GetMessages() ([]bank.SendMsg) {
	return []bank.SendMsg{
		// Simple address, 1 input, 1 output
		bank.SendMsg{
			Inputs: []bank.Input{
				{
					Address: crypto.Address([]byte("input")),
					Coins:   sdk.Coins{{"atom", 10}},
					//Sequence: 1,
				},
			},
			Outputs: []bank.Output{
				{
					Address: crypto.Address([]byte("output")),
					Coins:   sdk.Coins{{"atom", 10}},
				},
			},
		},

		// Real public key, 1 input, 1 output
		bank.SendMsg{
			Inputs: []bank.Input{
				{
					Address: crypto.Address(crypto.GenPrivKeySecp256k1().PubKey().Bytes()),
					Coins:   sdk.Coins{{"atom", 1000000}},
					//Sequence: 1,
				},
			},
			Outputs: []bank.Output{
				{
					Address: crypto.Address(crypto.GenPrivKeySecp256k1().PubKey().Bytes()),
					Coins:   sdk.Coins{{"atom", 1000000}},
				},
			},
		},

		// Simple address, 2 inputs, 2 outputs
		bank.SendMsg{
			Inputs: []bank.Input{
				{
					Address: crypto.Address([]byte("input")),
					Coins:   sdk.Coins{{"atom", 10}},
					//Sequence: 1,
				},
				{
					Address: crypto.Address([]byte("anotherinput")),
					Coins:   sdk.Coins{{"atom", 50}},
					//Sequence: 1,
				},
			},
			Outputs: []bank.Output{
				{
					Address: crypto.Address([]byte("output")),
					Coins:   sdk.Coins{{"atom", 10}},
				},
				{
					Address: crypto.Address([]byte("anotheroutput")),
					Coins:   sdk.Coins{{"atom", 50}},
				},
			},
		},

		// Simple address, 2 inputs, 2 outputs, 2 coins
		bank.SendMsg{
			Inputs: []bank.Input{
				{
					Address: crypto.Address([]byte("input")),
					Coins:   sdk.Coins{{"atom", 10}, {"bitcoin", 20}},
					//Sequence: 1,
				},
				{
					Address: crypto.Address([]byte("anotherinput")),
					Coins:   sdk.Coins{{"atom", 50}, {"bitcoin", 60}, {"ethereum", 70}},
					//Sequence: 1,
				},
			},
			Outputs: []bank.Output{
				{
					Address: crypto.Address([]byte("output")),
					Coins:   sdk.Coins{{"atom", 10}, {"bitcoin", 20}},
				},
				{
					Address: crypto.Address([]byte("anotheroutput")),
					Coins:   sdk.Coins{{"atom", 50}, {"bitcoin", 60}, {"ethereum", 70}},
				},
			},
		},
	}
}

func main() {
	ledger, err := ledger_goclient.FindLedger()

	if err != nil {
		fmt.Printf("Ledger NOT Found\n")
		fmt.Print(err.Error())
	} else {
		ledger.Logging = true

		fmt.Printf("\n************ Waiting for signature message 1..\n")

		messages := GetMessages()
		transactionData := messages[0].GetSignBytes()
		signedMsg, err := ledger.Sign(transactionData)

		if err == nil {
			fmt.Printf("Signed msg: %x\n", signedMsg)

		} else {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n************ Waiting for signature message 2..\n")

		transactionData = messages[1].GetSignBytes()
		signedMsg, err = ledger.Sign(transactionData)

		if err == nil {
			fmt.Printf("Signed msg: %x\n", signedMsg)
		} else {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n************ Waiting for signature message 3..\n")

		transactionData = messages[2].GetSignBytes()
		signedMsg, err = ledger.Sign(transactionData)

		if err == nil {
			fmt.Printf("Signed msg: %x\n", signedMsg)
		} else {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
}