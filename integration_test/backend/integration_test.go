package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/kitabisa/backend-takehome-test/api/v1/campaign"
	"github.com/kitabisa/backend-takehome-test/api/v1/payment"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type BackendTestSuite struct {
	suite.Suite
	Context            context.Context
	BackendTestCompose testcontainers.DockerCompose
	Score              int64            // Used in each test case to define weight of test score
	ScoreMap           map[string]int64 //key: test name, value: score. This describes the final score that candidate will get
	CumulativeMap      map[string]int64 //key: test name, value: score. This describes the cumulative score that candidate can get as test progress
}

// Run this before beginning test `Suite
func (suite *BackendTestSuite) SetupSuite() {
	//setup phase
	ctx := context.Background()
	suite.Context = ctx
	suite.ScoreMap = make(map[string]int64)
	suite.CumulativeMap = make(map[string]int64)

	compose, err := SetupDockerAppTest()
	if err != nil {
		panic(err)
	}

	suite.BackendTestCompose = compose

}

func Test_BackendTestSuite(t *testing.T) {
	suite.Run(t, new(BackendTestSuite))
}

func getScoreResultMap(scoreMap map[string]int64) []string {
	keys := make([]string, 0)
	for k, _ := range scoreMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func writeResult(fileName string, totalScore int64, score map[string]int64, cumulative map[string]int64) {
	pwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	dirPaths := strings.Split(pwd, "/integration_test/backend")

	f, err := os.Create(filepath.Join(dirPaths[0], fileName))

	if err != nil {
		log.Fatal(err)
	}

	// score result map contain value of test case ordered and the score that candidate get
	// each key will be used to retrieve the score and the max score that candidate can get for each question
	// Here we will write those result into the file
	for _, v := range getScoreResultMap(score) {
		result := fmt.Sprintf("Case : [%v]:  Cumulative Score/Max Score: [%d/%d]\n", v, score[v], cumulative[v])
		_, err2 := f.WriteString(result)
		if err2 != nil {
			log.Fatal(err2)
		}
	}

	_, err2 := f.WriteString(fmt.Sprintf("\nScore : %d", totalScore))
	if err2 != nil {
		log.Fatal(err2)
	}

	_, err2 = f.WriteString(fmt.Sprintf("\nResult : %v", judgePass(totalScore)))
	if err2 != nil {
		log.Fatal(err2)
	}

	f.Close()

}

func judgePass(score int64) string {
	if score > 90 {
		return "PASS"
	}

	return "FAIL"
}

func (suite *BackendTestSuite) TearDownSuite() {
	var fileOutputName = "BE-Test-Result.txt"

	// write result will write into a file with content of detail of each test cases and final result
	writeResult(fileOutputName, suite.Score, suite.ScoreMap, suite.CumulativeMap)
	if suite.BackendTestCompose != nil {
		compose := suite.BackendTestCompose
		compose.Down()
	}

	return
}

func (suite *BackendTestSuite) Test_02_1CreateCampaignOK() {
	var client = &http.Client{}

	body := apitest.New("CreateCampaign").
		EnableNetworking(client).
		Post("http://localhost:8081/v1/campaign/create").
		Header("Content-Type", "application/json").
		Body(`{"title": "campaign 1"}`).
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var campaignResponse campaign.CampaignResponse

	json.Unmarshal(responseBody, &campaignResponse)

	if assert.Equal(suite.T(), 201, body.StatusCode) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "BE-000", campaignResponse.Code) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "campaign created", campaignResponse.Message) {
		suite.Score += 1
	}

	suite.ScoreMap["Test_02_1CreateCampaignOK"] = suite.Score
	suite.CumulativeMap["Test_02_1CreateCampaignOK"] = 15
}
func (suite *BackendTestSuite) Test_01_1GetCampaignOK() {
	var client = &http.Client{}
	//insert first
	body := apitest.New("GetCampaign").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/campaign/0").
		Expect(suite.T()).
		End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var campaignResponse campaign.CampaignResponse

	json.Unmarshal(responseBody, &campaignResponse)
	dataResult := make(map[string]interface{})
	dataResult["id"] = float64(0)
	dataResult["title"] = "campaign 1"
	var arrayInterface []interface{}
	arrayInterface = append(arrayInterface, dataResult)

	//scoring
	if assert.Equal(suite.T(), 200, body.StatusCode) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), arrayInterface, campaignResponse.Campaign) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "BE-000", campaignResponse.Code) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "campaign id 0 successfully retrieved", campaignResponse.Message) {
		suite.Score += 1
	}

	suite.ScoreMap["Test_01_1GetCampaignOK"] = suite.Score
	suite.CumulativeMap["Test_01_1GetCampaignOK"] = 4
}

