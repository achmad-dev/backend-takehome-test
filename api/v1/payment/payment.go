package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/kitabisa/backend-takehome-test/internal/config"
	"github.com/rs/zerolog/log"
)

func Routes( /* any dependency injection comes here*/ ) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}", GetPaymentMethodByIdHandler)
	r.Post("/create", CreateNewPaymentMethodHandler)
	return r
}

type ResponseBodyResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func CreateNewPaymentMethodHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	var resp ResponseBodyResult
	var msg string
	//post payment method
	var newPaymentMethod PaymentMethods
	err := json.NewDecoder(r.Body).Decode(&newPaymentMethod)
	if err != nil {
		resp.Code = "BE-001"
		msg = "can't decode request"
	}
	var dbConn config.DbConnection
	dbConn.GetDbConnectionPool().AddTableWithName(PaymentMethods{}, "payment_methods").SetKeys(true, "id")
	err = dbConn.GetDbConnectionPool().Insert(&newPaymentMethod)
	if err != nil {
		resp.Code = "BE-001"
		msg = err.Error()
	} else {
		resp.Code = "BE-000"
		msg = "payment method created"
	}
	// this is dummy response, always OK. You need to implement the proper way
	resp.Message = msg
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Error().Err(err).Msg("Error while marshalling response")
	}
	rw.Write(jsonResp)
}

type PaymentMethods struct {
	Id   uint64 `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name"`
}

type PaymentMethodResponse struct {
	ResponseBodyResult
	PaymentMethod []interface{} `json:"data"`
}

func GetPaymentMethodByIdHandler(rw http.ResponseWriter, r *http.Request) {

	var resp PaymentMethodResponse

	idReq := chi.URLParam(r, "id")
	idReqInt, err := strconv.ParseInt(idReq, 0, 64)
	if err != nil {
		resp.Code = "BE-001"
		resp.Message = "payment method id must be a number"
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(resp)
		return
	}

	var msg string
	var httpStatus int

	var dbConn config.DbConnection
	dbConn.GetDbConnectionPool().AddTableWithName(PaymentMethods{}, "payment_methods").SetKeys(true, "id")
	paymentMethodResult, err := dbConn.GetDbConnectionPool().Get(PaymentMethods{}, idReqInt)
	if err != nil {
		resp.Code = "BE-001"

		msg = fmt.Sprintf("failed to retrieve payment data with id %s | %v", idReq, err)
		httpStatus = http.StatusInternalServerError
	} else if paymentMethodResult == nil {
		resp.Code = "BE-002"

		msg = fmt.Sprintf("payment method with id %s could not be found", idReq)
		httpStatus = http.StatusNotFound
	} else {
		resp.Code = "BE-000"
		resp.PaymentMethod = append(resp.PaymentMethod, buildPaymentResponse(paymentMethodResult))

		msg = fmt.Sprintf("payment method id %s successfully retrieved", idReq)
		httpStatus = http.StatusOK
	}

	log.Info().Msgf(msg)
	resp.Message = msg
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(httpStatus)
	json.NewEncoder(rw).Encode(resp)
}

type PaymentMethodData struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func buildPaymentResponse(paymentMethodResult interface{}) (resp PaymentMethodData) {
	result := paymentMethodResult.(*PaymentMethods)
	resp.Id = result.Id
	resp.Name = result.Name
	return resp
}
