package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type Align string

const (
	AlignLeft   Align = "left"
	AlignCenter Align = "center"
	AlignRight  Align = "right"
)

var AlignValues = []Align{
	AlignLeft,
	AlignCenter,
	AlignRight,
}

func (Align) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["Align"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "Align")
		schemaRef.Title = "Align"
		for _, v := range AlignValues {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["Align"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/Align"}
}

type LabelSpec struct {
	X      int32  `json:"x"     `
	Y      int32  `json:"y"     `
	Size   int32  `json:"size"  `
	Color  string `json:"color" `
	Weight int32  `json:"weight"`
	Align  Align  `json:"align"`
}

type WallpaperLayout struct {
	W         int32     `json:"w"        `
	H         int32     `json:"h"        `
	FontStack string    `json:"fontStack"`
	Teamname  LabelSpec `json:"teamname" `
	IP        LabelSpec `json:"ip"       `
}

func (w WallpaperLayout) Value() (driver.Value, error) {
	return json.Marshal(w)
}
func (w *WallpaperLayout) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, w)
	case string:
		return json.Unmarshal([]byte(v), w)
	default:
		return fmt.Errorf("cannot scan %T into WallpaperLayout", src)
	}
}

func DefaultWallpaperConfig() WallpaperLayout {
	return WallpaperLayout{
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
