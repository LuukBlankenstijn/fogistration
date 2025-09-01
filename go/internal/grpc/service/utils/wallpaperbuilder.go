package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/flopp/go-findfont"
	"github.com/fogleman/gg"
	drawx "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// ---- internal types ----

// Label describes one text overlay.
type Label struct {
	Text   string
	X, Y   int
	Size   int
	Weight int
	Color  string
	Align  float64 // 0=left, 0.5=center, 1=right
}

// wallpaperBuilder handles drawing with a fixed font + supersampling.
type wallpaperBuilder struct {
	w, h   int
	scale  int
	font   string
	canvas *gg.Context
}

func newWallpaperBuilder(w, h, scale int, fontStack string) *wallpaperBuilder {
	return &wallpaperBuilder{
		w:      w,
		h:      h,
		scale:  scale,
		font:   fontStack,
		canvas: gg.NewContext(w*scale, h*scale),
	}
}

func (b *wallpaperBuilder) SetBackground(base image.Image) {
	if base != nil {
		bg := resizeCover(base, b.w*b.scale, b.h*b.scale)
		b.canvas.DrawImage(bg, 0, 0)
		return
	}

	// fallback: black + watermark
	b.canvas.SetRGB(0, 0, 0)
	b.canvas.Clear()
	if face, _ := openFace(b.font, 800, float64(min(b.w, b.h))/12.0*float64(b.scale)); face != nil {
		b.canvas.SetFontFace(face)
		b.canvas.SetRGBA(1, 1, 1, 0.20)
		b.canvas.DrawStringAnchored("NO BACKGROUND FILE FOUND",
			float64(b.w*b.scale)/2, float64(b.h*b.scale)/2, 0.5, 0.5)
	}
}

func (b *wallpaperBuilder) DrawLabel(l Label) {
	if l.Size <= 0 {
		return
	}
	if face, _ := openFace(b.font, int32(l.Weight), float64(l.Size)*float64(b.scale)); face != nil {
		b.canvas.SetFontFace(face)
		setHexOrWhite(b.canvas, l.Color)
		b.canvas.DrawStringAnchored(l.Text,
			float64(l.X*b.scale),
			float64(l.Y*b.scale),
			l.Align, 0,
		)
	}
}

func (b *wallpaperBuilder) Result() ([]byte, error) {
	out := image.NewRGBA(image.Rect(0, 0, b.w, b.h))
	drawx.CatmullRom.Scale(out, out.Bounds(), b.canvas.Image(), b.canvas.Image().Bounds(), drawx.Over, nil)

	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}
	return buf.Bytes(), nil
}

func (b *wallpaperBuilder) SetBlackBackground() {
	b.canvas.SetRGB(0, 0, 0)
	b.canvas.Clear()
}

// scale-to-fill then center-crop to tw×th
func resizeCover(src image.Image, tw, th int) image.Image {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()

	sx := float64(tw) / float64(sw)
	sy := float64(th) / float64(sh)
	scale := math.Max(sx, sy)

	nw := int(math.Ceil(float64(sw) * scale))
	nh := int(math.Ceil(float64(sh) * scale))

	// 1) resize
	tmp := image.NewRGBA(image.Rect(0, 0, nw, nh))
	drawx.CatmullRom.Scale(tmp, tmp.Bounds(), src, sb, draw.Src, nil)

	// 2) center-crop
	ox := (nw - tw) / 2
	oy := (nh - th) / 2
	crop := image.Rect(ox, oy, ox+tw, oy+th).Intersect(tmp.Bounds())

	out := image.NewRGBA(image.Rect(0, 0, tw, th))
	draw.Draw(out, out.Bounds(), tmp.SubImage(crop), crop.Min, draw.Src)
	return out
}

func openFace(cssStack string, weight int32, sizePx float64) (font.Face, error) {
	fam := pickFamily(cssStack)
	if fam == "" {
		fam = "sans-serif"
	}
	wantBold := weight >= 600

	var names []string
	if wantBold {
		names = append(names, fam+" Bold", fam+"-Bold", fam+"Bold")
	}
	names = append(names, fam, "DejaVuSans", "NotoSans", "LiberationSans")

	var path string
	for _, n := range names {
		if p, err := findfont.Find(n); err == nil && p != "" {
			path = p
			break
		}
	}
	if path == "" {
		return nil, fmt.Errorf("no font found for stack %q", cssStack)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read font: %w", err)
	}
	ft, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse font: %w", err)
	}

	const dpi = 96.0
	pt := sizePx * 72.0 / dpi

	return opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    pt,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
}

func setHexOrWhite(dc *gg.Context, hex string) {
	if hex == "" {
		dc.SetColor(color.White)
		return
	}
	dc.SetHexColor(hex)
}

func pickFamily(cssStack string) string {
	var out []string
	cur := strings.Builder{}
	var quote rune
	for _, r := range strings.TrimSpace(cssStack) {
		if quote != 0 {
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
			continue
		}
		switch r {
		case '\'', '"':
			quote = r
		case ',':
			if s := strings.TrimSpace(cur.String()); s != "" {
				out = append(out, s)
			}
			cur.Reset()
		default:
			cur.WriteRune(r)
		}
	}
	if s := strings.TrimSpace(cur.String()); s != "" {
		out = append(out, s)
	}
	for _, f := range out {
		l := strings.ToLower(f)
		if l == "sans-serif" || l == "serif" || l == "monospace" || l == "system-ui" {
			continue
		}
		return f
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
