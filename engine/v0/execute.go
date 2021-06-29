package v0

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/task"
)

func (e *RehearsalEngine) Run() error {
	// todo: write definition

	type exitTask struct {
		err  error
		task task.Task
	}

	exit := make(chan exitTask)
	nTask := len(e.tasks)
	iTask := 0

	for name, t := range e.tasks {
		e.logger.SystemPrint(fmt.Sprint("initialize task (", 1+iTask, "/", nTask, " : ", name, ")..."))
		if err := t.RunInit(); err != nil {
			return errors.WithStack(err)
		}
		iTask++
	}

	iTask = 0

	for _, t := range e.tasks {

		e.logger.SystemPrint(fmt.Sprint("running start(", iTask+1, "/", nTask, " : ", t.Config().Name, ")..."))
		go func(t task.Task, exit chan exitTask) {
			exit <- exitTask{
				err:  t.RunWait(),
				task: t,
			}
		}(t, exit)
		iTask++
	}

	defer func() {
		for _, t := range e.tasks {
			e.logger.SystemPrint(fmt.Sprint("finalize: " + t.Config().Name))
			t.Finalize()
			e.logger.SystemPrint("finished")
		}
	}()

	iTask = 0

	for iTask < nTask {

		select {
		case exited := <-exit:
			if exited.err != nil {
				e.logger.SystemPrint(fmt.Sprint("error occered at ", exited.task.Config().Name, ": ", exited.err.Error()))
			}
			e.logger.SystemPrint(fmt.Sprint("running end(", iTask+1, "/", nTask, " : ", exited.task.Config().Name, ")"))
			iTask++
		}
	}

	return nil
}

func (e *RehearsalEngine) Kill() {
	// todo: write definition
}
