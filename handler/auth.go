package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dtusupbaev/goProject/internal/models"
	"github.com/dtusupbaev/goProject/internal/store"
	"github.com/dtusupbaev/goProject/token"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"time"
)

const (
	accessTokenTTL  = 2 * time.Hour
	refreshTokenTTL = 168 * time.Hour
)

type AuthResource struct {
	store       store.Store
	cache       *lru.TwoQueueCache
	tokenManger token.TokenManger
}

func NewAuthResource(store store.Store, cache *lru.TwoQueueCache, tokenManager token.TokenManger) *AuthResource {
	return &AuthResource{
		store:       store,
		cache:       cache,
		tokenManger: tokenManager,
	}
}
func (a *AuthResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", a.LoginUser)

	return r
}
func (a *AuthResource) LoginUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.LoginRequest)

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	u, err := a.store.Users().FindByEmail(r.Context(), user.Email)
	if err != nil || !u.ComparePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Incorrect email or password")
		return
	}

	tokens, err := a.createSession(&models.AuthorizedUserInfo{
		Id: u.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	render.JSON(w, r, tokens)
}
func (a *AuthResource) createSession(userInfo *models.AuthorizedUserInfo) (*models.Tokens, error) {
	var res models.Tokens
	var err error

	if res.AccessToken, err = a.tokenManger.NewJWT(userInfo, accessTokenTTL); err != nil {
		return nil, err
	}
	res.RefreshToken, err = a.tokenManger.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	session := models.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
	}

	a.cache.Add(userInfo.Id, session)

	return &res, err
}
