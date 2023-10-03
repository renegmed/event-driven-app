package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	pkgmodel "job-runner-app/pkg/model"
	"log"
	"net/http"
)

type client struct {
	baseAddress string
}

func newClient() *client {
	return &client{baseAddress: "http://localhost:4001"}
}

func (c *client) CreateWorker(name string, desc string) (pkgmodel.CreateWorkerResponse, error) {
	client := http.Client{}

	workerRequest := pkgmodel.CreateWorkerRequest{
		Name:        name,
		Description: desc,
	}

	b, err := json.Marshal(workerRequest)
	if err != nil {
		return pkgmodel.CreateWorkerResponse{}, err
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/workers", c.baseAddress),
		r,
	)

	if err != nil {
		return pkgmodel.CreateWorkerResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return pkgmodel.CreateWorkerResponse{}, err
	}

	defer resp.Body.Close()

	log.Println("...resp.Status:", resp.StatusCode, " http StatusCreated:", http.StatusCreated)

	if resp.StatusCode != http.StatusCreated {
		return pkgmodel.CreateWorkerResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cwresp pkgmodel.CreateWorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&cwresp); err != nil {
		return pkgmodel.CreateWorkerResponse{}, err
	}

	return cwresp, nil
}

func (c *client) GetWorker(id string) (pkgmodel.GetWorkerResponse, error) {
	client := http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/workers/%s", c.baseAddress, id),
		nil,
	)

	if err != nil {
		return pkgmodel.GetWorkerResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return pkgmodel.GetWorkerResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pkgmodel.GetWorkerResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var gwresp pkgmodel.GetWorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&gwresp); err != nil {
		return pkgmodel.GetWorkerResponse{}, err
	}

	return gwresp, nil

}

func (c *client) DeleteWorker(id string) (pkgmodel.DeleteWorkerResponse, error) {
	client := http.Client{}

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/workers/%s", c.baseAddress, id),
		nil,
	)

	if err != nil {
		return pkgmodel.DeleteWorkerResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return pkgmodel.DeleteWorkerResponse{}, err
	}

	defer resp.Body.Close()

	log.Println("...e2e.DeleteWorker() resp.StatusCode", resp.StatusCode, " http.StatusOK", http.StatusOK)

	if resp.StatusCode != http.StatusOK {
		return pkgmodel.DeleteWorkerResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dwresp pkgmodel.DeleteWorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&dwresp); err != nil {
		return pkgmodel.DeleteWorkerResponse{}, err
	}

	return dwresp, nil
}

func (c *client) LaunchJob(workerID string, body map[string]any) (pkgmodel.LaunchJobResponse, error) {
	client := http.Client{}

	r := pkgmodel.LaunchJobRequest{
		WorkerID: workerID,
		Input:    body,
	}

	b, err := json.Marshal(r)
	if err != nil {
		return pkgmodel.LaunchJobResponse{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/jobs", c.baseAddress),
		bytes.NewReader(b),
	)
	if err != nil {
		return pkgmodel.LaunchJobResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return pkgmodel.LaunchJobResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pkgmodel.LaunchJobResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ljresp pkgmodel.LaunchJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&ljresp); err != nil {
		return pkgmodel.LaunchJobResponse{}, err
	}

	return ljresp, nil

}
func (c *client) GetJobStatuses(jobID string) (pkgmodel.JobStatusesResponse, error) {
	client := http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/jobs/%s/statuses", c.baseAddress, jobID),
		nil,
	)

	if err != nil {
		return pkgmodel.JobStatusesResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return pkgmodel.JobStatusesResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pkgmodel.JobStatusesResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var gjsresp pkgmodel.JobStatusesResponse
	if err := json.NewDecoder(resp.Body).Decode(&gjsresp); err != nil {
		return pkgmodel.JobStatusesResponse{}, err
	}

	return gjsresp, nil
}
