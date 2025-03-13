package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/util"
)

// createUser
type createUserRequestBody struct {
	Email     string `json:"email" binding:"required,min=1,max=50"`
	FirstName string `json:"first_name" binding:"required,max=50"`
	LastName  string `json:"last_name" binding:"required,max=50"`
	Gender    string `json:"gender" binding:"required,max=50"`
	Birthdate string `json:"birth_date" binding:"required,date"`
	Country   string `json:"country" binding:"required,max=50"`
	Goal      string `json:"goals" binding:"required,max=500"`
	Platform  string `json:"platform" binding:"required,platform"`
	Password  string `json:"password" binding:"required,min=8"`
}

func (server *Server) createUser(ctx *gin.Context) {
	//verify request body
	var req createUserRequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// parse date
	format := util.GetDateFormat()
	parsedBirthDate, err := time.Parse(format, req.Birthdate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateUserTransactiontArgs{
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Gender:         req.Gender,
		Birthdate:      parsedBirthDate,
		Country:        req.Country,
		Goal:           req.Goal,
		Platform:       req.Platform,
		HashedPassword: hashedPassword,
	}

	userResult, err := server.store.CreateUserTransaction(ctx, args)

	if err != nil {
		println("err is not nil!")
		if strings.Contains(err.Error(), "unique_violation") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("email exists already")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, userResult)
}

// fetchuserInformation
type fetchUserInfoByIdParams struct {
	ID int64 `uri:"id"  binding:"required,min=1"`
}

type userInfoResult struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
	Country   string `json:"country"`
	// 1 = yes. 0 = no
	IsMrUser int16 `json:"is_mr_user"`
	// 1 = yes. 0 = no
	IsMobileUser int16  `json:"is_mobile_user"`
	Goals        string `json:"goals"`
}

func (server *Server) fetchUserInfoById(ctx *gin.Context) {
	var result userInfoResult
	//verify params
	var reqParam fetchUserInfoByIdParams
	if err := ctx.ShouldBindUri(&reqParam); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, reqParam.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user does not exists")))
		return
	}
	result.ID = user.ID
	result.Email = user.Email
	result.FirstName = user.FirstName
	result.LastName = user.LastName
	result.Gender = user.Gender
	result.BirthDate = user.BirthDate.Format(util.GetDateFormat())
	result.Country = user.Country
	result.IsMrUser = user.IsMrUser
	result.IsMobileUser = user.IsMobileUser
	result.Goals = user.Goals

	ctx.JSON(http.StatusOK, result)
}

// fetch the time spent by the user
type fetchUserTimeByIdParams struct {
	ID       int64  `uri:"id"  binding:"required,min=1"`
	Platform string `uri:"platform" binding:"required,platform"`
}

func (server *Server) fetchUserTime(ctx *gin.Context) {
	//verify params
	var reqParam fetchUserTimeByIdParams
	if err := ctx.ShouldBindUri(&reqParam); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var err error
	var userTime int64
	switch reqParam.Platform {
	case "mobile":
		userTime, err = server.store.GetUserProfileMobileTime(ctx, reqParam.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "mr":
		userTime, err = server.store.GetUserProfileMrTime(ctx, reqParam.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("platform not available. please check the platform")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total_time_spent_in_mins": userTime,
	})
}

// login user
type loginUserRequestBody struct {
	Email    string `json:"email" binding:"required,min=1,max=50"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginUserRequestResult struct {
	ID          int64  `json:"id"`
	AccessToken string `json:"access_token"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	//verify request body
	var req loginUserRequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		println("error in byscrypt")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		req.Email,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserRequestResult{
		AccessToken: accessToken,
		ID:          user.ID,
	}
	ctx.JSON(http.StatusCreated, rsp)
}
