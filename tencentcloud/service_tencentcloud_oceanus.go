package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type OceanusService struct {
	client *connectivity.TencentCloudClient
}

func (me *OceanusService) DescribeOceanusResourceRelatedJobByFilter(ctx context.Context, param map[string]interface{}) (ResourceRelatedJob []*oceanus.ResourceRefJobInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = oceanus.NewDescribeResourceRelatedJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ResourceId" {
			request.ResourceId = v.(*string)
		}

		if k == "DESCByJobConfigCreateTime" {
			request.DESCByJobConfigCreateTime = v.(*int64)
		}

		if k == "ResourceConfigVersion" {
			request.ResourceConfigVersion = v.(*int64)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeResourceRelatedJobs(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RefJobInfos) < 1 {
			break
		}

		ResourceRelatedJob = append(ResourceRelatedJob, response.Response.RefJobInfos...)
		if len(response.Response.RefJobInfos) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusSavepointListByFilter(ctx context.Context, param map[string]interface{}) (savepointList []*oceanus.Savepoint, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = oceanus.NewDescribeJobSavepointRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "JobId" {
			request.JobId = v.(*string)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeJobSavepoint(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Savepoint) < 1 {
			break
		}

		savepointList = append(savepointList, response.Response.Savepoint...)
		if len(response.Response.Savepoint) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusSystemResourceByFilter(ctx context.Context, param map[string]interface{}) (SystemResource []*oceanus.SystemResourceItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = oceanus.NewDescribeSystemResourcesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ResourceIds" {
			request.ResourceIds = v.([]*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*oceanus.Filter)
		}

		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}

		if k == "FlinkVersion" {
			request.FlinkVersion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeSystemResources(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ResourceSet) < 1 {
			break
		}

		SystemResource = append(SystemResource, response.Response.ResourceSet...)
		if len(response.Response.ResourceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusWorkSpacesByFilter(ctx context.Context, param map[string]interface{}) (WorkSpace []*oceanus.WorkSpaceSetItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = oceanus.NewDescribeWorkSpacesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}

		if k == "Filters" {
			request.Filters = v.([]*oceanus.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeWorkSpaces(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.WorkSpaceSetItem) < 1 {
			break
		}

		WorkSpace = append(WorkSpace, response.Response.WorkSpaceSetItem...)
		if len(response.Response.WorkSpaceSetItem) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusJobById(ctx context.Context, jobId string) (Job *oceanus.JobV1, errRet error) {
	logId := getLogId(ctx)

	request := oceanus.NewDescribeJobsRequest()
	request.JobIds = common.StringPtrs([]string{jobId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeJobs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	Job = response.Response.JobSet[0]
	return
}

func (me *OceanusService) DeleteOceanusJobById(ctx context.Context, jobId string) (errRet error) {
	logId := getLogId(ctx)

	request := oceanus.NewDeleteJobsRequest()
	request.JobIds = common.StringPtrs([]string{jobId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteJobs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusResourceById(ctx context.Context, resourceId string) (resource *oceanus.ResourceItem, errRet error) {
	logId := getLogId(ctx)

	request := oceanus.NewDescribeResourcesRequest()
	request.ResourceIds = common.StringPtrs([]string{resourceId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeResources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	resource = response.Response.ResourceSet[0]
	return
}

func (me *OceanusService) DeleteOceanusResourceById(ctx context.Context, resourceId string) (errRet error) {
	logId := getLogId(ctx)

	request := oceanus.NewDeleteResourcesRequest()
	request.ResourceIds = common.StringPtrs([]string{resourceId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteResources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusResourceConfigById(ctx context.Context, resourceId, version string) (resourceConfig *oceanus.ResourceConfigItem, errRet error) {
	logId := getLogId(ctx)

	request := oceanus.NewDescribeResourceConfigsRequest()
	request.ResourceId = &resourceId
	versionInt, _ := strconv.ParseInt(version, 10, 64)
	request.ResourceConfigVersions = common.Int64Ptrs([]int64{versionInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeResourceConfigs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	resourceConfig = response.Response.ResourceConfigSet[0]
	return
}

func (me *OceanusService) DeleteOceanusResourceConfigById(ctx context.Context, resourceId, version string) (errRet error) {
	logId := getLogId(ctx)

	request := oceanus.NewDeleteResourceConfigsRequest()
	request.ResourceId = &resourceId
	versionInt, _ := strconv.ParseInt(version, 10, 64)
	request.ResourceConfigVersions = common.Int64Ptrs([]int64{versionInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteResourceConfigs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusClustersByFilter(ctx context.Context, param map[string]interface{}) (Clusters []*oceanus.Cluster, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = oceanus.NewDescribeClustersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterIds" {
			request.ClusterIds = v.([]*string)
		}

		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}

		if k == "Filters" {
			request.Filters = v.([]*oceanus.Filter)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeClusters(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClusterSet) < 1 {
			break
		}

		Clusters = append(Clusters, response.Response.ClusterSet...)
		if len(response.Response.ClusterSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusTreeJobsByFilter(ctx context.Context, param map[string]interface{}) (treeJobs *oceanus.DescribeTreeJobsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = oceanus.NewDescribeTreeJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeTreeJobs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	treeJobs = response.Response
	return
}
