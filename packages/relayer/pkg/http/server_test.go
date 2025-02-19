package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/taikoxyz/taiko-mono/packages/relayer"
	"github.com/taikoxyz/taiko-mono/packages/relayer/pkg/mock"
	"github.com/taikoxyz/taiko-mono/packages/relayer/pkg/repo"
)

func newTestServer(url string) *Server {
	_ = godotenv.Load("../.test.env")

	srv := &Server{
		echo:            echo.New(),
		eventRepo:       mock.NewEventRepository(),
		suspendedTxRepo: mock.NewSuspendedTransactionRepository(),
	}

	srv.configureMiddleware([]string{"*"})
	srv.configureRoutes()

	return srv
}

func Test_NewServer(t *testing.T) {
	tests := []struct {
		name    string
		opts    NewServerOpts
		wantErr error
	}{
		{
			"success",
			NewServerOpts{
				Echo:            echo.New(),
				EventRepo:       &repo.EventRepository{},
				SuspendedTxRepo: &repo.SuspendedTransactionRepository{},
				CorsOrigins:     make([]string, 0),
				SrcEthClient:    &mock.EthClient{},
				DestEthClient:   &mock.EthClient{},
				BlockRepo:       &mock.BlockRepository{},
			},
			nil,
		},
		{
			"noSrcEthClient",
			NewServerOpts{
				Echo:            echo.New(),
				EventRepo:       &repo.EventRepository{},
				SuspendedTxRepo: &repo.SuspendedTransactionRepository{},
				CorsOrigins:     make([]string, 0),
				DestEthClient:   &mock.EthClient{},
				BlockRepo:       &mock.BlockRepository{},
			},
			relayer.ErrNoEthClient,
		},
		{
			"noDestEthClient",
			NewServerOpts{
				Echo:            echo.New(),
				EventRepo:       &repo.EventRepository{},
				SuspendedTxRepo: &repo.SuspendedTransactionRepository{},
				CorsOrigins:     make([]string, 0),
				SrcEthClient:    &mock.EthClient{},
				BlockRepo:       &mock.BlockRepository{},
			},
			relayer.ErrNoEthClient,
		},
		{
			"noBlockRepo",
			NewServerOpts{
				Echo:            echo.New(),
				EventRepo:       &repo.EventRepository{},
				SuspendedTxRepo: &repo.SuspendedTransactionRepository{},
				CorsOrigins:     make([]string, 0),
				SrcEthClient:    &mock.EthClient{},
				DestEthClient:   &mock.EthClient{},
			},
			relayer.ErrNoBlockRepository,
		},
		{
			"noEventRepo",
			NewServerOpts{
				Echo:            echo.New(),
				CorsOrigins:     make([]string, 0),
				SrcEthClient:    &mock.EthClient{},
				DestEthClient:   &mock.EthClient{},
				BlockRepo:       &mock.BlockRepository{},
				SuspendedTxRepo: &repo.SuspendedTransactionRepository{},
			},
			relayer.ErrNoEventRepository,
		},
		{
			"noCorsOrigins",
			NewServerOpts{
				Echo:            echo.New(),
				EventRepo:       &repo.EventRepository{},
				SrcEthClient:    &mock.EthClient{},
				DestEthClient:   &mock.EthClient{},
				BlockRepo:       &mock.BlockRepository{},
				SuspendedTxRepo: &repo.SuspendedTransactionRepository{},
			},
			relayer.ErrNoCORSOrigins,
		},
		{
			"noHttpFramework",
			NewServerOpts{
				EventRepo:       &repo.EventRepository{},
				CorsOrigins:     make([]string, 0),
				SrcEthClient:    &mock.EthClient{},
				DestEthClient:   &mock.EthClient{},
				BlockRepo:       &mock.BlockRepository{},
				SuspendedTxRepo: &repo.SuspendedTransactionRepository{},
			},
			ErrNoHTTPFramework,
		},
	}

	for _, tt := range tests {
		_, err := NewServer(tt.opts)
		assert.Equal(t, tt.wantErr, err)
	}
}

func Test_Health(t *testing.T) {
	srv := newTestServer("")

	req, _ := http.NewRequest(echo.GET, "/healthz", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Test_Health expected code %v, got %v", http.StatusOK, rec.Code)
	}
}

func Test_Root(t *testing.T) {
	srv := newTestServer("")

	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Test_Root expected code %v, got %v", http.StatusOK, rec.Code)
	}
}

func Test_StartShutdown(t *testing.T) {
	srv := newTestServer("")

	go func() {
		_ = srv.Start(":3928")
	}()
	assert.Nil(t, srv.Shutdown(context.Background()))
}
