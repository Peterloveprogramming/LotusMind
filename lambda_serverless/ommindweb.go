package lambdaServerless

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/lotusMind/meditation/chakareport"
	db "github.com/lotusMind/meditation/db/sqlc"
)

func (lambdaServerless *Lambda) Test(ctx context.Context, event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	fmt.Println("context", ctx)
	fmt.Println("event", event)
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Test is working!",
	}
	return response
}

type createUserEmailRequestBody struct {
	Email    string            `json:"email"`
	Answers  map[string]string `json:"answers"`
	Language string            `json:"language"`
	IP       string            `json:"ip"`
	Country  string            `json:"country"`
}

// type registEmailResponse struct {
// }
func buildBraceletURL(baseURL string, chakras []string) string {
	// Join chakras with comma
	chakraQuery := strings.Join(chakras, ",")
	// Escape query string
	escapedChakras := url.QueryEscape(chakraQuery)
	// Format final URL
	return fmt.Sprintf("%s?chakras=%s", baseURL, escapedChakras)
}

func (lambdaServerless *Lambda) RegisterEmail(ctx context.Context, event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// resp := registEmailResponse{}
	var req createUserEmailRequestBody

	// Parse JSON body
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "invalid JSON",
		}
	}

	// Manual Validation
	if req.Email == "" || len(req.Email) > 30 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "email is required and must be <= 30 characters",
		}
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "invalid email format",
		}
	}
	if len(req.Answers) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "answers are required",
		}
	}
	for key, val := range req.Answers {
		if val == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf("answer for '%s' cannot be empty", key),
			}
		}
	}
	if len(req.Language) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "language is required",
		}
	}
	if len(req.IP) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "ip is required",
		}
	}
	if len(req.Country) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "country is required",
		}
	}
	// 创建问题分数映射
	fmt.Printf("req.Answers: %+v\n", req.Answers)
	questionAnswers := make(map[string]string)

	uniqueCode := generateUniqueCode()

	for question, answer := range req.Answers {
		questionAnswers[question] = answer
	}
	// save answers
	// err := lambdaServerless.storageMaker.SaveChakaraReportAnswersAsText(req.Email, uniqueCode, req.Answers)

	saveAnswersUrl := lambdaServerless.config.ApiGateWayEndpoint + "/save-chakra-report-answers"
	fmt.Println("url is", saveAnswersUrl)

	data := map[string]interface{}{
		"email":    req.Email,
		"uniqueId": uniqueCode,
		"answers":  req.Answers, // assign map[string]string here
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create request with body
	saveAnsweReq, err := http.NewRequest("POST", saveAnswersUrl, bytes.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	saveAnsweReq.Header.Add("x-api-key", lambdaServerless.config.ApiGateWayApiKey)
	fmt.Println("api key is", lambdaServerless.config.ApiGateWayApiKey)

	_, err = http.DefaultClient.Do(saveAnsweReq)
	if err != nil {
		panic(err)
	}

	if err != nil {
		fmt.Println("error in saving answers", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "error in saving answers",
		}
	}

	// 创建一个切片用于存放所有 value
	values := make([]string, 0, len(questionAnswers))

	// 遍历 map，将每个 value 添加到切片中
	for _, value := range questionAnswers {
		values = append(values, value)
	}
	// 打印结果
	fmt.Println("Values:", values)
	fmt.Printf("Email: %s\n", req.Email)
	// fmt.Println("Question Scores Map:")

	var rootScores, sacralScores, solarPlexusScores, heartScores, throatScores, thirdEyeScores, crownScores float32

	for _, answer := range values {

		// fmt.Printf("index: %d,answer: %s\n", i, answer)

		parts := strings.Split(answer, "||")

		setIndex, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("转换失败:", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("转换失败: %v", err),
			}
		}
		score, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("转换失败:", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("转换失败: %v", err),
			}
		}

		switch setIndex {
		case 1:
			rootScores += float32(score)
		case 2:
			sacralScores += float32(score)
		case 3:
			solarPlexusScores += float32(score)
		case 4:
			heartScores += float32(score)
		case 5:
			throatScores += float32(score)
		case 6:
			thirdEyeScores += float32(score)
		case 7:
			crownScores += float32(score)
		}
	}
	rootScores = rootScores * 100 / (9 * 2)
	sacralScores = sacralScores * 100 / (8 * 2)
	solarPlexusScores = solarPlexusScores * 100 / (9 * 2)
	heartScores = heartScores * 100 / (9 * 2)
	throatScores = throatScores * 100 / (8 * 2)
	thirdEyeScores = thirdEyeScores * 100 / (9 * 2)
	crownScores = crownScores * 100 / (7 * 2)

	fmt.Printf("scoreRootScores: %f\n", rootScores)
	fmt.Printf("sacralScores : %f\n", sacralScores)
	fmt.Printf("solarPlexusScores: %f\n", solarPlexusScores)
	fmt.Printf("heartScores: %f\n", heartScores)
	fmt.Printf("throatScores: %f\n", throatScores)
	fmt.Printf("thirdEyeScores: %f\n", thirdEyeScores)
	fmt.Printf("CrownChakra Score: %f\n", crownScores)

	// 定义脉轮数据
	chakras := []struct {
		ChakraName   string `json:"chakra_name"`
		ChakraScore  int32  `json:"chakra_score"`
		ChakraStatus string `json:"chakra_status"`
	}{
		{"Root Chakra", int32(rootScores), getChakraStatus(rootScores)},
		{"Sacral Chakra", int32(sacralScores), getChakraStatus(sacralScores)},
		{"Solar Plexus Chakra", int32(solarPlexusScores), getChakraStatus(solarPlexusScores)},
		{"Heart Chakra", int32(heartScores), getChakraStatus(heartScores)},
		{"Throat Chakra", int32(throatScores), getChakraStatus(throatScores)},
		{"Third Eye Chakra", int32(thirdEyeScores), getChakraStatus(thirdEyeScores)},
		{"Crown Chakra", int32(crownScores), getChakraStatus(crownScores)},
	}

	// 将脉轮数据转换为 JSON
	chakraInfoJSON, err := json.Marshal(chakras)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("将脉轮数据转换为 JSON 失败: %v", err),
		}
	}

	emailRegistrationData := map[string]interface{}{
		"email":       req.Email,
		"chakra_info": chakraInfoJSON,
		"language":    req.Language,
		"unique_code": uniqueCode,
		"ip":          req.IP,
		"country":     req.Country,
	}

	// Marshal to JSON
	emailRegistrationjsonData, err := json.Marshal(emailRegistrationData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// save email registration
	url := lambdaServerless.config.ApiGateWayEndpoint + "/email-registration"
	fmt.Println("url is", url)

	emailRegistrationReq, _ := http.NewRequest("POST", url, bytes.NewReader(emailRegistrationjsonData))
	emailRegistrationReq.Header.Add("x-api-key", lambdaServerless.config.ApiGateWayApiKey)
	fmt.Println("api key is", lambdaServerless.config.ApiGateWayApiKey)

	_, err = http.DefaultClient.Do(emailRegistrationReq)
	if err != nil {
		panic(err)
	}

	// 创建用户注册邮箱
	// args := db.CreateUserEmailTransactiontArgs{
	// 	Email:      req.Email,
	// 	ChakraInfo: string(chakraInfoJSON),
	// 	Language:   req.Language,
	// 	UniqueCode: uniqueCode,
	// 	IP:         req.IP,
	// 	Country:    req.Country,
	// }

	// registration, err := lambdaServerless.store.CreateUserEmailTransaction(ctx, args)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "unique_violation") {
	// 		return events.APIGatewayProxyResponse{
	// 			StatusCode: 400,
	// 			Body:       "email exists already",
	// 		}
	// 	}
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 500,
	// 		Body:       "something went wrong while attempting registering email",
	// 	}
	// }

	//获取推荐的手串
	// 创建一个包含脉轮分数的切片
	chakraScores := []struct {
		name  string
		score float32
	}{
		{"Root Chakra", rootScores},
		{"Sacral Chakra", sacralScores},
		{"Solar Plexus Chakra", solarPlexusScores},
		{"Heart Chakra", heartScores},
		{"Throat Chakra", throatScores},
		{"Third Eye Chakra", thirdEyeScores},
		{"Crown Chakra", crownScores},
	}

	// 根据分数排序（从低到高）
	for i := 0; i < len(chakraScores)-1; i++ {
		for j := i + 1; j < len(chakraScores); j++ {
			if chakraScores[i].score > chakraScores[j].score {
				chakraScores[i], chakraScores[j] = chakraScores[j], chakraScores[i]
			}
		}
	}
	fmt.Println(("wtf???"))
	// 获取分数最低的两个脉轮
	lowestChakras := []string{chakraScores[0].name, chakraScores[1].name}
	fmt.Printf("分数最低的两个脉轮是: %v\n", lowestChakras)
	braceletApiUrl := buildBraceletURL(lambdaServerless.config.ApiGateWayEndpoint+"/bracelets", lowestChakras)
	fmt.Println(braceletApiUrl)

	fmt.Println("braceletApiUrl is", braceletApiUrl)
	bracetletGetRequest, _ := http.NewRequest("GET", braceletApiUrl, nil)
	bracetletGetRequest.Header.Add("x-api-key", lambdaServerless.config.ApiGateWayApiKey)
	resp, err := http.DefaultClient.Do(bracetletGetRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("error calling bracelet API: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("bracelet API returned status %d: %s", resp.StatusCode, string(bodyBytes)),
		}
	}

	var braceletAPIResponse struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&braceletAPIResponse); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("failed to decode bracelet API response: %v", err),
		}
	}
	bracelets := braceletAPIResponse.Data

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "something went wrong while attempting getting bracelets",
		}
	}

	// 发送电子邮件
	err = lambdaServerless.sendEmailMaker.SendChakaraResult([]string{req.Email}, uniqueCode, req.Language)
	if err != nil {
		println("error in sending email", err)
	}

	fmt.Println(("sending email "))
	fmt.Println("email is", req.Email)
	fmt.Println("unique code is", uniqueCode)

	response := map[string]interface{}{
		"email":          req.Email,
		"message":        "Email registered successfully",
		"chakra_results": chakras,
		"bracelets":      bracelets,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to marshal response JSON",
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseJSON),
	}
}

