<div id="top"></div>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://kitabisa.com">
    <img src="assets/homework.jpg" alt="Logo" width="160" height="160">
  </a>

<h3 align="center">Backend Homework Test</h3>
  <p align="center">
A repository that contains homework test for backend candidates.
<br />
  </p>
</div>

## Table of Contents
1. [Background Story](#background-story)
2. [Requirement](#requirement)
3. [How to run](#how-to-run)
4. [Your Job](#your-job)
5. [Request & Response Payload](#request--response-payload)
6. [Rules](#rules)

## Background Story
Kitabisa is planning to migrate from monolith to API based.
To begin with this project, we wanted to migrate only the core functionality.

Here is the progress so far:

- Payment
  - Create Payment (Dummy)
  - Get Payment
- Campaign
  - Create Campaign (Dummy)
  - Get Campaign
- Donation 
  - Create Donation (Not Yet Implemented)
  - Get Donation (Not Yet Implemented)

## Requirement
In order to run the application, you need
- Docker & Docker Compose v2++ (usually bundled together with Docker installation)
- Go at least v1.18

## How to run
1. Run the docker compose to start the database
```
docker-compose up -d
```

2. Run the application, it will start on port 8080
```
go mod download
go run cmd/app/main.go
```

3. Try curl or run postman request for health check
```
curl -I --request GET 'localhost:8080/health_check/db'
```

it should response with these result
```
HTTP/1.1 200 OK
Vary: Origin
Date: Fri, 25 Nov 2022 17:05:21 GMT
Content-Length: 0
```

## Your Job

Your tasks are 
- Implement [Create Payment API](https://github.com/kitabisa/backend-takehome-test/blob/75f23043f83885b6d7300703454df3efcab981b6/api/v1/payment/payment.go#L30)
  - ```
    func CreateNewPaymentMethodHandler(rw http.ResponseWriter, r *http.Request) {
    ```
- Implement [Create Campaign API](https://github.com/kitabisa/backend-takehome-test/blob/75f23043f83885b6d7300703454df3efcab981b6/api/v1/campaign/campaign.go#L30)
  - ```
    func GetCampaignByIdHandler(rw http.ResponseWriter, r *http.Request) {
    ```
- Implement [Create Donation API](https://github.com/kitabisa/backend-takehome-test/blob/3d5c99ef6e3f302518d43b0fab93868af613480c/api/v1/donation/donation.go#L15)
  - ```
    func CreateDonationHandler(rw http.ResponseWriter, r *http.Request) {
    ```
- Implement [Get Donation API](https://github.com/kitabisa/backend-takehome-test/blob/3d5c99ef6e3f302518d43b0fab93868af613480c/api/v1/donation/donation.go#L19)
  - ```
    func GetDonationByIdHandler(rw http.ResponseWriter, r *http.Request) {
    ```


If you are confused where to implement, just search it inside this codebase

### Request & Response Payload
Example API request and response format can be found in a json file inside `docs` folder.
You should use [Postman](https://www.postman.com/) and import the json file to see the request and response format.

### Rules
- Please use Git so we can see how your code evolve. 
- Please don't cheat, respect our code of conduct.
- You are free to use any kind of IDE. We recommend Goland from Jetbrains.
- You shouldn't change anything in Postman collection.
- We are not expecting a clean code but we want to see your ability to write a code that fulfills the requirement in Golang. The existing codebase should give you enough example to tinker around.
- Make sure you read the postman examples carefully because there are some business logic validations (e.g payment method not exist should be handled)
- Once you finish the task, please run the integration test (**make sure Docker is running in your machine**)
```
sh integration_test.sh

....

Example result

--- FAIL: Test_BackendTestSuite (13.77s)
    --- PASS: Test_BackendTestSuite/Test_01_1GetCampaignOK (0.01s)
    --- PASS: Test_BackendTestSuite/Test_01_2GetCampaignFailed_NotFound (0.00s)
    --- PASS: Test_BackendTestSuite/Test_01_3GetCampaignFailed_IdNotNumber (0.00s)
    --- PASS: Test_BackendTestSuite/Test_02_1CreateCampaignOK (0.00s)
    --- PASS: Test_BackendTestSuite/Test_03_1CreatePaymentMethodOK (0.00s)
    --- PASS: Test_BackendTestSuite/Test_04_1GetPaymentMethodOK (0.00s)
    --- PASS: Test_BackendTestSuite/Test_04_2GetPaymentMethod_NotFound (0.00s)
    --- PASS: Test_BackendTestSuite/Test_04_3GetPaymentMethod_IdNotNumber (0.00s)
    --- FAIL: Test_BackendTestSuite/Test_05_1GetDonationOK (0.00s)
    --- FAIL: Test_BackendTestSuite/Test_05_2GetDonation_NotFound (0.00s)
    --- FAIL: Test_BackendTestSuite/Test_05_3GetDonation_IdNotNumber (0.00s)
    --- FAIL: Test_BackendTestSuite/Test_06_1CreateDonationOK (0.00s)
    --- FAIL: Test_BackendTestSuite/Test_06_2CreateDonationCampaignNotFound (0.00s)
    --- FAIL: Test_BackendTestSuite/Test_06_3CreateDonationPaymentMethodNotFound (0.00s)
    --- FAIL: Test_BackendTestSuite/Test_06_4CreateDonationAmountLessThan10000 (0.00s)
FAIL
FAIL    github.com/kitabisa/backend-takehome-test/integration_test/backend      14.294s


```

- Integration Test result will generate result and your score
  - The objective is to fix the `FAIL` part so it turns into `PASS` by modifying the codebase
  - Score > 90 with PASS result means you are qualified to continue the next stage. 
  - Otherwise, your implementation of API request / response is different than the one we expected or you change something that shouldn't be changed.
  If this happen, verify your request / response payload with the postman example we gave to you. 
  - This integration test also serve as a reference that your code is working on your machine and should be working on our test pipeline.

Integration test will run in port, make sure these port are free 
- 8081 (app)
- 5433 (db)
  



## Copyright
Kitabisa™ 2024
Authored By : Sactio Swastioyono