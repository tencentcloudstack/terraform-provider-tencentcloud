package tke

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
			if vList[0] == nil {
				tmpList = append(tmpList, map[string]interface{}{})
			} else {
				tmpList = append(tmpList, vList[0].(map[string]interface{}))
			}
		}
	}
	if v, ok := d.GetOk("oidc_config"); ok {
		if vList, isList := v.([]interface{}); isList && len(vList) > 0 {
			if vList[0] == nil {
				tmpList = append(tmpList, map[string]interface{}{})
			} else {
				vMap := vList[0].(map[string]interface{})
				autoCreateClientID, _ := vMap["auto_create_client_id"].(*schema.Set)
				vMap["auto_create_client_id"] = autoCreateClientID.List()
				tmpList = append(tmpList, vMap)
			}
		}
	}

	return tmpList
}