func (lambdaServerless *Lambda) getEmailRegistrationByTestNum(ctx context.Context, email string, testNum int) (db.GetEmailRegistrationByTestNumRow, error) {
	// OFFSET 从 0 开始，所以需要减去 1
	arg := db.GetEmailRegistrationByTestNumParams{
		Email:  email,
		Offset: int32(testNum - 1), // Cast testNum-1 to int32 for the Offset field
	}
	registration, err := lambdaServerless.store.GetEmailRegistrationByTestNum(ctx, arg)
	if err != nil {
		return db.GetEmailRegistrationByTestNumRow{}, err
	}
	return registration, nil
}

func (lambdaServerless *Lambda) getLatestEmailRegistration(ctx context.Context, email string) (db.GetLatestEmailRegistrationRow, error) {
	latestRegistration, err := lambdaServerless.store.GetLatestEmailRegistration(ctx, email)
	if err != nil {
		return db.GetLatestEmailRegistrationRow{}, err
	}
	return latestRegistration, nil
}

type ChakraInfo struct {
	ChakraName   string `json:"chakra_name" binding:"required"`
	ChakraScore  int32  `json:"chakra_score" binding:"required"`
	ChakraStatus string `json:"chakra_status" binding:"required"`
}

type GetChakraReportRequest struct {
	Email      string                   `json:"email" binding:"required,email"`
	TestNum    string                   `json:"test_num"`
	Language   string                   `json:"language"`
	ChakraInfo []chakareport.ChakraInfo `json:"chakra_info"` // Use chakareport.ChakraInfo
}

