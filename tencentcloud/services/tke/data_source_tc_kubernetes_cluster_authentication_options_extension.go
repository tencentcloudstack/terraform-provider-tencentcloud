package tke

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func dataSourceTencentCloudKubernetesClusterAuthenticationOptionsReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	var tmpList []map[string]interface{}

	if v, ok := d.GetOk("service_accounts"); ok {
		if vList, isList := v.([]interface{}); isList && len(vList) > 0 {
			tmpList = append(tmpList, vList[0].(map[string]interface{}))
		}
	}
	if v, ok := d.GetOk("oidc_config"); ok {
		if vList, isList := v.([]interface{}); isList && len(vList) > 0 {
			tmpList = append(tmpList, vList[0].(map[string]interface{}))
		}
	}

	return tmpList
}
