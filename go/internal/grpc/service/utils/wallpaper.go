package utils

import (
	"bytes"
	"fmt"
	"image"
	"strings"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
)

var defaultStack = "system-ui, SegoeUI, Roboto, Ubuntu, Cantarell, NotoSans, HelveticaNeue, Arial, sans-serif"

// RenderCompleteBackground renders onto basePNG and returns the PNG bytes.
func RenderCompleteBackground(basePNG []byte, layout models.WallpaperLayout, teamName, ip string) ([]byte, error) {
	if len(basePNG) == 0 {
		return nil, fmt.Errorf("basePNG is empty")
	}

	// decode base
	base, _, err := image.Decode(bytes.NewReader(basePNG))
	if err != nil {
		return nil, fmt.Errorf("decode base: %w", err)
	}
	w, h := int(layout.W), int(layout.H)
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("invalid layout size %dx%d", w, h)
	}

	var fontStack string
	if layout.FontStack == "" {
		fontStack = defaultStack
	} else {
		fontStack = cssStackToFindfontName(layout.FontStack)
	}
	// builder
	builder := newWallpaperBuilder(w, h, 2, fontStack)
	builder.SetBackground(base)

	// team label
	if layout.Teamname.Size > 0 {
		builder.DrawLabel(Label{
			Text:   teamName,
			X:      int(layout.Teamname.X),
			Y:      int(layout.Teamname.Y),
			Size:   int(layout.Teamname.Size),
			Weight: int(layout.Teamname.Weight),
			Color:  layout.Teamname.Color,
			Align:  anchorX(layout.Teamname.Align),
		})
	}
	// ip label
	if layout.IP.Size > 0 {
		builder.DrawLabel(Label{
			Text:   ip,
			X:      int(layout.IP.X),
			Y:      int(layout.IP.Y),
			Size:   int(layout.IP.Size),
			Weight: int(layout.IP.Weight),
			Color:  layout.IP.Color,
			Align:  anchorX(layout.IP.Align),
		})
	}

	return builder.Result()
}

func RenderLayoutOnly(layout models.WallpaperLayout, teamName, ip string) ([]byte, error) {
	w, h := int(layout.W), int(layout.H)
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("invalid layout size %dx%d", w, h)
	}

	// font stack
	var fontStack string
	if layout.FontStack == "" {
		fontStack = defaultStack
	} else {
		fontStack = cssStackToFindfontName(layout.FontStack)
	}

	// builder
	builder := newWallpaperBuilder(w, h, 2, fontStack)
	builder.SetBlackBackground()

	// team label
	if layout.Teamname.Size > 0 {
		builder.DrawLabel(Label{
			Text:   teamName,
			X:      int(layout.Teamname.X),
			Y:      int(layout.Teamname.Y),
			Size:   int(layout.Teamname.Size),
			Weight: int(layout.Teamname.Weight),
			Color:  "#ffffff",
			Align:  anchorX(layout.Teamname.Align),
		})
	}
	// ip label
	if layout.IP.Size > 0 {
		builder.DrawLabel(Label{
			Text:   ip,
			X:      int(layout.IP.X),
			Y:      int(layout.IP.Y),
			Size:   int(layout.IP.Size),
			Weight: int(layout.IP.Weight),
			Color:  "#ffffff",
			Align:  anchorX(layout.IP.Align),
		})
	}

	return builder.Result()
}

// Black background + centered semi-transparent watermark only.
func RenderNoWallpaperWatermark() ([]byte, error) {
	const (
		w, h  = 1920, 1080
		scale = 2
	)
	fontStack := defaultStack // use built-in default

	// build
	builder := newWallpaperBuilder(w, h, scale, fontStack)
	builder.SetBlackBackground()

	// watermark text
	if face, _ := openFace(fontStack, 800, float64(min(w, h))/12.0*float64(scale)); face != nil {
		builder.canvas.SetFontFace(face)
	}
	builder.canvas.SetRGBA(1, 1, 1, 0.20)
	builder.canvas.DrawStringAnchored("NO WALLPAPER FOUND",
		float64(w*scale)/2, float64(h*scale)/2, 0.5, 0.5)

	return builder.Result()
}

// RenderNoTeamAssigned renders a 1920x1080 black background with a big IP
// label centered, and just below it smaller text "no team assigned".
func RenderNoTeamAssigned(ip string) ([]byte, error) {
	const (
		w, h  = 1920, 1080
		scale = 2
	)
	fontStack := defaultStack

	b := newWallpaperBuilder(w, h, scale, fontStack)
	b.SetBlackBackground()

	// IP label (large, bold)
	if face, _ := openFace(fontStack, 700, float64(h)/8.0*float64(scale)); face != nil {
		b.canvas.SetFontFace(face)
	}
	b.canvas.SetRGB(1, 1, 1)
	b.canvas.DrawStringAnchored(ip,
		float64(w*scale)/2, float64(h*scale)/4, 0.5, 0.5)

	// "no team assigned" label (smaller, lighter)
	if face, _ := openFace(fontStack, 400, float64(h)/16.0*float64(scale)); face != nil {
		b.canvas.SetFontFace(face)
	}
	b.canvas.SetRGBA(1, 1, 1, 0.7)
	b.canvas.DrawStringAnchored("no team assigned",
		float64(w*scale)/2, float64(h*scale)/2,
		0.5, 0.0)

	return b.Result()
}

// RenderInactiveContest renders a 1920x1080 black background with:
// - IP in the top-left corner
// - Team name centered
// - "(not in an active contest)" just below the team name
func RenderInactiveContest(ip, teamName string) ([]byte, error) {
	const (
		w, h  = 1920, 1080
		scale = 2
	)
	fontStack := defaultStack
	b := newWallpaperBuilder(w, h, scale, fontStack)
	b.SetBlackBackground()

	// IP top-left
	if face, _ := openFace(fontStack, 500, float64(h)/20.0*float64(scale)); face != nil {
		b.canvas.SetFontFace(face)
	}
	b.canvas.SetRGB(1, 1, 1)
	b.canvas.DrawStringAnchored(ip,
		20, 40, // a little padding
		0.0, 0.0)

	// Team name center
	if face, _ := openFace(fontStack, 700, float64(h)/8.0*float64(scale)); face != nil {
		b.canvas.SetFontFace(face)
	}
	b.canvas.SetRGB(1, 1, 1)
	b.canvas.DrawStringAnchored(teamName,
		float64(w*scale)/2, float64(h*scale)/2,
		0.5, 0.5)

	// "(not in an active contest)" below team name
	if face, _ := openFace(fontStack, 400, float64(h)/18.0*float64(scale)); face != nil {
		b.canvas.SetFontFace(face)
	}
	b.canvas.SetRGBA(1, 1, 1, 0.7)
	b.canvas.DrawStringAnchored("(not in an active contest)",
		float64(w*scale)/2, float64(h*scale)/2+float64(h)/12.0*float64(scale),
		0.5, 0.0)

	return b.Result()
}

// --- helpers ---

func anchorX(a models.Align) float64 {
	switch a {
	case models.AlignRight:
		return 1.0
	case models.AlignCenter:
		return 0.5
	default:
		return 0.0
	}
}

func cssStackToFindfontName(stack string) string {
	parts := strings.Split(stack, ",")
	if len(parts) == 0 {
		return ""
	}

	// trim quotes/whitespace
	fam := strings.TrimSpace(parts[0])
	fam = strings.Trim(fam, `"'`)

	return strings.NewReplacer(" ", "", "-", "").Replace(fam)
}
