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
	"golang.org/x/sync/errgroup"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
)

type RunnerRequest struct {
	config      *models.Config
	JobID       string
	ExecutionID string
	WorkerID    string
	BoxID       int
	Blocks      []utils.CodeBlock
}

type RunnerResponse struct {
	EndTime time.Time
	Error   *string
	Output  *string
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

func (r *Runner) NewRequest(job models.Job, executionID string) (*RunnerRequest, error) {
	doc, err := htmlparse.Parse(strings.NewReader(job.Payload))
	if err != nil {
		return nil, err
	}

	d := utils.DocumentData{}
	r.parser.ExtractCodeBlocks(doc, &d)

	return &RunnerRequest{
		config:      r.config,
		ExecutionID: executionID,
		JobID:       job.JobID,
		WorkerID:    r.config.WorkerID,
		Blocks:      d.CodeBlocks,
	}, nil
}

func (r *Runner) RunCode(in RunnerRequest, reporter *Reporter) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	g, ctx := errgroup.WithContext(context.Background())

	for i, block := range in.Blocks {
		name, err := WriteToTempFile([]byte(block.Content), block.Language, *r.config)
		if err != nil {
			return err
		}
		defer os.Remove(name)

		if _, err = reporter.StartExecution(in.ExecutionID); err != nil {
			return err
		}

		g.Go(func() error {
			return runBlock(ctx, r.config.TempDir, name, i, len(in.Blocks), in.ExecutionID, reporter)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func runBlock(ctx context.Context, tempDir string, fileName string, blockNo int, totalBlocks int, executionID string, reporter *Reporter) error {
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

	logger.Infof(cmd.String())

	cmd.WaitDelay = time.Second * 60

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorf("block %d of %d errored", blockNo+1, totalBlocks)
		return err
	}

	err = cmd.Start()
	if err != nil {
		logger.Errorf("block %d of %d errored", blockNo+1, totalBlocks)
		return err
	}

	res := RunnerResponse{
		EndTime: time.Now().UTC(),
	}

	scanner := bufio.NewScanner(stdoutpipe)
	for scanner.Scan() {
		// TODO: run response
		m := scanner.Text()
		formatted := *res.Output + "\n" + m
		res.Output = &formatted
		logger.Infof(*res.Output)
		logger.Infof("block %d of %d finished", blockNo+1, totalBlocks)
		return nil
	}

	err = cmd.Wait()
	if err != nil {
		if strings.Contains(stderr.String(), "box is currently in use by another process") {
			logger.Errorf("block %d of %d errored", blockNo+1, totalBlocks)
			return ErrSandboxBusy
		}

		if strings.Contains(err.Error(), "exit status 2") {
			logger.Errorf("block %d of %d errored", blockNo+1, totalBlocks)
			return fmt.Errorf("isolate error: %v", stderr.String())
		}

		// TODO: process error responses
		// for _, errStr := range strings.Split(stderr.String(), "\n") {
		// }
		logger.Errorf("block %d of %d errored", blockNo+1, totalBlocks)
	}

	// reporter.CompleteExecution(executionID)

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
	return fmt.Sprintf("%v/%v%v", dir, name, utils.GetSupportedLanguages()[lang])
}
