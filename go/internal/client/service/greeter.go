package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type configKey string

const (
	WallPaper  configKey = "greeter-background-image"
	ContestUrl configKey = "ccs-contest-api-url"
	Username   configKey = "ccs-autologin-username"
	Password   configKey = "ccs-autologin-password"
)

type GreeterConfigBuilder interface {
	SetValue(key configKey, value string)
	Commit() error
}

type configBuilder struct {
	values map[configKey]string
	path   string
}

func NewConfigBuilder(path string) GreeterConfigBuilder {
	return &configBuilder{
		values: make(map[configKey]string),
		path:   path,
	}
}

func (c *configBuilder) SetValue(key configKey, value string) {
	c.values[key] = value
}

func (c *configBuilder) Commit() error {
	if len(c.values) == 0 {
		return nil
	}
	if c.path == "" {
		return errors.New("config path not set")
	}

	var lines []string
	b, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(strings.NewReader(string(b)))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	done := make(map[configKey]bool, len(c.values))

	for i, raw := range lines {
		trim := strings.TrimLeft(raw, " \t")

		if trim == "" {
			continue
		}
		uc := strings.TrimLeft(trim, "; \t")
		kv := strings.SplitN(uc, "=", 2)
		if len(kv) < 1 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		if k == "" {
			continue
		}
		for wantK, val := range c.values {
			if k == string(wantK) {
				lines[i] = fmt.Sprintf("%s=%s", k, val)
				done[wantK] = true
				break
			}
		}
	}

	var bld strings.Builder
	for _, ln := range lines {
		bld.WriteString(ln)
		bld.WriteByte('\n')
	}
	for k, v := range c.values {
		if !done[k] {
			bld.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
	}

	tmp := c.path + ".tmp"
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(tmp, []byte(bld.String()), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, c.path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	c.values = make(map[configKey]string)
	return nil
}
