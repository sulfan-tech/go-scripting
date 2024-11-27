package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go-scripting/pkg/logger"
	"go-scripting/repositories/firestore"

	"github.com/joho/godotenv"
)

const (
	apiDraftURL = "https://web.svc.fithub.id/v1/transactions"
	// apiDraftURL = "https://web.svc.staging.fithubdev.com/v1/transactions"
	apiPaidURL = "https://web.svc.fithub.id/v1/transactions/%s/paid"
	// apiPaidURL        = "https://web.svc.staging.fithubdev.com/v1/transactions/%s/paid"
	apiKey            = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzEwNTU0NDQsImlhdCI6MTczMDQ1MDY0NCwibG9jYXRpb24iOiJGSVQgSFVCIEhFQUQgT0ZGSUNFIiwibG9jYXRpb25zIjoiZml0aHVic3VudGVyYWx0aXJhfGZpdGh1YmFydGVyaXBvbmRva2luZGFofGZpdGh1YnBhbmNvcmFufGZpdGh1YmJlbmhpbHxmaXRodWJtZW5hcmFkdXRhfGZpdGh1YmRhcm1vfGZpdGh1Ymd1YmVuZ3xmaXRodWJncmVlbnZpbGxlfGZpdGh1Ym5pcGFobWFsbHxmaXRodWJsZW5na29uZ3xmaXRodWJwbHVpdHxmaXRodWJtYW5hZG98Zml0aHVibWVydXlhfGZpdGh1YndpeXVuZ3xmaXRodWJhbGFtc3V0cmF8Zml0aHVia2VtYXlvcmFufGZpdGh1YmJzZHxmaXRodWJnYWRpbmdzZXJwb25nfGZpdGh1YmpnY3xmaXRodWJibG9rbXxmaXRodWJrb3BvfGZpdGh1YmthcmF3YWNpfGZpdGh1Yml0Y2t1bmluZ2FufGZpdGh1YnRlc3R8Zml0aHViYm9nb3J8Zml0aHVibWFueWFyfGZpdGh1YnJhd2FtYW5ndW58Zml0aHVidGViZXR8Zml0aHViYmludGFyb3xmaXRodWJzYWxlbWJhfGZpdGh1YmdyYWhhcGVuYXxmaXRodWJjaWxlZHVnfGZpdGh1YnBhbXVsYW5nfGZpdGh1YmZhdG1hd2F0aXxmaXRodWJyZW5vbnxmaXRodWJiZWthc2loYXJhcGFuaW5kYWh8Zml0aHViYmludGFyb2pvbWJhbmd8Zml0aHVic3VyeWFzdW1hbnRyaXxmaXRodWJjaXRyYXJheWF8Zml0aHVia2FyYW5nbXVseWF8Zml0aHViZGVwb2t0b3duc3F1YXJlfGZpdGh1YmRlcG9rdG93bmNlbnRlcnxmaXRodWJncm9nb2x8Zml0aHVia2VsYXBhZ2FkaW5nfGZpdGh1YmVsYW5nbGF1dHBpa3xmaXRodWJjaXB1dGF0bG90fGZpdGh1YmFkaXR5YXdhcm1hbnxmaXRodWJ0YW5nY2l0eXxmaXRodWJwdXJpaW5kYWh8Zml0aHViY2l0cmFsYW5kfGZpdGh1Ym1hbXBhbmd8Zml0aHVidGV1a3V1bWFyfGZpdGh1Ym1lcnJ8Zml0aHViZ2FqYWhtYWRhfGZpdGh1Ym1heW9yb2tpbmd8Zml0aHViam9namFjaXR5bWFsbHxmaXRodWJkLmlwYW5qYWl0YW58Zml0aHVic2VkYXl1Y2l0eXxmaXRodWJkdXJlbnNhd2l0fGZpdGh1YmNlbXBha2FwdXRpaHxmaXRodWJpdGNwZXJtYXRhaGlqYXV8Zml0aHVidGFtaW5pc3F1YXJlfGZpdGh1Ym1hbGxtZXRyb2tlYmF5b3JhbnxmaXRodWJrZW1hbmdnaXNhbnxmaXRodWJiaW50YXJvc2VrdG9yMXxmaXRodWJzdW5zZXRyb2FkfGZpdGh1YmdyZXNzbWFsbHxmaXRodWJrYWxpbWFsYW5nfGZpdGh1Yndpc21hYm5pNDZ8Zml0aHVic2V0aWFidWRpc2VtYXJhbmd8Zml0aHVic2V0aWFidWRpc2VtYXJhbmd8Zml0aHViYnNkZ29sZGZpbmNofGZpdGh1YnBhc2FybWluZ2d1fGZpdGh1YnBsYXphb2xlb3N8Zml0aHViY2luZXJlfGZpdGh1Ym1hc3Bpb25wbGF6YXxmaXRodWJ0YW1hbnBhbGVtbGVzdGFyaXxmaXRodWJsaXZpbmdwbGF6YWphYmFiZWthfGZpdGh1YnRyYW5zeW9naWNpYnVidXJ8Zml0aHVidHJhbnN5b2dpY2lidWJ1cnxmaXRodWJwYWhsYXdhbnNpZG9hcmpvfGZpdGh1YnBhaGxhd2Fuc2lkb2Fyam98Zml0aHViZ2F0c3V0aW11cnxmaXRodWJic2JzZW1hcmFuZ3xmaXRodWJzdW1tYXJlY29uYmVrYXNpfGZpdGh1YmJ1YWhiYXR1fGZpdGh1YnBhc2Fya290YXBvbmRva2dlZGV8Zml0aHVidWx1d2F0dXxmaXRodWJtYWphcGFoaXRzZW1hcmFuZ3xmaXRodWJiYWxvaXBlcnNlcm9iYXRhbXxmaXRodWJrZW5qZXJhbnxmaXRodWJqYXRpYXNpaHxmaXRodWJvdGlzdGFyYXlhfGZpdGh1Ynlhc21pbmJvZ29yfGZpdGh1YmRpZW5nbWFsYW5nfGZpdGh1YmNpYmFiYXRjaW1haGl8Zml0aHVicmF0dWxhbmdpbWFrYXNzYXJ8Zml0aHViZW50ZXJwcmlzZXxmaXRodWJzbGFtZXRyaXlhZGlzb2xvfGZpdGh1YnN1a2FtdG9wYWxlbWJhbmd8Zml0aHViaGVhZG9mZmljZXxmaXRodWJkYW1haWJhbGlrcGFwYW58Zml0aHViZGFtYWliYWxpa3BhcGFufGZpdGh1YnNvbG9iYXJ1fGZpdGh1YmthcnRpbmlkZXBva3xmaXRodWJqZW11cnNhcml8Zml0aHVicGFtdWxhcnNpaHNlbWFyYW5nfGZpdGh1Ym1hbnVrYW58Zml0aHViZ2F0c3ViYXJhdHxmaXRodWJjaXR5d2Fsa2xpcHBvY2lrYXJhbmd8Zml0aHViZ3JhbmR3aXNhdGFiZWthc2l8Iiwicm9sZXNNeVNRTCI6W10sInN1YiI6ImFkbWluLXByb2RAZml0aHViLmlkIiwidHlwZVVzZXIiOiJBZG1pbiJ9.4qtOfNJmgEkpvgaPavvZQ-OyzlqmRKq0H8BZdCd5fA0"
	customPackageName = "FREE SESSION CORPORATE - FUJIFILM"
	contentType       = "application/json"
	csvFilePath       = "input.csv"
	timeout           = 5 * time.Second
)

