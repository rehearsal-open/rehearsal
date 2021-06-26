package executer

import (
	"log"

	"github.com/pkg/errors"
)

func (manager *ExecuteManager) ExecuteAll() (err error) {

	type exitCode struct {
		name string
		err  error
	}

	nExecutes := len(manager.executes)
	exits := make(chan exitCode)
	err = nil

	outputsChs := make([]chan Packet, 0)
	errorsChs := make([]chan Packet, 0)
	for name, exec := range manager.executes {

		outputCh := make(chan Packet)
		errorsCh := make(chan Packet)
		exec.Exec.AppendOutPacketChannel(outputCh)
		exec.Exec.AppendErrPacketChannel(errorsCh)
		outputsChs = append(outputsChs, outputCh)
		errorsChs = append(errorsChs, errorsCh)

		go func(name string, execute *ExecuteItem, outCh chan Packet, errCh chan Packet) {
			exist := true
			for exist {
				select {
				case out, exist := <-outCh:
					if exist {
						manager.printOut(execute, out.StringData())
					}
				case err, exist := <-errCh:
					if exist {
						log.Println(err.StringData())
					}
				}
			}
		}(name, exec, outputCh, errorsCh)
	}

	defer func() {
		for _, outputCh := range outputsChs {
			close(outputCh)
		}
		for _, errorsCh := range errorsChs {
			close(errorsCh)
		}
	}()

	// Output log (initialize start)
	log.Println("Execute Initializing start....")

	// Initialize
	for name, exec := range manager.executes {
		go func(name string, execute *ExecuteItem) {
			exits <- exitCode{
				name: name,
				err:  execute.Exec.ExecuteInitialize(),
			}
		}(name, exec)
	}

	// Initialize wait
	for i := 0; i < nExecutes; i++ {
		if exit := <-exits; exit.err != nil {
			err = errors.WithMessage(exit.err, "Error occered when process initializing at "+exit.name+".")
		}
	}
	if err != nil {
		return
	}

	log.Println("Execute Initializing end....")
	log.Println("Execute Running Start....")

	// Running
	for name, exec := range manager.executes {
		go func(name string, execute *ExecuteItem) {
			exits <- exitCode{
				name: name,
				err:  execute.Exec.ExecuteWait(),
			}
		}(name, exec)
	}

	log.Println("Execute Running Starting end....")

	// Running wait
	for i := 0; i < nExecutes; i++ {
		if exit := <-exits; exit.err != nil {
			log.Println(exit.err)
			err = errors.WithMessage(exit.err, "Error occered when process initializing at "+exit.name+".")
		}
	}
	// if err != nil {
	// 	return
	// }

	log.Println("Execute Running end....")

	// Finalize Start

	log.Println("Finalizing....")
	for name, exec := range manager.executes {
		go func(name string, execute *ExecuteItem) {
			exits <- exitCode{
				name: name,
				err:  execute.Exec.ExecuteFinalize(),
			}
		}(name, exec)
	}
	for i := 0; i < nExecutes; i++ {
		if exit := <-exits; exit.err != nil {
			err = errors.WithMessage(exit.err, "Error occered when process finalizing at "+exit.name+".")
		}
	}

	return
}
