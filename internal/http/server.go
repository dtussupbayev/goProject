package http

import (
	"context"
	"github.com/Assyl00/goProject/handler"
	"github.com/Assyl00/goProject/internal/http/resource"
	"github.com/Assyl00/goProject/internal/message_broker"
	"github.com/Assyl00/goProject/internal/store"
	"github.com/Assyl00/goProject/token"
	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ctx          context.Context
	idleConnsCh  chan struct{}
	store        store.Store
	cache        *lru.TwoQueueCache
	Address      string
	broker       message_broker.MessageBroker
	tokenManager token.TokenManger
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	categoriesResource := resource.NewCategoryResource(s.store, s.broker, s.cache)
	r.Mount("/categories", categoriesResource.Routes())

	orderResource := resource.NewOrderResource(s.store)
	r.Mount("/orders", orderResource.Routes(s.userIdentity))

	productResource := resource.NewProductResource(s.store)
	r.Mount("/products", productResource.Routes())

	reviewResource := resource.NewReviewResource(s.store)
	r.Mount("/reviews", reviewResource.Routes(s.userIdentity))

	usersResource := resource.NewUsersResource(s.store, s.broker, s.cache)
	r.Mount("/users", usersResource.Routes(s.userIdentity))

	// Authentication
	authResource := handler.NewAuthResource(s.store, s.cache, s.tokenManager)
	r.Mount("/auth", authResource.Routes())

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // блокируемся, пока контекст приложения не отменен

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	// блок до записи или закрытия канала
	<-s.idleConnsCh
}
