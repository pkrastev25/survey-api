


func WithContext(operation func(ctx context) interface{}) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return operation(ctx)
}

func WithContextError(
	operation func(ctx context) (interface{}, error)
	) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return operation(ctx)
}