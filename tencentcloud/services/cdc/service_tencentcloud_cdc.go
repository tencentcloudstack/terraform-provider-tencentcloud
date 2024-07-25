package cdc

import (
	"context"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type CdcService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdcService) DescribeCdcSiteDetailById(ctx context.Context, siteId string) (siteDetail *cdc.SiteDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdc.NewDescribeSitesDetailRequest()
	request.SiteIds = helper.Strings([]string{siteId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCdcClient().DescribeSitesDetail(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.SiteDetailSet) < 1 {
		return
	}

	siteDetail = response.Response.SiteDetailSet[0]
	return
}

func (me *CdcService) DeleteCdcSiteById(ctx context.Context, siteId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdc.NewDeleteSitesRequest()
	request.SiteId = &siteId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdcClient().DeleteSites(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