func (suite *BackendTestSuite) Test_01_2GetCampaignFailed_NotFound() {
	var client = &http.Client{}

	body := apitest.New("GetCampaign").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/campaign/3").
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var campaignResponse campaign.CampaignResponse

	json.Unmarshal(responseBody, &campaignResponse)

	var arrayInterface []interface{}

	if assert.Equal(suite.T(), 404, body.StatusCode) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "BE-002", campaignResponse.Code) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "campaign with id 3 could not be found", campaignResponse.Message) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), arrayInterface, campaignResponse.Campaign) {
		suite.Score += 1
	}

	suite.ScoreMap["Test_01_2GetCampaignFailed_NotFound"] = suite.Score
	suite.CumulativeMap["Test_01_2GetCampaignFailed_NotFound"] = 8

}

func (suite *BackendTestSuite) Test_01_3GetCampaignFailed_IdNotNumber() {
	var client = &http.Client{}

	body := apitest.New("GetCampaign").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/campaign/test").
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var campaignResponse campaign.CampaignResponse

	json.Unmarshal(responseBody, &campaignResponse)

	var arrayInterface []interface{}

	if assert.Equal(suite.T(), 400, body.StatusCode) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "BE-001", campaignResponse.Code) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "campaign id must be a number", campaignResponse.Message) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), arrayInterface, campaignResponse.Campaign) {
		suite.Score += 1
	}

	suite.ScoreMap["Test_01_3GetCampaignFailed_IdNotNumber"] = suite.Score
	suite.CumulativeMap["Test_01_3GetCampaignFailed_IdNotNumber"] = 12

}

func (suite *BackendTestSuite) Test_03_1CreatePaymentMethodOK() {
	var client = &http.Client{}

	body := apitest.New("CreatePaymentMethod").
		EnableNetworking(client).
		Post("http://localhost:8081/v1/payment-method/create").
		Header("Content-Type", "application/json").
		Body(`{"name": "payment method 1"}`).
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var paymentMethodResponse payment.PaymentMethodResponse

	json.Unmarshal(responseBody, &paymentMethodResponse)

	if assert.Equal(suite.T(), 201, body.StatusCode) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "BE-000", paymentMethodResponse.Code) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "payment method created", paymentMethodResponse.Message) {
		suite.Score += 1
	}

	suite.ScoreMap["Test_03_1CreatePaymentMethodOK"] = suite.Score
	suite.CumulativeMap["Test_03_1CreatePaymentMethodOK"] = 18
}

func (suite *BackendTestSuite) Test_04_1GetPaymentMethodOK() {
	var client = &http.Client{}

	body := apitest.New("GetPaymentMethod").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/payment-method/0").
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var paymentMethodResponse payment.PaymentMethodResponse

	json.Unmarshal(responseBody, &paymentMethodResponse)

	dataResult := make(map[string]interface{})
	dataResult["id"] = float64(0)
	dataResult["name"] = "payment method 1"
	var arrayInterface []interface{}
	arrayInterface = append(arrayInterface, dataResult)
	if assert.Equal(suite.T(), body.StatusCode, 200) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "BE-000", paymentMethodResponse.Code) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "payment method id 0 successfully retrieved", paymentMethodResponse.Message) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), arrayInterface, paymentMethodResponse.PaymentMethod) {
		suite.Score += 1
	}

	suite.ScoreMap["Test_04_1GetPaymentMethodOK"] = suite.Score
	suite.CumulativeMap["Test_04_1GetPaymentMethodOK"] = 22
}

func (suite *BackendTestSuite) Test_04_2GetPaymentMethod_NotFound() {
	var client = &http.Client{}

	body := apitest.New("GetPaymentMethod").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/payment-method/3").
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var paymentMethodResponse payment.PaymentMethodResponse

	json.Unmarshal(responseBody, &paymentMethodResponse)

	var arrayInterface []interface{}

	if assert.Equal(suite.T(), 404, body.StatusCode) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "BE-002", paymentMethodResponse.Code) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "payment method with id 3 could not be found", paymentMethodResponse.Message) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), arrayInterface, paymentMethodResponse.PaymentMethod) {
		suite.Score += 1
	}

	suite.ScoreMap["Test_04_2GetPaymentMethod_NotFound"] = suite.Score
	suite.CumulativeMap["Test_04_2GetPaymentMethod_NotFound"] = 26

}

