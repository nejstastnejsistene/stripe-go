package charge

import (
	"testing"

	. "github.com/stripe/stripe-go"
	. "github.com/stripe/stripe-go/utils"
)

func init() {
	Key = GetTestKey()
}

func TestChargeCreate(t *testing.T) {
	chargeParams := &ChargeParams{
		Amount:   1000,
		Currency: USD,
		Card: &CardParams{
			Name:   "Stripe Tester",
			Number: "378282246310005",
			Month:  "06",
			Year:   "20",
		},
		Statement: "statement",
		Email:     "a@b.com",
	}

	target, err := Create(chargeParams)

	if err != nil {
		t.Error(err)
	}

	if target.Amount != chargeParams.Amount {
		t.Errorf("Amount %v does not match expected amount %v\n", target.Amount, chargeParams.Amount)
	}

	if target.Currency != chargeParams.Currency {
		t.Errorf("Currency %q does not match expected currency %q\n", target.Currency, chargeParams.Currency)
	}

	if target.Card.Name != chargeParams.Card.Name {
		t.Errorf("Card name %q does not match expected name %q\n", target.Card.Name, chargeParams.Card.Name)
	}

	if target.Statement != chargeParams.Statement {
		t.Errorf("Statement description %q does not match expected description %v\n", target.Statement, chargeParams.Statement)
	}

	if target.Email != chargeParams.Email {
		t.Errorf("Email %q does not match expected email %v\n", target.Email, chargeParams.Email)
	}
}

func TestChargeGet(t *testing.T) {
	chargeParams := &ChargeParams{
		Amount:   1001,
		Currency: USD,
		Card: &CardParams{
			Number: "378282246310005",
			Month:  "06",
			Year:   "20",
		},
	}

	res, _ := Create(chargeParams)

	target, err := Get(res.Id, nil)

	if err != nil {
		t.Error(err)
	}

	if target.Id != res.Id {
		t.Errorf("Charge id %q does not match expected id %q\n", target.Id, res.Id)
	}
}

func TestChargeUpdate(t *testing.T) {
	chargeParams := &ChargeParams{
		Amount:   1002,
		Currency: USD,
		Card: &CardParams{
			Number: "378282246310005",
			Month:  "06",
			Year:   "20",
		},
		Desc: "original description",
	}

	res, _ := Create(chargeParams)

	if res.Desc != chargeParams.Desc {
		t.Errorf("Original description %q does not match expected description %q\n", res.Desc, chargeParams.Desc)
	}

	updated := &ChargeParams{
		Desc: "updated description",
	}

	target, err := Update(res.Id, updated)

	if err != nil {
		t.Error(err)
	}

	if target.Desc != updated.Desc {
		t.Errorf("Updated description %q does not match expected description %q\n", target.Desc, updated.Desc)
	}
}

func TestChargeRefund(t *testing.T) {
	chargeParams := &ChargeParams{
		Amount:   1003,
		Currency: USD,
		Card: &CardParams{
			Number: "378282246310005",
			Month:  "06",
			Year:   "20",
		},
	}

	res, _ := Create(chargeParams)

	// full refund
	ref, err := RefundCharge(&RefundParams{Charge: res.Id})

	if err != nil {
		t.Error(err)
	}

	if ref.Charge != res.Id {
		t.Errorf("Refund charge %q does not match expected value %v\n", ref.Charge, res.Id)
	}

	target, _ := Get(res.Id, nil)

	if !target.Refunded || target.Refunds == nil {
		t.Errorf("Expected to have refunded this charge\n")
	}

	if len(target.Refunds.Values) != 1 {
		t.Errorf("Expected to have a refund, but instead have %v\n", len(target.Refunds.Values))
	}

	if target.Refunds.Values[0].Amount != target.AmountRefunded {
		t.Errorf("Refunded amount %v does not match amount refunded %v\n", target.Refunds.Values[0].Amount, target.AmountRefunded)
	}

	if target.Refunds.Values[0].Currency != target.Currency {
		t.Errorf("Refunded currency %q does not match charge currency %q\n", target.Refunds.Values[0].Currency, target.Currency)
	}

	if len(target.Refunds.Values[0].Tx.Id) == 0 {
		t.Errorf("Refund transaction not set\n")
	}

	if target.Refunds.Values[0].Charge != target.Id {
		t.Errorf("Refund charge %q does not match expected value %v\n", target.Refunds.Values[0].Charge, target.Id)
	}

	res, err = Create(chargeParams)

	// partial refund
	refund := &RefundParams{
		Charge: res.Id,
		Amount: 253,
	}

	RefundCharge(refund)

	target, _ = Get(res.Id, nil)

	if target.Refunded {
		t.Errorf("Partial refund should not be marked as Refunded\n")
	}

	if target.AmountRefunded != refund.Amount {
		t.Errorf("Refunded amount %v does not match expected amount %v\n", target.AmountRefunded, refund.Amount)
	}
}

func TestChargeCapture(t *testing.T) {
	chargeParams := &ChargeParams{
		Amount:   1004,
		Currency: USD,
		Card: &CardParams{
			Number: "378282246310005",
			Month:  "06",
			Year:   "20",
		},
		NoCapture: true,
	}

	res, _ := Create(chargeParams)

	if res.Captured {
		t.Errorf("The charge should not have been captured\n")
	}

	// full capture
	target, err := Capture(res.Id, nil)

	if err != nil {
		t.Error(err)
	}

	if !target.Captured {
		t.Errorf("Expected to have captured this charge after full capture\n")
	}

	res, err = Create(chargeParams)

	// partial capture
	capture := &CaptureParams{
		Amount: 554,
		Email:  "a@b.com",
	}

	target, err = Capture(res.Id, capture)

	if err != nil {
		t.Error(err)
	}

	if !target.Captured {
		t.Errorf("Expected to have captured this charge after partial capture\n")
	}

	if target.AmountRefunded != chargeParams.Amount-capture.Amount {
		t.Errorf("Refunded amount %v does not match expected amount %v\n", target.AmountRefunded, chargeParams.Amount-capture.Amount)
	}

	if target.Email != capture.Email {
		t.Errorf("Email %q does not match expected email %v\n", target.Email, capture.Email)
	}
}

func TestChargeList(t *testing.T) {
	params := &ChargeListParams{}
	params.Filters.AddFilter("include[]", "", "total_count")
	params.Filters.AddFilter("limit", "", "5")
	params.Single = true

	i := List(params)
	for !i.Stop() {
		target, err := i.Next()

		if err != nil {
			t.Error(err)
		}

		if target == nil {
			t.Error("No nil values expected")
		}

		if i.Meta() == nil {
			t.Error("No metadata returned")
		}
	}
}
