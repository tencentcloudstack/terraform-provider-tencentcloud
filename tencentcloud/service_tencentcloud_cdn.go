package tencentcloud

import (
	"context"
	"log"
	"strings"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CdnService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdnService) DescribeDomainsConfigByDomain(ctx context.Context, domain string) (domainConfig *cdn.DetailDomain, errRet error) {
	logId := getLogId(ctx)
	request := cdn.NewDescribeDomainsConfigRequest()
	request.Filters = make([]*cdn.DomainFilter, 0, 1)
	filter := &cdn.DomainFilter{
		Name:  helper.String("domain"),
		Value: []*string{&domain},
	}
	request.Filters = append(request.Filters, filter)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().DescribeDomainsConfig(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			if sdkErr.Code == CDN_HOST_NOT_FOUND {
				return
			}
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	if len(response.Response.Domains) > 0 {
		domainConfig = response.Response.Domains[0]
	}
	return
}

func (me *CdnService) DeleteDomain(ctx context.Context, domain string) error {
	logId := getLogId(ctx)
	request := cdn.NewDeleteCdnDomainRequest()
	request.Domain = &domain

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCdnClient().DeleteCdnDomain(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CdnService) StopDomain(ctx context.Context, domain string) error {
	logId := getLogId(ctx)
	request := cdn.NewStopCdnDomainRequest()
	request.Domain = &domain

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCdnClient().StopCdnDomain(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CdnService) StartDomain(ctx context.Context, domain string) error {
	logId := getLogId(ctx)
	request := cdn.NewStartCdnDomainRequest()
	request.Domain = &domain

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCdnClient().StartCdnDomain(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CdnService) DescribeDomainsConfigByFilters(ctx context.Context, filterMap map[string]interface{}) (domainConfig []*cdn.DetailDomain, errRet error) {
	logId := getLogId(ctx)
	request := cdn.NewDescribeDomainsConfigRequest()
	request.Filters = make([]*cdn.DomainFilter, 0, len(filterMap))

	for k, v := range filterMap {
		value := v.(string)

		filter := &cdn.DomainFilter{
			Name:  helper.String(underlineToHump([]byte(k))),
			Value: []*string{&value},
		}
		request.Filters = append(request.Filters, filter)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdnClient().DescribeDomainsConfig(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			if sdkErr.Code == CDN_HOST_NOT_FOUND {
				return
			}
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	if len(response.Response.Domains) > 0 {
		domainConfig = response.Response.Domains
	}
	return
}

func underlineToHump(underline []byte) (humpValue string) {
	lenUnderLine := len(underline)
	for i := 0; i < lenUnderLine; i++ {
		if string(underline[i]) == "_" {
			if i+1 < lenUnderLine {
				humpValue += strings.ToUpper(string(underline[i+1]))
				i++
			}
			continue
		}

		humpValue += string(underline[i])
	}

	return
}
