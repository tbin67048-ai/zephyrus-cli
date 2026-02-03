package utils

import (
	"fmt"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	cryptossh "golang.org/x/crypto/ssh"
)

func DeleteFile(vaultPath string, session *Session) error {
	repoURL := fmt.Sprintf("git@github.com:%s/.nexus.git", session.Username)
	storer := memory.NewStorage()
	fs := memfs.New()

	publicKeys, _ := ssh.NewPublicKeys("git", session.RawKey, "")
	publicKeys.HostKeyCallback = cryptossh.InsecureIgnoreHostKey()

	// 1. Clone to get current state
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:           repoURL,
		Auth:          publicKeys,
		ReferenceName: plumbing.ReferenceName("refs/heads/master"),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return err
	}

	w, _ := r.Worktree()

	// 2. Check local index for the file
	entry, exists := session.Index[vaultPath]
	if !exists {
		return fmt.Errorf("file '%s' not found in vault", vaultPath)
	}

	// 3. Remove physical hex file from Git
	_, err = w.Remove(entry.RealName)
	if err != nil {
		fmt.Printf("Warning: storage file %s already missing from remote\n", entry.RealName)
	}

	// 4. Update index and re-encrypt
	delete(session.Index, vaultPath)
	newIndexBytes, _ := session.Index.ToBytes(session.Password)

	idxFile, _ := fs.Create(".config/index")
	idxFile.Write(newIndexBytes)
	idxFile.Close()
	w.Add(".config/index")

	// 5. Commit and Push
	commit, _ := w.Commit("Nexus: Deleted "+vaultPath, &git.CommitOptions{
		Author: &object.Signature{Name: "Nexus", Email: "nexus@cli.io", When: time.Now()},
	})

	err = r.Push(&git.PushOptions{
		Auth:     publicKeys,
		RefSpecs: []config.RefSpec{config.RefSpec(fmt.Sprintf("%s:refs/heads/master", commit))},
	})
	if err != nil {
		return err
	}

	// 6. Sync local session
	return session.Save()
}
