// package refund provides the /refunds APIs
package refund

import (
	"fmt"
	"net/url"

	. "github.com/stripe/stripe-go"
)

// Client is used to invoke /refunds APIs.
type Client struct {
	B   Backend
	Tok string
}

// Get returns the details of a refund.
// For more details see https://stripe.com/docs/api#retrieve_refund.
func Get(id string, params *RefundParams) (*Refund, error) {
	return getC().Get(id, params)
}

func (c Client) Get(id string, params *RefundParams) (*Refund, error) {
	body := &url.Values{}

	params.AppendTo(body)

	refund := &Refund{}
	err := c.B.Call("GET", fmt.Sprintf("/charges/%v/refunds/%v", params.Charge, id), c.Tok, body, refund)

	return refund, err
}

// Update updates a refund's properties.
// For more details see https://stripe.com/docs/api#update_refund.
func Update(id string, params *RefundParams) (*Refund, error) {
	return getC().Update(id, params)
}

func (c Client) Update(id string, params *RefundParams) (*Refund, error) {
	body := &url.Values{}

	params.AppendTo(body)

	refund := &Refund{}
	err := c.B.Call("POST", fmt.Sprintf("/charges/%v/refunds/%v", params.Charge, id), c.Tok, body, refund)

	return refund, err
}

// List returns a list of refunds.
// For more details see https://stripe.com/docs/api#list_refunds.
func List(params *RefundListParams) *RefundIter {
	return getC().List(params)
}

func (c Client) List(params *RefundListParams) *RefundIter {
	body := &url.Values{}
	var lp *ListParams

	params.AppendTo(body)
	lp = &params.ListParams

	return &RefundIter{GetIter(lp, body, func(b url.Values) ([]interface{}, ListMeta, error) {
		list := &RefundList{}
		err := c.B.Call("GET", fmt.Sprintf("/charges/%v/refunds", params.Charge), c.Tok, &b, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

func getC() Client {
	return Client{GetBackend(), Key}
}
