/*
Copyright © 2024 Mathias Djärv <mathias.djarv@allbinary.se>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mdjarv/img2qmk/qmk"
	"github.com/spf13/cobra"
)

var (
	withType   bool
	namePrefix string
	frameRate  int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "img2qmk",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		if withType {
			qmk.PrintType()
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
		animation := qmk.Animation{
			Name:      namePrefix,
			FrameRate: frameRate,
		}
		// frames := make([]string, 0, len(args))
		for _, arg := range args {
			data, err := qmk.ImgToBytes(arg)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			animation.Frames = append(animation.Frames, data)
		}

		err := animation.Print()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// fmt.Printf("static const char* %s_frames[%d] = {\n", namePrefix, len(frames))
		// for _, frame := range frames {
		// 	fmt.Printf("\t%s,\n", frame)
		// }
		// fmt.Println("};")
		//
		// fmt.Printf("static const uint16_t %s_sizes[%d] = {\n", namePrefix, len(frames))
		// for _, frame := range frames {
		// 	fmt.Printf("\tsizeof(%s),\n", frame)
		// }
		// fmt.Println("};")
		//
		// fmt.Printf("static const uint16_t %s_delays[%d] = {\n", namePrefix, len(frames))
		// for _ = range frames {
		// 	fmt.Printf("\t%d,\n", frameRate)
		// }
		// fmt.Println("};")
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

	rootCmd.Flags().BoolVarP(&withType, "withType", "t", false, "generate with animation type")
	rootCmd.Flags().StringVarP(&namePrefix, "name", "n", "gfx", "name of the generated variable")
	rootCmd.Flags().IntVarP(&frameRate, "frameRate", "f", 400, "default frame rate in ms")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
