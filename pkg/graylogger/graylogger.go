package graylogger

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
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
	ID             string    `json:"_id"`
	Version        string    `json:"version"`
	Host           string    `json:"host"`
	Message        string    `json:"message"`
	Trace          *string   `json:"full_message"`
	Timestamp      time.Time `json:"timestamp"`
	Level          LogLevel  `json:"level"`
	Origin         string    `json:"_origin"`
	AdditionalData *string   `json:"_additional_data"`
}

type LogOptions struct {
	err            *error
	additionalData *map[string]any
}

func NewLogger(originator string) (*GrayLogger, error) {
	conf := config.GetConfig()
	conn, err := queue.GetConnection(conf.Rabbit.LoggingUsername, conf.Rabbit.LoggingPassword, conf)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
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

	log := l.formatGELFLog(msg, LogLevelInfo, additionalData, nil)

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	l.doLog(ctx, data)
}

func (l *GrayLogger) Error(msg string, trace *error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var errText string
	if trace != nil {
		err := *trace
		errText = err.Error()
	}
	log := l.formatGELFLog(msg, LogLevelError, utils.If(trace != nil, &errText, nil), nil)

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	l.doLog(ctx, data)
}

func (l *GrayLogger) ErrorWithData(msg string, trace *error, additionalData *map[string]any) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var errText string
	if trace != nil {
		err := *trace
		errText = err.Error()
	}

	var formattedData *string
	if additionalData != nil {
		d, err := json.Marshal(additionalData)
		if err != nil {
			panic(err)
		}
		formattedData = utils.StringPtr(string(d))
	}

	log := l.formatGELFLog(msg, LogLevelError, utils.If(trace != nil, &errText, nil), formattedData)

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	l.doLog(ctx, data)
}

func (l *GrayLogger) Warn(msg string, additionalData *string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log := l.formatGELFLog(msg, LogLevelWarning, additionalData, nil)

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	l.doLog(ctx, data)
}

func (l *GrayLogger) Debug(msg string, additionalData *string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log := l.formatGELFLog(msg, LogLevelDebug, additionalData, nil)

	data, err := json.Marshal(&log)
	if err != nil {
		panic(err)
	}

	l.doLog(ctx, data)
}

func (l *GrayLogger) doLog(ctx context.Context, body []byte) {
	err := LogCh.PublishWithContext(
		ctx,
		"logs",
		l.Originator,
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

func (l *GrayLogger) formatGELFLog(msg string, level LogLevel, trace *string, additionalData *string) Log {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	log := Log{
		ID:        uuid.New().String(),
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
	if additionalData != nil {
		log.AdditionalData = additionalData
	}

	return log
}
