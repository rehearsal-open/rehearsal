package v0

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/cli"
	"github.com/rehearsal-open/rehearsal/task/out"
)

func (e *RehearsalEngine) AssignConfig(conf *entity.Config) error {

	const (
		outTaskName string = "$std-out"
	)

	// assign
	e.config = conf

	// logger
	e.logger = &logger.Logger{}
	if err := e.logger.AssignConfig(conf); err != nil {
		return errors.WithStack(err)
	}

	// syncMs
	if e.config.SyncMs < 1 {
		e.config.SyncMs = 1
	}

	// option's tasks
	{

		if regex, err := regexp.Compile(entity.TaskNameRegexp); err != nil {
			return errors.WithStack(err)
		} else {
			e.tasks = map[string]task.Task{}
			for i, _ := range e.config.TaskConf {

				taskConf := &e.config.TaskConf[i]

				if taskConf.SyncMs < 1 {
					taskConf.SyncMs = e.config.SyncMs
				}

				if !regex.MatchString(taskConf.Name) {
					return errors.New("task's name is unmatch: use alphabet, number or underbar, don't begin with underbar")
				}
				var task task.Task

				switch taskConf.Type {
				case "CLI": // todo: define constant
					task = &cli.Task{}
				default:
					return errors.New("unsupported task's type: " + taskConf.Type + " (task's name is " + taskConf.Name + ")")
				}
				e.tasks[taskConf.Name] = task

				if err := task.AssignEngine(e, taskConf, taskConf.Name); err != nil {
					return errors.WithStack(err)
				}

				// append system task reciever
				if taskConf.ShowOut {
					taskConf.SendTo = append(taskConf.SendTo, outTaskName)
				}
			}
		}
	}

	// system task (and config)
	{

		// standard output task
		outTaskConf := entity.TaskConfig{
			Name: outTaskName,
		}

		e.config.TaskConf = append(e.config.TaskConf, outTaskConf)

		outTask := &out.Task{}
		outTask.AssignEngine(e, &e.config.TaskConf[len(e.config.TaskConf)-1], outTaskName)
		outTask.AssignLogger(e.logger)
		e.tasks[outTaskName] = task.Task(outTask)
	}

	for i, _ := range e.config.TaskConf {

		taskConf := &e.config.TaskConf[i]

		if len(taskConf.SendTo) > 0 {

			if t, ok := e.tasks[taskConf.Name].(task.OutTask); !ok {
				return errors.New("sendTo property is used only task which can send: " + taskConf.Name)
			} else {
				for _, sendTo := range taskConf.SendTo {
					if reciever, ok := e.tasks[sendTo].(task.RecieverTask); !ok {
						return errors.New("selected task which selected as sendTo property is cannot recieve: " + sendTo)
					} else {
						t.AppendTaskAsOut(reciever)
					}
				}
			}

		}
	}
	return nil
}
