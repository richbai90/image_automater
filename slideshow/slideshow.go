package slideshow // Define the package name

import (
	"fmt"           // Package for formatted I/O
	"io/fs"         // Package for file system interface
	"path/filepath" // Package for file path manipulation
	"time"          // Package for time manipulation

	"fyne.io/fyne/v2"                             // GUI toolkit
	"fyne.io/fyne/v2/app"                         // Create and manage GUI applications
	"fyne.io/fyne/v2/canvas"                      // Package for drawing on the screen
	"fyne.io/fyne/v2/container"                   // Package for organizing GUI elements
	"github.com/richbai90/image_automater/serial" // Package for serial communication
	log "github.com/sirupsen/logrus"              // Package for logging
)

// Define a custom type named mode_t, which is an integer
type mode_t int

// Define two constants of type mode_t
const (
	INTERVAL mode_t = iota // The slideshow is controlled by an interval
	TRIGGER  mode_t = 1    // The slideshow is controlled by a trigger signal from the serial port
)

// Define a struct named Slideshow
type Slideshow struct {
	App         fyne.App       // An instance of the fyne.App struct
	imageFolder string         // The path of the folder containing images
	images      []string       // A slice containing the names of image files
	index       int            // The index of the currently displayed image
	window      fyne.Window    // The window displaying the slideshow
	image       *canvas.Image  // The image displayed in the window
	interval    time.Duration  // The time duration between image changes
	paused      bool           // Whether the slideshow is currently paused
	mode        mode_t         // The mode of operation (interval or trigger)
	s           *serial.Serial // An instance of the serial.Serial struct
}

// Define a function named NewSlideshow that creates a new Slideshow instance
func NewSlideshow(imageFolder string, recursive bool, interval time.Duration, triggerMode bool, s *serial.Serial) (*Slideshow, error) {
	// Call the getImages function to retrieve the names of image files
	images, err := getImages(imageFolder, recursive)
	if err != nil {
		return nil, err
	}

	// If no images were found, return an error
	if len(images) == 0 {
		return nil, fmt.Errorf("no images found in %s", imageFolder)
	}

	// Create a new Slideshow instance with the specified parameters
	return &Slideshow{
		imageFolder: imageFolder,
		images:      images,
		// If triggerMode is true, set the interval to 0; otherwise, use the specified interval
		interval: func() time.Duration {
			if triggerMode {
				return 0
			} else {
				return interval
			}
		}(),
		App: app.New(), // Create a new fyne.App instance
		// If triggerMode is true, set the mode to TRIGGER; otherwise, use INTERVAL
		mode: func() mode_t {
			if triggerMode {
				return TRIGGER
			} else {
				return INTERVAL
			}
		}(),
		s: s, // Set the serial instance
	}, nil
}

func (s *Slideshow) Start() {
	log.Info("Starting slideshow")
	// Create a new fyne window with the title "Slideshow"
	s.window = s.App.NewWindow("Slideshow")

	// Set the window to full screen
	s.window.SetFullScreen(true)
	// Create an image from the first image in the `s.images` array
	s.image = canvas.NewImageFromFile(s.images[s.index])
	// Set the content of the window to the image
	s.window.SetContent(container.NewMax(s.image))
	// Display the window
	log.Debug("Displaying window")
	s.window.Show()
	// Start a goroutine to handle the slideshow's logic
	go func() {
		// Create a buffer with length 5 to read from the serial port
		buf := make([]byte, 5)
		for {
			if !s.paused {
				if s.mode == TRIGGER {
					log.Info("mode is trigger. waiting for next command")
					// Keep reading from the serial port until the "next" command is received if the slideshow is in trigger mode
					for {
						n, err := s.s.Read(buf)
						if err != nil {
							log.Fatal("Error reading from serial port: ", err)
						}
						if n > 0 && string(buf[:n]) == "next" {
							log.Info("next command received")
							break
						}
					}
				}

				// Wait for the specified interval
				// If the slideshow is in trigger mode, this will be skipped because the interval is 0
				log.Debug("waiting for interval")
				<-time.Tick(s.interval)
				log.Debug("interval elapsed")
				// Set the index to the next image in the `s.images` array
				s.index = (s.index + 1) % len(s.images)

				// Create a new image from the next image in the `s.images` array
				img := canvas.NewImageFromFile(s.images[s.index])

				// Set the fill mode of the image to ImageFillOriginal
				img.FillMode = canvas.ImageFillOriginal

				// Set the content of the window to the new image
				s.image = img
				s.window.SetContent(container.NewMax(s.image))

				// If the slideshow is in trigger mode, send the "capture" command to the serial port
				if s.mode == TRIGGER {
					log.Info("sending capture command")
					s.s.Write([]byte("capture"))
				}

			}
			// Restart from the beginning of the loop
		}
	}()
}

// TogglePause toggles the paused state of the slideshow
func (s *Slideshow) TogglePause() {
	log.Info("Slideshow paused")
	s.paused = !s.paused
}

// Stop stops the slideshow
func (s *Slideshow) Stop() {
	log.Info("Stopping slideshow gracefully")
	s.paused = true
	s.window.Close()
}

// AddKeyHandlers adds key event handlers to the window
func (s *Slideshow) AddKeyHandlers(callback func(event *fyne.KeyEvent)) {
	// Add key event handlers
	s.window.Canvas().SetOnTypedKey(callback)
}

// AddShortcut adds a shortcut to the window
func (s *Slideshow) AddShortcut(shortcut fyne.Shortcut, callback func(shortcut fyne.Shortcut)) {
	// Add a shortcut
	s.window.Canvas().AddShortcut(shortcut, callback)
}

// getImages is a helper function that returns a slice of strings, containing the names of the image files in the specified folder
// It takes a single argument, the path to the image folder and returns two values: the list of image names, and an error if there was any
func getImages(imageFolder string, recursive bool) ([]string, error) {
	// Read the directory at imageFolder and return any error encountered
	var images []string
	filepath.WalkDir(imageFolder, func(path string, d fs.DirEntry, err error) error {

		// Check if the file is not a directory and is an image
		if !d.IsDir() && isImage(d.Name()) {
			// Append the full file path to the images slice
			images = append(images, path)
		} else if d.IsDir() && recursive {
			// recursively call getImages on the subdirectory
			subImages, err := getImages(path, recursive)
			if err != nil {
				return err
			}
			images = append(images, subImages...)
		}
		return nil
	})
	// Return the list of image names and no error
	return images, nil
}

// isImage is a helper function to check if the given filename represents an image.
// It takes a single argument, the filename, and returns a boolean value indicating whether the file is an image or not.
func isImage(filename string) bool {
	// Get the extension of the filename
	ext := filepath.Ext(filename)
	// Check if the extension is one of the supported image types
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp"
}
