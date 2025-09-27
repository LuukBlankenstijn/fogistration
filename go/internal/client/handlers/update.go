package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"reflect"

	"github.com/LuukBlankenstijn/fogistration/internal/client/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

type UpdateHandler struct {
	baseHandler
}

func (u *UpdateHandler) MessageType() reflect.Type {
	return reflect.TypeOf(&pb.ServerMessage_Update{})
}

func (u *UpdateHandler) HandleMessage(m *pb.ServerMessage) {
	msg := m.GetUpdate()
	if msg == nil {
		return
	}
	configBuilder := service.NewConfigBuilder(u.config.GreeterConfigPath)

	if contestUrl := msg.GetContestUrl(); contestUrl != nil {
		configBuilder.SetValue(service.ContestUrl, contestUrl.GetUrl())
	}

	if wallpaper := msg.GetWallpaper(); wallpaper != nil {
		filename, err := u.saveWallpaperPNG(wallpaper)
		if err != nil {
			logging.Error("failed to save wallpaper", err)
		} else {
			configBuilder.SetValue(service.WallPaper, filename)
		}
	}

	changes, err := configBuilder.Commit()
	if err != nil {
		logging.Error("failed to commit new config for wallpaper", err)
		return
	}

	if changes {
		err = service.RestartLightDM()
		if err != nil {
			logging.Error("failed to restart lightDM", err)
			return
		}
	}
}

func (u *UpdateHandler) saveWallpaperPNG(wallpaper *pb.Wallpaper) (string, error) {
	if wallpaper == nil || len(wallpaper.Data) == 0 {
		return "", errors.New("empty wallpaper data")
	}
	if wallpaper.Size != 0 && wallpaper.Size != uint64(len(wallpaper.Data)) {
		return "", fmt.Errorf("size mismatch: header=%d data=%d", wallpaper.Size, len(wallpaper.Data))
	}

	// Validate PNG (optional)
	if _, err := png.Decode(bytes.NewReader(wallpaper.Data)); err != nil {
		return "", fmt.Errorf("invalid PNG: %w", err)
	}

	// Hash content -> hex filename
	sum := sha256.Sum256(wallpaper.Data)
	hashHex := hex.EncodeToString(sum[:])

	dir := "/tmp" // or wherever you want to store deduped wallpapers
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	target := filepath.Join(dir, hashHex+".png")

	// If it already exists, reuse it
	if _, err := os.Stat(target); err == nil {
		abs, _ := filepath.Abs(target)
		return abs, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	// Create exclusively to avoid races with parallel writers
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		// If another process created it between Stat and OpenFile
		if os.IsExist(err) {
			abs, _ := filepath.Abs(target)
			return abs, nil
		}
		return "", err
	}
	defer f.Close()

	if _, err := f.Write(wallpaper.Data); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}

	abs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	return abs, nil
}