type GetChakraReportResponse struct {
	Report string `json:"report"`
}

type GetEmailRegistrationByTestNumRow struct {
	UniqueID   uuid.UUID `json:"unique_id"`
	Email      string    `json:"email"`
	Language   string    `json:"language"`
	ChakraInfo string    `json:"chakra_info"`
	UniqueCode string    `json:"unique_code"`
	CreatedAt  string    `json:"created_at"` // <- string instead of time.Time
	DeletedAt  string    `json:"deleted_at"` // <- string instead of time.Time
}

type EmailRegistrationResponse struct {
	Data GetEmailRegistrationByTestNumRow `json:"data"`
}

func (lambdaServerless *Lambda) GetChakraReport(ctx context.Context, event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var req GetChakraReportRequest
	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid JSON body",
		}
	}
	fmt.Println("test num received is", req.TestNum)
	// save a copy of testnum incase its 3J32J2J...
	// reportUniqueIdentifier := req.TestNum
	testNum, err := strconv.Atoi(req.TestNum)
	isNumericTestNum := err == nil

	// 打印接收到的 email, testNum, language 和 Chakra Info 的值
	fmt.Printf("Received Email: %s\n", req.Email)
	fmt.Printf("Received Test Number: %d\n", testNum)
	fmt.Printf("Received Language: %s\n", req.Language)
	fmt.Printf("Received Chakra Info: %+v\n", req.ChakraInfo)

	var report []byte
	var language string
	var uniqueCode string
	if !isNumericTestNum && req.TestNum != "" {
		fmt.Println("testNum is alphanumeric")
		// registration, err := lambdaServerless.store.GetReportByCode(ctx, reportUniqueIdentifier)

		// use http to get email registration
		url := lambdaServerless.config.ApiGateWayEndpoint + "/email-registration?unique_code=" + req.TestNum
		fmt.Println("url is", url)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("x-api-key", lambdaServerless.config.ApiGateWayApiKey)
		fmt.Println("api key is", lambdaServerless.config.ApiGateWayApiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var parsed EmailRegistrationResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			panic(err)
		}

		// ✅ Now you can access:
		fmt.Println("UniqueCode:", parsed.Data.UniqueCode)
		fmt.Println("Language:", parsed.Data.Language)
		fmt.Println("Email:", parsed.Data.Email)

		if err != nil {
			if err == sql.ErrNoRows {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusNotFound,
					Body:       "no reports found",
				}
			}

			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("something went wrong in getting report by code: %v", err),
			}
		}

		// 使用 chakra_report
		// var reportString string
		// reportString, err = lambdaServerless.storageMaker.GetChakaraReportByUniqueCode(parsed.Data.Email, parsed.Data.UniqueCode)

		// use http to get the report
		reportEndPoint := lambdaServerless.config.ApiGateWayEndpoint + "/chakra-report?unique_code=" + parsed.Data.UniqueCode + "&email=" + parsed.Data.Email
		fmt.Println("reportEndPoint ", reportEndPoint)

		fetchReportreq, _ := http.NewRequest("GET", reportEndPoint, nil)
		fetchReportreq.Header.Add("x-api-key", lambdaServerless.config.ApiGateWayApiKey)
		fmt.Println("api key is", lambdaServerless.config.ApiGateWayApiKey)

		fetchReportResp, err := http.DefaultClient.Do(fetchReportreq)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		report, _ = io.ReadAll(fetchReportResp.Body)

		// Convert the string report to a byte slice
		// report = []byte(reportString)
		fmt.Println("report is", string(report))

		// 获取 language
		language = parsed.Data.Language
		uniqueCode = parsed.Data.UniqueCode

	} else if testNum > 0 {
		fmt.Println("testNum is numeric")
		// 获取特定的 email_registrations 记录
		registration, err := lambdaServerless.getEmailRegistrationByTestNum(ctx, req.Email, testNum)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("something went wrong while checking test num: %v", err),
		}
		// 使用 chakra_report
		var reportString string
		reportString, err = lambdaServerless.storageMaker.GetChakaraReportByUniqueCode(req.Email, registration.UniqueCode)
		if err != nil {

			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("failed to get report from storage: %v", err),
			}
		}
		// Convert the string report to a byte slice
		report = []byte(reportString)

		// 获取 language
		language = registration.Language
		uniqueCode = registration.UniqueCode
	} else {
		// 使用请求中的 language
		language = req.Language

		// 生成报告的逻辑
		// report, err = lambdaServerless.chakaraReportMaker.GenerateChakaraReport(req.ChakraInfo, language)
		generateReportUrl := lambdaServerless.config.ApiGateWayEndpoint + "/chakra-report"
		fmt.Println("generate report url is", generateReportUrl)

		generateReportData := map[string]interface{}{
			"language":    language,
			"chakra_info": req.ChakraInfo, // assign map[string]string here
		}

		// // Marshal to JSON
		generateReportDataJson, err := json.Marshal(generateReportData)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}

		// // Create request with body
		generateReportReq, err := http.NewRequest("POST", generateReportUrl, bytes.NewReader(generateReportDataJson))
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		generateReportReq.Header.Add("x-api-key", lambdaServerless.config.ApiGateWayApiKey)
		fmt.Println("api key is", lambdaServerless.config.ApiGateWayApiKey)

		rep, err := http.DefaultClient.Do(generateReportReq)
		// if err != nil {
		// 	panic(err)
		// }

		defer rep.Body.Close()

		report, err = io.ReadAll(rep.Body)

		if err != nil {
			println("there is error in the server", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("error occurred while calling GenerateChakaraReport: %v", err),
			}
		}
		// 获取最新的 email_registrations 记录
		// latestRegistration, err := lambdaServerless.getLatestEmailRegistration(ctx, req.Email)
		LatestRegistrationurl := lambdaServerless.config.ApiGateWayEndpoint + "/latest-registration?email=" + req.Email
		fmt.Println("uLatestRegistrationurlrl is", LatestRegistrationurl)

		req, _ := http.NewRequest("GET", LatestRegistrationurl, nil)
		req.Header.Add("x-api-key", lambdaServerless.config.ApiGateWayApiKey)
		fmt.Println("api key is", lambdaServerless.config.ApiGateWayApiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var parsed EmailRegistrationResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			panic(err)
		}

		// ✅ Now you can access:
		fmt.Println("UniqueCode:", parsed.Data.UniqueCode)
		fmt.Println("Language:", parsed.Data.Language)
		fmt.Println("Email:", parsed.Data.Email)

		if err != nil {
			println("error occurred while attempting getLatestEmailRegistration", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("error occurred while attempting getLatestEmailRegistration: %v", err),
			}
		}
		fmt.Printf("Latest registration for email %s: %+v\n", parsed.Data.Email)
		println("the report is ", string(report))

		saveReportUrl := lambdaServerless.config.ApiGateWayEndpoint + "/save-report"
		fmt.Println("saveReportUrl is", saveReportUrl)

		saveReportdata := map[string]interface{}{
			"email":    parsed.Data.Email,
			"uniqueId": parsed.Data.UniqueCode,
			"content":  string(report), // assign map[string]string here
		}
		fmt.Println("saving report...")

		// // Marshal to JSON
		saveReportJsonData, err := json.Marshal(saveReportdata)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}

		// // Create request with body
		saveReportReq, err := http.NewRequest("POST", saveReportUrl, bytes.NewReader(saveReportJsonData))
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		saveReportReq.Header.Add("x-api-key", lambdaServerless.config.ApiGateWayApiKey)
		fmt.Println("api key is", lambdaServerless.config.ApiGateWayApiKey)

		_, err = http.DefaultClient.Do(saveReportReq)
		// if err != nil {
		// 	panic(err)
		// }

		if err != nil {
			fmt.Println("error in saving report", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "error in saving chakra report",
			}
		}

		// 更新 chakra_report 字段
		// if err := lambdaServerless.storageMaker.SaveChakaraReportAsText(parsed.Data.Email, parsed.Data.UniqueCode, string(report)); err != nil {
		// 	return events.APIGatewayProxyResponse{
		// 		StatusCode: http.StatusInternalServerError,
		// 		Body:       fmt.Sprintf("error occurred while attempting SaveChakaraReportAsText: %v", err),
		// 	}
		// }
		uniqueCode = parsed.Data.UniqueCode
	}

	// 直接返回报告
	// ctx.Data(http.StatusOK, "application/json", report)

	fmt.Println("getting ready to return report")
	var data map[string]interface{}
	if err := json.Unmarshal(report, &data); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Invalid report data: %v", err),
		}
	}

	fmt.Println("data:", data)
	// 添加 uniqueCode
	data["uniqueCode"] = uniqueCode

	// 重新编码
	newReport, err := json.Marshal(data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to encode data: %v", err),
		}
	}
	fmt.Println("newReport111:", string(newReport))
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(newReport), // make sure it's string, not []byte
	}

}

