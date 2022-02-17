package config

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

const (
	ENG = "en"
	IDN = "id"
)

type Validator struct {
	Driver     *validator.Validate
	Uni        *ut.UniversalTranslator
	Translator ut.Translator
}

func SetupValidator(c *AppConfig) *Validator {
	en := en.New()
	id := id.New()
	uni := ut.New(en, id)

	transEN, _ := uni.GetTranslator(ENG)
	transID, _ := uni.GetTranslator(IDN)

	validatorDriver := validator.New()

	_ = enTranslations.RegisterDefaultTranslations(validatorDriver, transEN)
	_ = idTranslations.RegisterDefaultTranslations(validatorDriver, transID)

	var translator ut.Translator
	switch c.Locale {
	case ENG:
		translator = transEN
	case IDN:
		translator = transID
	}

	return &Validator{Driver: validatorDriver, Uni: uni, Translator: translator}
}
