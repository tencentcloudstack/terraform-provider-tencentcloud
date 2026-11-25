package oceanus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type OceanusService struct {
	client *connectivity.TencentCloudClient
}

func (me *OceanusService) DescribeOceanusResourceRelatedJobByFilter(ctx context.Context, param map[string]interface{}) (ResourceRelatedJob []*oceanus.ResourceRefJobInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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

		if k == "Filters" {
			request.Filters = v.([]*oceanus.Filter)
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

func (me *OceanusService) DescribeOceanusTreeResourcesByFilter(ctx context.Context, param map[string]interface{}) (treeResources *oceanus.DescribeTreeResourcesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeTreeResourcesRequest()
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

	response, err := me.client.UseOceanusClient().DescribeTreeResources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	treeResources = response.Response
	return
}

func (me *OceanusService) DescribeOceanusJobSubmissionLogByFilter(ctx context.Context, param map[string]interface{}) (jobSubmissionLog *oceanus.DescribeJobSubmissionLogResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeJobSubmissionLogRequest()
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

		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}

		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}

		if k == "RunningOrderId" {
			request.RunningOrderId = v.(*int64)
		}

		if k == "Keyword" {
			request.Keyword = v.(*string)
		}

		if k == "OrderType" {
			request.OrderType = v.(*string)
		}

		if k == "Cursor" {
			request.Cursor = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	request.Limit = common.Int64Ptr(100)
	response, err := me.client.UseOceanusClient().DescribeJobSubmissionLog(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	jobSubmissionLog = response.Response

	return
}

func (me *OceanusService) DescribeOceanusCheckSavepointByFilter(ctx context.Context, param map[string]interface{}) (CheckSavepoint *oceanus.CheckSavepointResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewCheckSavepointRequest()
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

		if k == "SerialId" {
			request.SerialId = v.(*string)
		}

		if k == "RecordType" {
			request.RecordType = v.(*int64)
		}

		if k == "SavepointPath" {
			request.SavepointPath = v.(*string)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().CheckSavepoint(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	CheckSavepoint = response.Response
	return
}

func (me *OceanusService) DescribeOceanusWorkSpaceById(ctx context.Context, workSpaceName string) (WorkSpace *oceanus.WorkSpaceSetItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeWorkSpacesRequest()
	request.Filters = []*oceanus.Filter{
		{
			Name:   common.StringPtr("WorkSpaceName"),
			Values: common.StringPtrs([]string{workSpaceName}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeWorkSpaces(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.WorkSpaceSetItem) != 1 {
		return
	}

	WorkSpace = response.Response.WorkSpaceSetItem[0]
	return
}

func (me *OceanusService) DeleteOceanusWorkSpaceById(ctx context.Context, workSpaceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteWorkSpaceRequest()
	request.WorkSpaceId = &workSpaceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteWorkSpace(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusJobConfigById(ctx context.Context, jobId, version string) (JobConfig *oceanus.JobConfig, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeJobConfigsRequest()
	request.JobId = &jobId
	versionInt, _ := strconv.ParseUint(version, 10, 64)
	request.JobConfigVersions = common.Uint64Ptrs([]uint64{versionInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeJobConfigs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.JobConfigSet) != 1 {
		return
	}

	JobConfig = response.Response.JobConfigSet[0]
	return
}

func (me *OceanusService) DeleteOceanusJobConfigById(ctx context.Context, jobId, version string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteJobConfigsRequest()
	request.JobId = &jobId
	versionInt, _ := strconv.ParseInt(version, 10, 64)
	request.JobConfigVersions = common.Int64Ptrs([]int64{versionInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteJobConfigs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusFolderById(ctx context.Context, workSpaceId, folderId, folderType string) (Folder *oceanus.DescribeFolderResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeFolderRequest()
	request.FolderId = common.StringPtr(folderId)
	folderTypeInt, _ := strconv.ParseInt(folderType, 10, 64)
	request.FolderType = common.Int64Ptr(folderTypeInt)
	request.WorkSpaceId = common.StringPtr(workSpaceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeFolder(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	Folder = response.Response
	return
}

func (me *OceanusService) DeleteOceanusFolderById(ctx context.Context, workSpaceId, folderId, folderType string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteFoldersRequest()
	request.FolderIds = common.StringPtrs([]string{folderId})
	folderTypeInt, _ := strconv.ParseInt(folderType, 10, 64)
	request.FolderType = common.Int64Ptr(folderTypeInt)
	request.WorkSpaceId = common.StringPtr(workSpaceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteFolders(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusJobEventsByFilter(ctx context.Context, param map[string]interface{}) (JobEvents *oceanus.DescribeJobEventsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeJobEventsRequest()
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

		if k == "StartTimestamp" {
			request.StartTimestamp = v.(*uint64)
		}

		if k == "EndTimestamp" {
			request.EndTimestamp = v.(*uint64)
		}

		if k == "Types" {
			request.Types = v.([]*string)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeJobEvents(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	JobEvents = response.Response
	return
}

func (me *OceanusService) DescribeOceanusMetaTableByFilter(ctx context.Context, param map[string]interface{}) (MetaTable *oceanus.GetMetaTableResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewGetMetaTableRequest()
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

		if k == "Catalog" {
			request.Catalog = v.(*string)
		}

		if k == "Database" {
			request.Database = v.(*string)
		}

		if k == "Table" {
			request.Table = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().GetMetaTable(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	MetaTable = response.Response
	return
}

// Oceanus cluster status values reported by DescribeClusters. Only 1
// (not initialized), 3 (initializing) and 2 (running) are documented by the
// cloud API; any other value is treated as an unknown transition state and
// keeps the caller waiting until the poll timeout.
const (
	OceanusClusterStatusNotInitialized = 1
	OceanusClusterStatusInitializing   = 3
	OceanusClusterStatusRunning        = 2

	// OceanusClusterTaskSuccess is the TaskExecResult value the write APIs
	// report once the operation was accepted.
	OceanusClusterTaskSuccess = "success"
)

// CreateOceanusCluster creates a cluster and returns the id assigned by the
// cloud side.
func (me *OceanusService) CreateOceanusCluster(ctx context.Context, request *oceanus.CreateOceanusClusterRequest) (clusterId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *oceanus.CreateOceanusClusterResponse
	err := sdkresource.Retry(tccommon.WriteRetryTimeout, func() *sdkresource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseOceanusClient().CreateOceanusClusterWithContext(ctx, request)
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
		errRet = fmt.Errorf("CreateOceanusCluster returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ClusterId == nil || *response.Response.ClusterId == "" {
		errRet = fmt.Errorf("CreateOceanusCluster returned an empty cluster id")
		return
	}

	clusterId = *response.Response.ClusterId
	return
}

// DescribeOceanusClusterById returns the cluster with the given id.
//
// A nil result with a nil error means "queried successfully, no such cluster";
// an API failure is always reported as an error so the caller never confuses
// a failed request with a deleted resource.
func (me *OceanusService) DescribeOceanusClusterById(ctx context.Context, clusterId string) (cluster *oceanus.Cluster, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeClustersRequest()
	request.ClusterIds = []*string{helper.String(clusterId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *oceanus.DescribeClustersResponse
	err := sdkresource.Retry(tccommon.ReadRetryTimeout, func() *sdkresource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseOceanusClient().DescribeClustersWithContext(ctx, request)
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
		errRet = fmt.Errorf("DescribeClusters returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, item := range response.Response.ClusterSet {
		if item == nil {
			continue
		}
		if item.ClusterId == nil || *item.ClusterId == clusterId {
			cluster = item
			return
		}
	}
	return
}

// ScaleOceanusCluster resizes an existing cluster to the new CU count. The
// cloud API is asynchronous: the operation is only accepted when it reports
// a successful TaskExecResult.
func (me *OceanusService) ScaleOceanusCluster(ctx context.Context, request *oceanus.ScaleOceanusClusterRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *oceanus.ScaleOceanusClusterResponse
	err := sdkresource.Retry(tccommon.WriteRetryTimeout, func() *sdkresource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseOceanusClient().ScaleOceanusClusterWithContext(ctx, request)
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
		errRet = fmt.Errorf("ScaleOceanusCluster returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.TaskExecResult == nil || *response.Response.TaskExecResult != OceanusClusterTaskSuccess {
		result := "<nil>"
		if response.Response.TaskExecResult != nil {
			result = *response.Response.TaskExecResult
		}
		errRet = fmt.Errorf("ScaleOceanusCluster did not report a successful task execution result, got %q", result)
	}
	return
}

// RenewOceanusCluster renews a prepaid cluster.
func (me *OceanusService) RenewOceanusCluster(ctx context.Context, request *oceanus.RenewOceanusClusterRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *oceanus.RenewOceanusClusterResponse
	err := sdkresource.Retry(tccommon.WriteRetryTimeout, func() *sdkresource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseOceanusClient().RenewOceanusClusterWithContext(ctx, request)
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
		errRet = fmt.Errorf("RenewOceanusCluster returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

// DeleteOceanusCluster deletes a cluster.
func (me *OceanusService) DeleteOceanusCluster(ctx context.Context, request *oceanus.DeleteOceanusClusterRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *oceanus.DeleteOceanusClusterResponse
	err := sdkresource.Retry(tccommon.WriteRetryTimeout, func() *sdkresource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseOceanusClient().DeleteOceanusClusterWithContext(ctx, request)
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
		errRet = fmt.Errorf("DeleteOceanusCluster returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.TaskExecResult == nil || *response.Response.TaskExecResult != OceanusClusterTaskSuccess {
		result := "<nil>"
		if response.Response.TaskExecResult != nil {
			result = *response.Response.TaskExecResult
		}
		errRet = fmt.Errorf("DeleteOceanusCluster did not report a successful task execution result, got %q", result)
	}
	return
}

// WaitOceanusClusterReady polls DescribeClusters until the cluster reports
// status 2 (running).
func (me *OceanusService) WaitOceanusClusterReady(ctx context.Context, clusterId string, timeout, delay, minTimeout time.Duration) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"INITIALIZING"},
		Target:  []string{"RUNNING"},
		Refresh: func() (interface{}, string, error) {
			cluster, err := me.DescribeOceanusClusterById(ctx, clusterId)
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return nil, "MISSING", fmt.Errorf("cluster %s not found", clusterId)
			}
			if cluster.Status == nil {
				return cluster, "INITIALIZING", nil
			}
			switch *cluster.Status {
			case OceanusClusterStatusRunning:
				return cluster, "RUNNING", nil
			default:
				// Undocumented transition state: keep waiting instead of
				// guessing an extra enum value.
				return cluster, "INITIALIZING", nil
			}
		},
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: minTimeout,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
