package response

import (
	"dev.azure.com/coderollers/almeria/go-shared-noversion/http/model"
	"fmt"
	"github.com/danbordeanu/go-utils"
	"github.com/gin-gonic/gin"
	"jwt-sign/configuration"
	"math"
	"net/http"
)

// SuccessResponse sends a successful JSON response to the client.
//
// Parameters:
//   - c *gin.Context: Gin context for handling the response
//   - data interface{}: Data to be included in the response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, model.JSONSuccessResult{
		Code:          http.StatusOK,
		Data:          data,
		Message:       "Success",
		CorrelationId: c.MustGet("correlation_id").(string),
	})
}

// FailureResponse sends a failure JSON response to the client.
//
// Parameters:
//   - c *gin.Context: Gin context for handling the response
//   - data interface{}: Data to be included in the response
//   - err utils.HttpError: HTTP error details
func FailureResponse(c *gin.Context, data interface{}, err utils.HttpError) {
	if err.Err == nil {
		err = utils.HttpError{Code: int(math.Max(float64(err.Code), 500)), Err: fmt.Errorf("FailureResponse was called with a nil error (%s)", err.Message)}
	}
	var errorString, stackString string
	conf := configuration.AppConfig()
	if conf.Development {
		errorString = err.Error()
		stackString = err.StackTrace()
	} else {
		// we want to see error reason but not stackstring
		errorString = err.Error()
	}
	c.JSON(err.Code, model.JSONFailureResult{
		Code:          err.Code,
		Data:          data,
		Error:         errorString,
		Stack:         stackString,
		CorrelationId: c.MustGet("correlation_id").(string),
	})
}

func RegistrationHtmlResponse(c *gin.Context, page, login, status, testSignature string) {
	PutBody := map[string]interface{}{
		"login":         login,
		"status":        status,
		"testSignature": testSignature, // Add the testSignature field
	}
	c.HTML(http.StatusOK, page, gin.H(PutBody))
}
