package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/lotusMind/meditation/db/sqlc"
)

// var chakraQuestions = map[int]string{
// 	1:  "You usually feel present in the moment and grounded in life.",
// 	2:  "You always feel a strong sense of security.",
// 	3:  "You worry about your financial situation and the safety of your home.",
// 	4:  "You feel comfortable no matter where you are.",
// 	5:  "You feel comfortable with intimacy and physical desires.",
// 	6:  "You can express your feelings about sexuality.",
// 	7:  "You are an emotional and passionate person.",
// 	8:  "You have a strong need to establish emotional connections with others.",
// 	9:  "You express yourself through some form of artistic creation (music, painting, singing, or other).",
// 	10: "You often cultivate self-discipline.",
// 	11: "You can stand firm and confident when necessary.",
// 	12: "You have a strong desire to be in control of situations.",
// 	13: "You feel capable of influencing the course of events in a team.",
// 	14: "You take action toward what you want.",
// 	15: "You are a confident person.",
// 	16: "You tend to plan ahead rather than go with the flow.",
// 	17: "You genuinely like most people.",
// 	18: "You feel at ease working in a team.",
// 	19: "You trust most people.",
// 	20: "You strive for harmony in your relationships.",
// 	21: "You easily show compassion to both yourself and others.",
// 	22: "When conflicts arise, you consider others' feelings.",
// 	23: "You give a lot to others, sometimes even neglecting your own needs.",
// 	24: "You enjoy talking.",
// 	25: "You are good at communication, both listening to others and expressing yourself.",
// 	26: "Your voice is loud and clear when you speak.",
// 	27: "You express your emotions openly and without hesitation.",
// 	28: "You are skilled at writing as a form of communication.",
// 	29: "You are good at developing awareness and insight.",
// 	30: "You often have a sense of what will happen in the future.",
// 	31: "You believe coincidences usually have meaning rather than being purely random.",
// 	32: "You rely heavily on your intuition.",
// 	33: "You can easily recall your dreams.",
// 	34: "You instinctively perceive the deeper connections between all things.",
// 	35: "You frequently engage in daydreaming or imagination.",
// 	36: "You are a creative person.",
// 	37: "You are aware of your likes, dislikes, and needs.",
// 	38: "You accept whatever happens to you with ease.",
// 	39: "You see life experiences as opportunities to learn.",
// 	40: "You think effectively using words, symbols, and abstract concepts.",
// 	41: "You feel a deep connection with everything, from the vast universe to the small things around you.",
// }

// createUser
type createUserEmailRequestBody struct {
	Email   string            `json:"email" binding:"required,min=1,max=50"`
	Answers map[string]string `json:"answers"`
}

func (server *Server) registEmail(ctx *gin.Context) {
	var req createUserEmailRequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 创建问题分数映射
	questionAnswers := make(map[string]string)
	for question, answer := range req.Answers {
		questionAnswers[question] = answer
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

	var rootScores int
	var sacralScores int
	var solarPlexusScores int
	var heartScores int
	var throatScores int
	var thirdEyeScores int
	var crownScores int

	for i, answer := range values {

		fmt.Printf("index: %d,answer: %s\n", i, answer)

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
			rootScores += score
		case 2:
			sacralScores += score
		case 3:
			solarPlexusScores += score
		case 4:
			heartScores += score
		case 5:
			throatScores += score
		case 6:
			thirdEyeScores += score
		case 7:
			crownScores += score
		}
	}

	rootScores = (rootScores - 3*8) / (8 * 2) * 100
	sacralScores = (sacralScores - 3*8) / (8 * 2) * 100
	solarPlexusScores = (solarPlexusScores - 3*8) / (8 * 2) * 100
	heartScores = (heartScores - 3*9) / (9 * 2) * 100
	throatScores = (throatScores - 3*8) / (8 * 2) * 100
	thirdEyeScores = (thirdEyeScores - 3*9) / (9 * 2) * 100
	crownScores = (crownScores - 3*7) * 100 / (7 * 2)
	fmt.Printf("RootScores: %d\n", rootScores)
	fmt.Printf("sacralScores : %d\n", sacralScores)
	fmt.Printf("solarPlexusScores: %d\n", solarPlexusScores)
	fmt.Printf("heartScores: %d\n", heartScores)
	fmt.Printf("throatScores: %d\n", throatScores)
	fmt.Printf("thirdEyeScores: %d\n", thirdEyeScores)
	fmt.Printf("CrownChakra Score: %d\n", crownScores)

	args := db.CreateUserEmailTransactiontArgs{
		Email: req.Email,
	}

	err := server.store.CreateUserEmailTransaction(ctx, args)

	if err != nil {
		println("err is not nil!")
		if strings.Contains(err.Error(), "unique_violation") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("email exists already")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 添加返回结果
	result := gin.H{
		"email":   req.Email,
		"message": "Email registered successfully",
	}

	ctx.JSON(http.StatusCreated, result)
}

// // fetchuserInformation
// type fetchUserInfoByIdParams struct {
// 	ID int64 `uri:"id"  binding:"required,min=1"`
// }

// type userInfoResult struct {
// 	ID        int64  `json:"id"`
// 	Email     string `json:"email"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Gender    string `json:"gender"`
// 	BirthDate string `json:"birth_date"`
// 	Country   string `json:"country"`
// 	// 1 = yes. 0 = no
// 	IsMrUser int16 `json:"is_mr_user"`
// 	// 1 = yes. 0 = no
// 	IsMobileUser int16  `json:"is_mobile_user"`
// 	Goals        string `json:"goals"`
// }

// func (server *Server) fetchUserInfoById(ctx *gin.Context) {
// 	var result userInfoResult
// 	//verify params
// 	var reqParam fetchUserInfoByIdParams
// 	if err := ctx.ShouldBindUri(&reqParam); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	user, err := server.store.GetUserById(ctx, reqParam.ID)

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user does not exists")))
// 		return
// 	}
// 	result.ID = user.ID
// 	result.Email = user.Email
// 	result.FirstName = user.FirstName
// 	result.LastName = user.LastName
// 	result.Gender = user.Gender
// 	result.BirthDate = user.Gender
// 	result.Country = user.Country
// 	result.IsMrUser = user.IsMrUser
// 	result.IsMobileUser = user.IsMobileUser
// 	result.Goals = user.Goals

// 	ctx.JSON(http.StatusOK, result)
// }

type chakraTestResult struct {
	UniqueID     string    `json:"unique_id"`
	ChakraName   string    `json:"chakra_name"`
	ChakraScore  int32     `json:"chakra_score"`
	ChakraStatus string    `json:"chakra_status"`
	CreatedAt    time.Time `json:"created_at"`
}

type getChakraTestResultsRequest struct {
	Email string `uri:"email" binding:"required,email"`
}

func (server *Server) getChakraTestResults(ctx *gin.Context) {
	var req getChakraTestResultsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 创建参数结构体
	params := db.GetChakraTestResultsParams{
		Email: req.Email,
	}

	results, err := server.store.GetChakraTestResults(ctx, params) // 传入结构体
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []chakraTestResult
	for _, r := range results {
		response = append(response, chakraTestResult{
			UniqueID:     r.UniqueId.String(),
			ChakraName:   r.ChakraName,
			ChakraScore:  r.ChakraScore,
			ChakraStatus: r.ChakraStatus,
			CreatedAt:    r.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, response)
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
