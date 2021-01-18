package request

type CasbinInfo struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// Casbin structure for input parameters
type CasbinInReceive struct {
	RoleId      string       `json:"roleId"`
	CasbinInfos []CasbinInfo `json:"casbinInfos"`
}
