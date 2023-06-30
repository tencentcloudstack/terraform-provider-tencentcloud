package tencentcloud

const (
	COS_ACL_GRANTEE_TYPE_USER      = "CanonicalUser"
	COS_ACL_GRANTEE_TYPE_ANONYMOUS = "Group"
)

var COSACLPermissionMap = map[string]string{
	"PERMISSON_READ":         "READ",
	"PERMISSON_WRITE":        "WRITE",
	"PERMISSON_FULL_CONTROL": "FULL_CONTROL",
	"PERMISSON_WRITE_ACP":    "WRITE_ACP",
	"PERMISSON_READ_ACP":     "READ_ACP",
}