func (suite *BackendTestSuite) Test_04_3GetPaymentMethod_IdNotNumber() {
	var client = &http.Client{}

	body := apitest.New("GetPaymentMethod").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/payment-method/test").
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var paymentMethodResponse payment.PaymentMethodResponse

	json.Unmarshal(responseBody, &paymentMethodResponse)
	var arrayInterface []interface{}

	if assert.Equal(suite.T(), 400, body.StatusCode) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "BE-001", paymentMethodResponse.Code) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), "payment method id must be a number", paymentMethodResponse.Message) {
		suite.Score += 1
	}

	if assert.Equal(suite.T(), arrayInterface, paymentMethodResponse.PaymentMethod) {
		suite.Score += 1
	}

	suite.ScoreMap["Test_04_3GetPaymentMethod_IdNotNumber"] = suite.Score
	suite.CumulativeMap["Test_04_3GetPaymentMethod_IdNotNumber"] = 30

}

type ResponseBodyResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type DonationResponse struct {
	ResponseBodyResult
	Donation []interface{} `json:"data"`
}

// below result should be success after candidate create donation endpoint implementation
func (suite *BackendTestSuite) Test_05_1GetDonationOK() {
	var client = &http.Client{}

	body := apitest.New("GetDonation").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/donation/0").
		Header("Content-Type", "application/json").
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)

	var donationResponse DonationResponse

	json.Unmarshal(responseBody, &donationResponse)

	dataResult := make(map[string]interface{})
	dataResult["id"] = float64(0)
	dataResult["payment_method_id"] = float64(0)
	dataResult["campaign_id"] = float64(0)
	dataResult["amount"] = float64(10000)
	var arrayInterface []interface{}
	arrayInterface = append(arrayInterface, dataResult)

	if assert.Equal(suite.T(), 200, body.StatusCode) {
		suite.Score += 3
	}

	if assert.Equal(suite.T(), "BE-000", donationResponse.Code) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), "donation id 0 is successfully retrieved", donationResponse.Message) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), arrayInterface, donationResponse.Donation) {
		suite.Score += 3
	}

	suite.ScoreMap["Test_05_1GetDonationOK"] = suite.Score
	suite.CumulativeMap["Test_05_1GetDonationOK"] = 40

}

func (suite *BackendTestSuite) Test_05_2GetDonation_NotFound() {
	var client = &http.Client{}

	body := apitest.New("GetDonation").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/donation/3").
		Header("Content-Type", "application/json").
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var donationResponse DonationResponse

	json.Unmarshal(responseBody, &donationResponse)
	var arrayInterface []interface{}
	if assert.Equal(suite.T(), 404, body.StatusCode) {
		suite.Score += 3
	}

	if assert.Equal(suite.T(), "BE-002", donationResponse.Code) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), "donation id 3 could not be found", donationResponse.Message) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), arrayInterface, donationResponse.Donation) {
		suite.Score += 3
	}

	suite.ScoreMap["Test_05_2GetDonation_NotFound"] = suite.Score
	suite.CumulativeMap["Test_05_2GetDonation_NotFound"] = 50

}

func (suite *BackendTestSuite) Test_05_3GetDonation_IdNotNumber() {
	var client = &http.Client{}

	body := apitest.New("GetDonation").
		EnableNetworking(client).
		Get("http://localhost:8081/v1/donation/test").
		Header("Content-Type", "application/json").
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var donationResponse DonationResponse

	json.Unmarshal(responseBody, &donationResponse)
	var arrayInterface []interface{}
	if assert.Equal(suite.T(), 400, body.StatusCode) {
		suite.Score += 3
	}

	if assert.Equal(suite.T(), "BE-001", donationResponse.Code) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), "donation id must be a number", donationResponse.Message) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), arrayInterface, donationResponse.Donation) {
		suite.Score += 3
	}

	suite.ScoreMap["Test_05_3GetDonation_IdNotNumber"] = suite.Score
	suite.CumulativeMap["Test_05_3GetDonation_IdNotNumber"] = 60

}

