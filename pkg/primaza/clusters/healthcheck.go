package clusters

import (
	"context"
	"fmt"
	"time"

	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ErrProbeNotFound = fmt.Errorf("probe not found")

type HealthcheckResult string

type Healthchecker interface {
	AddProbe(primazaiov1alpha1.ClusterEnvironment)
	RemoveProbe(string) error
}

func NewHealthchecker(ctx context.Context, cli client.Client) Healthchecker {
	return &healthchecker{
		ctx:    ctx,
		cli:    cli,
		checks: map[string]func(){},
	}
}

type healthchecker struct {
	ctx context.Context
	cli client.Client

	checks map[string]func()
}

func (h *healthchecker) AddProbe(ce primazaiov1alpha1.ClusterEnvironment) {
	ctx, cancel := context.WithCancel(h.ctx)
	c := newHealthcheck(ctx, ce)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(120 * time.Second):
				c.Probe()
			}
		}
	}()

	h.checks[ce.Name] = cancel
}

func (h *healthchecker) RemoveProbe(ceName string) error {
	cancel, ok := h.checks[ceName]
	if !ok {
		return fmt.Errorf("error removing probe '%s': %w", ceName, ErrProbeNotFound)
	}

	cancel()
	return nil
}

type Healthcheck interface {
	Probe()
}

func newHealthcheck(context.Context, primazaiov1alpha1.ClusterEnvironment) Healthcheck {
	panic("To be implemented")
}
