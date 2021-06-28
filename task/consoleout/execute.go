package consoleout

const (
	Successfully string = "Successfully"
)

func (t *Task) Initialize() error {
	t.kill = make(chan string)
	t.progress = make(chan string)

	go t.execute()
	<-t.waiter
	if t.finalized {
		close(t.kill)
		close(t.progress)
	}
	return nil
}

func (t *Task) Wait() error {
	if !t.finalized {
		t.progress <- ""
	}
	return nil
}

func (t *Task) Finalize() error {
	if !t.finalized {
		close(t.kill)
	}
	close(t.progress)
	close(t.in)
	return nil
}

func (t *Task) Kill() error {
	if !t.finalized {
		close(t.kill)
	}
	return nil
}

func (t *Task) execute() {

	// internal func
	CheckKill := func() bool {
		select {
		case <-t.kill:
			return true
		default:
			return false
		}
	}

	// check kill
	if exit := CheckKill(); exit {
		t.finalized = true
		return
	}

	// call initialize ended
	t.waiter <- Successfully

	// wait for next
	for isContinue := true; isContinue; {
		select {
		case <-t.kill:
			t.finalized = true
			return
		case <-t.progress:
			isContinue = false
			break
		}
	}

	for isContinue := true; isContinue; {
		select {
		case <-t.kill:
			isContinue = true
			break
		case input := <-t.in:
			t.PacketLogger.PacketLog(input)
		}
	}
}
