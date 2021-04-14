package lfc

func FastRecoverGoroutineFunc(f func()) {
	go func() {
		defer recover()
		f()
	}()
}
