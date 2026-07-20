package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTeoBindSecurityTemplateCreateStateRefreshFunc_0_0(ctx context.Context, zoneId string, templateId string, entity string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}

		service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resp, err := service.DescribeTeoBindSecurityTemplateById(ctx, zoneId, templateId, entity)
		if err != nil {
			return nil, "", err
		}
		if resp == nil {
			return nil, "", nil
		}
		if resp.Status == nil {
			return resp, "", nil
		}
		return resp, *resp.Status, nil
	}
}
