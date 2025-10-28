package order_balancer

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/value"
	"strconv"
	"strings"
	"time"
)

func checkExecutorParams(executor *aggregate.ExecutorWithParams, parameters []entity.OrderParameter) bool {
	if len(parameters) == 0 {
		return true
	}

	mappedParams := mapParams(parameters)
	paramsLen := len(parameters)

	var count int
	for _, executorParam := range executor.Parameters {
		if param, ok := mappedParams[executorParam.Id]; ok && matchParam(executorParam, param) {
			count++
			if count == paramsLen {
				return true
			}
		}
	}

	return false
}

func mapParams(params []entity.OrderParameter) map[int]string {
	m := make(map[int]string)
	for _, param := range params {
		m[param.Id] = param.Value
	}
	return m
}

func matchParam(param entity.ExecutorParameter, val string) bool {
	switch param.Type {
	case value.ParameterTypeText, value.ParameterTypeBool:
		return param.Mask == val
	case value.ParameterTypeInt, value.ParameterTypeFloat, value.ParameterTypeDatetime:
		parsedMask := parseMask(param.Mask, val)
		if len(parsedMask) == 1 {
			return param.Mask == val
		}

		switch param.Type {
		case value.ParameterTypeInt:
			return compareInt(parsedMask)
		case value.ParameterTypeFloat:
			return compareFloat(parsedMask)
		case value.ParameterTypeDatetime:
			return compareDatetime(parsedMask)
		}
	}

	return false
}

// TODO use interface
func compareInt(s []string) bool {
	lastN, _ := strconv.Atoi(s[0])
	for _, v := range s[1:] {
		n, _ := strconv.Atoi(v)
		if lastN >= n {
			return false
		}
		lastN = n
	}
	return true
}
func compareFloat(s []string) bool {
	lastN, _ := strconv.ParseFloat(s[0], 64)
	for _, v := range s[1:] {
		n, _ := strconv.ParseFloat(v, 64)
		if lastN >= n {
			return false
		}
		lastN = n
	}
	return true
}
func compareDatetime(s []string) bool {
	lastTime, _ := time.Parse(time.DateTime, s[0])
	for _, v := range s[1:] {
		t, _ := time.Parse(time.DateTime, v)
		if lastTime.After(t) {
			return false
		}
		lastTime = t
	}
	return true
}

func parseMask(mask string, val string) []string {
	xIdx := strings.Index(mask, "x")
	if xIdx == -1 {
		return []string{mask}
	}

	if xIdx == 0 {
		return []string{val, mask[1:]}
	}
	if xIdx == len(mask)-1 {
		return []string{mask[:len(mask)-1], val}
	}

	return []string{mask[:xIdx], val, mask[xIdx+1:]}
}
