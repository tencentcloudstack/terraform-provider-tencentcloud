package makers

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	makersv20261113 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewMakersService(client *connectivity.TencentCloudClient) MakersService {
	return MakersService{client: client}
}

type MakersService struct {
	client *connectivity.TencentCloudClient
}

// CreateZoneCustomDomain binds a custom domain to a Makers zone and returns
// the domain echoed back by the cloud side. The cloud-side id is not used as
// a resource identity: all follow-up APIs address the domain by the composite
// key (zone id, project id, domain name).
func (me *MakersService) CreateZoneCustomDomain(ctx context.Context, request *makersv20261113.CreateZoneCustomDomainRequest) (domain string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makersv20261113.CreateZoneCustomDomainResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersV20261113Client().CreateZoneCustomDomainWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateZoneCustomDomain returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Domain == nil || *response.Response.Domain == "" {
		errRet = fmt.Errorf("CreateZoneCustomDomain returned an empty domain")
		return
	}

	domain = *response.Response.Domain
	return
}

// DescribeZoneCustomDomain queries a single custom domain by the composite
// key (zone id, project id, domain name). Only `Custom` type domains are
// considered managed; `Preset` domains are platform built-ins.
//
// A nil result with a nil error means "queried successfully, no such domain";
// an API failure is always reported as an error so the caller never confuses
// a failed request with a deleted resource.
func (me *MakersService) DescribeZoneCustomDomain(ctx context.Context, zoneId, projectId, domain string) (zoneDomain *makersv20261113.ZoneDomain, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := makersv20261113.NewDescribeZoneCustomDomainsRequest()
	request.ZoneId = &zoneId
	request.ProjectId = &projectId
	request.CustomDomains = []*makersv20261113.CustomDomainInput{
		{
			Domain: &domain,
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makersv20261113.DescribeZoneCustomDomainsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersV20261113Client().DescribeZoneCustomDomainsWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeZoneCustomDomains returned an empty response for domain %s", domain)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, item := range response.Response.PagesDomains {
		if item == nil || item.Domain == nil || *item.Domain != domain {
			continue
		}
		if item.Type != nil && *item.Type != "Custom" {
			continue
		}
		zoneDomain = item
		return
	}

	return
}

// ModifyZoneCustomDomainEnv switches the environment a custom domain is
// bound to. This is the only in-place update the cloud API supports; any
// other change is handled by replacement in the Terraform layer.
func (me *MakersService) ModifyZoneCustomDomainEnv(ctx context.Context, request *makersv20261113.ModifyZoneCustomDomainEnvRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makersv20261113.ModifyZoneCustomDomainEnvResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersV20261113Client().ModifyZoneCustomDomainEnvWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("ModifyZoneCustomDomainEnv returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// DeleteProjectCustomDomains unbinds the given custom domain from the
// project. The cloud API is batch-shaped; the Terraform layer always sends
// exactly one element so destroy only affects the managed domain.
func (me *MakersService) DeleteProjectCustomDomain(ctx context.Context, zoneId, projectId, domain string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := makersv20261113.NewDeleteProjectCustomDomainsRequest()
	request.ZoneId = &zoneId
	request.ProjectId = &projectId
	request.CustomDomains = []*string{&domain}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makersv20261113.DeleteProjectCustomDomainsResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersV20261113Client().DeleteProjectCustomDomainsWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DeleteProjectCustomDomains returned an empty response for domain %s", domain)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
