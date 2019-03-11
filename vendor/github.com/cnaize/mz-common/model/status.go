package model

type CenterStatus struct {
	MinCoreVersion *string `json:"minCoreVersion,omitempty" form:"minCoreVersion"`
	Error          *Error  `json:"error,omitempty" form:"error"`
}

type CoreStatus struct {
	Version *string `json:"version,omitempty" form:"version"`
	Error   *Error  `json:"error,omitempty" form:"error"`
}
