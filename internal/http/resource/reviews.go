package resource

import (
	"encoding/json"
	"fmt"
	"github.com/Assyl00/goProject/internal/models"
	"github.com/Assyl00/goProject/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"
)

type ReviewResource struct {
	store store.Store
}

func NewReviewResource(store store.Store) *ReviewResource {
	return &ReviewResource{store: store}
}
func (p *ReviewResource) Routes(auth func(handler http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	//r.Get("/{id}", p.ByCategoryID)
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Post("/", p.CreateReview)
		r.Put("/", p.UpdateReview)
		r.Delete("/{id}", p.DeleteReview)
		r.Get("/", p.GetAllReviews)
		r.Get("/{id}", p.ById)

	})

	return r
}

func (p *ReviewResource) CreateReview(w http.ResponseWriter, r *http.Request) {
	Review := new(models.Review)
	if err := json.NewDecoder(r.Body).Decode(Review); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) //invalid json
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := p.store.Reviews().Create(r.Context(), Review); err != nil {
		w.WriteHeader(http.StatusInternalServerError) //error with db
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func (p *ReviewResource) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ReviewsFilter{}

	if searchQuery := queryValues.Get("query"); searchQuery != "" {
		filter.Query = &searchQuery
	}
	Reviews, err := p.store.Reviews().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, Reviews)
}

func (p *ReviewResource) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	Review, err := p.store.Reviews().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, Review)
}
func (p *ReviewResource) UpdateReview(w http.ResponseWriter, r *http.Request) {
	Review := new(models.Review)
	if err := json.NewDecoder(r.Body).Decode(Review); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		Review,
		validation.Field(&Review.ID, validation.Required),
		validation.Field(&Review.Body, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := p.store.Reviews().Update(r.Context(), Review); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}
func (p *ReviewResource) DeleteReview(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := p.store.Reviews().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}
