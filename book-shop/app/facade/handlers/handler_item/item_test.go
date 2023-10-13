package handler_item

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/biz-demo/book-shop/app/facade/model"
	"github.com/cloudwego/biz-demo/book-shop/kitex_gen/cwg/bookshop/item"
	"github.com/issue9/assert"
)

const host = "http://localhost:8080"

func getToken() string {
	// get the token
	loginBody := map[string]string{
		"password": "emate",
		"username": "emate",
	}
	loginBytes, _ := json.Marshal(loginBody)
	bodyReader := bytes.NewReader(loginBytes)
	req, _ := http.NewRequest("POST", host+"/user/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	respData, _ := io.ReadAll(resp.Body)
	respMap := make(map[string]interface{})
	json.Unmarshal(respData, &respMap)
	token := respMap["token"].(string)
	return token
}

func TestEditProductPic(t *testing.T) {
	httpClient := http.DefaultClient
	token := fmt.Sprintf("Bearer %s", getToken())

	// add a product
	nowTime := time.Now()
	req := model.AddProductRequest{
		Name: fmt.Sprintf("%d", nowTime.UnixMilli()),
		Pic:  "default",
	}
	reqBytes, _ := json.Marshal(req)
	bodyReader := bytes.NewReader(reqBytes)
	httpReq, _ := http.NewRequest("POST", host+"/item2b/add", bodyReader)
	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-Type", "application/json")
	httpResp, _ := httpClient.Do(httpReq)
	httpRespBytes, _ := io.ReadAll(httpResp.Body)
	itemAddResp := model.Response{}
	json.Unmarshal(httpRespBytes, &itemAddResp)
	resp := itemAddResp.Data.(map[string]interface{})
	productID := resp["product_id"].(string)

	// update pic field
	newPic := "new pic"
	editReq := model.EditProductRequest{
		ProductId: productID,
		Pic:       &newPic,
	}
	reqBytes, _ = json.Marshal(editReq)
	bodyReader = bytes.NewReader(reqBytes)
	httpReq, _ = http.NewRequest("POST", host+"/item2b/edit", bodyReader)
	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-Type", "application/json")
	httpResp, _ = httpClient.Do(httpReq)
	httpRespBytes, _ = io.ReadAll(httpResp.Body)
	itemAddResp = model.Response{}
	json.Unmarshal(httpRespBytes, &itemAddResp)
	assert.Equal(t, 0, itemAddResp.Code)

	// check if the pic is updated
	httpReq, _ = http.NewRequest("GET", host+"/item2b/get?product_id="+productID, nil)
	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-Type", "application/json")
	httpResp, _ = httpClient.Do(httpReq)
	httpRespBytes, _ = io.ReadAll(httpResp.Body)
	itemGetResp := model.Response{}
	json.Unmarshal(httpRespBytes, &itemGetResp)
	getBytes, _ := json.Marshal(itemGetResp.Data)
	product := item.Product{}
	json.Unmarshal(getBytes, &product)
	assert.Equal(t, 0, itemAddResp.Code)
	assert.Equal(t, newPic, product.Pic)

	// delete the product for clean env
	delReq := model.OperateProductReq{
		ProductId: productID,
	}
	reqBytes, _ = json.Marshal(delReq)
	bodyReader = bytes.NewReader(reqBytes)
	httpReq, _ = http.NewRequest("POST", host+"/item2b/del", bodyReader)
	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-Type", "application/json")
	httpClient.Do(httpReq)
}
