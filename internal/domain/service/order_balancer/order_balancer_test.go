package order_balancer

import (
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/value"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMatchParam(t *testing.T) {
	var data = []struct {
		ExecutorParam entity.ExecutorParameter
		OrderParam    string

		Result bool
	}{
		{
			ExecutorParam: entity.ExecutorParameter{
				Type: value.ParameterTypeFloat,
				Mask: "2.5x",
			},
			OrderParam: "3.3",

			Result: true,
		},
		{
			ExecutorParam: entity.ExecutorParameter{
				Type: value.ParameterTypeText,
				Mask: "asdsad",
			},
			OrderParam: "asdsad",

			Result: true,
		},
		{
			ExecutorParam: entity.ExecutorParameter{
				Type: value.ParameterTypeInt,
				Mask: "x12",
			},
			OrderParam: "10",

			Result: true,
		},
	}

	for _, d := range data {
		t.Logf("mask '%s', value '%s' %s", d.ExecutorParam.Mask, d.OrderParam, d.ExecutorParam.Type)
		require.Equal(t, d.Result, matchParam(d.ExecutorParam, d.OrderParam), "ExecutorParam %+v OrderParam %+v", d.ExecutorParam, d.OrderParam)
	}

}
