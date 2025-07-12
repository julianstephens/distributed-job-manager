package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/config"
	internalssm "github.com/julianstephens/distributed-task-scheduler/backend/pkg/aws/ssm"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/httputil"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("X-API-Key")

		conf := config.GetConfig()

		if key == "" {
			httputil.NewError(c, http.StatusUnauthorized, errors.New("No API key provided"))
			c.Abort()
			return
		}

		ssmClient, err := internalssm.GetSSMClient()
		if err != nil {
			httputil.NewError(c, http.StatusInternalServerError, errors.New("Unable to check key 1"))
			c.Abort()
			return
		}
		paramRes, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{Name: aws.String(conf.APIKeyParam), WithDecryption: aws.Bool(true)})
		if err != nil {
			// httputil.NewError(c, http.StatusInternalServerError, errors.New("Unable to check key 2"))
			httputil.NewError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		if key != *paramRes.Parameter.Value {
			httputil.NewError(c, http.StatusUnauthorized, errors.New("Unauthorized"))
			c.Abort()
			return
		}

		c.Next()
	}
}
