package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var verbose bool

var (
	rootCmd = &cobra.Command{
		Use:   "sops-sm",
		Short: "Sops-SM is a fork of Mozilla's Security Operations tools that supports AWS Parameter store and SecrestsManager",
		Long:  `Sops-SM decrypts a yaml or json file that contain references to secrets stored in AWS SecretsManager or Parameter Store`,
	}
	decrypt = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt input file",
		Args: func(cmd *cobra.Command, args []string) error {
			SetLogLevel(cmd.Flags())
			return validateFile(args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			threads, _ := cmd.Flags().GetInt("threads")
			overwrite, _ := cmd.Flags().GetBool("overwrite")
			h := NewHandler(threads, overwrite)
			h.HandleDecrypt(args[0])
		},
	}
	encrypt = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt input file - requires .cache file generated from the decryption phase",
		Args: func(cmd *cobra.Command, args []string) error {
			SetLogLevel(cmd.Flags())
			return validateFile(args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			overwrite, _ := cmd.Flags().GetBool("overwrite")
			h := NewHandler(1, overwrite)
			h.HandleEncrypt(args[0])
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	decrypt.PersistentFlags().Int("threads", 1, "Parallelize the decryption process. Consider for files with more than 30 secrets. Careful of AWS throttling.")

	rootCmd.AddCommand(encrypt)
	rootCmd.AddCommand(decrypt)
	rootCmd.PersistentFlags().Bool("verbose", false, "Show info messages")
	rootCmd.PersistentFlags().Bool("overwrite", false, "Overwrite input file")
}

func validateFile(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Input file is required")
	}
	filePath := args[0]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("input file %s does not exist", filePath)
	}
	return nil
}

func SetLogLevel(flags *pflag.FlagSet) {
	verbose, _ := flags.GetBool("verbose")

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.WarnLevel)
	if verbose {
		log.SetLevel(log.InfoLevel)
	}
}
