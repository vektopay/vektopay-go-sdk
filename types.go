package vektopay

type PaymentMethodInput struct {
  Type         string `json:"type"`
  Token        *string `json:"token,omitempty"`
  CardID       *string `json:"card_id,omitempty"`
  CvcToken     *string `json:"cvc_token,omitempty"`
  Installments *int    `json:"installments,omitempty"`
}

type PaymentCustomerInput struct {
  ExternalID string  `json:"external_id"`
  Name       *string `json:"name,omitempty"`
  Email      *string `json:"email,omitempty"`
  DocType    string  `json:"doc_type"`
  DocNumber  string  `json:"doc_number"`
}

type PaymentItemInput struct {
  PriceID  string `json:"price_id"`
  Quantity int    `json:"quantity"`
}

type PaymentInput struct {
  CustomerID    *string              `json:"customer_id,omitempty"`
  Customer      *PaymentCustomerInput `json:"customer,omitempty"`
  Items         []PaymentItemInput   `json:"items,omitempty"`
  Amount        *int                 `json:"amount,omitempty"`
  Currency      *string              `json:"currency,omitempty"`
  CouponCode    *string              `json:"coupon_code,omitempty"`
  Mode          *string              `json:"mode,omitempty"`
  WebhookURL    *string              `json:"webhook_url,omitempty"`
  PaymentMethod PaymentMethodInput   `json:"payment_method"`
}

type PaymentCreateResponse struct {
  PaymentID      string             `json:"payment_id"`
  Status         string             `json:"status"`
  PaymentStatus  *string            `json:"payment_status,omitempty"`
  SubscriptionID *string            `json:"subscription_id,omitempty"`
  Amount         *int               `json:"amount,omitempty"`
  Currency       *string            `json:"currency,omitempty"`
  Challenge      map[string]string  `json:"challenge,omitempty"`
}

type PaymentStatusResponse struct {
  ID     string `json:"id"`
  Status string `json:"status"`
}

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

type TransactionItemInput struct {
  PriceID  string `json:"price_id"`
  Quantity int    `json:"quantity"`
}

type TransactionPaymentMethodInput struct {
  Type         string `json:"type"`
  Token        string `json:"token"`
  Installments int    `json:"installments"`
}

type TransactionInput struct {
  CustomerID    string                        `json:"customer_id"`
  Items         []TransactionItemInput        `json:"items"`
  CouponCode    *string                       `json:"coupon_code,omitempty"`
  PaymentMethod TransactionPaymentMethodInput `json:"payment_method"`
}

type TransactionResponse struct {
  ID            string  `json:"id"`
  Status        string  `json:"status"`
  PaymentStatus *string `json:"paymentStatus,omitempty"`
  MerchantID    *string `json:"merchantId,omitempty"`
  Amount        *int    `json:"amount,omitempty"`
  Currency      *string `json:"currency,omitempty"`
}

type CustomerCreateInput struct {
  MerchantID string  `json:"merchant_id"`
  ExternalID string  `json:"external_id"`
  Name       *string `json:"name,omitempty"`
  Email      *string `json:"email,omitempty"`
  DocType    string  `json:"doc_type"`
  DocNumber  string  `json:"doc_number"`
}

type CustomerUpdateInput struct {
  MerchantID *string `json:"merchant_id,omitempty"`
  ExternalID *string `json:"external_id,omitempty"`
  Name       *string `json:"name,omitempty"`
  Email      *string `json:"email,omitempty"`
  DocType    *string `json:"doc_type,omitempty"`
  DocNumber  *string `json:"doc_number,omitempty"`
}

type CustomerCreateResponse struct {
  ID string `json:"id"`
}

type CustomerResponse struct {
  ID         string  `json:"id"`
  MerchantID *string `json:"merchantId,omitempty"`
  ExternalID *string `json:"externalId,omitempty"`
  Name       *string `json:"name,omitempty"`
  Email      *string `json:"email,omitempty"`
  DocType    *string `json:"docType,omitempty"`
  DocNumber  *string `json:"docNumber,omitempty"`
  CreatedAt  *string `json:"createdAt,omitempty"`
  UpdatedAt  *string `json:"updatedAt,omitempty"`
}

type CustomerListParams struct {
  MerchantID *string
  Limit      *int
  Offset     *int
}

type CheckoutSessionInput struct {
  CustomerID       string `json:"customer_id"`
  Amount           int    `json:"amount"`
  Currency         string `json:"currency"`
  PriceID          *string `json:"price_id,omitempty"`
  Quantity         *int    `json:"quantity,omitempty"`
  ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
  SuccessURL       *string `json:"success_url,omitempty"`
  CancelURL        *string `json:"cancel_url,omitempty"`
}

type CheckoutSessionResponse struct {
  ID        string `json:"id"`
  Token     string `json:"token"`
  ExpiresAt any    `json:"expires_at"`
}
