package handlers

import (
	"encoding/json"
	"mvp-2-spms/services/manage-universities/inputdata"
	"mvp-2-spms/web_server/handlers/interfaces"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UniversityHandler struct {
	uniInteractor interfaces.IUniversityInteractor
}

func InitUniversityHandler(uInteractor interfaces.IUniversityInteractor) UniversityHandler {
	return UniversityHandler{
		uniInteractor: uInteractor,
	}
}

func (h *UniversityHandler) GetAllUniEdProgrammes(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	uniId, err := strconv.ParseUint(chi.URLParam(r, "uniID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	input := inputdata.GetUniEducationalProgrammes{
		ProfessorId:  uint(id),
		UniversityId: uint(uniId),
	}
	result, err := h.uniInteractor.GetUniEdProgrammes(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
