package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/georgemessiha22/georgemessiha22/internal/render"
)

// Write persists artifacts under dir, creating subdirectories as needed.
func Write(dir string, artifacts []render.Artifact) error {
	for _, a := range artifacts {
		dest := filepath.Join(dir, a.RelPath)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("create dir for %s: %w", a.RelPath, err)
		}
		if err := os.WriteFile(dest, a.Bytes, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", a.RelPath, err)
		}
	}
	return nil
}
