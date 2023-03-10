/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/ 

// Declare the package name and import necessary packages
package cmd

import (
	"os"
	"time"

	// Fyne is a Go-based GUI toolkit and driver is the package that manages driver-based interactions with Fyne
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	// Import custom packages for serial communication and slideshow creation
	"github.com/richbai90/image_automater/serial"
	"github.com/richbai90/image_automater/slideshow"

	// Import a popular logging package for Go
	log "github.com/sirupsen/logrus"

	// Import a popular command-line library for Go
	"github.com/spf13/cobra"
)

// Declare variables for the root command
var interval float32
var background string
var triggerMode bool
var device string
var baud int
var Serial *serial.Serial

// Declare the root command
var rootCmd = &cobra.Command{
	Use:   "image_automater <folder>",
	Short: "Start a slideshow of images in a folder",
	Long:  `Start a slideshow of images in a folder`,
	Args:  cobra.MinimumNArgs(1),

	// PreRunE is a function that is run before the command is executed
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if trigger mode is enabled
		if triggerMode {
			// If trigger mode is enabled, ignore the interval flag
			if cmd.Flags().Changed("interval") {
				log.Warn("Interval flag is ignored in trigger mode")
			}
			// Initialize the Serial variable for serial communication
			var err error
			Serial, err = serial.NewSerial(device, baud)
			if err != nil {
				return err
			}
		}
		return nil
	},

	// RunE is the main function that is executed when the command is called
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the folder argument
		folder := args[0]
		// Initialize the slideshow
		ss, err := slideshow.NewSlideshow(folder, cmd.Flags().Changed("recursive"), time.Duration(interval)*time.Millisecond, triggerMode, Serial)
		if err != nil {
			return err
		}
		// Start the slideshow
		ss.Start()

		// Add key event handlers
		ss.AddKeyHandlers(func(event *fyne.KeyEvent) {
			if event.Name == fyne.KeySpace {
				ss.TogglePause()
			}
		})

		// Add a shortcut to close the application when the user presses Ctrl+K
		ctrlk := &desktop.CustomShortcut{KeyName: fyne.KeyK, Modifier: desktop.ControlModifier}
		ss.AddShortcut(ctrlk, func(shortcut fyne.Shortcut) {
			ss.App.Quit()
		})

		// Set the background color of the application
		ss.App.Settings().SetTheme(slideshow.NewCustomTheme(background))
		// Run the application
		ss.App.Run()

		return nil
	},
}

// Execute is called by the main function to execute the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Float32VarP(&interval, "interval", "i", 1e3, "Interval between images in milliseconds")
	rootCmd.Flags().StringVarP(&background, "background", "b", "black", "Background color")
	rootCmd.Flags().BoolVarP(&triggerMode, "trigger-mode", "t", false, "Enable trigger mode")
	rootCmd.Flags().StringVarP(&device, "Device", "d", "", "Device to use for trigger mode")
	rootCmd.Flags().IntVarP(&baud, "Baud", "b", 9600, "Baud rate for trigger mode")
	rootCmd.Flags().BoolP("recursive", "r", false, "Enable recursive search for images")
}
