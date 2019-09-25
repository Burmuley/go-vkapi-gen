package responses

import (
	"encoding/json"
	"gitlab.com/Burmuley/go-vkapi/errors"
)

type RequestParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ApiRawResponse struct {
	Error         *errors.ApiError `json:"error"`
	RequestParams []RequestParams  `json:"request_params"`
	Response      json.RawMessage  `json:"response"`
}
