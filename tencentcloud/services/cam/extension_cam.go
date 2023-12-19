package cam

const PAGE_ITEM = 200

const (
	CAM_POLICY_CREATE_STRATEGY_CUSTOM = "User"
	CAM_POLICY_CREATE_STRATEGY_PRESET = "QCS"
	CAM_POLICY_CREATE_STRATEGY_NULL   = ""
)

var CAM_POLICY_CREATE_STRATEGY = []string{
	CAM_POLICY_CREATE_STRATEGY_CUSTOM,
	CAM_POLICY_CREATE_STRATEGY_PRESET,
	CAM_POLICY_CREATE_STRATEGY_NULL,
}

type Principal struct {
	Service []string `json:"service"`
}
type Statement struct {
	Principal Principal `json:"principal"`
}
type Document struct {
	Version   string      `json:"version"`
	Statement []Statement `json:"statement"`
}
