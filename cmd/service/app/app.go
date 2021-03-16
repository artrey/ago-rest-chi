package app

import (
	"encoding/json"
	"github.com/artrey/ago-rest-chi/cmd/service/rest"
	"github.com/artrey/ago-rest-chi/pkg/offers"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	"net/http"
)

type Server struct {
	offersSvc *offers.Service
	router    chi.Router
}

func NewServer(offersSvc *offers.Service, router chi.Router) *Server {
	return &Server{offersSvc: offersSvc, router: router}
}

func (s *Server) Init() error {
	s.router.Get("/offers", s.handleGetOffers)
	s.router.Get("/offers/{id}", s.handleGetOfferByID)
	s.router.Post("/offers", s.handleSaveOffer)
	s.router.Delete("/offers/{id}", s.handleRemoveOfferByID)

	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleGetOffers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.offersSvc.All(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_ = rest.WriteAsJson(writer, items)
}

func (s *Server) handleGetOfferByID(writer http.ResponseWriter, request *http.Request) {
	id, err := rest.ExtractID(request)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.offersSvc.ByID(request.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	_ = rest.WriteAsJson(writer, item)
}

func (s *Server) handleSaveOffer(writer http.ResponseWriter, request *http.Request) {
	itemToSave := &offers.Offer{}
	err := json.NewDecoder(request.Body).Decode(&itemToSave)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.offersSvc.Save(request.Context(), itemToSave)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_ = rest.WriteAsJson(writer, item)
}

func (s *Server) handleRemoveOfferByID(writer http.ResponseWriter, request *http.Request) {
	id, err := rest.ExtractID(request)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	deletedItem, err := s.offersSvc.DeleteByID(request.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	_ = rest.WriteAsJson(writer, deletedItem)
}
