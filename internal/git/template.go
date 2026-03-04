package git

import (
	"os"
	"path/filepath"
)

func FindPRTemplate(root string) (string, error) {
	paths := []string{
		filepath.Join(root, ".github", "pull_request_template.md"),
		filepath.Join(root, "docs", "pull_request_template.md"),
		filepath.Join(root, "PULL_REQUEST_TEMPLATE.md"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			content, err := os.ReadFile(p)
			if err != nil {
				return "", err
			}

			return string(content), nil
		}
	}

	return "", nil
}
