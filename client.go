package vektopay

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "time"
)

type Client struct {
  APIKey        string
  BaseURL       string
  DefaultHeader map[string]string
  httpClient    *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
  return &Client{
    APIKey:        apiKey,
    BaseURL:       trimSlash(baseURL),
    DefaultHeader: map[string]string{},
    httpClient:    &http.Client{Timeout: 30 * time.Second},
  }
}

func (c *Client) CreateCharge(input ChargeInput) (*ChargeResponse, error) {
  body := map[string]any{
    "customer_id": input.CustomerID,
    "card_id":     input.CardID,
    "amount":      input.Amount,
    "currency":    input.Currency,
    "installments": input.Installments,
    "country":     input.Country,
    "metadata":    input.Metadata,
    "price_id":    input.PriceID,
  }
  var resp ChargeResponse
  err := c.post("/v1/charges", body, input.IdempotencyKey, &resp)
  return &resp, err
}

func (c *Client) CreateCheckoutSession(input CheckoutSessionInput) (*CheckoutSessionResponse, error) {
  body := map[string]any{
    "customer_id":        input.CustomerID,
    "amount":             input.Amount,
    "currency":           input.Currency,
    "expires_in_seconds": input.ExpiresInSeconds,
    "success_url":        input.SuccessURL,
    "cancel_url":         input.CancelURL,
  }
  var resp struct {
    ID        string `json:"id"`
    Token     string `json:"token"`
    ExpiresAt string `json:"expires_at"`
  }
  if err := c.post("/v1/checkout-sessions", body, "", &resp); err != nil {
    return nil, err
  }
  return &CheckoutSessionResponse{ID: resp.ID, Token: resp.Token, ExpiresAt: resp.ExpiresAt}, nil
}

func (c *Client) PollChargeStatus(transactionID string, interval, timeout time.Duration) (*ChargeStatusResponse, error) {
  started := time.Now()
  for {
    if time.Since(started) > timeout {
      return nil, errors.New("poll_timeout")
    }
    var resp ChargeStatusResponse
    if err := c.get("/v1/charges/"+transactionID+"/status", &resp); err != nil {
      return nil, err
    }
    if resp.Status == "PAID" || resp.Status == "FAILED" {
      return &resp, nil
    }
    time.Sleep(interval)
  }
}

func (c *Client) post(path string, body map[string]any, idempotencyKey string, out any) error {
  payload, _ := json.Marshal(body)
  req, _ := http.NewRequest("POST", c.BaseURL+path, bytes.NewReader(payload))
  req.Header.Set("content-type", "application/json")
  req.Header.Set("x-api-key", c.APIKey)
  if idempotencyKey != "" {
    req.Header.Set("idempotency-key", idempotencyKey)
  }
  for k, v := range c.DefaultHeader {
    req.Header.Set(k, v)
  }
  res, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }
  defer res.Body.Close()
  if res.StatusCode >= 300 {
    var errPayload map[string]any
    _ = json.NewDecoder(res.Body).Decode(&errPayload)
    return fmt.Errorf("request_failed_%d", res.StatusCode)
  }
  return json.NewDecoder(res.Body).Decode(out)
}

func (c *Client) get(path string, out any) error {
  req, _ := http.NewRequest("GET", c.BaseURL+path, nil)
  req.Header.Set("x-api-key", c.APIKey)
  for k, v := range c.DefaultHeader {
    req.Header.Set(k, v)
  }
  res, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }
  defer res.Body.Close()
  if res.StatusCode >= 300 {
    return fmt.Errorf("request_failed_%d", res.StatusCode)
  }
  return json.NewDecoder(res.Body).Decode(out)
}

func trimSlash(base string) string {
  for len(base) > 0 && base[len(base)-1] == '/' {
    base = base[:len(base)-1]
  }
  return base
}
