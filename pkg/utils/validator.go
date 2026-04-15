package utils

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var Trans ut.Translator

func init() {
	if err := InitTrans("zh"); err != nil {
		fmt.Println("初始化翻译器错误")
		panic("初始化翻译器错误")
	}
}
func RemoveTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// InitTrans 定义翻译的方法
func InitTrans(locale string) (err error) {
	//修改gin框架中的validator引擎属性, 实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}
		switch locale {
		case "en":
			enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			zhTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			enTranslations.RegisterDefaultTranslations(v, Trans)
		}
		return
	}
	return
}

// ShowErrorMessage 翻译错误
func ShowErrorMessage(err error) string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		fmt.Println(err.Error(), "类型错误")
		return err.Error()
	}
	// 直接返回第一个就可以
	for _, val := range RemoveTopStruct(errs.Translate(Trans)) {
		return val
	}
	return ""
}