func (suite *BackendTestSuite) Test_06_1CreateDonationOK() {
	var client = &http.Client{}

	body := apitest.New("CreateDonation").
		EnableNetworking(client).
		Post("http://localhost:8081/v1/donation/create").
		Header("Content-Type", "application/json").
		Body(`{"payment_method_id":1,"campaign_id":1,"amount":10000}`).
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var donationResponse DonationResponse

	json.Unmarshal(responseBody, &donationResponse)
	dataResult := make(map[string]interface{})
	// JSON numbers always unmarshalled into float64
	dataResult["id"] = float64(1)
	dataResult["payment_method_id"] = float64(1)
	dataResult["campaign_id"] = float64(1)
	dataResult["amount"] = float64(10000)
	var arrayInterface []interface{}
	arrayInterface = append(arrayInterface, dataResult)
	if assert.Equal(suite.T(), body.StatusCode, 200) {
		suite.Score += 3
	}

	if assert.Equal(suite.T(), "BE-000", donationResponse.Code) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), "donation created successfully", donationResponse.Message) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), arrayInterface, donationResponse.Donation) {
		suite.Score += 3
	}

	suite.ScoreMap["Test_06_1CreateDonationOK"] = suite.Score
	suite.CumulativeMap["Test_06_1CreateDonationOK"] = 70
}

func (suite *BackendTestSuite) Test_06_2CreateDonationCampaignNotFound() {
	var client = &http.Client{}

	body := apitest.New("CreateDonation").
		EnableNetworking(client).
		Post("http://localhost:8081/v1/donation/create").
		Header("Content-Type", "application/json").
		Body(`{"payment_method_id":1,"campaign_id":999,"amount":10000}`).
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var donationResponse DonationResponse

	json.Unmarshal(responseBody, &donationResponse)

	var arrayInterface []interface{}
	if assert.Equal(suite.T(), donationResponse.Donation, arrayInterface) {
		suite.Score += 3
	}

	if assert.Equal(suite.T(), 500, body.StatusCode) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), "BE-001", donationResponse.Code) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), "failed to create donation because campaign id 999 does not exist", donationResponse.Message) {
		suite.Score += 3
	}

	suite.ScoreMap["Test_06_2CreateDonationCampaignNotFound"] = suite.Score
	suite.CumulativeMap["Test_06_2CreateDonationCampaignNotFound"] = 80
}

func (suite *BackendTestSuite) Test_06_3CreateDonationPaymentMethodNotFound() {
	var client = &http.Client{}

	body := apitest.New("CreateDonation").
		EnableNetworking(client).
		Post("http://localhost:8081/v1/donation/create").
		Header("Content-Type", "application/json").
		Body(`{"payment_method_id":999,"campaign_id":1,"amount":10000}`).
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var donationResponse DonationResponse

	json.Unmarshal(responseBody, &donationResponse)

	var arrayInterface []interface{}

	if assert.Equal(suite.T(), 500, body.StatusCode) {
		suite.Score += 3
	}

	if assert.Equal(suite.T(), "BE-001", donationResponse.Code) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), "failed to create donation because payment method id 999 does not exist", donationResponse.Message) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), arrayInterface, donationResponse.Donation) {
		suite.Score += 3
	}

	suite.ScoreMap["Test_06_3CreateDonationPaymentMethodNotFound"] = suite.Score
	suite.CumulativeMap["Test_06_3CreateDonationPaymentMethodNotFound"] = 90
}

func (suite *BackendTestSuite) Test_06_4CreateDonationAmountLessThan10000() {
	var client = &http.Client{}

	body := apitest.New("CreateDonation").
		EnableNetworking(client).
		Post("http://localhost:8081/v1/donation/create").
		Header("Content-Type", "application/json").
		Body(`{"payment_method_id":1,"campaign_id":1,"amount":9000}`).
		Expect(suite.T()).End().Response

	responseBody, _ := io.ReadAll(body.Body)
	var donationResponse DonationResponse

	json.Unmarshal(responseBody, &donationResponse)

	var arrayInterface []interface{}

	if assert.Equal(suite.T(), 400, body.StatusCode) {
		suite.Score += 3
	}

	if assert.Equal(suite.T(), "BE-001", donationResponse.Code) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), "failed to create donation because amount is less than 10000", donationResponse.Message) {
		suite.Score += 2
	}

	if assert.Equal(suite.T(), arrayInterface, donationResponse.Donation) {
		suite.Score += 3
	}

	suite.ScoreMap["Test_06_4CreateDonationAmountLessThan10000"] = suite.Score
	suite.CumulativeMap["Test_06_4CreateDonationAmountLessThan10000"] = 100
}
