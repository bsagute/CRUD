package logger

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
)

func ConfigureSentry() error {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              viper.GetString("SENTRY_DSN"),
		EnableTracing:    viper.GetBool("SENTRY_ENABLE_TRACING"),
		TracesSampleRate: viper.GetFloat64("SENTRY_TRACES_SAMPLE_RATE"),
		SampleRate:       viper.GetFloat64("SENTRY_SAMPLE_RATE"),
		AttachStacktrace: true,
		Environment:      os.Getenv("APP_ENV"),
	}); err != nil {
		return err
	}
	return nil
}

func NotifySentry(err error) {
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureException(err)
}

func FlushSentry() {
	defer sentry.Flush(2 * time.Second)
}
