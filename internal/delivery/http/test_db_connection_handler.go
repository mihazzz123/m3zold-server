package http

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mihazzz123/m3zold-server/internal/usecase/health"
)

type TestDBConnectionHandler struct {
	TestDBConnectionUC *health.TestDBConnectionUseCase
}

func NewTestDBConnectionHandler(testDBConnectionUC *health.TestDBConnectionUseCase) *TestDBConnectionHandler {
	return &TestDBConnectionHandler{TestDBConnectionUC: testDBConnectionUC}
}

// TestDBConnection использует инфраструктурный слой
func (h *TestDBConnectionHandler) TestDBConnection(dbURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	return h.TestDBConnectionUC.TestDBConnection(ctx)
}
