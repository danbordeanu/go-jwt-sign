package model

import (
	"fmt"
)

// SignatureValidation represents the structure for validating user signatures.
//
// swagger:model
type SignatureValidation struct {
	Request   `json:"-" swaggerignore:"true"`
	User      string `json:"user" example:"JonnyBoy"`
	Signature string `json:"signature" example:"test-signature-JonnyBoy"`
}

// Validate checks if the required fields in SignatureValidation are present.
//
// Returns:
//   - error: Validation error, nil if validation passes
func (r *SignatureValidation) Validate() error {
	if r.User == "" {
		return fmt.Errorf("missing parameter: Username")
	}
	if r.Signature == "" {
		return fmt.Errorf("missing parameter: signature")
	}
	return nil
}

// JwtValidation represents the structure for validating JWTs.
//
// swagger:model
type JwtValidation struct {
	Request   `json:"-" swaggerignore:"true"`
	Jwt       string   `json:"jwt" example:"your_jwt_here"`
	Questions []string `json:"questions" example:"question1,question2"`
	Answers   []string `json:"answers" example:"answer1,answer2"`
}

// Validate checks if the required fields in JwtValidation are present.
//
// Returns:
//   - error: Validation error, nil if validation passes
func (r *JwtValidation) Validate() error {
	if r.Jwt == "" {
		return fmt.Errorf("missing parameter: JWT")
	}
	if len(r.Questions) == 0 {
		return fmt.Errorf("missing parameter: Questions")
	}
	if len(r.Answers) == 0 {
		return fmt.Errorf("missing parameter: Answers")
	}
	return nil
}
