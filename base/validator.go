package base

import (
	"github.com/SAIKAII/skResk-Infra"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
)

var validate *validator.Validate
var translator ut.Translator

func Validate() *validator.Validate {
	return validate
}

func Translate() ut.Translator {
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	validate = validator.New()
	// 创建消息国际化通用翻译器
	cn := zh.New()
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh")
	if found {
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		logrus.Error("Not found translator: zh")
	}
}

func ValidateStruct(s interface{}) (err error) {
	err = Validate().Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logrus.Error("验证错误", err)
		}
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				logrus.Error(e.Translate(Translate()))
			}
		}
	}
	return err
}
