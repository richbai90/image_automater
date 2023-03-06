/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"time"
	"github.com/spf13/cobra"
	"github.com/richbai90/image_automater/slideshow"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	
)

var interval float32;
var background string;


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "image_automater <folder>",
	Short: "Start a slideshow of images in a folder",
	Long: `Start a slideshow of images in a folder`,
	Args: cobra.MinimumNArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		folder := args[0]
		ss, err := slideshow.NewSlideshow(folder, time.Duration(interval)*time.Second)
		if err != nil {
			return err
		}
		ss.Start()

	// Add key event handlers
	ss.AddKeyHandlers(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeySpace {
			ss.TogglePause()
		}
	})

	// // Close on Ctrl+K
	ctrlk := &desktop.CustomShortcut{KeyName: fyne.KeyK, Modifier: desktop.ControlModifier}
	ss.AddShortcut(ctrlk, func(shortcut fyne.Shortcut) {
		ss.App.Quit()
	})

	ss.App.Settings().SetTheme(slideshow.NewCustomTheme(background))
	ss.App.Run()

		return nil
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.image_automater.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Float32VarP(&interval, "interval", "i", 1, "Interval between images in seconds")
	rootCmd.Flags().StringVarP(&background, "background", "b", "black", "Background color")
}


