package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"nyeh-back/internal/core"
	"nyeh-back/internal/domain"
)

// HTTP DTOs
type ContactDTO struct {
	Email     string `json:"email"`
	Github    string `json:"github"`
	Linkedin  string `json:"linkedin"`
	Instagram string `json:"instagram"`
}

type BioRequest struct {
	Name    string     `json:"name"`
	AboutMe string     `json:"about_me"`
	Address string     `json:"address"`
	Contact ContactDTO `json:"contact"`
}

type BioResponse struct {
	Data *BioRequest `json:"data"`
}

type BioHandler struct {
	bioUsecase domain.BioUsecase
}

func NewBioHandler(u domain.BioUsecase) *BioHandler {
	return &BioHandler{bioUsecase: u}
}

// GetBioHandler godoc
//
//	@Summary	Get Bio
//	@Accept		json
//	@Produce	json
//	@Router		/bio [get]
func (h *BioHandler) GetBioHandler(w http.ResponseWriter, r *http.Request) {
	bio, err := h.bioUsecase.Get(r.Context(), core.Settings.ALLOWED_EMAIL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bio == nil {
		http.Error(w, "Bio not found", http.StatusNotFound)
		return
	}

	resp := BioResponse{
		Data: &BioRequest{
			Name:    bio.Name,
			AboutMe: bio.AboutMe,
			Address: bio.Address,
			Contact: ContactDTO{
				Email:     bio.Contact.Email,
				Github:    bio.Contact.Github,
				Linkedin:  bio.Contact.Linkedin,
				Instagram: bio.Contact.Instagram,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpsertBioHandler godoc
//
//	@Summary	Upsert Bio (create and update)
//	@Accept		json
//	@Produce	json
//	@Param		request	body	BioRequest	true	"Bio update payload"
//	@Router		/bio [post]
func (h *BioHandler) UpsertBioHandler(w http.ResponseWriter, r *http.Request) {

	claimEmail, _ := r.Context().Value(core.JWT_CLAIM_USER_EMAIL_KEY).(string)

	var req BioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	bio := &domain.Bio{
		Name:    req.Name,
		AboutMe: req.AboutMe,
		Address: req.Address,
		Contact: domain.Contact{
			Email:     req.Contact.Email,
			Github:    req.Contact.Github,
			Instagram: req.Contact.Instagram,
			Linkedin:  req.Contact.Linkedin,
		},
	}

	if err := h.bioUsecase.Upsert(r.Context(), claimEmail, bio); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update bio %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// DeleteBioHandler godoc
//
//	@Summary	Delete Bio
//	@Router		/bio [delete]
func (h *BioHandler) DeleteBioHandler(w http.ResponseWriter, r *http.Request) {
	claimEmail, _ := r.Context().Value(core.JWT_CLAIM_USER_EMAIL_KEY).(string)

	if err := h.bioUsecase.Delete(r.Context(), claimEmail); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete bio %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
