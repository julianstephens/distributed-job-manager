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
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
)

type RunnerRequest struct {
	config            *models.Config
	JobID             string
	ExecutionID       string
	WorkerID          string
	BoxID             int
	ExpectedStartTime time.Time
	Blocks            []utils.CodeBlock
}

type RunnerResponse struct {
	StartTime time.Time
	EndTime   time.Time
	Error     *string
	Output    *string
}

type Result struct {
	Value RunnerResponse
	Err   error
}

type Runner struct {
	config *models.Config
	log    *graylogger.GrayLogger
	box    *Sandbox
	parser utils.Parser
}

const Timeout = 120 * time.Second

func NewRunner(log *graylogger.GrayLogger) *Runner {
	return &Runner{
		config: config.GetConfig(),
		parser: utils.Parser{},
		log:    log,
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
		config:            r.config,
		ExecutionID:       executionID,
		JobID:             job.JobID,
		ExpectedStartTime: job.ExecutionTime,
		WorkerID:          r.config.WorkerID,
		Blocks:            d.CodeBlocks,
	}, nil
}

func (r *Runner) RunCode(in RunnerRequest, reporter *Reporter) error {
	block := in.Blocks[0]

	name, err := writeToTempFile([]byte(block.Content), block.Language, *r.config)
	if err != nil {
		return err
	}
	defer os.Remove(name)

	timeToStart := time.Until(in.ExpectedStartTime) - (10 * time.Second)
	if timeToStart > 0 {
		r.log.Info(fmt.Sprintf("waiting %v before starting job %s", timeToStart, in.JobID), nil)
		time.Sleep(timeToStart)
	}

	r.log.Info(fmt.Sprintf("starting execution %s for job %s", in.ExecutionID, in.JobID), nil)

	if _, err = reporter.StartExecution(in.ExecutionID); err != nil {
		return err
	}

	results := make(chan Result)

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	go runBlock(ctx, in.BoxID, r.config.TempDir, name, results)

	result := <-results

	if result.Err != nil {
		logger.Errorf("execution %s failed with error: %v", in.ExecutionID, result.Err)
		r.log.Error(fmt.Sprintf("execution %s failed", in.ExecutionID), &result.Err)
		return result.Err
	}

	r.log.Info(fmt.Sprintf("execution %s completed successfully", in.ExecutionID), utils.StringPtr(string(utils.MustMarshalJson(result))))
	if _, err = reporter.CompleteExecution(in.ExecutionID, result.Value); err != nil {
		return err
	}

	return nil
}

func runBlock(ctx context.Context, boxId int, tempDir string, fileName string, results chan<- Result) {
	cmd := exec.CommandContext(ctx,
		"isolate",
		fmt.Sprintf("--box-id=%v", boxId),
		// max size (in KB) of files that can be created per execution = 5MB
		"--fsize=5120",
		// makes directory visible in the sandbox
		fmt.Sprintf("--dir=%v", tempDir),
		// give read write access to the go cache dir as it needs to be cleaned
		"--dir=/root/.cache/go-build:rw",
		// if sandbox is busy, wait instead of returning error right away
		// instead of serving 25/100 requests in 10 sandbox, it's gonna serve all
		"--wait",
		// to keep the child process in parent’s network namespace and communicate with the outside world
		"--share-net",
		"--processes=100",
		// unlimited open files
		"--open-files=0",
		"--env=GOROOT",
		"--env=GOPATH",
		"--env=GO111MODULE=on",
		"--env=HOME",
		// makes commands visible in the sandbox e.g. 'ls', 'echo' or other installed command
		"--env=PATH",
		// log package writes to stderr instead of stdout, so we need to redirect this to stdout.
		// only exit code determines if the program ran successfully or not
		"--stderr-to-stdout",
		"--run",
		"--",
		"/usr/local/go/bin/go",
		"run",
		fileName,
	)

	cmd.WaitDelay = Timeout - (10 * time.Second) // give some time to isolate to clean up the sandbox

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		results <- Result{Value: RunnerResponse{}, Err: err}
		return
	}

	res := RunnerResponse{
		StartTime: time.Now().UTC(),
	}
	err = cmd.Start()
	if err != nil {
		results <- Result{Value: res, Err: err}
		return
	}

	scanner := bufio.NewScanner(stdoutpipe)
	for scanner.Scan() {
		m := scanner.Text()
		res.Output = utils.StringPtr(m)
	}

	err = cmd.Wait()
	if err != nil {
		if strings.Contains(stderr.String(), "box is currently in use by another process") {
			results <- Result{Value: res, Err: ErrSandboxBusy}
			return
		}

		if strings.Contains(err.Error(), "exit status 2") {
			results <- Result{Value: res, Err: fmt.Errorf("isolate error: %v", stderr.String())}
			return
		}

		// TODO: process error responses
		// for _, errStr := range strings.Split(stderr.String(), "\n") {
		// }
	}

	res.EndTime = time.Now().UTC()
	results <- Result{Value: res, Err: nil}
}

func writeToTempFile(b []byte, lang string, conf models.Config) (string, error) {
	unscaped := html.UnescapeString(string(b))

	id := fmt.Sprintf("%v", time.Now().UnixNano())
	fileName := codeFileName(conf.TempDir, id, lang)
	err := os.WriteFile(fileName, []byte(unscaped), 0777)

	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(conf.TempDir, os.ModePerm)
		if err != nil {
			return fileName, err
		}
		err = os.WriteFile(fileName, []byte(unscaped), 0777)
	}

	return fileName, err
}

func codeFileName(dir string, name string, lang string) string {
	return fmt.Sprintf("%v/%v%v", dir, name, utils.GetSupportedLanguages()[lang])
}
