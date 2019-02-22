package validator

import (
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// ValidateStruct make sure the struct is valid and translate errors to human-readable messages.
func ValidateStruct(s interface{}) error {
	if err := defaults.Set(s); err != nil {
		return err
	}

	var uni *ut.UniversalTranslator
	var validate *validator.Validate

	en := en.New()
	uni = ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Validate Config struct
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		// Translate errors to human-readable english messages
		errs := err.(validator.ValidationErrors)
		transErrs := errs.Translate(trans)

		var errStrings []string
		for _, transErr := range transErrs {
			errStrings = append(errStrings, transErr)
		}

		validationError := fmt.Errorf(strings.Join(errStrings, "\n"))
		return validationError
	}

	return nil
}
