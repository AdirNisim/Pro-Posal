package integrationtests

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pro-posal/webserver/models"
	"github.com/stretchr/testify/assert"
)

func TestPostAndGetCompany(t *testing.T) {
	req := map[string]string{
		"name":    gofakeit.BeerName(),
		"address": gofakeit.Address().Address,
	}

	var resp models.Company
	client.Post(t, "/companies", req, http.StatusCreated, &resp)

	assert.Equal(t, req["name"], resp.Name)
	assert.Equal(t, req["address"], resp.Address)
	assert.NotEmpty(t, resp.ID)
	assert.NotEmpty(t, resp.CreatedAt)
	assert.NotEmpty(t, resp.UpdatedAt)

	// // TODO: Add Get company
	// var getResp models.Company
	// client.Get(t, "/companies?id="+resp.ID, http.StatusOK, &getResp)
}
