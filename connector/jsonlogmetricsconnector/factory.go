package jsonlogmetricsconnector

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
)

var (
	Type = component.MustNewType("jsonlogmetricsconnector")
)

func NewFactory() connector.Factory {
	return connector.NewFactory(
		Type,
		createDefaultConfig,
		connector.WithLogsToMetrics(createLogsToMetricsConnector, component.StabilityLevelAlpha))
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createLogsToMetricsConnector(
	_ context.Context,
	set connector.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (connector.Logs, error) {
	return newMetricsConnector(set, cfg, nextConsumer), nil
}
