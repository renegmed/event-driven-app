package e2e_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteShouldReturnDeleted(t *testing.T) {
	client := newClient()

	r, err := client.CreateWorker(
		fmt.Sprintf("test-%s", uuid.NewString()),
		fmt.Sprintf("test-%s", uuid.NewString()),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.DeleteWorker(r.ID)

	assert.NoError(t, err)
	assert.True(t, resp.Deleted)
}

func TestDeleteShouldReturnNotDeleted(t *testing.T) {
	client := newClient()

	resp, err := client.DeleteWorker(uuid.NewString()) // non-existing worker

	log.Println(".....err:", err)

	assert.NoError(t, err)
	assert.False(t, resp.Deleted)

}
