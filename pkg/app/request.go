package app

import (
	"github.com/astaxie/beego/validation"

	"gitee.com/muzipp/Distribution/pkg/logging"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}

	return
}
