package vektopay

type ChargeInput struct {
  CustomerID     string                 `json:"customer_id"`
  CardID         string                 `json:"card_id"`
  Amount         int                    `json:"amount"`
  Currency       string                 `json:"currency"`
  Installments   *int                   `json:"installments,omitempty"`
  Country        *string                `json:"country,omitempty"`
  PriceID        *string                `json:"price_id,omitempty"`
  Metadata       map[string]any         `json:"metadata,omitempty"`
  IdempotencyKey string                 `json:"-"`
}

type ChargeResponse struct {
  ID        string            `json:"id"`
  Status    string            `json:"status"`
  Error     map[string]any    `json:"error,omitempty"`
  Challenge map[string]string `json:"challenge,omitempty"`
}

type ChargeStatusResponse struct {
  ID     string `json:"id"`
  Status string `json:"status"`
}

type CheckoutSessionInput struct {
  CustomerID       string `json:"customer_id"`
  Amount           int    `json:"amount"`
  Currency         string `json:"currency"`
  ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
  SuccessURL       *string `json:"success_url,omitempty"`
  CancelURL        *string `json:"cancel_url,omitempty"`
}

type CheckoutSessionResponse struct {
  ID        string `json:"id"`
  Token     string `json:"token"`
  ExpiresAt string `json:"expires_at"`
}
