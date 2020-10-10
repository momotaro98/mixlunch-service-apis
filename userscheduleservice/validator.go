package userscheduleservice

import (
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
	//if err := validate.RegisterValidation("custom_tag_name", func(fl validator.FieldLevel) bool {
	//	return false
	//}); err != nil {
	//	panic(err)
	//}
}

func ValidateUserSchedule(schedule *UserScheduleForCommand) error {
	if err := validate.Struct(schedule); err != nil {
		return err
	}
	return nil
}
