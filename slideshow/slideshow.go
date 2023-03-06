package slideshow

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)


type Slideshow struct {
	App 	   fyne.App
	imageFolder string
	images      []string
	index       int
	window      fyne.Window
	image       *canvas.Image
	interval    time.Duration
	paused      bool
}

func NewSlideshow(imageFolder string, interval time.Duration) (*Slideshow, error) {
	images, err := getImages(imageFolder)
	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no images found in %s", imageFolder)
	}

	return &Slideshow{
		imageFolder: imageFolder,
		images:      images,
		interval:    interval,
		App: 	   app.New(),
	}, nil
}

func (s *Slideshow) Start() {
	s.window = s.App.NewWindow("Slideshow")
	s.window.SetFullScreen(true)
	s.image = canvas.NewImageFromFile(s.images[s.index])
	s.window.SetContent(container.NewMax(s.image))
	s.window.Show()

	go func() {
		for {
			if !s.paused {
				time.Sleep(s.interval)
				s.index = (s.index + 1) % len(s.images)
				img := canvas.NewImageFromFile(s.images[s.index])
				img.FillMode = canvas.ImageFillOriginal
				s.image = img
				s.window.SetContent(container.NewMax(s.image))
			}
		}
	}()
}

func (s *Slideshow) TogglePause() {
	s.paused = !s.paused
}

func (s *Slideshow) Stop() {
	s.paused = true
	s.window.Close()
}

func(s *Slideshow) AddKeyHandlers(callback func (event *fyne.KeyEvent)) {
	// Add key event handlers
	s.window.Canvas().SetOnTypedKey(callback)
}

func (s *Slideshow) AddShortcut(shortcut fyne.Shortcut, callback func (shortcut fyne.Shortcut)) {
	s.window.Canvas().AddShortcut(shortcut, callback)
}

func getImages(imageFolder string) ([]string, error) {
	files, err := ioutil.ReadDir(imageFolder)
	if err != nil {
		return nil, err
	}

	var images []string
	for _, file := range files {
		if !file.IsDir() && isImage(file.Name()) {
			images = append(images, filepath.Join(imageFolder, file.Name()))
		}
	}

	return images, nil
}

func isImage(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp"
}

// func main() {
// 	if len(os.Args) < 3 {
// 		fmt.Println("Usage: image_automater <folder> <interval>")
// 		return
// 	}

// 	imageFolder := os.Args[1]
// 	interval, err := strconv.Atoi(os.Args[2])
// 	if err != nil {
// 		fmt.Println("Invalid interval argument")
// 		return
// 	}

// 	slideshow, err := NewSlideshow(imageFolder, time.Duration(interval)*time.Second)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	slideshow.Start()

// 	//Add key event handlers
// 	slideshow.window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
// 		if event.Name == fyne.KeySpace {
// 			slideshow.TogglePause()
// 		}
// 	})

// 	// Close on Ctrl+K
// 	ctrlk := &desktop.CustomShortcut{KeyName: fyne.KeyK, Modifier: desktop.ControlModifier}
// 	slideshow.window.Canvas().AddShortcut(ctrlk, func(shortcut fyne.Shortcut) {
// 		slideshow.app.Quit()
// 	})

// 	slideshow.app.Settings().SetTheme(&CustomTheme{})
// 	slideshow.app.Run()
// }
