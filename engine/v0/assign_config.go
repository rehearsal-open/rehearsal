package v0

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/rehearsal-open/rehearsal/entity"
	color "github.com/rehearsal-open/rehearsal/entity/cli-color"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/cli"
)

func (e *RehearsalEngine) AssignConfig(config *entity.Conf) error {
	e.Conf = config

	// directory
	{
		var dir string

		if !filepath.IsAbs(e.Conf.Dir) {
			dir = filepath.Join(e.Conf.ConfigDir, e.Conf.Dir)
		} else {
			dir = e.Conf.Dir
		}
		if f, err := os.Stat(dir); os.IsNotExist(err) || !f.IsDir() {
			return errors.New("directory cannot found: " + dir)
		} else {
			e.Dir = dir
		}
	}

	// io sync
	if e.Conf.SyncMs < 1 {
		e.Conf.SyncMs = 1
	}

	e.Conf.TasksMap = make(map[string]*entity.TaskConf)
	for iTask, _ := range e.Conf.Tasks {

		t := &e.Conf.Tasks[iTask]

		// name
		if t.Name == "" {
			return errors.New("task's name is required, but cannot found: No." + strconv.Itoa(iTask) + " task")
		}
		// path
		if t.Path == "" {
			return errors.New("task's path is required, but cannot found: " + t.Name)
		}
		// io sync
		if t.SyncMs < 1 {
			e.Conf.SyncMs = 1
		}
		// max length
		if lenName := len(t.Name); lenName > e.Conf.MaxNameLen {
			e.Conf.MaxNameLen = lenName
		}

		// color
		t.Color, _ = color.FromString(t.ColorStr)
		if t.Color == color.Default {
			switch iTask % 3 {
			case 0:
				t.Color = color.Green
			case 1:
				t.Color = color.Syan
			case 2:
				t.Color = color.Magenta
			}
		}

		// task
		t.Type = strings.ToUpper(t.Type)
		switch t.Type {
		case string(task.CLI):
			e.tasks[t.Name] = &cli.Task{}
		default:
			return errors.New("task's kind unsupported: " + t.Type)
		}
		if err := e.tasks[t.Name].AssignEngine(e, t.Name); err != nil {
			return errors.WithStack(err)
		}

		e.TasksMap[t.Name] = t
	}

	return nil
}
