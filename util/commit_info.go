package util

import (
	"bytes"
	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
)

type Server struct {
}

func CheckErr(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}

func (s *Server) Translate(ctx context.Context, in *CommitInfo) (*ServerResponse, error) {
	log.Printf("Receive message body from client: %s", in.HeadHash)
	// вот тут мне не то, чтобы все понятно, потому что у нас получается есть всегда на сервере
	// master версия репозитория. Да, есть, нужно прописать отдельную логику, которая его будет обновлять просто
	// пока что загружаю захардкоженный url

	dir := "hehe/"

	file, err := os.Create("branch_patch.diff")

	_, err = file.WriteString(in.CommitDiff)
	CheckErr(err, "Error while file writing")
	err = file.Close()
	CheckErr(err, "Error while file saving")

	patch, err := os.Open("branch_patch.diff")
	files, _, err := gitdiff.Parse(patch)
	err = patch.Close()
	CheckErr(err, "Failed patch reading")

	for _, f := range files {
		//TODO это надо потом убрать точно
		err := os.RemoveAll(dir)
		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL: "https://github.com/polupanovaanna/github_actions_test_project.git",
		})
		CheckErr(err, "Error when uploading git repository: %s")

		file, err := os.OpenFile(dir+f.OldName, os.O_CREATE|os.O_APPEND, os.ModePerm)
		CheckErr(err, "Error while opening "+f.OldName)

		var output bytes.Buffer
		err = gitdiff.Apply(&output, file, f)
		CheckErr(err, "Error while applying changes "+f.OldName)

		err = file.Close()
		CheckErr(err, "Error while closing "+f.OldName)

		err = ioutil.WriteFile(dir+f.OldName, output.Bytes(), 0)
		CheckErr(err, "Error while writing to file "+f.OldName)

		//TODO github actions ??

	}

	return &ServerResponse{Response: "Evetything is ok!"}, nil
}
