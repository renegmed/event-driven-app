package e2e_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetShouldBeOK(t *testing.T) {
	client := newClient()

	name := fmt.Sprintf("test-%s", uuid.NewString())
	desc := fmt.Sprintf("test-%s", uuid.NewString())

	r, err := client.CreateWorker(
		name,
		desc,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.GetWorker(r.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, r.ID, resp.Worker.ID)
	assert.Equal(t, name, resp.Worker.Name)
	assert.Equal(t, desc, resp.Worker.Description)
}

func TestGetShouldReturnNotFound(t *testing.T) {
	client := newClient()

	_, err := client.GetWorker(uuid.NewString())

	assert.Error(t, err)
	assert.Equal(t, "unexpected status code: 404", err.Error())
}
