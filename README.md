# Vektopay Go SDK

MVP: charges + checkout sessions + polling.

## Install

```bash
go get github.com/vektopay/vektopay-go-sdk
```

## Usage

```go
client := vektopay.NewClient(os.Getenv("VEKTOPAY_API_KEY"), "https://api.vektopay.com")

session, err := client.CreateCheckoutSession(vektopay.CheckoutSessionInput{
  CustomerID: "cust_123",
  Amount: 1000,
  Currency: "BRL",
})
```
