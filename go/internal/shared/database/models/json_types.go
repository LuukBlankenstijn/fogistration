package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Align string

const (
	AlignLeft   Align = "left"
	AlignCenter Align = "center"
	AlignRight  Align = "right"
)

type LabelSpec struct {
	X      int    `json:"x"      binding:"required"`
	Y      int    `json:"y"      binding:"required"`
	Size   int    `json:"size"   binding:"required"`
	Color  string `json:"color"  binding:"required"`
	Weight int    `json:"weight" binding:"required"`
	Align  Align  `json:"align"  binding:"required,oneof=left center right"`
}

type WallpaperConfigJSON struct {
	W         int       `json:"w"         binding:"required"`
	H         int       `json:"h"         binding:"required"`
	FontStack string    `json:"fontStack" binding:"required"`
	Teamname  LabelSpec `json:"teamname"  binding:"required"`
	IP        LabelSpec `json:"ip"        binding:"required"`
}

func (w WallpaperConfigJSON) Value() (driver.Value, error) {
	return json.Marshal(w)
}
func (w *WallpaperConfigJSON) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, w)
	case string:
		return json.Unmarshal([]byte(v), w)
	default:
		return fmt.Errorf("cannot scan %T into WallpaperConfigJSON", src)
	}
}

func DefaultWallpaperConfig() WallpaperConfigJSON {
	return WallpaperConfigJSON{
		W:         1920,
		H:         1080,
		FontStack: "Inter, system-ui, Arial, sans-serif",
		Teamname: LabelSpec{
			X:      1920 / 2,
			Y:      200,
			Size:   88,
			Color:  "#ffffff",
			Weight: 700,
			Align:  AlignCenter,
		},
		IP: LabelSpec{
			X:      20,
			Y:      40,
			Size:   48,
			Color:  "#ffffff",
			Weight: 500,
			Align:  AlignLeft,
		},
	}
}
