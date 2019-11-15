/*
Copyright 2019 Konstantin Vasilev (burmuley@gmail.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package go_vkapi

import (
	"encoding/json"
	"fmt"
	"github.com/Burmuley/go-vkapi/responses"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	apiVersion = "5.101"
	apiUrl     = "https://api.vk.com/method/"
)

type VKApi struct {
	userToken  string
	apiVersion string
	apiUrl     string
}

// SendAPIRequest calls defined method of the VK API with the defined parameters
// Returns slice of bytes with API response
func (vk *VKApi) SendAPIRequest(method string, parameters map[string]interface{}) ([]byte, error) {
	//Format API endpoint
	u, err := url.Parse(apiUrl + method)

	if err != nil {
		return nil, err
	}

	//Fill mandatory parameters
	parameters["access_token"] = vk.userToken
	parameters["v"] = apiVersion

	// Format URL-encoded key-value parameters
	request := url.Values{}
	for k, v := range parameters {
		switch v.(type) {
		case string:
			request.Add(k, v.(string))
		default:
			request.Add(k, fmt.Sprint(v))
		}
	}

	// Send request and read response
	resp, err := http.PostForm(u.String(), request)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	//Read response body and check for errors
	rBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	var apiResp responses.ApiRawResponse

	if err := json.Unmarshal(rBody, &apiResp); err != nil {
		return []byte{}, err
	}

	if apiResp.Error.GetCode() != 0 {
		return []byte{}, apiResp.Error
	}

	return apiResp.Response, nil
}

func (vk *VKApi) SendObjRequest(method string, params map[string]interface{}, object interface{}) error {
	info, err := vk.SendAPIRequest(method, params)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(info, &object); err != nil {
		return err
	}

	return nil
}

func NewApiWithToken(token string) *VKApi {
	envApiUrl := os.Getenv("VK_API_URL")
	locApiUrl := apiUrl

	if len(envApiUrl) > 0 {
		locApiUrl = envApiUrl
	}

	return &VKApi{userToken: token,
		apiVersion: apiVersion,
		apiUrl:     locApiUrl,
	}
}