var usersRepository firestore.UsersRepo
var trainerRepo firestore.PersonalTrainer

type RequestDraftPayload struct {
	LocationOn       string `json:"locationOn"`
	UserId           string `json:"userId"`
	SalesBy          string `json:"salesBy"`
	PtId             string `json:"ptId"`
	PackageId        int    `json:"packageId"`
	Sessions         int    `json:"sessions"`
	Price            int    `json:"price"`
	Expiry           int    `json:"expiry"`
	PaymentMethodId1 string `json:"paymentMethodId1"`
	Payment1         int    `json:"payment1"`
	PaymentID1       string `json:"paymentID1"`
	PaymentMethodId2 string `json:"paymentMethodId2"`
	PromoCode        string `json:"promoCode"`
	StartDate        string `json:"startDate"`
	TypeTransaction  string `json:"typeTransaction"`
	PromoCodeInput   string `json:"promoCodeInput"`
	SalesType        string `json:"salesType"`
	IsVerified       bool   `json:"isVerified"`
}

type TransactionRequest struct {
	TypeTransaction        string `json:"typeTransaction"`
	TransactionId          string `json:"transId"`
	UserId                 string `json:"userId" validate:"required"`
	DateTransaction        string `json:"dateTransaction" validate:"omitempty"`
	UserName               string `json:"userName" validate:"omitempty,max=100"`
	Phone                  string `json:"phone" validate:"required,max=18"`
	PackageId              string `json:"packageId" validate:"required"`
	PackageLocationType    string `json:"packageLocationType"`
	PromoCode              string `json:"promoCode"`
	PromoCodeInfo          string `json:"promoCodeInfo"`
	Payment1               string `json:"payment1"`
	Price                  string `json:"price"`
	Total                  string `json:"total"`
	PaymentMethod          string `json:"paymentMethod"`
	PaymentID1             string `json:"paymentID1" validate:"required"`
	PaymentMethodId1       string `json:"paymentMethodId1" validate:"required"`
	PaymentMethodDetailId1 string `json:"paymentMethodDetailId1" validate:"omitempty"`
	Payment2               int64  `json:"payment2"`
	PaymentID2             string `json:"paymentID2"`
	PaymentMethodId2       string `json:"paymentMethodId2"`
	PaymentMethodDetailId2 string `json:"paymentMethodDetailId2" validate:"omitempty"`
	Remarks                string `json:"remarks"`
	HomeClub               string `json:"homeClub" validate:"required"`
	StartDate              string `json:"startDate" validate:"required,date=2006-01-02"`
	ExpiredDate            string `json:"expiredDate" validate:"required"`
	StartFrom              string `json:"startFrom" validate:"required,date=2006-01-02"`
	SalesBy                string `json:"salesBy" validate:"required"`
	SalesType              string `json:"salesType" validate:"omitempty"`
	StatusTransaction      string `json:"statusTransaction"`
	UpdatedBy              string `json:"updatedBy"`
	UpdatedDate            string `json:"updatedDate"`
	PTId                   string `json:"ptId"`
	PTLevel                string `json:"ptLevel"`
	PTName                 string `json:"ptName"`
	Session                string `json:"session"`
	SessionUser            string `json:"sessionUser"`
	LeadsId                string `json:"leadsId" validate:"omitempty"`
	TransferType           string `json:"transferType"`
	IsPTtransferred        bool   `json:"isPTtransferred"`
	Name                   string `json:"name"`
	LocationOn             string `json:"locationOn" validate:"required"`
	Months                 int    `json:"months"`
	Sessions               int    `json:"sessions"`
}

