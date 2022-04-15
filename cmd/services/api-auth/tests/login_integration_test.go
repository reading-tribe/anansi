package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/reading-tribe/anansi/pkg/nettypes"
)

var (
	TEST_EMAIL    = "test@user.com"
	TEST_PASSWORD = "test"
)

func TestLogin(t *testing.T) {
	t.Log("\tAPI-Auth Login Test")
	t.Log("\tBuilding Request")
	req := nettypes.LoginRequest{
		EmailAddress: TEST_EMAIL,
		Password:     TEST_PASSWORD,
	}

	t.Log("\tMarshalling JSON")

	reqBytes, jsonMarshalErr := json.Marshal(req)
	if jsonMarshalErr != nil {
		t.Fatal(jsonMarshalErr)
	}
	t.Log("\tPOST Request")

	res, postErr := http.Post(
		"https://81qrzgok36.execute-api.eu-central-1.amazonaws.com/auth/login",
		"Content-Type: application/json",
		bytes.NewReader(reqBytes),
	)
	t.Log("\tCheck POST Request Error")
	if postErr != nil {
		t.Fatal(postErr)
	}

	t.Log("\tCheck Response Code")
	if res.StatusCode != http.StatusOK {
		t.Errorf("Got non 200 status code, %d", res.StatusCode)
	}

	t.Log("\tParse Response Body")
	var resParsed nettypes.LoginResponse
	decodeErr := json.NewDecoder(res.Body).Decode(&resParsed)
	t.Log("\tCheck Parsing Response Body Error")
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	t.Log("\tCheck Parsed Body Message")
	if resParsed.Message != "Successfully logged user in" {
		t.Errorf("Unexpected message received: %s", resParsed.Message)
	}
	t.Log("\tPass")
}
