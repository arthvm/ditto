package git_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arthvm/ditto/internal/git"
)

func SetupGitRepo(ctx context.Context) (string, error) {
	path, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	cmd := exec.CommandContext(ctx, "git", "init")
	cmd.Dir = path

	if err := cmd.Run(); err != nil {
		os.RemoveAll(path)
		return path, fmt.Errorf("init git repo: %w", err)
	}

	return path, nil
}

func TestStagedDiff(t *testing.T) {
	ctx := context.Background()

	repo, err := SetupGitRepo(ctx)
	require.NoError(t, err)
	defer os.RemoveAll(repo)

	helloFile, err := os.Create(fmt.Sprintf("%s/hello.txt", repo))
	require.NoError(t, err)
	defer helloFile.Close()
	helloFile.WriteString("Hello, World!\n")

	byeFile, err := os.Create(fmt.Sprintf("%s/bye.txt", repo))
	require.NoError(t, err)
	defer byeFile.Close()
	byeFile.WriteString("Goodbye, World!\n")

	cmd := exec.CommandContext(ctx, "git", "add", "hello.txt")
	cmd.Dir = repo
	require.NoError(t, cmd.Run())

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	os.Chdir(repo)
	diff, err := git.Diff(ctx, git.Staged)

	assert.NoError(t, err)
	assert.Equal(t, `diff --git a/hello.txt b/hello.txt
new file mode 100644
index 0000000..8ab686e
--- /dev/null
+++ b/hello.txt
@@ -0,0 +1 @@
+Hello, World!
`, diff)
}
