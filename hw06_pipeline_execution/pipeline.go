package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func createChannel(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			default:
				out <- v
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := createChannel(in, done)
	for _, stage := range stages {
		out = stage(out)
		out = createChannel(out, done)
	}

	return out
}
