package v0

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
	"github.com/rehearsal-open/rehearsal/task"
)

func (e RehearsalEngine) AssignConfig(conf *entity.Config) error {

	// assign
	e.config = conf

	// logger
	e.logger = &logger.Logger{}
	if err := e.logger.AssignConfig(conf); err != nil {
		return errors.WithStack(err)
	}

	e.logger.SystemPrint("Happy New Year")
	e.logger.SystemPrint(fmt.Sprintln("Hpppy world"))
	e.logger.SystemPrint(fmt.Sprintln("Hpppy\nworld"))

	// option's tasks
	{

		if regex, err := regexp.Compile(entity.TaskNameRegexp); err != nil {
			return errors.WithStack(err)
		} else {
			e.tasks = map[string]task.Task{}
			for _, taskConf := range e.config.TaskConf {
				if !regex.MatchString(taskConf.Name) {
					return errors.New("task's name is unmatch: use alphabet, number or underbar, don't begin with underbar")
				}
				var task task.Task
				switch taskConf.Type {
				case "CLI": // todo: define constant

				default:
					return errors.New("unsupported task's type: " + taskConf.Type + " (task's name is " + taskConf.Name + ")")
				}
				e.tasks[taskConf.Name] = task
			}
		}
	}

	// system's task (and config)
	{

		// output task

	}
	return nil
}
