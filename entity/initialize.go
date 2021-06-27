package entity

import (
	"os"

	"github.com/pkg/errors"

	color "github.com/rehearsal-open/rehearsal/entity/cli-color"
)

func (c *Conf) Initialize() error {

	// global config

	// dir
	if c.Dir == "" {
		c.Dir, _ = os.Getwd()
	}

	// syncMs
	if c.SyncMs < 1 {
		c.SyncMs = 1
	}

	// maxnamelen
	c.MaxNameLen = 0

	// tasks
	for name, task := range c.Tasks {

		// maxnamelen
		if len(name) > c.MaxNameLen {
			c.MaxNameLen = len(name)
		}

		switch task.Type {
		case CLITask:
			// task.path
			if task.Path == "" {
				errors.New("task is required, but cannot found it")
			}
		default:
			return errors.New("invalid task type: " + task.Type)
		}

		// task.name
		task.Name = name

		// task.colorstr
		if task.ColorStr == "" {
			task.Color = color.Default
		} else if col, err := color.FromString(task.ColorStr); err != nil {
			return errors.WithStack(err)
		} else {
			task.Color = col
		}

		// task.sendto
		for _, sendTo := range task.SendTo {
			if _, exist := c.Tasks[sendTo]; !exist {
				return errors.New("cannot found reciever task: " + sendTo + " (send from: " + name + ")")
			}
		}

		// task.syncms
		if task.SyncMs < 1 {
			task.SyncMs = c.SyncMs
		}
	}
	return nil
}
