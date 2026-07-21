package mariadb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

// resourceTencentCloudMariadbParametersReadPreHandleResponse0 filters the API
// response so that only the parameters configured by the user are kept. This
// avoids drifting state with the full parameter list returned by the API.
func resourceTencentCloudMariadbParametersReadPreHandleResponse0(ctx context.Context, resp *mariadb.DescribeDBParametersResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if v, ok := d.GetOk("params"); ok {
		params := []*mariadb.ParamDesc{}
		for _, item := range v.(*schema.Set).List() {
			paramsMap := item.(map[string]interface{})
			if name, ok := paramsMap["param"].(string); ok && name != "" {
				if resp.Params != nil {
					for _, paramDesc := range resp.Params {
						if paramDesc.Param != nil && *paramDesc.Param == name {
							params = append(params, paramDesc)
						}
					}
				}
			}
		}
		resp.Params = params
	}

	return nil
}
