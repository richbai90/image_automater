<h1 align="center">Welcome to image_automater 👋</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-2.0.0-blue.svg?cacheSeconds=2592000" />
  <a href="https://opensource.org/license/mit/" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> Image Automater is a command-line tool that allows you to start a slideshow of images in a folder. You can customize the interval between images, background color, and other options.

## Install

### From source

1. Install Go on your system. You can download it from the [Official Go Website](https://golang.org/dl/).
2. Clone the repository
3. Navigate to the repository folder
4. Run the following command to build the binary
```sh
go build -o ./bin/image_automater ./main.go
```

### Pre-built binaries
Pre-built binaries are available for Windows and Linux. You can download them from the [Releases Page](https://github.com/richbai90/image_automater/releases/latest)

## Usage

```sh
image_automater /path/to/folder [flags]
```

The command accepts the following flags:

* -i, --interval: Interval between images in milliseconds (default 1000)
* -b, --background: Background color (default "black")
* -t, --trigger-mode: Enable trigger mode
* -d, --Device: Device to use for trigger mode
* --Baud: Baud rate for trigger mode (default 9600)
* -r, --recursive: Enable recursive search for images

### Modes

The program has two modes: interval and trigger mode. In interval mode, the program will display images at a specified interval. In trigger mode, the program will wait for a trigger signal from a serial device before displaying the next image. This is useful for displaying images in a slideshow when you want to control when the next image is displayed. The exxpected trigger signal is "next". Once the next image is displayed, the program will send a "capture" message over serial and will wait for another trigger signal.

This mode is specifically designed for synchronizing the display of images with a camera. The camera will send a trigger signal over serial when it captures an image. The program will then display the next image in the slideshow. This allows you to display images in a slideshow that are captured by a camera. Without needing to synchronize the camera and computer clocks.

## Contributing
Contributions are welcome! If you find a bug or have a feature request, please open an issue on the [GitHub repository](https://github.com/richbai90/image_automater/issues). If you want to contribute code, please fork the repository and submit a pull request.

## Author

👤 **Rich Baird**

* Github: [@richbai90](https://github.com/richbai90)
* LinkedIn: [@richbai90](https://linkedin.com/in/richbai90)

## Show your support

Give a ⭐️ if this project helped you!

## 📝 License

Copyright © 2023 [Rich Baird](https://github.com/richbai90).<br />
This project is [MIT](https://opensource.org/license/mit/) licensed.

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_