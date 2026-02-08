package vektopay

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "net/url"
  "time"
)

type Client struct {
  APIKey        string
  BaseURL       string
  DefaultHeader map[string]string
  BearerToken   string
  httpClient    *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
  return &Client{
    APIKey:        apiKey,
    BaseURL:       trimSlash(baseURL),
    DefaultHeader: map[string]string{},
    BearerToken:   "",
    httpClient:    &http.Client{Timeout: 30 * time.Second},
  }
}

func (c *Client) CreatePayment(input PaymentInput) (*PaymentCreateResponse, error) {
  var resp PaymentCreateResponse
  err := c.post("/v1/payments", input, "", &resp)
  return &resp, err
}

func (c *Client) GetPaymentStatus(id string) (*PaymentStatusResponse, error) {
  var resp PaymentStatusResponse
  err := c.get("/v1/payments/"+id+"/status", &resp)
  return &resp, err
}

func (c *Client) PollPaymentStatus(paymentID string, interval, timeout time.Duration) (*PaymentStatusResponse, error) {
  started := time.Now()
  for {
    if time.Since(started) > timeout {
      return nil, errors.New("poll_timeout")
    }
    status, err := c.GetPaymentStatus(paymentID)
    if err != nil {
      return nil, err
    }
    if status.Status == "PAID" || status.Status == "FAILED" || status.Status == "CANCELED" {
      return status, nil
    }
    time.Sleep(interval)
  }
}

func (c *Client) CreateCharge(input ChargeInput) (*ChargeResponse, error) {
  // Legacy alias: `/v1/charges` is deprecated; map to `/v1/payments`.
  payment, err := c.CreatePayment(PaymentInput{
    CustomerID: &input.CustomerID,
    Amount:     &input.Amount,
    Currency:   &input.Currency,
    PaymentMethod: PaymentMethodInput{
      Type:         "credit_card",
      CardID:       &input.CardID,
      Installments: input.Installments,
    },
  })
  if err != nil {
    return nil, err
  }
  return &ChargeResponse{
    ID:     payment.PaymentID,
    Status: payment.Status,
    Challenge: payment.Challenge,
  }, nil
}

func (c *Client) CreateTransaction(input TransactionInput) (*TransactionResponse, error) {
  // Legacy alias: `/v1/transactions` is deprecated; map to `/v1/payments`.
  items := make([]PaymentItemInput, 0, len(input.Items))
  for _, it := range input.Items {
    items = append(items, PaymentItemInput{PriceID: it.PriceID, Quantity: it.Quantity})
  }
  payment, err := c.CreatePayment(PaymentInput{
    CustomerID: &input.CustomerID,
    Items:      items,
    CouponCode: input.CouponCode,
    PaymentMethod: PaymentMethodInput{
      Type:         input.PaymentMethod.Type,
      Token:        &input.PaymentMethod.Token,
      Installments: &input.PaymentMethod.Installments,
    },
  })
  if err != nil {
    return nil, err
  }
  return &TransactionResponse{
    ID:            payment.PaymentID,
    Status:        payment.Status,
    PaymentStatus: payment.PaymentStatus,
    Amount:        payment.Amount,
    Currency:      payment.Currency,
  }, nil
}

func (c *Client) CreateCustomer(input CustomerCreateInput) (*CustomerCreateResponse, error) {
  if c.BearerToken == "" {
    return nil, errors.New("bearer_token_required")
  }
  body := map[string]any{
    "merchant_id": input.MerchantID,
    "external_id": input.ExternalID,
    "name":        input.Name,
    "email":       input.Email,
    "doc_type":    input.DocType,
    "doc_number":  input.DocNumber,
  }
  var resp CustomerCreateResponse
  err := c.postBearer("/v1/customers", body, &resp)
  return &resp, err
}

func (c *Client) UpdateCustomer(id string, input CustomerUpdateInput) (*CustomerResponse, error) {
  if c.BearerToken == "" {
    return nil, errors.New("bearer_token_required")
  }
  body := map[string]any{}
  if input.MerchantID != nil {
    body["merchant_id"] = *input.MerchantID
  }
  if input.ExternalID != nil {
    body["external_id"] = *input.ExternalID
  }
  if input.Name != nil {
    body["name"] = *input.Name
  }
  if input.Email != nil {
    body["email"] = *input.Email
  }
  if input.DocType != nil {
    body["doc_type"] = *input.DocType
  }
  if input.DocNumber != nil {
    body["doc_number"] = *input.DocNumber
  }
  var resp CustomerResponse
  err := c.putBearer("/v1/customers/"+id, body, &resp)
  return &resp, err
}

