package jsonlogmetricsconnector

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type connectorMetrics struct {
	config          Config
	metricsConsumer consumer.Metrics
	logger          *zap.Logger

	component.StartFunc
	component.ShutdownFunc
}

// newMetricsConnector is a function to create a new connector for metrics extraction
func newMetricsConnector(set connector.Settings, config component.Config, metricsConsumer consumer.Metrics) *connectorMetrics {
	set.TelemetrySettings.Logger.Info("Building jsonlogmetricsconnector connector for metrics")
	cfg := config.(*Config)

	return &connectorMetrics{
		config:          *cfg,
		logger:          set.TelemetrySettings.Logger,
		metricsConsumer: metricsConsumer,
	}
}

// Capabilities implements the consumer interface.
func (c *connectorMetrics) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeLogs method is called for each instance of a log sent to the connector
func (c *connectorMetrics) ConsumeLogs(ctx context.Context, pl plog.Logs) error {
	// // loop through the levels of logs
	metricsUnmarshaler := &pmetric.JSONUnmarshaler{}
	c.logger.Info("Total resource logs inside:", zap.Int("length", pl.ResourceLogs().Len()))
	for i := 0; i < pl.ResourceLogs().Len(); i++ {
		li := pl.ResourceLogs().At(i)
		for j := 0; j < li.ScopeLogs().Len(); j++ {
			logRecord := li.ScopeLogs().At(j)

			for k := 0; k < logRecord.LogRecords().Len(); k++ {
				lRecord := logRecord.LogRecords().At(k)
				token := lRecord.Body()
				var m pmetric.Metrics
				m, err := metricsUnmarshaler.UnmarshalMetrics([]byte(token.AsString()))
				if err != nil {
					c.logger.Error("could extract metrics from json", zap.Error(err))
					continue
				}
				err = c.metricsConsumer.ConsumeMetrics(ctx, m)
				if err != nil {
					c.logger.Error("could not consume metrics from json", zap.Error(err))
				}
			}
		}
	}
	return nil
}
