package middlewares

import "github.com/sadlil/gologger"

/*
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] %s %s %d %s \n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC822),
			params.Method,
			params.Path,
			params.StatusCode,
			params.Latency,
		)
	})
}*/

func Logger() gologger.GoLogger {
	logger := gologger.GetLogger(gologger.CONSOLE, gologger.ColoredLog)
	return logger
}
