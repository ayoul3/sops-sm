package cmd

import (
	"os"

	"github.com/ayoul3/sops-sm/stores"
	"github.com/ayoul3/sops-sm/stores/json"
	"github.com/ayoul3/sops-sm/stores/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var formats = map[string]stores.StoreAPI{
	"yaml": yaml.NewStore(),
	"json": json.NewStore(),
}

var (
	rootCmd = &cobra.Command{
		Use:   "sops-sm",
		Short: "Sops-SM is a fork of Mozilla's Security Operations tools that supports AWS Parameter store and SecrestsManager",
		Long:  `Sops-SM decrypts a yaml or json file that contain references to secrets stored in AWS SecretsManager or Parameter Store`,
	}
	decrypt = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt input file",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatalf("Input file is required")
			}
			validateFile(args[0])
			HandleDecrypt(args[0])
		},
	}
	encrypt = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt input file - requires .cache file generated from the decryption phase",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatalf("Input file is required")
			}
			validateFile(args[0])
			HandleEncrypt(args[0])
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(encrypt)
	rootCmd.AddCommand(decrypt)
}

func validateFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("input file %s does not exist", filePath)
	}
}
