/*
Copyright © 2024 Mathias Djärv <mathias.djarv@allbinary.se>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mdjarv/img2qmk/qmk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var namePrefix string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "img2qmk",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		// Single image
		if len(args) == 1 {
			err := qmk.ParseImage(args[0], namePrefix)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			os.Exit(0)
		}

		// Animation
		frames := make([]string, 0, len(args))
		for i, arg := range args {
			name := namePrefix
			if namePrefix != "" && len(args) > 1 {
				name = fmt.Sprintf("%s%d", namePrefix, i+1)
			}
			frames = append(frames, name)

			err := qmk.ParseImage(arg, name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		fmt.Printf("static const char* %s_frames[%d] = {\n", namePrefix, len(frames))
		for _, frame := range frames {
			fmt.Printf("\t%s,\n", frame)
		}
		fmt.Println("};")

		fmt.Printf("static const uint16_t %s_sizes[%d] = {\n", namePrefix, len(frames))
		for _, frame := range frames {
			fmt.Printf("\tsizeof(%s),\n", frame)
		}
		fmt.Println("};")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.img2qmk.yaml)")
	rootCmd.Flags().StringVarP(&namePrefix, "name", "n", "", "name of the generated variable")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".img2qmk" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".img2qmk")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
