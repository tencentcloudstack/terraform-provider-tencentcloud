package tke

import (
	"context"
	"fmt"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKubernetesClusterCommonNamesReadPreRequest0(ctx context.Context, req *tke.DescribeClusterCommonNamesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	if v, ok := d.GetOk("subaccount_uins"); ok {
		req.SubaccountUins = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("role_ids"); ok {
		req.RoleIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	return nil
}
