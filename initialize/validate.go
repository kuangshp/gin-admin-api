package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var Trans ut.Translator

// NewValidator Wire Provider：初始化自定义校验器和中文翻译
func NewValidator() (ut.Translator, error) {
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		if err := zhTranslations.RegisterDefaultTranslations(v, trans); err != nil {
			return nil, err
		}
	}

	Trans = trans
	return trans, nil
}
