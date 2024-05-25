package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pro-posal/webserver/models"
	"github.com/pro-posal/webserver/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var TEST_COMPANY_ID = uuid.NewSHA1(uuid.Nil, []byte("company_id")).String()
var TEST_CONTRACT_ID = uuid.NewSHA1(uuid.Nil, []byte("test")).String()

func TestPostCompanies_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	companyServiceMock := services.NewMockCompanyManagementService(ctrl)

	api := NewAPI(nil, nil, nil, companyServiceMock, nil, nil, nil, nil)

	companyServiceMock.EXPECT().
		CreateCompany(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, req services.CreateCompanyRequest) {
			assert.Equal(t, "Test Company", req.Name)
			assert.Equal(t, TEST_CONTRACT_ID, req.ContactID)
			assert.Equal(t, "Seattle, WA", req.Address)
			assert.Equal(t, "base64-example", req.LogoBase64)
		}).
		Return(&models.Company{
			ID:   TEST_COMPANY_ID,
			Name: "Acme Enterprises",
		}, nil)

	req, err := http.NewRequest("POST", "/companies", buildPostCompaniesPayload(t))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	api.PostCompanies(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp models.Company
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, TEST_COMPANY_ID, resp.ID)
	assert.Equal(t, "Acme Enterprises", resp.Name)
}

func TestPostCompanies_ErrorIsReturned(t *testing.T) {
	ctrl := gomock.NewController(t)
	companyServiceMock := services.NewMockCompanyManagementService(ctrl)

	api := NewAPI(nil, nil, nil, companyServiceMock, nil, nil, nil, nil)

	companyServiceMock.EXPECT().
		CreateCompany(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("failed to create company"))

	req, err := http.NewRequest("POST", "/companies", buildPostCompaniesPayload(t))
	require.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	api.PostCompanies(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func buildPostCompaniesPayload(t *testing.T) io.Reader {
	t.Helper()

	requestBody := map[string]string{
		"name":        "Test Company",
		"address":     "Seattle, WA",
		"contact_id":  TEST_CONTRACT_ID,
		"logo_base64": "base64-example",
	}

	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	return bytes.NewBuffer(bodyBytes)
}
