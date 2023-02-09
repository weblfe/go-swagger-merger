// Package cmd /*
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weblfe/go-swagger-merger/swagger"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var outputFile string
var cfgFile string
var suffix = new([]string)
var beautify bool
var version = "0.1.2"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "go-swagger-merger",
	Short:   "swagger merger",
	Version: version,
	Long:    `swagger merger Author: weblfe`,
	Example: `<files> go-swagger-merger -o swagger.json b.swagger.json c.swagger.json
<dir> go-swagger-merger -o swagger.yaml api/bff/ 
<suffix> go-swagger-merger -o swagger.json -s .swagger.json -s .swagger.yml  api/bff
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New(`miss input files`)
		}
		if outputFile == "" {
			outputFile = "swagger.json"
		}
		var (
			merger = swagger.NewMerger()
			filter = merger.CreatePatternFilter(*suffix)
		)
		for _, file := range args {
			if err := merger.AddFile(file, filter); err != nil {
				return err
			}
		}
		return merger.Save(outputFile, beautify)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-swagger-merger.
	//yaml)")
	defaultSuffix := []string{".swagger.json", ".swagger.yaml"}
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "merge swagger.json to save file")
	rootCmd.Flags().StringArrayVarP(suffix, "suffix", "s", defaultSuffix, "suffix filter for swagger matcher")
	rootCmd.Flags().BoolVarP(&beautify, "beautify", "b", false, "merge swagger.json unzip beautify format")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-swagger-merger" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-swagger-merger")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
