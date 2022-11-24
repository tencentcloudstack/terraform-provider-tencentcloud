package tencentcloud

import (
	"context"
	"log"
	"strconv"

	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type RumService struct {
	client *connectivity.TencentCloudClient
}

func (me *RumService) DescribeRumTawInstance(ctx context.Context, instanceId string) (tawInstance *rum.RumInstanceInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
