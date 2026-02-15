package batch

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zerjioang/rooftop/v2/analysis"
)

func TestProcessCSV(t *testing.T) {
	t.Run("2017_eth_sorted", func(t *testing.T) {
		err := ProcessCSV("/Users/sergio/Downloads/unique_contracts_2017_sorted.csv", 0, 546)
		require.NoError(t, err)
	})
	t.Run("2018_eth_sorted", func(t *testing.T) {
		err := ProcessCSV("/Users/sergio/Downloads/unique_contracts_2018_sorted.csv", 10, 0)
		require.NoError(t, err)
	})
}

func TestProcessSample(t *testing.T) {
	t.Run("stuck", func(t *testing.T) {
		name := "0x0969a6faae9cc596a42ead5ac9f780360c2f3248"
		code := "0x606060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806361461954146044575b600080fd5b3415604e57600080fd5b6054606a565b6040518082815260200191505060405180910390f35b6000806073606a565b8101905080806001019150915050905600a165627a7a723058201cff09a7222fbd72f9f18386a0a03a1a1f02313950b8306cbdb5ce84ed7749c40029"
		cli := analysis.NewCLI()
		if err := cli.Run(context.Background(), name, code, false, time.Now()); err != nil {
			log.Fatal(err)
		}
	})
}
