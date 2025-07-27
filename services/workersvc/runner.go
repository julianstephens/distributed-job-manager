package workersvc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html"
	"os"
	"os/exec"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
)

func (w *Worker) RunCode(in utils.CodeBlock) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	name, err := WriteToTempFile([]byte(in.Content), in.Language, *w.config)
	if err != nil {
		return err
	}
	defer os.Remove(name)

	cmd := exec.CommandContext(ctx,
		"isolate",
		"--fsize=5120",
		fmt.Sprintf("--dir=%v", w.config.TempDir),
		"--dir=/root/.cache/go-build:rw",
		"--wait",
		"--processes=100",
		"--open-files=0",
		"--env=GOROOT",
		"--env=GOPATH",
		"--env=GO111MODULE=on",
		"--env=PATH",
		"--stderr-to-stdout",
		"--run",
		"--",
		name,
	)

	cmd.WaitDelay = time.Second * 60

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// stdoutpipe, err := cmd.StdoutPipe()
	// if err != nil {
	// 	return err
	// }

	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func WriteToTempFile(b []byte, lang string, conf models.Config) (string, error) {
	unscaped := html.UnescapeString(string(b))

	id := fmt.Sprintf("%v", time.Now().UnixNano())
	name := CodeFileName(conf.TempDir, id, lang)
	err := os.WriteFile(name, []byte(unscaped), 0777)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(conf.TempDir, os.ModePerm)
		if err != nil {
			return name, err
		}
	}

	err = os.WriteFile(name, []byte(unscaped), 0777)

	return name, err
}

func CodeFileName(dir string, name string, lang string) string {
	return fmt.Sprintf("%v/%v.%v", dir, name, utils.GetSupportedLanguages()[lang])
}
