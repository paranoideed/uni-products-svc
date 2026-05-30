package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var statusSuccess = attribute.String("status", "success")
var statusError = attribute.String("status", "error")

type Metrics struct {
	created metric.Int64Counter
	deleted metric.Int64Counter
	active  metric.Int64UpDownCounter
}

func New() (*Metrics, error) {
	meter := otel.Meter("uni-products-svc")

	created, err := meter.Int64Counter("products_created_total",
		metric.WithDescription("Total number of created products"),
	)
	if err != nil {
		return nil, err
	}

	deleted, err := meter.Int64Counter("products_deleted_total",
		metric.WithDescription("Total number of deleted products"),
	)
	if err != nil {
		return nil, err
	}

	active, err := meter.Int64UpDownCounter("products_active",
		metric.WithDescription("Current number of active products"),
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		created: created,
		deleted: deleted,
		active:  active,
	}, nil
}

func (m *Metrics) RecordCreated(ctx context.Context, err error) {
	if err != nil {
		m.created.Add(ctx, 1, metric.WithAttributes(statusError))
		return
	}
	m.created.Add(ctx, 1, metric.WithAttributes(statusSuccess))
	m.active.Add(ctx, 1)
}

func (m *Metrics) RecordDeleted(ctx context.Context, err error) {
	if err != nil {
		m.deleted.Add(ctx, 1, metric.WithAttributes(statusError))
		return
	}
	m.deleted.Add(ctx, 1, metric.WithAttributes(statusSuccess))
	m.active.Add(ctx, -1)
}