type chakraTestResult struct {
	UniqueID     string    `json:"unique_id"`
	ChakraName   string    `json:"chakra_name"`
	ChakraScore  int32     `json:"chakra_score"`
	ChakraStatus string    `json:"chakra_status"`
	CreatedAt    time.Time `json:"created_at"`
}

type getChakraTestResultsRequest struct {
	Email   string `uri:"email" binding:"required,email"`
	TestNum string `uri:"testNum" binding:"required"`
}

func (lambdaServerless *Lambda) GetChakraTestResults(ctx context.Context, event events.APIGatewayProxyRequest, inputEmail string, inputTestNum string) events.APIGatewayProxyResponse {
	// Manual Validation
	if inputEmail == "" || len(inputEmail) <= 5 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "email is empty or less than 5 characters",
		}
	}
	if len(inputTestNum) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "test_num is required",
		}
	}

	// 将 TestNum 从字符串转换为整数
	testNum, err := strconv.Atoi(inputTestNum)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("invalid test_num: %v", err),
		}
	}
	// 获取特定的 email_registrations 记录
	registration, err := lambdaServerless.getEmailRegistrationByTestNum(ctx, inputEmail, testNum)
	if err != nil {
		println("error occured while attempting getEmailRegistrationByTestNum ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("something went wrong while attempting getEmailRegistrationByTestNum: %v", err),
		}

	}

	var reportData []ChakraInfo
	if registration.ChakraInfo.Valid {
		if err := json.Unmarshal([]byte(registration.ChakraInfo.String), &reportData); err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("something went wrong while attempting unmarshalling chakra info: %v", err),
			}
		}
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "No report available",
		}
	}

	// 获取手串
	// 创建一个切片存储分数最低的两个脉轮名称
	type chakraScore struct {
		name  string
		score int32
	}
	scores := make([]chakraScore, len(reportData))
	for i, r := range reportData {
		scores[i] = chakraScore{
			name:  r.ChakraName,
			score: r.ChakraScore,
		}
	}

	// 根据分数排序（从低到高）
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score < scores[j].score
	})

	// 获取分数最低的两个脉轮
	lowestChakras := []string{scores[0].name, scores[1].name}

	// 获取推荐的手串
	var bracelets []db.GetChakraBraceletRow
	bracelets, err = lambdaServerless.store.GetChakraBracelet(ctx, lowestChakras)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("something went wrong while attempting getting bracelets: %v", err),
		}
	}

	var response []chakraTestResult
	for _, r := range reportData {
		response = append(response, chakraTestResult{
			ChakraName:   r.ChakraName,
			ChakraScore:  r.ChakraScore,
			ChakraStatus: r.ChakraStatus,
		})
	}

	//is Go's way of saying: "I’m creating a map where the keys are strings,
	//and the values can be anything — string, int, struct, slice, bool, whatever."

	result := map[string]interface{}{
		"chakra_results": response,
		"language":       registration.Language,
		"bracelets":      bracelets,
	}

	jsonBody, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to marshal JSON response",
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonBody),
	}
}

func (lambdaServerless *Lambda) RequestNotFound(ctx context.Context, event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	fmt.Println("context", ctx)
	fmt.Println("event", event)
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       "Not found!",
	}
	return response
}
