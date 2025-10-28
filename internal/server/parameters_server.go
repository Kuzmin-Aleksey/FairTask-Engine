package server

import (
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/service/parameters_service"
	"FairTask_Engine/pkg/failure"
	"encoding/json"
	"net/http"
	"strconv"
)

type ParametersServer struct {
	parameters *parameters_service.ParametersService
}

func NewParametersServer(parameters *parameters_service.ParametersService) *ParametersServer {
	return &ParametersServer{
		parameters: parameters,
	}
}

func (s *ParametersServer) ApiHandleCreateParameter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param := new(entity.Parameter)

	if err := json.NewDecoder(r.Body).Decode(param); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.parameters.CreateParameter(ctx, param); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, IdResponse{param.Id}, http.StatusOK)
}

func (s *ParametersServer) ApiHandleGetParameters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := s.parameters.GetAll(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, params, http.StatusOK)
}

func (s *ParametersServer) ApiHandleDeleteParameter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
	}

	if err := s.parameters.DeleteParameter(ctx, id); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}
}
