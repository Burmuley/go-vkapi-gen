package go_vkapi

type VK interface {
	SendAPIRequest(method string, parameters map[string]string) ([]byte, error)
	SendObjRequest(method string, params map[string]string, object interface{}) error
}
