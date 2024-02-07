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
)

// ValidateJwt godoc
// @Summary Validate jwt
// @Description Validate Jwt
// @ID  validateJwt
// @Produce html
// @Param model.JwtValidation body model.JwtValidation true "validate signature"
// @Success 200 {object} model.JSONSuccessResult "The request was validated and has been processed successfully (sync)"
// @Router /v1/validate-jwt [post]
func ValidateJwt(c *gin.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()
	log := logger.SugaredLogger().WithContextCorrelationId(c).With("package", "handlers", "action", "ValidateJwt")

	var (
		err           error
		e             error
		rr            model.JwtValidation
		ctx           = c.Request.Context()
		correlationId = c.MustGet("correlation_id").(string)
	)

	ctx, span := tracer.Start(ctx, "JWT user validation",
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

	span.AddEvent("New Jwt Provider")

	// decode jwt token received in post request
	span.AddEvent("Decode Jwt Token")

	// query db by jwt and get ravendb ID
	span.AddEvent("Query Db after jwt")

	// load document by ravenID
	span.AddEvent("Load Document By Id")

	// decode jwt stored in ravendb
	span.AddEvent("Decode JWT stored in db")

	// decide what type of jwt we have based on presence absence of password key
	span.AddEvent("Unmarshal JWT")

	// Retrieve the questions and answers from the query parameters
	questions := rr.Questions
	answers := rr.Answers

	// Sign the answers using your logic
	testSignature, err := SignAnswers(c, questions, answers)
	if err != nil {
		e = fmt.Errorf("failed to sign answers: %s", err)
		log.Errorf("%s", e)
		span.SetStatus(codes.Error, e.Error())
		span.RecordError(err)
		response.RegistrationHtmlResponse(c, configuration.HtmlJwtValidationSuccessPage, "", "failed", "")
		return
	}

	log.Debugf("we got signature:%s", testSignature)

	concurrency.GlobalWaitGroup.Add(1)
	go func() {
		defer log.Debugf("onboarding subroutine proccess finished")
		defer concurrency.GlobalWaitGroup.Done()
		_, span = tracer.Start(ctx, "On boarding subroutine process")
		defer span.End()
		log.Debugf("start doing things")
		span.AddEvent("we do some stuff here")
	}()

	response.RegistrationHtmlResponse(c, configuration.HtmlJwtValidationSuccessPage, "", "successfully", testSignature)

}

// SignAnswers signs the provided answers based on the given questions.
//
// Parameters:
//   - c *gin.Context: Gin context for logging purposes
//   - questions []string: List of questions for which answers are provided
//   - answers []string: List of answers corresponding to the questions
//
// Returns:
//   - string: The generated signature for the provided answers
//   - error: An error, if any, encountered during the signing process
func SignAnswers(c *gin.Context, questions, answers []string) (string, error) {
	log := logger.SugaredLogger().WithContextCorrelationId(c).With("package", "routine", "action", "doSignature")
	defer log.Debugf("sign answer proccess finished")
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()
	// Implement logic for signing answers here
	log.Debugf("question:%s, answer:%s", questions, answers)
	// This is just a placeholder, replace it with your actual signing logic
	return "test-signature-JonnyBoy", nil
}
