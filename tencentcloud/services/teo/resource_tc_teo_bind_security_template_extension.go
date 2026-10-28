package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
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

func resourceTeoBindSecurityTemplateDeleteStateRefreshFunc_0_0(ctx context.Context, zoneId string, templateId string, entity string) resource.StateRefreshFunc {
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
		// Only when the query cannot find the binding is the unbind complete.
		// A non-nil result must be returned together with the "deleted" state,
		// otherwise the SDK treats a nil result as NotFound and keeps polling
		// until it raises NotFoundError/TimeoutError.
		if resp == nil {
			return &teov20220901.EntityStatus{}, "deleted", nil
		}
		// The binding still exists. If Status is absent, keep polling with a
		// pending state instead of prematurely treating it as deleted.
		if resp.Status == nil {
			return resp, "process", nil
		}
		return resp, *resp.Status, nil
	}
}
