package common

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func init() {
	// 注册中文翻译
	if err := transInit("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}
}

// local 通常取决于 http 请求头的 'Accept-Language'
func transInit(local string) (err error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)
		var o bool
		trans, o = uni.GetTranslator(local)
		if !o {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		// 注册翻译器
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return nil
	}
	return nil
}

// 从错误中提取文字提示
func ShowErrorMessage(err error) string {
	// 获取validator.ValidationErrors类型的errors,比如要求的是数字,你传递字符串
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		fmt.Println(err.Error(), "类型错误")
		return err.Error()
	}
	// validator.ValidationErrors类型错误则进行翻译
	for _, val := range errs.Translate(trans) {
		return val
	}
	return ""
}
