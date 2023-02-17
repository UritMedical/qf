package action

import (
	"errors"
	"net/http"
	"qf/helper/content"
)

func NewActionByWebApi() *ByWebApi {
	return &ByWebApi{
		client: &http.Client{},
	}
}

type ByWebApi struct {
	client *http.Client
}

func (b *ByWebApi) Call(bllId string, method, relative string, content content.Content) (string, error) {
	//if method == "Get" {
	//	url := fmt.Sprintf("http://localhost/api/%s/%s", bllId, relative)
	//	if params != "" {
	//		p := map[string]interface{}{}
	//		e := json.Unmarshal([]byte(params), &p)
	//		if e == nil {
	//			query := ""
	//			for k, v := range p {
	//				query += fmt.Sprintf("%s=%s", k, v)
	//			}
	//			url += "?" + query
	//		}
	//	}
	//	rep, err := b.client.Get(url)
	//	if err == nil {
	//		data, err := ioutil.ReadAll(rep.Body)
	//		rep.Body.Close()
	//		return string(data), err
	//	}
	//	return "", err
	//} else if method == "Post" {
	//
	//} else if method == "Put" {
	//
	//} else if method == "Delete" {
	//
	//}
	return "", errors.New("not find")
}
