package slideshow

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color" 
	"fmt"
)

// declare a theme
type CustomTheme struct {
	BackgroundColor color.Color
}
var _ fyne.Theme = (*CustomTheme)(nil)

// implement the theme interface
func (t CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return t.BackgroundColor
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func NewCustomTheme(strColor string) fyne.Theme {
	trueColor := strToColor(strColor);
	return &CustomTheme{
		BackgroundColor: trueColor,
	}
}

func strToColor(strColor string) color.Color {
	switch strColor {
	case "black":
		return color.Black
	case "white":
		return color.White
	case "red":
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	case "green":
		return color.RGBA{R: 0, G: 255, B: 0, A: 255}
	case "blue":
		return color.RGBA{R: 0, G: 0, B: 255, A: 255}
	case "yellow":
		return color.RGBA{R: 255, G: 255, B: 0, A: 255}
	case "cyan":
		return color.RGBA{R: 0, G: 255, B: 255, A: 255}
	case "magenta":
		return color.RGBA{R: 255, G: 0, B: 255, A: 255}
	default:
		// if the color is a hex code
		if(len(strColor) == 7 && strColor[0] == '#') {
			var r, g, b uint64
			fmt.Sscanf(strColor[1:3], "%x", &r)
			fmt.Sscanf(strColor[3:5], "%x", &g)
			fmt.Sscanf(strColor[5:7], "%x", &b)
			return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
		}
		// otherwise return black
		fmt.Println("Failed to parse color, using black instead")
		return color.Black
	}
}


