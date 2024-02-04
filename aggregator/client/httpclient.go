package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Naman15032001/tolling/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invReq := types.GetInvoiceRequest{
		ObuId: int32(id),
	}
	b, err := json.Marshal(&invReq)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/%s?obu=%d", c.Endpoint, "invoice", id)
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responsded with non 200 status code %d", resp.StatusCode)
	}
	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &inv, nil
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregrateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responsded with non 200 status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}
