package service

import (
	"bufio"
	"bytes"
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
	Commit() (bool, error)
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

func (c *configBuilder) Commit() (bool, error) {
	if len(c.values) == 0 {
		return false, nil
	}
	if c.path == "" {
		return false, errors.New("config path not set")
	}

	lines, err := readLinesIfExists(c.path)
	if err != nil {
		return false, err
	}

	changed, lines := applyUpdates(lines, c.values)

	if !changed {
		c.values = make(map[configKey]string)
		return false, nil
	}

	if err := writeAtomically(c.path, render(lines)); err != nil {
		return false, err
	}
	c.values = make(map[configKey]string)
	return true, nil
}

/* ---------- helpers ---------- */

func readLinesIfExists(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

// parseLine returns (key, value, isKV, isCommented).
func parseLine(raw string) (k, v string, isKV, commented bool) {
	trim := strings.TrimLeft(raw, " \t")
	if trim == "" {
		return "", "", false, false
	}
	commented = strings.HasPrefix(trim, ";")
	uc := strings.TrimLeft(trim, "; \t") // uncommented-content for parsing
	kv := strings.SplitN(uc, "=", 2)
	if len(kv) == 0 || strings.TrimSpace(kv[0]) == "" {
		return "", "", false, commented
	}
	k = strings.TrimSpace(kv[0])
	if len(kv) == 2 {
		v = strings.TrimSpace(kv[1])
	}
	return k, v, true, commented
}

func applyUpdates(lines []string, updates map[configKey]string) (bool, []string) {
	changed := false
	done := make(map[configKey]bool, len(updates))

	for i, raw := range lines {
		k, existing, isKV, commented := parseLine(raw)
		if !isKV {
			continue
		}
		if val, ok := updates[configKey(k)]; ok {
			done[configKey(k)] = true

			// desired empty string => ensure commented entry (keep old value if present)
			if val == "" {
				if !commented {
					// comment the existing line (preserve its value)
					lines[i] = commentify(k, existing)
					changed = true
				}
				continue
			}

			// desired non-empty => ensure uncommented k=val
			newLine := fmt.Sprintf("%s=%s", k, val)
			if commented || existing != val || !isExactLine(raw, newLine) {
				lines[i] = newLine
				changed = true
			}
		}
	}

	// append missing keys
	for k, v := range updates {
		if done[k] {
			continue
		}
		if v == "" {
			lines = append(lines, commentify(string(k), "")) // commented placeholder
		} else {
			lines = append(lines, fmt.Sprintf("%s=%s", k, v))
		}
		changed = true
	}

	return changed, lines
}

func commentify(k, existing string) string {
	// keep existing value if we have it; otherwise empty
	if existing != "" {
		return fmt.Sprintf("; %s=%s", k, existing)
	}
	return fmt.Sprintf("; %s=", k)
}

func isExactLine(raw, want string) bool {
	// treat exact same (including whitespace) as unchanged
	return strings.TrimRight(raw, "\r\n") == want
}

func render(lines []string) []byte {
	var b strings.Builder
	for _, ln := range lines {
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	return []byte(b.String())
}

func writeAtomically(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
