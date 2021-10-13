package config

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/zkcrescent/chaos/dingtalk"

	"github.com/sirupsen/logrus"
	"github.com/zkcrescent/chaos/logrotate"
)

// Logrus config
type Logrus struct {
	Env             string           `json:"env" yaml:"env" toml:"env"`
	Service         string           `json:"service" yaml:"service" toml:"service"`
	Level           string           `json:"level" yaml:"level" toml:"level"`
	FileHook        *FileHook        `json:"file_hook" yaml:"file_hook" toml:"file_hook"`
	DingtalkBotHook *DingtalkBotHook `json:"dingtalk_bot_hook" yaml:"dingtalk_bot_hook" toml:"dingtalk_bot_hook"`
}

// Logrus config
type FileHook struct {
	w io.Writer

	File      string `json:"file" yaml:"file" toml:"file"`
	LimitSize int    `json:"limit_size" yaml:"limit_size" toml:"limit_size"`
	Backups   int    `json:"backups" yaml:"backups" toml:"backups"`
	BufSize   int    `json:"buf_size" yaml:"buf_size" toml:"buf_size"`
}

func (h *FileHook) Init() error {
	if err := os.MkdirAll(path.Dir(h.File), os.ModePerm); err != nil {
		return err
	}

	f, err := logrotate.FileWriter(
		h.File,
		logrotate.Rotating(h.LimitSize, h.Backups),
		logrotate.Buffer(h.BufSize),
	)
	if err != nil {
		return err
	}

	h.w = f
	return nil
}

func (h *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *FileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel:
		_, err := h.w.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Printf("================write error================\n%v", err)
		}
		return nil
	case logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel:
		_, err := h.w.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Printf("================write error================\n%v", err)
		}
	default:
		return nil
	}
	return nil
}

type DingtalkBotHook struct {
	dingtalk.Bot `json:",inline" yaml:",inline"`

	Level string `json:"level" yaml:"level"`

	levels []logrus.Level
}

func (h *DingtalkBotHook) Init() error {
	l, err := logrus.ParseLevel(h.Level)
	if err != nil {
		return err
	}

	var levels []logrus.Level
	for i := logrus.PanicLevel; i <= l; i++ {
		levels = append(levels, i)
	}
	h.levels = levels

	return nil
}
func (h *DingtalkBotHook) Levels() []logrus.Level {
	return h.levels
}

func (h *DingtalkBotHook) Fire(entry *logrus.Entry) error {
	var mds []string
	mds = append(mds, dingtalk.MarkdownHeader(3, entry.Level.String()))

	var fields []string
	for k, v := range entry.Data {
		fields = append(fields, fmt.Sprintf("%v: %v", k, v))
	}

	mds = append(mds, dingtalk.MarkdownList(fields...))

	mds = append(mds, dingtalk.MarkdownInline(entry.Message))

	rs, err := h.Send(dingtalk.NewMarkdownMsg(entry.Level.String(), nil, false, mds...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "dingtalk bot send error, %v", err)
		return err
	}

	if !rs.IsSuccess() {
		fmt.Fprintf(os.Stderr, "dingtalk bot send failed, %+v", rs)
		return err
	}

	return nil
}

func (conf *Logrus) Init() error {
	if conf.Level != "" {
		level, err := logrus.ParseLevel(conf.Level)
		if err != nil {
			return err
		}
		logrus.SetLevel(level)
	}

	// clean old hooks
	logrus.StandardLogger().ReplaceHooks(make(logrus.LevelHooks))

	if conf.FileHook != nil {
		if err := conf.FileHook.Init(); err != nil {
			return err
		}

		logrus.AddHook(conf.FileHook)
		logrus.SetFormatter(
			&logrus.JSONFormatter{
				DisableTimestamp: false,
				TimestampFormat:  "2006-01-02 15:04:05.000",
				FieldMap:         nil,
				CallerPrettyfier: nil,
				PrettyPrint:      true,
			})

	}
	if conf.DingtalkBotHook != nil {
		if err := conf.DingtalkBotHook.Init(); err != nil {
			return err
		}

		logrus.AddHook(conf.DingtalkBotHook)
	}

	return nil
}
