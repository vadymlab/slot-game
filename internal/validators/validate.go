package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	log "github.com/public-forge/go-logger"
	"strings"
)

// Validate runs struct-level validation on the provided struct `s`.
// It applies custom validators, such as the UUID validator, and returns a slice of error messages if validation fails.
// Each error message follows the format: "field::tag::param" in lowercase.
func Validate(s interface{}) []string {
	validate := validator.New(
		validator.WithRequiredStructEnabled(),
		WithUUIDValidator())
	err := validate.Struct(s)
	if err != nil {
		var errs = make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, strings.ToLower(err.Field()+"::"+err.Tag()+"::"+err.Param()))
		}
		return errs
	}
	return nil
}

// WithUUIDValidator adds a custom UUID validator to the validator instance.
// This function registers the custom `uuid` validation with the validator.
func WithUUIDValidator() validator.Option {
	return func(v *validator.Validate) {
		err := v.RegisterValidation("uuid", uuidValidation)
		if err != nil {
			log.FromDefaultContext().Error(err)
		}
	}
}

// uuidValidation checks if the given field contains a valid UUID.
// It validates both standard UUID formats and byte representations, returning true if the UUID is valid.
func uuidValidation(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	if uuid.Validate(id) == nil {
		return true
	}
	_, err := uuid.FromBytes(fl.Field().Bytes())
	return err == nil
}
