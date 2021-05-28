package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type gitClient interface {
	CheckoutFileFromRepository(repository, commitHash, path string) ([]byte, error)
}

type basicGitClient struct {
	auth *ssh.PublicKeys
	mu   *sync.Mutex
}

func newBasicGitClient(sshPemFile string) (basicGitClient, error) {
	auth, err := ssh.NewPublicKeysFromFile("git", sshPemFile, "")
	if err != nil {
		return basicGitClient{}, err
	}

	return basicGitClient{
		auth: auth,
		mu:   &sync.Mutex{},
	}, nil
}

func (g basicGitClient) CheckoutFileFromRepository(repository, commitHash, path string) ([]byte, error) {
	filePath := filepath.Join(os.TempDir(), repository)

	g.mu.Lock()
	defer g.mu.Unlock()

	var repo *git.Repository

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// TODO: use context version and make depth configurable
		repo, err = git.PlainClone(filePath, false, &git.CloneOptions{
			URL:  repository,
			Auth: g.auth,
		})
		if err != nil {
			return []byte{}, err
		}
	} else {
		repo, err = git.PlainOpen(filePath)
		if err != nil {
			return []byte{}, err
		}
	}

	w, err := repo.Worktree()
	if err != nil {
		return []byte{}, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})
	if err != nil {
		return []byte{}, err
	}

	pathToManifest := filepath.Join(filePath, path)
	fileStat, err := os.Stat(pathToManifest)
	if err != nil {
		return []byte{}, err
	}

	if fileStat.IsDir() {
		return []byte{}, fmt.Errorf("path provided is not a file '%s'", path)
	}

	file, err := os.Open(pathToManifest)
	if err != nil {
		return []byte{}, err
	}

	fileContents := make([]byte, fileStat.Size())
	_, err = file.Read(fileContents)

	return fileContents, err
}