type RequestPaidPayload struct {
	Promo            string
	StartDate        string
	UserId           string
	SalesBy          string
	PtId             string
	PackageId        int
	PaymentMethodId1 string
	PaymentID1       string
	ExpiredDate      string
	DateTransaction  string
	TypeTransaction  string
	ImageReceiptPath string
	ImageOtherPath   string
}

func init() {
	env := "../../.env"
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}
	usersRepository = firestore.NewFirestoreUsersRepository()
	trainerRepo = firestore.NewFirestoreTrainerRepository()
}

func main() {

	client := &http.Client{
		Timeout: timeout,
	}
	csvData, err := readCSV(csvFilePath)
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	for i, dataCSV := range csvData {
		fmt.Printf("Processing row %d: %v\n", i+1, dataCSV)

		// query dataCSV data on db
		memberFS, err := usersRepository.QueryUserByUserAppId(dataCSV.UserId)
		if err != nil {
			fmt.Println("Error quering member:", err)
			continue
		}

		trainerData, err := trainerRepo.GetPT(context.Background(), dataCSV.SalesBy)
		if err != nil {
			fmt.Println("Error quering trainer:", err)
			continue
		}

		// Create RequestDraftPayload from dataCSV data
		payload := RequestDraftPayload{
			UserId:           memberFS.Uid,
			SalesBy:          dataCSV.SalesBy,
			LocationOn:       dataCSV.LocationOn,
			PtId:             dataCSV.SalesBy,
			PackageId:        111822,
			Sessions:         1,
			Price:            1,
			Expiry:           60,
			PaymentMethodId1: "4",
			Payment1:         1,
			PaymentID1:       "002",
			PaymentMethodId2: "",
			PromoCode:        "",
			StartDate:        "2024-11-04",
			TypeTransaction:  "pt",
			PromoCodeInput:   "",
			SalesType:        "Personal Trainer",
			IsVerified:       false,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			continue
		}

		// First request to create a draft transaction
		req, err := createRequest(apiDraftURL, jsonPayload)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			printResponse(resp, i+1, dataCSV.UserId)
			continue
		}

		var draftResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&draftResponse)
		if err != nil {
			fmt.Println("Error decoding draft response:", err)
			continue
		}

		data, ok := draftResponse["data"].(map[string]interface{})
		if !ok {
			fmt.Println("Error: transaction_id not found in draft response")
			continue
		}

		transactionId := data["transId"].(string)

		// Second request to mark the transaction as paid
		paidURL := fmt.Sprintf(apiPaidURL, transactionId)
		fmt.Println(paidURL, "PAID URL")

		var trainerLevel string

		if trainerData.Level == 0 {
			trainerLevel = "BASIC"
		} else {
			trainerLevel = "ADVANCE"
		}

		//from data
		payloadFromData := TransactionRequest{
			UserId:              memberFS.Uid,
			TransactionId:       transactionId,
			LocationOn:          dataCSV.LocationOn,
			Name:                customPackageName,
			PackageId:           "111822",
			PackageLocationType: "ALL_CLUB",
			PaymentID1:          "002",
			Payment1:            "1",
			PaymentMethod:       "Transfer",
			PaymentMethodId1:    "4",
			Phone:               memberFS.Phone,
			Price:               "1",
			PTId:                trainerData.Email,
			PTLevel:             trainerLevel,
			PTName:              trainerData.Name,
			SalesBy:             trainerData.Email,
			Session:             "1",
			StartDate:           "2024-11-04T00:00:00Z",
			StatusTransaction:   "Draft",
			Total:               "1",
			TypeTransaction:     "PT",
			UpdatedBy:           "admin-stg@fithub.id",
			UpdatedDate:         "2024-11-04T00:00:00.822Z",
			UserName:            memberFS.Name,
			ExpiredDate:         "18 Jun 2024",
			DateTransaction:     "2024-11-04",
		}

		payloadMap := map[string]string{
			"userId":              payloadFromData.UserId,
			"transId":             payloadFromData.TransactionId,
			"locationOn":          payloadFromData.LocationOn,
			"name":                payloadFromData.Name,
			"packageId":           payloadFromData.PackageId,
			"packageLocationType": payloadFromData.PackageLocationType,
			"paymentMethod":       payloadFromData.PaymentMethod,
			"paymentMethod1":      payloadFromData.PaymentMethodId1,
			"paymentMethodId1":    payloadFromData.PaymentMethodId1,
			"phone":               payloadFromData.Phone,
			"price":               payloadFromData.Price,
			"ptId":                payloadFromData.PTId,
			"ptLevel":             payloadFromData.PTLevel,
			"ptName":              payloadFromData.PTName,
			"salesBy":             payloadFromData.SalesBy,
			"sessions":            payloadFromData.Session,
			"startDate":           payloadFromData.StartDate,
			"total":               payloadFromData.Total,
			"typeTransaction":     payloadFromData.TypeTransaction,
			"updatedBy":           payloadFromData.UpdatedBy,
			"updatedDate":         payloadFromData.UpdatedDate,
			"userName":            payloadFromData.UserName,
			"expiredDate":         payloadFromData.ExpiredDate,
			"dateTransaction":     payloadFromData.DateTransaction,
			"imageReceipt":        "",
		}

		req, err = createFormRequest(paidURL, payloadMap)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("Error sending paid request:", err)
			continue
		}
		defer resp.Body.Close()

		printResponse(resp, i+1, dataCSV.UserId)
	}
}

func readCSV(filePath string) ([]RequestDraftPayload, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var csvData []RequestDraftPayload
	for i, record := range records {
		if i == 0 { // Skip header row
			continue
		}

		dataCSV := RequestDraftPayload{
			UserId:     record[0],
			SalesBy:    record[1],
			LocationOn: record[2],
			PtId:       record[2],
		}
		csvData = append(csvData, dataCSV)
	}

	return csvData, nil
}

func createRequest(url string, jsonPayload []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return req, nil
}

func createFormRequest(urlPath string, formData map[string]string) (*http.Request, error) {
	// Create form data string
	data := url.Values{}
	for key, value := range formData {
		data.Set(key, value)
	}

	// Create a new POST request with form data
	req, err := http.NewRequest("POST", urlPath, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Set appropriate headers for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return req, nil
}

func printResponse(resp *http.Response, rowNumber int, userId string) {
	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println("Response Status Code:", resp.Status)

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		successMessage := fmt.Sprintf("Successfully processed row %d with userId %s", rowNumber, userId)
		logger.LogInfo(successMessage)
	} else {
		fmt.Printf("Error processing row %d: %v\n", rowNumber, responseBody)
	}
}
