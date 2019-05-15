// Git functions based on go-git https://github.com/src-d/go-git/
package git

import (
	git "gopkg.in/src-d/go-git.v4"
	//	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"fmt"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"log"
	"os"
	"time"
)

func Clone(source string, userName string, accessToken string) (fileystem billy.Filesystem, repo *git.Repository) {
	fs := memfs.New()
	repository, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: source,
		Auth: &http.BasicAuth{
			Username: userName, // This just has to be not null if we are using a token
			Password: accessToken,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Unable to clone: %v", err)
	}

	return fs, repository
}

func Commit(repo *git.Repository, fileName string, branch string) {
	w, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Unable to commit %v", err)
	}

	_, err = w.Add(fileName)
	if err != nil {
		log.Fatalf("Unable to commit %v", err)
	}

	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	fmt.Println(commit)
	if err != nil {
		log.Fatalf("Unable to commit %v", err)
	}

}

func Push(repo *git.Repository) {
	err := repo.Push(&git.PushOptions{})
	if err != nil {
		log.Fatalf("Unable to push branch %v", err)
	}
}

func Branch(Name string, repo *git.Repository) {
	headRef, err := repo.Head()
	ref := plumbing.NewHashReference("refs/heads/my-branch", headRef.Hash())
	err = repo.Storer.SetReference(ref)
	if err != nil {
		log.Fatalf("Unable to create branch %v", err)
	}

}

func MergeRequest(repo *git.Repository) {

}
