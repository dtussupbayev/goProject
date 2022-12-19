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

type ProductResource struct {
	store store.Store
}

func NewProductResource(store store.Store) *ProductResource {
	return &ProductResource{store: store}
}
func (p *ProductResource) Routes() chi.Router {
	r := chi.NewRouter()

	//r.Get("/{id}", p.ByCategoryID)
	r.Group(func(r chi.Router) {
		//r.Use(auth)
		r.Post("/", p.CreateProduct)
		r.Put("/", p.UpdateProduct)
		r.Delete("/{id}", p.DeleteProduct)
		r.Get("/", p.GetAllProducts)
		r.Get("/{id}", p.ById)

	})

	return r
}

func (p *ProductResource) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(models.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) //invalid json
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := p.store.Products().Create(r.Context(), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError) //error with db
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func (p *ProductResource) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ProductsFilter{}

	if searchQuery := queryValues.Get("query"); searchQuery != "" {
		filter.Query = &searchQuery
	}
	products, err := p.store.Products().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, products)
}

func (p *ProductResource) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	product, err := p.store.Products().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, product)
}
func (p *ProductResource) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(models.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		product,
		validation.Field(&product.ID, validation.Required),
		validation.Field(&product.Name, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := p.store.Products().Update(r.Context(), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}
func (p *ProductResource) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := p.store.Products().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}

//func (p *ProductResource) ByCategoryID(w http.ResponseWriter, r *http.Request) {
//	idStr := chi.URLParam(r, "id")
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprintf(w, "Unknown err: %v", err)
//		return
//	}
//
//	product, err := p.store.Products().ByID(r.Context(), id)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(w, "DB err: %v", err)
//		return
//	}
//
//	render.JSON(w, r, product)
//}
