package rum

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewRumService(client *connectivity.TencentCloudClient) RumService {
	return RumService{client: client}
}

type RumService struct {
	client *connectivity.TencentCloudClient
}

func (me *RumService) DescribeRumTawInstance(ctx context.Context, instanceId string) (tawInstance *rum.RumInstanceInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeTawInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&rum.Filter{
			Name:   helper.String("InstanceIDs"),
			Values: []*string{&instanceId},
		},
	)

	response, err := me.client.UseRumClient().DescribeTawInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.InstanceSet) < 1 {
		return
	}
	tawInstance = response.Response.InstanceSet[0]
	return
}

func (me *RumService) DeleteRumTawInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := rum.NewDeleteInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DeleteInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *RumService) DescribeRumProject(ctx context.Context, id string) (project *rum.RumProject, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeProjectsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&rum.Filter{
			Name:   helper.String("ProjectID"),
			Values: []*string{&id},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances := make([]*rum.RumProject, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseRumClient().DescribeProjects(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ProjectSet) < 1 {
			break
		}
		instances = append(instances, response.Response.ProjectSet...)
		if len(response.Response.ProjectSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	project = instances[0]

	return

}

func (me *RumService) DeleteRumProjectById(ctx context.Context, id string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := rum.NewDeleteProjectRequest()

	projectId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("[DEBUG]%s api[%s] sting to uint64 error, err [%s]", logId, request.GetAction(), err)
		errRet = err
		return
	}

	request.ID = helper.Uint64(uint64(projectId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DeleteProject(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *RumService) DescribeRumWhitelist(ctx context.Context, instanceID, id string) (whitelist *rum.Whitelist, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeWhitelistsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceID = &instanceID

	response, err := me.client.UseRumClient().DescribeWhitelists(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.WhitelistSet) < 1 {
		return
	}
	if response != nil && len(response.Response.WhitelistSet) > 0 {
		for _, v := range response.Response.WhitelistSet {
			if *v.ID == id {
				whitelist = v
				break
			}
		}
	}
	return
}

func (me *RumService) DeleteRumWhitelistById(ctx context.Context, instanceID, id string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := rum.NewDeleteWhitelistRequest()

	request.InstanceID = &instanceID
	request.ID = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DeleteWhitelist(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *RumService) DescribeRumOfflineLogConfigAttachment(ctx context.Context, projectKey, uniqueId string) (offlineLogConfig *rum.DescribeOfflineLogConfigsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeOfflineLogConfigsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ProjectKey = &projectKey

	response, err := me.client.UseRumClient().DescribeOfflineLogConfigs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	offlineLogConfig = response.Response
	return
}

func (me *RumService) DeleteRumOfflineLogConfigAttachmentById(ctx context.Context, projectKey, uniqueId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := rum.NewDeleteOfflineLogConfigRequest()

	request.ProjectKey = &projectKey
	request.UniqueID = &uniqueId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DeleteOfflineLogConfig(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *RumService) DescribeRumOfflineLogConfigByFilter(ctx context.Context, param map[string]interface{}) (configs *rum.DescribeOfflineLogConfigsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeOfflineLogConfigsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "project_key" {
			request.ProjectKey = v.(*string)
		}
	}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseRumClient().DescribeOfflineLogConfigs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		configs = response.Response
	}

	return
}

func (me *RumService) DescribeRumProjectByFilter(ctx context.Context, param map[string]interface{}) (project []*rum.RumProject, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeProjectsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.Filters = append(
				request.Filters,
				&rum.Filter{
					Name:   helper.String("InstanceID"),
					Values: []*string{v.(*string)},
				},
			)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseRumClient().DescribeProjects(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ProjectSet) < 1 {
			break
		}
		project = append(project, response.Response.ProjectSet...)
		if len(response.Response.ProjectSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *RumService) DescribeRumWhitelistByFilter(ctx context.Context, param map[string]interface{}) (whitelist []*rum.Whitelist, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeWhitelistsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceID = v.(*string)
		}
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeWhitelists(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.WhitelistSet) > 0 {
		whitelist = response.Response.WhitelistSet
	}

	return
}

func (me *RumService) DescribeRumTawInstanceByFilter(ctx context.Context, param map[string]interface{}) (tawInstance []*rum.RumInstanceInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeTawInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "charge_statuses" {
			request.ChargeStatuses = v.([]*int64)
		}
		if k == "charge_types" {
			request.ChargeTypes = v.([]*int64)
		}
		if k == "area_ids" {
			request.AreaIds = v.([]*int64)
		}
		if k == "instance_statuses" {
			request.InstanceStatuses = v.([]*int64)
		}
		if k == "instance_ids" {
			request.Filters = append(
				request.Filters,
				&rum.Filter{
					Name:   helper.String("InstanceIDs"),
					Values: v.([]*string),
				},
			)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseRumClient().DescribeTawInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		tawInstance = append(tawInstance, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *RumService) DescribeRumCustomUrlByFilter(ctx context.Context, param map[string]interface{}) (customUrl *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataCustomUrlRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Url" {
			request.Url = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataCustomUrl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	customUrl = response.Response.Result

	return
}

func (me *RumService) DescribeRumSetUrlStatisticsByFilter(ctx context.Context, param map[string]interface{}) (setUrlStatistics *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataSetUrlStatisticsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
		if k == "PackageType" {
			request.PackageType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataSetUrlStatistics(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	setUrlStatistics = response.Response.Result

	return
}

func (me *RumService) DescribeRumEventUrlByFilter(ctx context.Context, param map[string]interface{}) (eventUrl *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataEventUrlRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "Name" {
			request.Name = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseRumClient().DescribeDataEventUrl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	eventUrl = response.Response.Result

	return
}

func (me *RumService) DescribeRumFetchUrlByFilter(ctx context.Context, param map[string]interface{}) (fetchUrl *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataFetchUrlRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Url" {
			request.Url = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "Ret" {
			request.Ret = v.(*string)
		}
		if k == "NetStatus" {
			request.NetStatus = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataFetchUrl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	fetchUrl = response.Response.Result

	return
}
func (me *RumService) DescribeRumFetchUrlInfoByFilter(ctx context.Context, param map[string]interface{}) (fetchUrlInfo *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataFetchUrlInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Url" {
			request.Url = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataFetchUrlInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	fetchUrlInfo = response.Response.Result

	return
}
func (me *RumService) DescribeRumLogUrlInfoByFilter(ctx context.Context, param map[string]interface{}) (logUrlInfo *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataLogUrlInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataLogUrlInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	logUrlInfo = response.Response.Result

	return
}
func (me *RumService) DescribeRumLogUrlStatisticsByFilter(ctx context.Context, param map[string]interface{}) (logUrlStatistics *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataLogUrlStatisticsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataLogUrlStatistics(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	logUrlStatistics = response.Response.Result

	return
}
func (me *RumService) DescribeRumPerformancePageByFilter(ctx context.Context, param map[string]interface{}) (performancePage *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataPerformancePageRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
		if k == "NetStatus" {
			request.NetStatus = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataPerformancePage(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	performancePage = response.Response.Result

	return
}
func (me *RumService) DescribeRumPvUrlInfoByFilter(ctx context.Context, param map[string]interface{}) (pvUrlInfo *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataPvUrlInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataPvUrlInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	pvUrlInfo = response.Response.Result

	return
}
func (me *RumService) DescribeRumPvUrlStatisticsByFilter(ctx context.Context, param map[string]interface{}) (pvUrlStatistics *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataPvUrlStatisticsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
		if k == "GroupByType" {
			request.GroupByType = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataPvUrlStatistics(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	pvUrlStatistics = response.Response.Result

	return
}
func (me *RumService) DescribeRumReportCountByFilter(ctx context.Context, param map[string]interface{}) (reportCount *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataReportCountRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ReportType" {
			request.ReportType = v.(*string)
		}
		if k == "InstanceID" {
			request.InstanceID = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataReportCount(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	reportCount = response.Response.Result

	return
}

func (me *RumService) DescribeRumStaticProjectByFilter(ctx context.Context, param map[string]interface{}) (staticProject *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataStaticProjectRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Url" {
			request.Url = v.([]*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataStaticProject(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	staticProject = response.Response.Result

	return
}
func (me *RumService) DescribeRumStaticResourceByFilter(ctx context.Context, param map[string]interface{}) (staticResource *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataStaticResourceRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Url" {
			request.Url = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataStaticResource(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	staticResource = response.Response.Result

	return
}
func (me *RumService) DescribeRumStaticUrlByFilter(ctx context.Context, param map[string]interface{}) (staticUrl *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataStaticUrlRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Url" {
			request.Url = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataStaticUrl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	staticUrl = response.Response.Result

	return
}
func (me *RumService) DescribeRumWebVitalsPageByFilter(ctx context.Context, param map[string]interface{}) (webVitalsPage *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeDataWebVitalsPageRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "ExtSecond" {
			request.ExtSecond = v.(*string)
		}
		if k == "Engine" {
			request.Engine = v.(*string)
		}
		if k == "Isp" {
			request.Isp = v.(*string)
		}
		if k == "From" {
			request.From = v.(*string)
		}
		if k == "Level" {
			request.Level = v.(*string)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "Brand" {
			request.Brand = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "VersionNum" {
			request.VersionNum = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
		if k == "ExtThird" {
			request.ExtThird = v.(*string)
		}
		if k == "ExtFirst" {
			request.ExtFirst = v.(*string)
		}
		if k == "NetType" {
			request.NetType = v.(*string)
		}
		if k == "Device" {
			request.Device = v.(*string)
		}
		if k == "IsAbroad" {
			request.IsAbroad = v.(*string)
		}
		if k == "Os" {
			request.Os = v.(*string)
		}
		if k == "Browser" {
			request.Browser = v.(*string)
		}
		if k == "CostType" {
			request.CostType = v.(*string)
		}
		if k == "Env" {
			request.Env = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeDataWebVitalsPage(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	webVitalsPage = response.Response.Result

	return
}

func (me *RumService) DescribeRumGroupLogByFilter(ctx context.Context, param map[string]interface{}) (groupLog *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeRumGroupLogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "Query" {
			request.Query = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "GroupField" {
			request.GroupField = v.(*string)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 20
	)

	request.Page = &offset
	request.Limit = &limit

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeRumGroupLog(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	groupLog = response.Response.Result

	return
}

func (me *RumService) DescribeRumLogListByFilter(ctx context.Context, param map[string]interface{}) (logList *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeRumLogListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "Query" {
			request.Query = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 20
	)

	request.Page = &offset
	request.Limit = &limit

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeRumLogList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	logList = response.Response.Result

	return
}
func (me *RumService) DescribeRumLogStatsLogListByFilter(ctx context.Context, param map[string]interface{}) (logStatsLogList *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeRumStatsLogListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "Query" {
			request.Query = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeRumStatsLogList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	logStatsLogList = response.Response.Result

	return
}
func (me *RumService) DescribeRumSignByFilter(ctx context.Context, param map[string]interface{}) (sign *rum.DescribeReleaseFileSignResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeReleaseFileSignRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Timeout" {
			request.Timeout = v.(*int64)
		}
		if k == "FileType" {
			request.FileType = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeReleaseFileSign(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	sign = response.Response

	return
}
func (me *RumService) DescribeRumScoresByFilter(ctx context.Context, param map[string]interface{}) (scores []*rum.ScoreInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeScoresRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "ID" {
			request.ID = v.(*int64)
		}
		if k == "IsDemo" {
			request.IsDemo = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeScores(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.ScoreSet) < 1 {
		return
	}

	scores = response.Response.ScoreSet

	return
}

func (me *RumService) DescribeRumTawAreaByFilter(ctx context.Context, param map[string]interface{}) (tawArea []*rum.RumAreaInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeTawAreasRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "AreaIds" {
			request.AreaIds = v.([]*int64)
		}
		if k == "AreaKeys" {
			request.AreaKeys = v.([]*string)
		}
		if k == "AreaStatuses" {
			request.AreaStatuses = v.([]*int64)
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
		response, err := me.client.UseRumClient().DescribeTawAreas(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AreaSet) < 1 {
			break
		}
		tawArea = append(tawArea, response.Response.AreaSet...)
		if len(response.Response.AreaSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *RumService) DescribeRumInstanceStatusConfigById(ctx context.Context, instanceId string) (instanceStatusConfig *rum.RumInstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := rum.NewDescribeTawInstancesRequest()
	request.Filters = append(
		request.Filters,
		&rum.Filter{
			Name:   helper.String("InstanceId"),
			Values: []*string{&instanceId},
		},
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseRumClient().DescribeTawInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceSet) < 1 {
		return
	}

	instanceStatusConfig = response.Response.InstanceSet[0]
	return
}

func (me *RumService) DescribeRumProjectStatusConfigById(ctx context.Context, projectId string) (projectStatusConfig *rum.RumProject, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := rum.NewDescribeProjectsRequest()
	request.Filters = append(
		request.Filters,
		&rum.Filter{
			Name:   helper.String("ID"),
			Values: []*string{&projectId},
		},
	)

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)

	request.Offset = &offset
	request.Limit = &limit

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseRumClient().DescribeProjects(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ProjectSet) < 1 {
		return
	}

	projectStatusConfig = response.Response.ProjectSet[0]
	return
}

func (me *RumService) DescribeRumReleaseFileById(ctx context.Context, projectID, releaseFileId int64) (releaseFile *rum.ReleaseFile, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := rum.NewDescribeReleaseFilesRequest()
	request.ProjectID = &projectID

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseRumClient().DescribeReleaseFiles(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Files) < 1 {
		return
	}

	for _, v := range response.Response.Files {
		if *v.ID == releaseFileId {
			releaseFile = v
			return
		}
	}

	return
}

func (me *RumService) DeleteRumReleaseFileById(ctx context.Context, releaseFileId int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := rum.NewDeleteReleaseFileRequest()
	request.ID = &releaseFileId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseRumClient().DeleteReleaseFile(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *RumService) DescribeRumLogExportByFilter(ctx context.Context, param map[string]interface{}) (logExport *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeRumLogExportRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Name" {
			request.Name = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "Query" {
			request.Query = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "ProjectId" {
			request.ID = v.(*int64)
		}
		if k == "Fields" {
			request.Fields = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeRumLogExport(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	logExport = response.Response.Result

	return
}

func (me *RumService) DescribeRumLogExportListByFilter(ctx context.Context, param map[string]interface{}) (logExportList *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = rum.NewDescribeRumLogExportsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProjectId" {
			request.ID = v.(*int64)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 20
	)

	request.PageNum = &offset
	request.PageSize = &limit

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseRumClient().DescribeRumLogExports(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	logExportList = response.Response.Result

	return
}
