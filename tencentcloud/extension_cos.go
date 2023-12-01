package tencentcloud

const (
	COS_ACL_GRANTEE_TYPE_USER      = "CanonicalUser"
	COS_ACL_GRANTEE_TYPE_ANONYMOUS = "Group"
)

var COSACLGranteeTypeSeq = []string{
	COS_ACL_GRANTEE_TYPE_USER,
	COS_ACL_GRANTEE_TYPE_ANONYMOUS,
}

var COSACLPermissionSeq = []string{
	"READ",
	"WRITE",
	"FULL_CONTROL",
	"WRITE_ACP",
	"READ_ACP",
}