func (c *Client) ListCustomers(params CustomerListParams) (*[]CustomerResponse, error) {
  if c.BearerToken == "" {
    return nil, errors.New("bearer_token_required")
  }
  query := buildQuery(map[string]string{
    "merchant_id": valueOrEmpty(params.MerchantID),
    "limit":       intOrEmpty(params.Limit),
    "offset":      intOrEmpty(params.Offset),
  })
  var resp []CustomerResponse
  err := c.getBearer("/v1/customers"+query, &resp)
  return &resp, err
}

func (c *Client) GetCustomer(id string) (*CustomerResponse, error) {
  if c.BearerToken == "" {
    return nil, errors.New("bearer_token_required")
  }
  var resp CustomerResponse
  err := c.getBearer("/v1/customers/"+id, &resp)
  return &resp, err
}

func (c *Client) DeleteCustomer(id string) error {
  if c.BearerToken == "" {
    return errors.New("bearer_token_required")
  }
  return c.delBearer("/v1/customers/" + id)
}

func (c *Client) CreateCheckoutSession(input CheckoutSessionInput) (*CheckoutSessionResponse, error) {
  body := map[string]any{
    "customer_id":        input.CustomerID,
    "amount":             input.Amount,
    "currency":           input.Currency,
    "price_id":           input.PriceID,
    "quantity":           input.Quantity,
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
  status, err := c.PollPaymentStatus(transactionID, interval, timeout)
  if err != nil {
    return nil, err
  }
  return &ChargeStatusResponse{ID: status.ID, Status: status.Status}, nil
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

func (c *Client) postBearer(path string, body map[string]any, out any) error {
  payload, _ := json.Marshal(body)
  req, _ := http.NewRequest("POST", c.BaseURL+path, bytes.NewReader(payload))
  req.Header.Set("content-type", "application/json")
  req.Header.Set("authorization", "Bearer "+c.BearerToken)
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

func (c *Client) put(path string, body map[string]any, out any) error {
  payload, _ := json.Marshal(body)
  req, _ := http.NewRequest("PUT", c.BaseURL+path, bytes.NewReader(payload))
  req.Header.Set("content-type", "application/json")
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
    var errPayload map[string]any
    _ = json.NewDecoder(res.Body).Decode(&errPayload)
    return fmt.Errorf("request_failed_%d", res.StatusCode)
  }
  return json.NewDecoder(res.Body).Decode(out)
}

func (c *Client) putBearer(path string, body map[string]any, out any) error {
  payload, _ := json.Marshal(body)
  req, _ := http.NewRequest("PUT", c.BaseURL+path, bytes.NewReader(payload))
  req.Header.Set("content-type", "application/json")
  req.Header.Set("authorization", "Bearer "+c.BearerToken)
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

func (c *Client) del(path string) error {
  req, _ := http.NewRequest("DELETE", c.BaseURL+path, nil)
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
  return nil
}

func (c *Client) delBearer(path string) error {
  req, _ := http.NewRequest("DELETE", c.BaseURL+path, nil)
  req.Header.Set("authorization", "Bearer "+c.BearerToken)
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
  return nil
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

func (c *Client) getBearer(path string, out any) error {
  req, _ := http.NewRequest("GET", c.BaseURL+path, nil)
  req.Header.Set("authorization", "Bearer "+c.BearerToken)
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

func buildQuery(values map[string]string) string {
  query := ""
  first := true
  for key, value := range values {
    if value == "" {
      continue
    }
    if first {
      query += "?"
      first = false
    } else {
      query += "&"
    }
    query += key + "=" + url.QueryEscape(value)
  }
  return query
}

func valueOrEmpty(value *string) string {
  if value == nil {
    return ""
  }
  return *value
}

func intOrEmpty(value *int) string {
  if value == nil {
    return ""
  }
  return fmt.Sprintf("%d", *value)
}
