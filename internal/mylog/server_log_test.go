package mylog

import (
	"testing"
)

func TestServerLogger(t *testing.T) {
	CommonLogger.Info().
		Msg("开始登录...")
	CommonLogger.Info().
		Msg("开始登录...")
}
