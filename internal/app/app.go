package app

import (
	"context"

	"firebase.google.com/go"
	"github.com/core-go/health"
	"github.com/core-go/health/firestore"
	"github.com/core-go/log/zap"
	"google.golang.org/api/option"

	"go-service/internal/user"
)

type ApplicationContext struct {
	Health *health.Handler
	User   user.UserTransport
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	opts := option.WithCredentialsJSON([]byte(cfg.Credentials))
	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	logError := log.LogError

	userHandler, err := user.NewUserHandler(client, logError)
	if err != nil {
		return nil, err
	}

	firestoreChecker, err := firestore.NewHealthChecker(ctx, []byte(cfg.Credentials), cfg.ProjectId)
	if err != nil {
		return nil, err
	}
	healthHandler := health.NewHandler(firestoreChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
