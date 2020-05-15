package nologr

func ExampleLogger() {
	logger := New()
	logger.Info("test log", "hello", "world")
	// Output:
}
