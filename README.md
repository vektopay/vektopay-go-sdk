# Vektopay Go SDK

Go SDK for Vektopay API (server-side). Supports payments, checkout sessions, and payment status polling.

## Install

```bash
go get github.com/vektopay/vektopay-go-sdk
```

## Setup

```go
client := vektopay.NewClient(os.Getenv("VEKTOPAY_API_KEY"), "https://api.vektopay.com")
```

## Create Payment (Recommended)

```go
payment, err := client.CreatePayment(vektopay.PaymentInput{
  CustomerID: func() *string { v := "cust_123"; return &v }(),
  Items: []vektopay.PaymentItemInput{
    { PriceID: "price_basic", Quantity: 1 },
  },
  CouponCode: func() *string { v := "OFF10"; return &v }(),
  PaymentMethod: vektopay.PaymentMethodInput{
    Type: "credit_card",
    Token: func() *string { v := "ev:tk_123"; return &v }(),
    Installments: func() *int { v := 1; return &v }(),
  },
})
```

## Create Customer

Customers must exist before creating transactions or charges.

```go
customer, err := client.CreateCustomer(vektopay.CustomerCreateInput{
  MerchantID: "mrc_123",
  ExternalID: "cust_ext_123",
  Name: func() *string { v := "Ana Silva"; return &v }(),
  Email: func() *string { v := "ana@example.com"; return &v }(),
  DocType: "CPF",
  DocNumber: "12345678901",
})
```

## Update Customer

```go
updated, err := client.UpdateCustomer("cust_123", vektopay.CustomerUpdateInput{
  Name: func() *string { v := "Ana Maria Silva"; return &v }(),
  Email: func() *string { v := "ana.maria@example.com"; return &v }(),
})
```

## Get Customer

```go
detail, err := client.GetCustomer("cust_123")
```

## List Customers

```go
customers, err := client.ListCustomers(vektopay.CustomerListParams{
  MerchantID: func() *string { v := "mrc_123"; return &v }(),
  Limit: func() *int { v := 50; return &v }(),
  Offset: func() *int { v := 0; return &v }(),
})
```

## Delete Customer

```go
err := client.DeleteCustomer("cust_123")
```

## Create Charge (Card)

```go
charge, err := client.CreateCharge(vektopay.ChargeInput{
  CustomerID: "cust_123",
  CardID: "card_123",
  Amount: 1000,
  Currency: "BRL",
})
```

## Create Checkout Session (Frontend)

```go
session, err := client.CreateCheckoutSession(vektopay.CheckoutSessionInput{
  CustomerID: "cust_123",
  Amount: 1000,
  Currency: "BRL",
})
```

## Poll Payment Status

```go
status, err := client.PollPaymentStatus(payment.PaymentID, 3*time.Second, 2*time.Minute)
```

## Notes
- Never expose your API key in the browser.
