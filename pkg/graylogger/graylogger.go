package graylogger

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/rabbitmq/amqp091-go"
)

var (
	LogCh *amqp091.Channel
)

type GrayLogger struct {
	LogCh      *amqp091.Channel
	LogConn    *amqp091.Connection
	Originator string
}

type LogLevel int64

const (
	LogLevelEmergency LogLevel = iota
	LogLevelAlert
	LogLevelCritical
	LogLevelError
	LogLevelWarning
	LogLevelNotice
	LogLevelInfo
	LogLevelDebug
)

type Log struct {
	Version   string    `json:"version"`
	Host      string    `json:"host"`
	Message   string    `json:"message"`
	Trace     *string   `json:"full_message"`
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Origin    string    `json:"_origin"`
}

func NewLogger(originator string) (*GrayLogger, error) {
	conf := config.GetConfig()
	conn, ch, err := queue.GetConnection(conf)
	if err != nil {
		return nil, err
	}

	LogCh = ch

	return &GrayLogger{
		LogConn:    conn,
		LogCh:      ch,
		Originator: originator,
	}, nil
}

func (l *GrayLogger) Info(msg string, additionalData *string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log := l.formatGELFLog(msg, LogLevelInfo, additionalData)

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	go l.doLog(ctx, data)
}

func (l *GrayLogger) Error(msg string, trace *error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var errText string
	if trace != nil {
		err := *trace
		errText = err.Error()
	}
	log := l.formatGELFLog(msg, LogLevelError, httputil.If(trace != nil, &errText, nil))

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	go l.doLog(ctx, data)
}

func (l *GrayLogger) Warn(msg string, additionalData *string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log := l.formatGELFLog(msg, LogLevelWarning, additionalData)

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	go l.doLog(ctx, data)
}

func (l *GrayLogger) Debug(msg string, additionalData *string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log := l.formatGELFLog(msg, LogLevelDebug, additionalData)

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	go l.doLog(ctx, data)
}

func (l *GrayLogger) doLog(ctx context.Context, body []byte) {
	err := LogCh.PublishWithContext(
		ctx,
		"logs",
		"jobsvc",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		panic(err)
	}
}

func (l *GrayLogger) formatGELFLog(msg string, level LogLevel, trace *string) Log {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	log := Log{
		Version:   "1.1",
		Host:      hostname,
		Message:   msg,
		Timestamp: time.Now().UTC(),
		Level:     level,
		Origin:    l.Originator,
	}
	if trace != nil {
		log.Trace = trace
	}

	return log
}
