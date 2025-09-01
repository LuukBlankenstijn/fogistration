package client

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"reflect"

	"github.com/LuukBlankenstijn/fogistration/internal/client/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

type WallpaperHandler struct {
	config config.ClientConfig
}

func (w *WallpaperHandler) MessageType() reflect.Type {
	return reflect.TypeOf(&pb.ServerMessage_Wallpaper{})
}

func (w *WallpaperHandler) HandleMessage(m *pb.ServerMessage) {
	msg := m.GetWallpaper()
	if msg == nil {
		return
	}

	logging.Info("wallpaper message")
	filename, err := w.saveWallpaperPNG(msg)
	if err != nil {
		logging.Error("failed to save wallpaper", err)
		return
	}

	configBuilder := service.NewConfigBuilder(w.config.GreeterConfigPath)
	configBuilder.SetValue(service.WallPaper, filename)

	err = configBuilder.Commit()
	if err != nil {
		logging.Error("failed to commit new config for wallpaper", err)
		return
	}

	err = service.RestartLightDM()
	if err != nil {
		logging.Error("failed to restart lightDM", err)
		return
	}
}

func (w *WallpaperHandler) SetConfig(config config.ClientConfig) {
	w.config = config
}

func (w *WallpaperHandler) saveWallpaperPNG(wallpaper *pb.Wallpaper) (string, error) {
	if wallpaper == nil || len(wallpaper.Data) == 0 {
		return "", errors.New("empty wallpaper data")
	}
	if wallpaper.Size != 0 && wallpaper.Size != uint64(len(wallpaper.Data)) {
		return "", fmt.Errorf("size mismatch: header=%d data=%d", wallpaper.Size, len(wallpaper.Data))
	}

	// (Optional) validate it's a PNG; remove if you don't want the extra check.
	if _, err := png.Decode(bytes.NewReader(wallpaper.Data)); err != nil {
		return "", fmt.Errorf("invalid PNG: %w", err)
	}

	f, err := os.CreateTemp("/tmp", "wallpaper-*.png")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write(wallpaper.Data); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil { // ensure bytes hit disk
		return "", err
	}

	if err := f.Chmod(0o644); err != nil {
		return "", err
	}

	abs, err := filepath.Abs(f.Name())
	if err != nil {
		return "", err
	}
	return abs, nil
}
