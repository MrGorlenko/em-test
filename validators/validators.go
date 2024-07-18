package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatePassportNumberFormat(fl validator.FieldLevel) bool {
    re := regexp.MustCompile(`^\d{4} \d{6}$`)
    return re.MatchString(fl.Field().String())
}
