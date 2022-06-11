package gatherer

import "go.uber.org/zap"

var zlog *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	zlog, _ = config.Build()
	// zlog, _ = zap.NewProduction()
}
func SetLogger(logger *zap.Logger) {
	zlog = logger
}
