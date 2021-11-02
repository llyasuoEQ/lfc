package lfc

func FastRecoverGoroutineFunc(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// TODO
			}
		}()
		f()
	}()
}
