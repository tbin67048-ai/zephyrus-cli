package utils

import (
	"fmt"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	cryptossh "golang.org/x/crypto/ssh"
)

func PurgeVault(session *Session) error {
	repoURL := fmt.Sprintf("git@github.com:%s/.nexus.git", session.Username)
	storer := memory.NewStorage()
	fs := memfs.New()

	publicKeys, _ := ssh.NewPublicKeys("git", session.RawKey, "")
	publicKeys.HostKeyCallback = cryptossh.InsecureIgnoreHostKey()

	r, _ := git.Init(storer, fs)
	w, _ := r.Worktree()

	commit, _ := w.Commit("Nexus: PURGE VAULT", &git.CommitOptions{
		Author:            &object.Signature{Name: "Nexus", Email: "nexus@cli.io", When: time.Now()},
		AllowEmptyCommits: true,
	})

	_, _ = r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{repoURL}})

	err := r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       publicKeys,
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("%s:refs/heads/master", commit))},
		Force:      true,
	})
	if err != nil {
		return err
	}

	// Reset local index to empty
	session.Index = NewIndex()
	return session.Save()
}
