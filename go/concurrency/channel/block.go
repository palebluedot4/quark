package channel

func BlockOnEmptySelect() {
	select {}
}

func BlockOnNilChannel() {
	var ch chan struct{}
	<-ch
}
