package api

// createSession
type createUserRequestBody struct {
	Email     string `json:"email" binding:"required,min=1,max=50"`
	FirstName string `json:"first_name" binding:"required,max=50"`
	LastName  string `json:"last_name" binding:"required,max=50"`
	Gender    string `json:"gender" binding:"required,max=50"`
	Birthdate string `json:"birth_date" binding:"required,max=50"`
	Country   string `json:"country" binding:"required,max=50"`
	Goal      string `json:"goal" binding:"required,max=500"`
	Platform  string `json:"platform" binding:"required,max=20"`
	Password  string `json:"password" binding:"required,min=8"`
}

// func (server *Server) createSession(ctx *gin.Context) {
// 	//verify params
// 	var req createSessionRequestParams
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 	}
// 	// verify session platform
// 	sessionPlatform := util.GetPlatformTypeBasedOnSessionType(req.SessionType)
// 	if len(sessionPlatform) <= 1 {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(InvalidSessionType))
// 	}

// 	createSessionLogArgs := db.CreateSessionLogParams{
// 		UserID:          req.ID,
// 		SessionType:     req.SessionType,
// 		SessionPlatform: sessionPlatform,
// 	}
// 	sessionLogResult, err := server.store.CreateSessionLogTransaction(ctx, createSessionLogArgs)

// 	if err != nil {
// 		if strings.Contains(err.Error(), "foreign_key_violation") {
// 			ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		}
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 	}

// 	ctx.JSON(http.StatusCreated, sessionLogResult)
// }
