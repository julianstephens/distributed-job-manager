package worker

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"html"
	"os"
	"os/exec"
	"strings"
	"time"

	htmlparse "golang.org/x/net/html"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
)

type RunnerRequest struct {
	config   *models.Config
	JobID    string
	WorkerID string
	BoxID    int
	Blocks   []utils.CodeBlock
}

type Runner struct {
	config *models.Config
	box    *Sandbox
	parser utils.Parser
}

func NewRunner() *Runner {
	return &Runner{
		config: config.GetConfig(),
		parser: utils.Parser{},
	}
}

func (r *Runner) NewRequest(job models.Job) (*RunnerRequest, error) {
	doc, err := htmlparse.Parse(strings.NewReader(job.Payload))
	if err != nil {
		return nil, err
	}

	d := utils.DocumentData{}
	r.parser.ExtractCodeBlocks(doc, &d)

	return &RunnerRequest{
		config:   r.config,
		JobID:    job.JobID,
		WorkerID: r.config.WorkerID,
		Blocks:   d.CodeBlocks,
	}, nil
}

func (r *Runner) RunCode(in RunnerRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	for i, block := range in.Blocks {
		name, err := WriteToTempFile([]byte(block.Content), block.Language, *r.config)
		if err != nil {
			return err
		}
		defer os.Remove(name)

		err = runBlock(ctx, r.config.TempDir, name, i, len(in.Blocks))
		if err != nil {
			return err
		}
	}

	return nil
}

func runBlock(ctx context.Context, tempDir string, fileName string, blockNo int, totalBlocks int) error {
	cmd := exec.CommandContext(ctx,
		"isolate",
		"--fsize=5120",
		fmt.Sprintf("--dir=%v", tempDir),
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
		fileName,
	)

	cmd.WaitDelay = time.Second * 60

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorf("block %d of %d errored", blockNo, totalBlocks)
		return err
	}

	err = cmd.Start()
	if err != nil {
		logger.Errorf("block %d of %d errored", blockNo, totalBlocks)
		return err
	}

	scanner := bufio.NewScanner(stdoutpipe)
	for scanner.Scan() {
		// TODO: run response
		// m := scanner.Text()
		logger.Infof("block %d of %d finished", blockNo, totalBlocks)
	}

	err = cmd.Wait()
	if err != nil {
		if strings.Contains(stderr.String(), "box is currently in use by another process") {
			logger.Errorf("block %d of %d errored", blockNo, totalBlocks)
			return ErrSandboxBusy
		}

		if strings.Contains(err.Error(), "exit status 2") {
			logger.Errorf("block %d of %d errored", blockNo, totalBlocks)
			return fmt.Errorf("isolate error: %v", err.Error())
		}

		// TODO: process error responses
		// for _, errStr := range strings.Split(stderr.String(), "\n") {
		// }
		logger.Errorf("block %d of %d errored", blockNo, totalBlocks)
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
