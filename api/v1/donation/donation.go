package donation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kitabisa/backend-takehome-test/internal/config"
	"github.com/rs/zerolog/log"
)

func Routes( /* any dependency injection comes here*/ ) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}", GetDonationByIdHandler)
	r.Post("/create", CreateDonationHandler)
	return r
}

type ResponseBodyResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Donations struct {
	Id              *uint64 `db:"id" json:"id,omitempty"`
	PaymentMethodId uint64  `db:"payment_method_id" json:"payment_method_id"`
	CampaignId      uint64  `db:"campaign_id" json:"campaign_id"`
	Amount          int64   `db:"amount" json:"amount"`
}

type DonationResponse struct {
	ResponseBodyResult
	Donation []interface{} `json:"data"`
}

type DonationData struct {
	Id              uint64 `json:"id,omitempty"`
	PaymentMethodId uint64 `json:"payment_method_id"`
	CampaignId      uint64 `json:"campaign_id"`
	Amount          int64  `json:"amount"`
}

type Campaigns struct {
	Id    uint64 `db:"id"`
	Title string `db:"title"`
}
type PaymentMethods struct {
	Id   uint64 `db:"id"`
	Name string `db:"name"`
}

func GetCampaignById(id uint64) (interface{}, error) {
	var dbConn config.DbConnection
	dbConn.GetDbConnectionPool().AddTableWithName(Campaigns{}, "campaigns").SetKeys(true, "id")
	return dbConn.GetDbConnectionPool().Get(Campaigns{}, id)
}

func GetPaymentMethodByID(id uint64) (interface{}, error) {
	var dbConn config.DbConnection
	dbConn.GetDbConnectionPool().AddTableWithName(PaymentMethods{}, "payment_methods").SetKeys(true, "id")
	return dbConn.GetDbConnectionPool().Get(PaymentMethods{}, id)
}

func CreateDonationHandler(rw http.ResponseWriter, r *http.Request) {
	// TODO: implement this
	rw.Header().Set("Content-Type", "application/json")
	var resp ResponseBodyResult
	var msg string
	var httpStatus int
	//post donation method
	var newDonation Donations
	err := json.NewDecoder(r.Body).Decode(&newDonation)
	if err != nil {
		resp.Code = "BE-001"
		httpStatus = http.StatusBadRequest
		msg = "can't decode request"
	}
	//validate is campaign and pm exist
	paymentMethod, err := GetPaymentMethodByID(newDonation.PaymentMethodId)
	if err != nil || paymentMethod == nil {
		resp.Code = "BE-001"
		httpStatus = http.StatusBadRequest
		msg = fmt.Sprintf("failed to create donation because payment method id %v does not exist", newDonation.PaymentMethodId)
	}
	campaign, err := GetCampaignById(newDonation.CampaignId)
	if err != nil || campaign == nil {
		resp.Code = "BE-001"
		httpStatus = http.StatusBadRequest
		msg = fmt.Sprintf("failed to create donation because campaign id %v does not exist", newDonation.CampaignId)
	}
	if newDonation.Amount < 10000 {
		resp.Code = "BE-001"
		httpStatus = http.StatusBadRequest
		msg = "failed to create donation because amount is less than 10000"
	}
	if paymentMethod != nil && campaign != nil && newDonation.Amount >= 10000 {
		var dbConn config.DbConnection
		dbConn.GetDbConnectionPool().AddTableWithName(Donations{}, "donations").SetKeys(true, "id")
		err = dbConn.GetDbConnectionPool().Insert(&newDonation)
		if err != nil {
			resp.Code = "BE-002"
			httpStatus = http.StatusInternalServerError
			msg = fmt.Sprintf("can't create donation because %s", err.Error())
		} else {
			resp.Code = "BE-000"
			httpStatus = http.StatusCreated
			msg = "donation created successfully"
		}
	}
	resp.Message = msg
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Error().Err(err).Msg("Error while marshalling response")
	}
	rw.WriteHeader(httpStatus)
	rw.Write(jsonResp)
}

func GetDonationByIdHandler(rw http.ResponseWriter, r *http.Request) {
	// TODO: implement this
}
