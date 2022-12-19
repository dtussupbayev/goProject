package resource

import (
	"encoding/json"
	"fmt"
	"github.com/dtusupbaev/goProject/internal/message_broker"
	"github.com/dtusupbaev/goProject/internal/models"
	"github.com/dtusupbaev/goProject/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"strconv"
)

type CategoryResource struct {
	store  store.Store
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewCategoryResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *CategoryResource {
	return &CategoryResource{
		store:  store,
		broker: broker,
		cache:  cache,
	}

}
func (cr *CategoryResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", cr.GetAllCategories)
	r.Get("/{id}", cr.ById)

	r.Group(func(r chi.Router) {
		//r.Use(auth)
		r.Post("/", cr.CreateCategory)
		r.Put("/", cr.UpdateCategory)
		r.Delete("/{id}", cr.DeleteCategory)

	})
	return r
}

func (cr *CategoryResource) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := new(models.Category)
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) //invalid json
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Categories().Create(r.Context(), category); err != nil {
		w.WriteHeader(http.StatusInternalServerError) //error with db
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	//cr.cache.Purge()
	cr.broker.Cache().Purge()

	w.WriteHeader(http.StatusCreated)
}
func (cr *CategoryResource) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.CategoriesFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		categoriesFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, categoriesFromCache)
			return
		}
		filter.Query = &searchQuery
	}

	categories, err := cr.store.Categories().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	if searchQuery != "" {
		cr.cache.Add(searchQuery, categories)
	}
	render.JSON(w, r, categories)
}

func (cr *CategoryResource) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	categoryFromCache, ok := cr.cache.Get(id)
	if ok {
		render.JSON(w, r, categoryFromCache)
		return
	}

	category, err := cr.store.Categories().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, category)
}
func (cr *CategoryResource) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	category := new(models.Category)
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		category,
		validation.Field(&category.ID, validation.Required),
		validation.Field(&category.Name, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Categories().Update(r.Context(), category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(category.ID)
}
func (cr *CategoryResource) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Categories().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(id)
}
