package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTeoBindSecurityTemplateCreateStateRefreshFunc_0_0(ctx context.Context, zoneId string, templateId string, entity string) resource.StateRefreshFunc {
	var req *teov20220901.DescribeSecurityTemplateBindingsRequest
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)

		service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}
		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}
			_ = d
			req = teov20220901.NewDescribeSecurityTemplateBindingsRequest()
		}
		resp, err := service.DescribeTeoBindSecurityTemplateById(ctx, zoneId, templateId, entity)
		if err != nil {
			return nil, "", err
		}
		if resp == nil {
			return nil, "", nil
		}
		return resp, *resp.Status, nil
	}
}
