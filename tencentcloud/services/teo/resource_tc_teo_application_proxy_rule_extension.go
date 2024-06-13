package teo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoApplicationProxyRuleReadPostFillRequest0(ctx context.Context, req *teo.DescribeApplicationProxiesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	req.Filters = []*teo.Filter{{
		Name:   helper.String("zone-id"),
		Values: []*string{&zoneId},
	}, {
		Name:   helper.String("proxy-id"),
		Values: []*string{&proxyId},
	}}

	return nil
}

func resourceTencentCloudTeoApplicationProxyRuleCreateRequestOnError0(ctx context.Context, req *teo.CreateApplicationProxyRuleRequest, e error) *resource.RetryError {
	return tccommon.RetryError(e, tccommon.InternalError)
}

func resourceTencentCloudTeoApplicationProxyRuleCreatePostHandleResponse0(ctx context.Context, resp *teo.CreateApplicationProxyRuleResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	var proxyId string
	if v, ok := d.GetOk("proxy_id"); ok {
		proxyId = v.(string)
	}
	var ruleId string
	if resp.Response.RuleId != nil {
		ruleId = *resp.Response.RuleId
	}

	return waitForTeoApplicationProxyRuleOnline(ctx, zoneId, proxyId, ruleId)
}

func resourceTencentCloudTeoApplicationProxyRuleUpdateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]
	ruleId := idSplit[2]

	return waitForTeoApplicationProxyRuleOnline(ctx, zoneId, proxyId, ruleId)
}

func resourceTencentCloudTeoApplicationProxyRuleDeletePostHandleResponse0(ctx context.Context, resp *teo.ModifyApplicationProxyRuleStatusResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]
	ruleId := idSplit[2]

	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	return resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if *instance.Status == "offline" {
			return nil
		}
		if *instance.Status == "stopping" {
			return resource.RetryableError(fmt.Errorf("applicationProxyRule stopping"))
		}

		return resource.RetryableError(fmt.Errorf("setting applicationProxyRule `status` to offline"))
	})
}

func waitForTeoApplicationProxyRuleOnline(ctx context.Context, zoneId, proxyId, ruleId string) error {
	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	return resource.Retry(60*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if *instance.Status == "online" {
			return nil
		}
		if *instance.Status == "fail" {
			return resource.NonRetryableError(fmt.Errorf("applicationProxyRule status is %v, operate failed.", *instance.Status))
		}
		return resource.RetryableError(fmt.Errorf("applicationProxyRule status is %v, retry...", *instance.Status))
	})
}
