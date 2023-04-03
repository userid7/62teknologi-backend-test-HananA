package v1

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}

func parseError(err error) []string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make([]string, len(validationErrs))
		for i, e := range validationErrs {
			switch e.Tag() {
			case "required_without":
				errorMessages[i] = fmt.Sprintf("The field %s is required if %s is not supplied", e.Field(), e.Param())
			}
		}
		return errorMessages
	} else if marshallingErr, ok := err.(*json.UnmarshalTypeError); ok {
		return []string{fmt.Sprintf("The field %s must be a %s", marshallingErr.Field, marshallingErr.Type.String())}
	}
	return nil
}
