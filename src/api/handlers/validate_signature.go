package handlers

import (
	"fmt"
	"github.com/danbordeanu/go-logger"
	"github.com/danbordeanu/go-stats/concurrency"
	"github.com/danbordeanu/go-utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"jwt-sign/api/response"
	"jwt-sign/configuration"
	"jwt-sign/model"
	"strings"
)

// VerifySignature godoc
// @Summary Verify signature
// @Description Verify signature for a given user
// @ID verifySignature
// @Accept json
// @Produce html
// @Param model.SignatureValidation body model.SignatureValidation true "validate signature"
// @Success 200 {object} model.JSONSuccessResult "The request was validated and has been processed successfully (sync)"
// @Failure 400 {object} model.JSONFailureResult "The payload is invalid"
// @Router /v1/verify-signature [post]
func VerifySignature(c *gin.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()
	log := logger.SugaredLogger().WithContextCorrelationId(c).With("package", "handlers", "action", "VerifySignature")

	var (
		e             error
		err           error
		rr            model.SignatureValidation
		ctx           = c.Request.Context()
		correlationId = c.MustGet("correlation_id").(string)
	)
	ctx, span := tracer.Start(ctx, "Signature Validation",
		oteltrace.WithAttributes(attribute.String("CorrelationId", correlationId)))
	defer span.End()

	// validate params
	if err = c.ShouldBindJSON(&rr); err != nil {
		e = fmt.Errorf("error while parsing request: %s", err.Error())
		span.SetStatus(codes.Error, e.Error())
		span.RecordError(err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 400, Err: e})
		return
	}
	if err = rr.Validate(); err != nil {
		e = fmt.Errorf("error while validating request: %s", err.Error())
		span.SetStatus(codes.Error, e.Error())
		span.RecordError(err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 400, Err: e})
		return
	}

	// Get user and signature from query parameters
	user := rr.User
	signature := rr.Signature

	log.Debugf("user:%s, signature:%s", user, signature)

	span.AddEvent("Validate signature")
	// Perform signature verification logic here
	// Check if the user is present in the signature
	if !ValidateUserSignature(c, signature, user) {
		e = fmt.Errorf("user is not present in the signature")
		span.SetStatus(codes.Error, e.Error())
		span.RecordError(err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 400, Err: e})
		return
	}

	// Assuming a successful verification, construct the response
	response.RegistrationHtmlResponse(c, configuration.HtmlJwtValidationSuccessPage, "", "OK if signature belongs to user,", "")

}

// ValidateUserSignature validates if the given user is present in the provided signature.
//
// Parameters:
//   - c *gin.Context: Gin context for logging purposes
//   - signature string: The signature to be validated
//   - user string: The user to check within the signature
//
// Returns:
//   - bool: True if the user is present in the signature, otherwise false
func ValidateUserSignature(c *gin.Context, signature, user string) bool {
	log := logger.SugaredLogger().WithContextCorrelationId(c).With("package", "routine", "action", "doValidate")
	defer log.Debugf("validate proccess finished")
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()
	log.Debugf("signature:%s, user:%s", signature, user)
	return strings.Contains(signature, user)
}
