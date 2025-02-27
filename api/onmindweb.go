package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/lotusMind/meditation/db/sqlc"
)

// createUser
type createUserEmailRequestBody struct {
	Email string `json:"email" binding:"required,min=1,max=50"`
}

func (server *Server) registEmail(ctx *gin.Context) {
	//verify request body
	var req createUserEmailRequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

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
