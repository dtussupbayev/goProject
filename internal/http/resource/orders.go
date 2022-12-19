package resource

import (
	"encoding/json"
	"fmt"
	"github.com/dtusupbaev/goProject/internal/models"
	"github.com/dtusupbaev/goProject/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type OrderResource struct {
	store store.Store
}

func NewOrderResource(store store.Store) *OrderResource {
	return &OrderResource{store: store}
}
func (o *OrderResource) Routes(auth func(handler http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", o.GetAllOrders)
	r.Get("/{id}", o.ByIdOrder)

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Post("/", o.CreateOrder)
		r.Delete("/{id}", o.DeleteOrder)

	})

	return r
}

func (o *OrderResource) CreateOrder(w http.ResponseWriter, r *http.Request) {
	order := new(models.Order)
	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := o.store.Orders().Create(r.Context(), order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func (o *OrderResource) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := o.store.Orders().All(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	render.JSON(w, r, orders)
}

func (o *OrderResource) ByIdOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	order, err := o.store.Orders().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, order)
}

func (o *OrderResource) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := o.store.Orders().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}
