package value

type ParameterType string

const (
	ParameterTypeText     ParameterType = "text"
	ParameterTypeInt      ParameterType = "int"
	ParameterTypeFloat    ParameterType = "float"
	ParameterTypeBool     ParameterType = "bool"
	ParameterTypeDatetime ParameterType = "datetime"
)
