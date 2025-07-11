package middleware

// func AuthGuard(api controller.Controller, scopes ...string) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		apiKey := ctx.GetHeader("X-API-KEY")
// 		if apiKey == "" {
// 			httputil.NewError(ctx, http.StatusUnauthorized, errors.New("no API key provided"))
// 			ctx.Abort()
// 			return
// 		}

// 		isAuthed, userId, err := api.Authorize(apiKey, scopes...)
// 		if err != nil {
// 			httputil.NewError(ctx, http.StatusUnauthorized, err)
// 			ctx.Abort()
// 			return
// 		}

// 		if !isAuthed {
// 			httputil.NewError(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set("user", userId)
// 		ctx.Next()
// 	}
// }
