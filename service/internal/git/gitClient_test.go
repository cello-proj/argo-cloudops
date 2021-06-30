package git

import (
	"io/fs"
	"os"
	"sync"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/google/go-cmp/cmp"
)

type mockGitSvc struct{}

func (g mockGitSvc) PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	return nil, nil
}

func (g mockGitSvc) PlainOpen(path string) (*git.Repository, error) {
	return nil, nil
}

func (g mockGitSvc) Fetch(r *git.Repository, o *git.FetchOptions) error {
	return nil
}

func (g mockGitSvc) Worktree(r *git.Repository) (*git.Worktree, error) {
	return nil, nil
}

func (g mockGitSvc) Checkout(w *git.Worktree, opts *git.CheckoutOptions) error {
	return nil
}

type mockOsSvc struct{}

func (o mockOsSvc) Stat(name string) (fs.FileInfo, error) {
	return nil, nil
}

func (o mockOsSvc) Open(name string) (*os.File, error) {
	return nil, nil
}

func (o mockOsSvc) IsDir(filestat fs.FileInfo) bool {
	return false
}

func (o mockOsSvc) Size(filestat fs.FileInfo) int64 {
	return 8
}

func (o mockOsSvc) Read(file *os.File, b []byte) (n int, err error) {
	copy(b, "my bytes")
	return 0, nil
}

func newGitClient() BasicClient {
	return BasicClient{
		auth: nil,
		mu:   &sync.Mutex{},
		git:  mockGitSvc{},
		os:   mockOsSvc{},
	}
}

func TestGetManifestFile(t *testing.T) {
	tests := []struct {
		name       string
		repository string
		commitHash string
		path       string
		errResult  bool
		res        string
	}{
		{
			name:       "get manifest success",
			repository: "myrepo",
			commitHash: "123",
			path:       "path/to/manifest.yaml",
			errResult:  false,
			res:        "my bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitClient := newGitClient()
			res, err := gitClient.GetManifestFile(tt.repository, tt.commitHash, tt.path)
			if err != nil {
				if !tt.errResult {
					t.Errorf("\ndid not expect error, got: %v", err)
				}
			} else {
				if tt.errResult {
					t.Errorf("\nexpected error")
				}
				if !cmp.Equal(string(res), tt.res) {
					t.Errorf("\nwant: %v\n got: %v", tt.res, string(res))
				}
			}
		})
	}
}
