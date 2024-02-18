package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	f, errFile := os.OpenFile("./test.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, fs.ModeAppend)

	if errFile != nil {
		fmt.Print("cant create file")
		fmt.Println(errFile)
		os.Exit(-1)
	}

	f.WriteString("test 1\n")

	wd, errWd := os.Getwd()

	if errWd != nil {
		fmt.Println("cant get wd")
		os.Exit(-2)
	}

	r, err := git.PlainOpen(wd)
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	Info("git add .")
	_, err = w.Add(".")
	CheckIfError(err)

	Info("git status --porcelain")
	status, err := w.Status()
	CheckIfError(err)

	fmt.Println(status)

	Info("git commit -m \"push new configs \"")

	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Amir Hemmateenejad",
			Email: "amirhemmateenejad@gmail.com",
			When:  time.Now(),
		},
	})

	CheckIfError(err)

	// Prints the current HEAD to verify that all worked well.
	Info("git show -s")
	obj, err := r.CommitObject(commit)
	CheckIfError(err)

	fmt.Println(obj)

	Info("git push")

	errPush := r.Push(&git.PushOptions{})
	CheckIfError(errPush)

}
