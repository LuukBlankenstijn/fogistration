package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := flag.String("in", "docs/swagger.json", "input Swagger 2.0 JSON file")
	out := flag.String("out", "docs/swagger.json", "output JSON file (with basePath merged into paths)")
	flag.Parse()

	raw, err := os.ReadFile(*in)
	if err != nil {
		panic(fmt.Errorf("read input: %w", err))
	}

	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		panic(fmt.Errorf("unmarshal: %w", err))
	}

	// Get basePath & paths
	bp, _ := doc["basePath"].(string)
	paths, _ := doc["paths"].(map[string]any)

	// If nothing to do, just write through
	if bp == "" || len(paths) == 0 {
		_ = os.WriteFile(*out, raw, 0o644)
		fmt.Printf("No basePath or no paths; wrote input through to %s\n", *out)
		return
	}

	// Normalize basePath: ensure single leading slash, no trailing slash (except "/")
	bp = "/" + strings.Trim(strings.TrimSpace(bp), "/")
	if bp == "/" {
		// basePath is root; nothing to change
		delete(doc, "basePath")
	} else {
		// Build new paths with prefix
		newPaths := make(map[string]any, len(paths))
		for p, v := range paths {
			np := joinURL(bp, p)
			// Avoid double-prefixing if already present
			if strings.HasPrefix(p, bp+"/") || p == bp {
				np = p
			}
			newPaths[np] = v
		}
		doc["paths"] = newPaths
		delete(doc, "basePath")
	}

	outBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(fmt.Errorf("marshal: %w", err))
	}
	if err := os.WriteFile(*out, outBytes, 0o644); err != nil {
		panic(fmt.Errorf("write output: %w", err))
	}
	fmt.Printf("Wrote %s with basePath merged into path keys\n", *out)
}

// joinURL joins basePath and a path key safely.
func joinURL(base, p string) string {
	if p == "" || p == "/" {
		return base
	}
	// Swagger path keys always start with '/', but be defensive
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	// path.Join would clean "//" but also remove trailing slash semantics; build manually
	return strings.TrimRight(base, "/") + p
}
