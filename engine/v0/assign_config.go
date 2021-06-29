package v0

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/out"
)

func (e *RehearsalEngine) AssignConfig(conf *entity.Config) error {

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

	stdoutTasks := make(map[string]task.Task, 0)

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

				// append system task reciever
				if taskConf.ShowOut {
					stdoutTasks[taskConf.Name] = task
				}
			}
		}
	}

	// system task (and config)
	{

		const (
			outTaskName string = "$std-out"
		)

		// standard output task
		outTaskConf := entity.TaskConfig{
			Name: outTaskName,
		}

		e.config.TaskConf = append(e.config.TaskConf, outTaskConf)

		outTask := &out.Task{}
		outTask.AssignEngine(e, outTaskName)
		outTask.AssignLogger(e.logger)
		e.tasks[outTaskName] = task.Task(outTask)
		for name, t := range stdoutTasks {
			if outtask, ok := t.(task.OutTask); !ok {
				return errors.New("cannot standard out, this is reciever only: " + name)
			} else if err := outtask.AppendTaskAsOut(outTask); err != nil {
				return errors.WithStack(err)
			}
		}

		e.logger.SystemPrint(fmt.Sprint(e))

	}
	return nil
}
