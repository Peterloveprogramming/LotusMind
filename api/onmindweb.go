package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lotusMind/meditation/chakareport"
	db "github.com/lotusMind/meditation/db/sqlc"
)

// createUser
type createUserEmailRequestBody struct {
	Email    string            `json:"email" binding:"required,min=1,max=50"`
	Answers  map[string]string `json:"answers"`
	Language string            `json:"language"`
	IP       string            `json:"ip"`
	Country  string            `json:"country"`
}

func (server *Server) registEmail(ctx *gin.Context) {
	var req createUserEmailRequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 创建问题分数映射
	fmt.Printf("req.Answers: %+v\n", req.Answers)
	questionAnswers := make(map[string]string)

	// 准备批量插入的参数
	var uniqueIds []uuid.UUID
	var emails []string
	var uniqueCodes []string
	var questions []string
	var answers []string
	uniqueCode := generateUniqueCode()

	for question, answer := range req.Answers {
		questionAnswers[question] = answer
		uniqueIds = append(uniqueIds, uuid.New())
		emails = append(emails, req.Email)
		uniqueCodes = append(uniqueCodes, uniqueCode)
		questions = append(questions, question)
		answers = append(answers, answer)
	}

	// 批量插入选项答案
	_, err := server.store.CreateChakraTestOptionAnswersBatch(ctx, db.CreateChakraTestOptionAnswersBatchParams{
		UniqueIds:   uniqueIds,
		Emails:      emails,
		UniqueCodes: uniqueCodes,
		Questions:   questions,
		Answers:     answers,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
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
			return
		}
		score, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("转换失败:", err)
			return
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

	// rootScores = (rootScores - 3*8) / (8 * 2) * 100
	// sacralScores = (sacralScores - 3*8) / (8 * 2) * 100
	// solarPlexusScores = (solarPlexusScores - 3*8) / (8 * 2) * 100
	// heartScores = (heartScores - 3*9) / (9 * 2) * 100
	// throatScores = (throatScores - 3*8) / (8 * 2) * 100
	// thirdEyeScores = (thirdEyeScores - 3*9) / (9 * 2) * 100
	// crownScores = (crownScores - 3*7) * 100 / (7 * 2)
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 创建用户注册邮箱
	args := db.CreateUserEmailTransactiontArgs{
		Email:      req.Email,
		ChakraInfo: string(chakraInfoJSON),
		Language:   req.Language,
		UniqueCode: uniqueCode,
		IP:         req.IP,
		Country:    req.Country,
	}

	err = server.store.CreateUserEmailTransaction(ctx, args)
	if err != nil {
		if strings.Contains(err.Error(), "unique_violation") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("email exists already")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

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

	// 获取分数最低的两个脉轮
	lowestChakras := []string{chakraScores[0].name, chakraScores[1].name}
	fmt.Printf("分数最低的两个脉轮是: %v\n", lowestChakras)

	// 获取推荐的手串
	bracelets, err := server.store.GetChakraBracelet(ctx, lowestChakras)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// 添加返回结果
	result := gin.H{
		"email":          req.Email,
		"message":        "Email registered successfully",
		"chakra_results": chakras,   // 直接返回解析后的脉轮结果
		"bracelets":      bracelets, // 添加手串信息到返回结果
	}

	ctx.JSON(http.StatusCreated, result)
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateUniqueCode() string {
	result := make([]byte, 11)
	for i := 0; i < 11; i++ {
		num := rand.Intn(len(charset))
		result[i] = charset[num]
	}
	return string(result)
}

// 根据分数确定脉轮状态
func getChakraStatus(score float32) string {
	if score >= 80 && score <= 100 {
		return "Overactive"
	} else if score >= 20 && score < 80 {
		return "Open"
	} else if score >= 0 && score < 20 {
		return "Underactive"
	} else if score >= -50 && score < 0 {
		return "Partially Blocked"
	} else {
		return "Severely Blocked"
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

func (server *Server) getChakraTestResults(ctx *gin.Context) {
	var req getChakraTestResultsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 将 TestNum 从字符串转换为整数
	testNum, err := strconv.Atoi(req.TestNum)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid test_num: %v", err)))
		return
	}

	// 获取特定的 email_registrations 记录
	registration, err := server.getEmailRegistrationByTestNum(ctx, req.Email, testNum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var reportData []ChakraInfo
	if registration.ChakraInfo.Valid {
		if err := json.Unmarshal([]byte(registration.ChakraInfo.String), &reportData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("No report available")))
		return
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
	var bracelets []db.ChakraBracelet
	bracelets, err = server.store.GetChakraBracelet(ctx, lowestChakras)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []chakraTestResult
	for _, r := range reportData {
		response = append(response, chakraTestResult{
			ChakraName:   r.ChakraName,
			ChakraScore:  r.ChakraScore,
			ChakraStatus: r.ChakraStatus,
		})
	}

	// 修改返回结果的结构
	result := gin.H{
		"chakra_results": response,
		"language":       registration.Language,
		"bracelets":      bracelets,
	}

	ctx.JSON(http.StatusOK, result)
}

type chakraData struct {
	ChakraName   string `json:"chakra_name" binding:"required"`
	ChakraScore  int32  `json:"chakra_score" binding:"required,min=0,max=100"`
	ChakraStatus string `json:"chakra_status" binding:"required,oneof=active inactive balanced"`
}

type createChakraTestResultsRequest struct {
	Email   string       `json:"email" binding:"required,email"`
	Chakras []chakraData `json:"chakras" binding:"required,min=7,max=7"` // 确保正好7个脉轮
}

func (server *Server) createChakraTestResults(ctx *gin.Context) {
	var req createChakraTestResultsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var results []db.ChakraTestResult
	err := server.store.ExecTx(ctx, func(q *db.Queries) error {
		for _, chakra := range req.Chakras {
			arg := db.CreateChakraTestResultParams{
				UniqueId:     uuid.New(),
				Email:        req.Email,
				ChakraName:   chakra.ChakraName,
				ChakraScore:  chakra.ChakraScore,
				ChakraStatus: chakra.ChakraStatus,
			}
			result, err := q.CreateChakraTestResult(ctx, arg)
			if err != nil {
				return err
			}
			results = append(results, result)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "foreign_key_violation") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("email not registered")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, results)
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
	ChakraInfo []chakareport.ChakraInfo `json:"chakra_info" binding:"required"` // Use chakareport.ChakraInfo
}

type GetChakraReportResponse struct {
	Report string `json:"report"`
}

func (server *Server) getChakraReport(ctx *gin.Context) {
	var req GetChakraReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 如果 TestNum 是空字符串，设置默认值
	if req.TestNum == "" {
		req.TestNum = "0" // 设置默认值为 0
	}

	// 将 TestNum 从字符串转换为整数
	testNum, err := strconv.Atoi(req.TestNum)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid test_num: %v", err)))
		return
	}

	// 打印接收到的 email, testNum, language 和 Chakra Info 的值
	fmt.Printf("Received Email: %s\n", req.Email)
	fmt.Printf("Received Test Number: %d\n", testNum)
	fmt.Printf("Received Language: %s\n", req.Language)
	fmt.Printf("Received Chakra Info: %+v\n", req.ChakraInfo)

	var report []byte
	var language string
	var uniqueCode string
	if testNum > 0 {
		// 获取特定的 email_registrations 记录
		registration, err := server.getEmailRegistrationByTestNum(ctx, req.Email, testNum)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		fmt.Printf("Registration for email %s: %+v\n", req.Email, registration)

		// 使用数据库中的 chakra_report
		if registration.ChakraReport.Valid {
			report = []byte(registration.ChakraReport.String)
		} else {
			report = []byte("No report available")
		}

		// 获取 language
		language = registration.Language
		uniqueCode = registration.UniqueCode
	} else {
		// 使用请求中的 language
		language = req.Language

		// 生成报告的逻辑
		report, err = server.chakaraReportMaker.GenerateChakaraReport(req.ChakraInfo, language)
		if err != nil {
			println("there is error in the server", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// 获取最新的 email_registrations 记录
		latestRegistration, err := server.getLatestEmailRegistration(ctx, req.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		fmt.Printf("Latest registration for email %s: %+v\n", req.Email, latestRegistration)

		// get uniqueid
		// save it like ommind/chakara-report/y.peter998@gmail.com/WZFZ5T6PSIX.txt
		// we want to pass email, and uniqueid, report  s3 - perhaps has a function called, save chakarareport?
		// function needs to have crud
		// 更新 chakra_report 字段
		if err := server.updateChakraReportByUniqueId(ctx, latestRegistration.UniqueId, report); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		uniqueCode = latestRegistration.UniqueCode
	}

	// 直接返回报告
	// ctx.Data(http.StatusOK, "application/json", report)

	var data map[string]interface{}
	if err := json.Unmarshal(report, &data); err != nil {
		ctx.String(http.StatusInternalServerError, "Invalid report data")
		return
	}

	// 添加 uniqueCode
	data["uniqueCode"] = uniqueCode

	// 重新编码
	newReport, err := json.Marshal(data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to encode data")
		return
	}
	fmt.Println("newReport111:", string(newReport))
	ctx.Data(http.StatusOK, "application/json", newReport)

}

func (server *Server) getEmailRegistrationByTestNum(ctx *gin.Context, email string, testNum int) (db.EmailRegistrations, error) {
	// OFFSET 从 0 开始，所以需要减去 1
	registration, err := server.store.GetEmailRegistrationByTestNum(ctx, email, testNum-1)
	if err != nil {
		return db.EmailRegistrations{}, err
	}
	return registration, nil
}

func (server *Server) updateChakraReportByUniqueId(ctx *gin.Context, uniqueId uuid.UUID, report []byte) error {
	err := server.store.UpdateChakraReportByUniqueId(ctx, string(report), uniqueId)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) getLatestEmailRegistration(ctx *gin.Context, email string) (db.EmailRegistrations, error) {
	latestRegistration, err := server.store.GetLatestEmailRegistration(ctx, email)
	if err != nil {
		return db.EmailRegistrations{}, err
	}
	return latestRegistration, nil
}

func (server *Server) getReportByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	registration, err := server.store.GetReportByCode(ctx, code)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 解析 chakra_info
	var chakraInfo []ChakraInfo
	if registration.ChakraInfo.Valid {
		if err := json.Unmarshal([]byte(registration.ChakraInfo.String), &chakraInfo); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse chakra info"})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No chakra info available"})
		return
	}

	// 获取推荐的手串
	// 创建一个切片存储分数最低的两个脉轮名称
	type chakraScore struct {
		name  string
		score int32
	}
	scores := make([]chakraScore, len(chakraInfo))
	for i, r := range chakraInfo {
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
	var bracelets []db.ChakraBracelet
	bracelets, err = server.store.GetChakraBracelet(ctx, lowestChakras)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建返回结果
	var response []chakraTestResult
	for _, r := range chakraInfo {
		response = append(response, chakraTestResult{
			ChakraName:   r.ChakraName,
			ChakraScore:  r.ChakraScore,
			ChakraStatus: r.ChakraStatus,
		})
	}

	// 返回结果
	result := gin.H{
		"chakra_results": response,
		"language":       registration.Language,
		"bracelets":      bracelets,
	}

	ctx.JSON(http.StatusOK, result)
}
