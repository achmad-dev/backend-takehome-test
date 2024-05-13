package campaign

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
	r.Get("/{id}", GetCampaignByIdHandler)
	r.Post("/create", CreateNewCampaignHandler)
	return r
}

type ResponseBodyResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func CreateNewCampaignHandler(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	//create campaign
	var resp ResponseBodyResult
	var msg string
	//post payment method
	var newCampaign Campaigns
	err := json.NewDecoder(r.Body).Decode(&newCampaign)
	if err != nil {
		resp.Code = "BE-001"
		msg = "can't decode request"
	}
	var dbConn config.DbConnection
	dbConn.GetDbConnectionPool().AddTableWithName(Campaigns{}, "campaigns").SetKeys(true, "id")
	err = dbConn.GetDbConnectionPool().Insert(&newCampaign)
	if err != nil {
		resp.Code = "BE-002"
		msg = fmt.Sprintf("can't create campaign because %s", err.Error())
	} else {
		resp.Code = "BE-000"
		msg = "campaign created"
	}
	// this is dummy response, always OK. You need to implement the proper way
	resp.Message = msg
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Error().Err(err).Msg("Error while marshalling response")
	}

	rw.Write(jsonResp)
}

type Campaigns struct {
	Id    *uint64 `db:"id" json:"id,omitempty"`
	Title string  `db:"title" json:"title"`
}

type CampaignResponse struct {
	ResponseBodyResult
	Campaign []interface{} `json:"data"`
}

func GetCampaignByIdHandler(rw http.ResponseWriter, r *http.Request) {
	var resp CampaignResponse

	idReq := chi.URLParam(r, "id")
	idReqInt, err := strconv.ParseInt(idReq, 0, 64)
	if err != nil {
		resp.Code = "BE-001"
		resp.Message = "campaign id must be a number"
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(resp)
		return
	}

	var msg string
	var httpStatus int

	var dbConn config.DbConnection
	dbConn.GetDbConnectionPool().AddTableWithName(Campaigns{}, "campaigns").SetKeys(true, "id")
	campaignResult, err := dbConn.GetDbConnectionPool().Get(Campaigns{}, idReqInt)
	if err != nil {
		resp.Code = "BE-001"

		msg = fmt.Sprintf("failed to retrieve campaign with id %s", idReq)
		httpStatus = http.StatusInternalServerError
		log.Info().Msgf(msg+" | %s", err.Error())
	} else if campaignResult == nil {
		resp.Code = "BE-002"

		msg = fmt.Sprintf("campaign with id %s could not be found", idReq)
		httpStatus = http.StatusNotFound
		log.Info().Msgf(msg)
	} else {
		resp.Code = "BE-000"

		resp.Campaign = append(resp.Campaign, buildCampaignResponse(campaignResult))

		msg = fmt.Sprintf("campaign id %s successfully retrieved", idReq)
		httpStatus = http.StatusOK
		log.Info().Msgf(msg)
	}

	resp.Message = msg
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(httpStatus)
	json.NewEncoder(rw).Encode(resp)
}

type CampaignData struct {
	Id    uint64 `json:"id,omitempty"`
	Title string `json:"title"`
}

func buildCampaignResponse(campaignResult interface{}) (resp CampaignData) {
	result := campaignResult.(*Campaigns)
	resp.Id = *result.Id
	resp.Title = result.Title
	return resp
}
