package cmd

import (
	"fmt"
	"os"

	"github.com/marotpam/gostripeplanner/pkg"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gostripeplanner",
	Short: "A tool to keep products and plans in Sync between multiple Stripe environments",
	Long: `Avoid manually creating products and plans between live and test modes in the same
account, and between different Stripe accounts.`,
}

var svc *pkg.Service

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initService)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gostripeplanner.json)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".gostripeplanner")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatalf("Error reading from config file %s, %s", cfgFile, err)
	}
}

func initService() {
	var stripeEnvs []*pkg.Environment

	envConfig := viper.GetStringMapString("envs")
	fmt.Printf("Found %d envs\n", len(envConfig))
	for name, key := range envConfig {
		env := pkg.NewStripeEnvironment(name, key)

		stripeEnvs = append(stripeEnvs, env)
	}

	svc = pkg.NewService(stripeEnvs)
}
