package v0

func (e *RehearsalEngine) Finalize() {
	if e.logger != nil {
		e.logger.Finalize()
	}
}
