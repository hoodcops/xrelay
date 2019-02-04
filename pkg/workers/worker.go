package workers

// Worker represents any object that listens for message on
// a queue and acts on them
type Worker interface {
	Run()
}

// VerificationWorker consumes messages on the verification queue
// and sends verification codes to msisds to verify their authenticity
type VerificationWorker struct {
}

// Run ...
func (vw VerificationWorker) Run() {

}
