package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// NationalCodeRequest represents the request payload for the API
type NationalCodeRequest struct {
	NationalCode string `json:"national_code"`
	APIKey       string `json:"api_key"`
}

// NationalCodeResponse represents the response payload from the API
type NationalCodeResponse struct {
	Valid     bool   `json:"valid"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Error     string `json:"error,omitempty"`
}

func CheckNationalCode(nationalCode, apiKey string) (*NationalCodeResponse, error) {
	// API URL
	apiURL := "https://api.sabteahval.ir/national-code"

	// Prepare the payload
	requestData := NationalCodeRequest{
		NationalCode: nationalCode,
		APIKey:       apiKey,
	}
	payload, _ := json.Marshal(requestData)

	// Make the HTTP POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("something went wrong: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Read the response
	body, _ := io.ReadAll(resp.Body)

	// Parse the JSON response
	var response NationalCodeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("something went wrong: %v", err)
	}

	return &response, nil
}

// ValidateIranianNationalCode validates the structure of an Iranian National Code
func ValidateIranianNationalCode(nationalCode string) error {
	// Ensure it's exactly 10 digits
	if matched, _ := regexp.MatchString(`^\d{10}$`, nationalCode); !matched {
		return errors.New("the national code most contain exactly 10 characters")
	}

	// Extract digits and calculate check digit
	checkDigit := int(nationalCode[9] - '0')
	sum := 0
	for i := 0; i < 9; i++ {
		digit := int(nationalCode[i] - '0')
		sum += digit * (10 - i)
	}

	// Calculate the expected check digit
	remainder := sum % 11
	expectedCheckDigit := remainder
	if remainder >= 2 {
		expectedCheckDigit = 11 - remainder
	}

	// Validate the check digit
	if checkDigit != expectedCheckDigit {
		return errors.New("the national code is not valid")
	}

	return nil
}
