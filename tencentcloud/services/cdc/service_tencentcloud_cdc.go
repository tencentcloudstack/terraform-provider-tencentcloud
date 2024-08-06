package cdc

import (
	"context"
	"log"

	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
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
	logId := tccommon.GetLogId(ctx)

	request := cdc.NewDeleteSitesRequest()
	request.SiteIds = helper.Strings([]string{siteId})

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

func (me *CdcService) DescribeCdcDedicatedClusterById(ctx context.Context, dedicatedClusterId string) (dedicatedCluster *cdc.DedicatedCluster, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdc.NewDescribeDedicatedClustersRequest()
	request.DedicatedClusterIds = helper.Strings([]string{dedicatedClusterId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdcClient().DescribeDedicatedClusters(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.DedicatedClusterSet) < 1 {
		return
	}

	dedicatedCluster = response.Response.DedicatedClusterSet[0]
	return
}

func (me *CdcService) DeleteCdcDedicatedClusterById(ctx context.Context, dedicatedClusterId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdc.NewDeleteDedicatedClustersRequest()
	request.DedicatedClusterIds = helper.Strings([]string{dedicatedClusterId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdcClient().DeleteDedicatedClusters(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdcService) DescribeCdcHostByFilter(ctx context.Context, param map[string]interface{}) (hostList []*cdc.HostInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdc.NewDescribeDedicatedClusterHostsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DedicatedClusterId" {
			request.DedicatedClusterId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdcClient().DescribeDedicatedClusterHosts(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.HostInfoSet) < 1 {
			break
		}

		hostList = append(hostList, response.Response.HostInfoSet...)
		offset += limit
	}

	return
}

func (me *CdcService) DescribeCdcDedicatedClusterInstanceTypesByFilter(ctx context.Context, param map[string]interface{}) (DedicatedClusterInstanceTypes []*cdc.DedicatedClusterInstanceType, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdc.NewDescribeDedicatedClusterInstanceTypesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DedicatedClusterId" {
			request.DedicatedClusterId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdcClient().DescribeDedicatedClusterInstanceTypes(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DedicatedClusterInstanceTypeSet) < 1 {
		return
	}

	DedicatedClusterInstanceTypes = response.Response.DedicatedClusterInstanceTypeSet
	return
}

func (me *CdcService) DescribeCdcDedicatedClusterOrdersByFilter(ctx context.Context, param map[string]interface{}) (dedicatedClusterOrders []*cdc.DedicatedClusterOrder, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdc.NewDescribeDedicatedClusterOrdersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DedicatedClusterIds" {
			request.DedicatedClusterIds = v.([]*string)
		}

		//if k == "DedicatedClusterOrderIds" {
		//	request.DedicatedClusterOrderIds = v.(*string)
		//}

		if k == "Status" {
			request.Status = v.(*string)
		}

		if k == "ActionType" {
			request.ActionType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdcClient().DescribeDedicatedClusterOrders(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DedicatedClusterOrderSet) < 1 {
		return
	}

	dedicatedClusterOrders = response.Response.DedicatedClusterOrderSet
	return
}
