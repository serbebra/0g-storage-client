package cmd

import (
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var (
	deployArgs struct {
		url            string
		key            string
		bytecodeOrFile string
	}

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy ZeroGStorage contract to specified blockchain",
		Run:   deploy,
	}
)

func init() {
	deployCmd.Flags().StringVar(&deployArgs.url, "url", "", "Fullnode URL to interact with blockchain")
	if err := deployCmd.MarkFlagRequired("url"); err != nil {
		logrus.WithError(err).Fatal("Failed to mark 'url' flag as required")
	}
	deployCmd.Flags().StringVar(&deployArgs.key, "key", "", "Private key to create smart contract")
	if err := deployCmd.MarkFlagRequired("key"); err != nil {
		logrus.WithError(err).Fatal("Failed to mark 'key' flag as required")
	}
	deployCmd.Flags().StringVar(&deployArgs.bytecodeOrFile, "bytecode", "", "ZeroGStorage smart contract bytecode or path to bytecode file")
	if err := deployCmd.MarkFlagRequired("bytecode"); err != nil {
		logrus.WithError(err).Fatal("Failed to mark 'bytecode' flag as required")
	}

	// Ensure rootCmd is properly initialized
	if rootCmd == nil {
		logrus.Fatal("rootCmd is not initialized")
	}
	rootCmd.AddCommand(deployCmd)
}

func deploy(cmd *cobra.Command, args []string) {
	client := blockchain.MustNewWeb3(deployArgs.url, deployArgs.key)

	bytecode := deployArgs.bytecodeOrFile
	if isFile(bytecode) {
		data, err := ioutil.ReadFile(bytecode)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to read bytecode file")
		}
		bytecode = string(data)
	}

	contract, err := blockchain.Deploy(client, bytecode)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to deploy smart contract")
	}

	logrus.WithField("contract", contract).Info("Smart contract deployed")
}

func isFile(path string) bool {
	return strings.HasSuffix(path, ".bin") || strings.HasSuffix(path, ".sol")
}
