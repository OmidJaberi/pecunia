package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"github.com/OmidJaberi/pecunia/internal/db"
	"github.com/OmidJaberi/pecunia/internal/domain"
)

type API struct {
	DB           *sqlx.DB
	UserRepo     *db.UserRepo
	AssetRepo    *db.AssetRepo
	CurrencyRepo *db.CurrencyRepo
	RateRepo     *db.ExchangeRateRepo
}

func NewAPI(database *sqlx.DB) *API {
	return &API{
		DB:           database,
		UserRepo:     db.NewUserRepo(database),
		AssetRepo:    db.NewAssetRepo(database),
		CurrencyRepo: db.NewCurrencyRepo(database),
		RateRepo:     db.NewExchangeRateRepo(database),
	}
}

// POST /users
type createUserReq struct {
	Name string `json:"name"`
}
type createUserResp struct {
	ID string `json:"id"`
}

func (api *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	u := &domain.User{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}
	if err := api.UserRepo.Create(u); err != nil {
		http.Error(w, "failed create user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createUserResp{ID: u.ID.String()})
}

// GET /users/{id}/assets
func (api *API) ListAssets(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	as, err := api.AssetRepo.ListByUserID(uid)
	if err != nil {
		http.Error(w, "failed to list assets", http.StatusInternalServerError)
		return
	}
	// convert domain.Asset -> response-friendly struct
	type assetResp struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Currency  string `json:"currency"`
		Amount    string `json:"amount"`
		Category  string `json:"category"`
		CreatedAt int64  `json:"created_at"`
	}
	out := make([]assetResp, 0, len(as))
	for _, a := range as {
		out = append(out, assetResp{
			ID:        a.ID.String(),
			Name:      a.Name,
			Currency:  a.Value.Currency.Code,
			Amount:    a.Value.Amount.String(),
			Category:  a.Category,
			CreatedAt: a.CreatedAt.Unix(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// POST /users/{id}/assets
type createAssetReq struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"` // decimal as string
	Category string `json:"category"`
}

func (api *API) CreateAsset(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	var req createAssetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	amt, err := decimal.NewFromString(req.Amount)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}
	// Todo: FindCurrency Function
	cur := &domain.Currency{
		Code: req.Currency,
	}
	val := &domain.Money{
		Amount:   amt,
		Currency: *cur,
	}
	a := &domain.Asset{
		ID:        uuid.New(),
		UserID:    uid,
		Name:      req.Name,
		Value:     *val,
		Category:  req.Category,
		CreatedAt: time.Now(),
	}
	if err := api.AssetRepo.Insert(*a); err != nil {
		http.Error(w, "failed to create asset", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": a.ID.String()})
}
