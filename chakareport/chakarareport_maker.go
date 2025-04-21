package chakareport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ChakaraReportMaker struct {
	apiUrl string
}

func NewChakaraReportMaker(apiUrl string) (Maker, error) {
	if apiUrl == "" {
		return nil, fmt.Errorf("API URL cannot be empty for ChakaraReportMaker")
	}
	maker := &ChakaraReportMaker{
		apiUrl: apiUrl,
	}
	return maker, nil
}

func (maker *ChakaraReportMaker) GenerateChakaraReport(chakaraInfo []ChakraInfo, language string) ([]byte, error) {
	fmt.Printf("current language: %+v\n", language)

	// 将 chakraInfo 和 language 转换为 JSON
	jsonData, err := json.Marshal(map[string]interface{}{
		"chakra_info": chakaraInfo,
		"language":    language,
	})
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return []byte("Error generating report about json.Marshal"), nil
	}

	// 发送 POST 请求到外部 API
	resp, err := http.Post(maker.apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return []byte("Error generating report about http.Post"), nil
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return []byte("Error generating report about ioutil.ReadAll"), nil
	}
	fmt.Println("body:", body)

	// 直接返回响应内容
	return body, nil
}
