package workers

import "sync"

func FanIn(inputChs ...chan interface{}) chan interface{} {
	outCh := make(chan interface{})

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)

			go func(inputCh chan interface{}) {
				defer wg.Done()
				for item := range inputCh {
					outCh <- item
				}
			}(inputCh)
		}

		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func FanOut(inputCh chan interface{}, n int) []chan interface{} {
	chs := make([]chan interface{}, 0, n)
	for i := 0; i < n; i++ {
		ch := make(chan interface{})
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan interface{}) {
			for _, ch := range chs {
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}

			num, ok := <-inputCh
			if !ok {
				return
			}

			ch := chs[i]
			ch <- num
		}
	}()

	return chs
}
