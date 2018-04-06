package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
)

var srcEnv, destEnv string

var copyProductsCmd = &cobra.Command{
	Use:   "copy-products",
	Short: "Copies all products from a src into a dest environment",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if srcEnv == "" {
			return errors.New("src cannot be empty")
		}
		if destEnv == "" {
			return errors.New("dest cannot be empty")
		}

		cmd.SilenceUsage = true
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Copying all products from %s to %s...\n", srcEnv, destEnv)
		if err := svc.CopyAllProducts(srcEnv, destEnv); err != nil {
			return err
		}
		log.Println("All done")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyProductsCmd)

	copyProductsCmd.Flags().StringVarP(&srcEnv, "src", "s", "", "The name of the Stripe environment to copy products from")
	copyProductsCmd.Flags().StringVarP(&destEnv, "dest", "d", "", "The name of the Stripe environment to copy products to")
}
