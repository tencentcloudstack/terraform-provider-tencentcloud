package tencentcloud

import (
	"context"
	"log"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ClsService struct {
	client *connectivity.TencentCloudClient
}

// cls logset
func (me *ClsService) DescribeClsLogsetByFilter(ctx context.Context, filters map[string]string) (instances []*cls.LogsetInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cls.NewDescribeLogsetsRequest()
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = make([]*cls.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cls.Filter{
			Key:    helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset int64 = 0
	var pageSize int64 = 100
	instances = make([]*cls.LogsetInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClsClient().DescribeLogsets(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Logsets) < 1 {
			break
		}
		instances = append(instances, response.Response.Logsets...)
		if len(response.Response.Logsets) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ClsService) DescribeClsLogsetById(ctx context.Context, logSetId string) (logset *cls.LogsetInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cls.NewDescribeLogsetsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	filter := make(map[string]string, 0)
	filter["logsetId"] = logSetId
	request.Filters = append(
		request.Filters,
		&cls.Filter{
			Key:    helper.String("logsetId"),
			Values: []*string{&logSetId},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*cls.LogsetInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClsClient().DescribeLogsets(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Logsets) < 1 {
			break
		}
		instances = append(instances, response.Response.Logsets...)
		if len(response.Response.Logsets) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	logset = instances[0]

	return
}

func (me *ClsService) DeleteClsLogset(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := cls.NewDeleteLogsetRequest()
	request.LogsetId = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DeleteLogset(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// cls topic
func (me *ClsService) DescribeClsTopicByFilter(ctx context.Context, filters map[string]string) (instances []*cls.TopicInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cls.NewDescribeTopicsRequest()
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = make([]*cls.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cls.Filter{
			Key:    helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset int64 = 0
	var pageSize int64 = 100
	instances = make([]*cls.TopicInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClsClient().DescribeTopics(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Topics) < 1 {
			break
		}
		instances = append(instances, response.Response.Topics...)
		if len(response.Response.Topics) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ClsService) DescribeClsTopicById(ctx context.Context, topicId string) (topic *cls.TopicInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cls.NewDescribeTopicsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = []*cls.Filter{
		{
			Key:    common.StringPtr("topicId"),
			Values: []*string{&topicId},
		},
	}
	ratelimit.Check(request.GetAction())
	var (
		offset    int64 = 0
		pageSize  int64 = 100
		instances       = make([]*cls.TopicInfo, 0)
	)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClsClient().DescribeTopics(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Topics) < 1 {
			break
		}
		instances = append(instances, response.Response.Topics...)
		if len(response.Response.Topics) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	topic = instances[0]
	return
}

func (me *ClsService) DeleteClsTopic(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := cls.NewDeleteTopicRequest()
	request.TopicId = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DeleteTopic(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// cls machine group
func (me *ClsService) DescribeClsMachineGroupByFilter(ctx context.Context, filters map[string]string) (instances []*cls.MachineGroupInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cls.NewDescribeMachineGroupsRequest()
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = make([]*cls.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cls.Filter{
			Key:    helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset int64 = 0
	var pageSize int64 = 100
	instances = make([]*cls.MachineGroupInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClsClient().DescribeMachineGroups(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MachineGroups) < 1 {
			break
		}
		instances = append(instances, response.Response.MachineGroups...)
		if len(response.Response.MachineGroups) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ClsService) DescribeClsMachineGroupById(ctx context.Context, id string) (machineGroup *cls.MachineGroupInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cls.NewDescribeMachineGroupsRequest()
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = []*cls.Filter{
		{
			Key:    common.StringPtr("machineGroupId"),
			Values: []*string{&id},
		},
	}

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*cls.MachineGroupInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClsClient().DescribeMachineGroups(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MachineGroups) < 1 {
			break
		}
		instances = append(instances, response.Response.MachineGroups...)
		if len(response.Response.MachineGroups) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	machineGroup = instances[0]
	return
}

func (me *ClsService) DeleteClsMachineGroup(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := cls.NewDeleteMachineGroupRequest()
	request.GroupId = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DeleteMachineGroup(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// cls cos shipper
func (me *ClsService) DescribeClsCosShippersByFilter(ctx context.Context, filters map[string]string) (instances []*cls.ShipperInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cls.NewDescribeShippersRequest()
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = make([]*cls.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cls.Filter{
			Key:    helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances = make([]*cls.ShipperInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClsClient().DescribeShippers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Shippers) < 1 {
			break
		}
		instances = append(instances, response.Response.Shippers...)
		if len(response.Response.Shippers) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ClsService) DescribeClsCosShipperById(ctx context.Context, shipperId string) (instance *cls.ShipperInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cls.NewDescribeShippersRequest()
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = []*cls.Filter{
		{
			Key:    common.StringPtr("shipperId"),
			Values: []*string{&shipperId},
		},
	}

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances := make([]*cls.ShipperInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClsClient().DescribeShippers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Shippers) < 1 {
			break
		}
		instances = append(instances, response.Response.Shippers...)
		if len(response.Response.Shippers) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	instance = instances[0]
	return
}

func (me *ClsService) DeleteClsCosShipper(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := cls.NewDeleteShipperRequest()
	request.ShipperId = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DeleteShipper(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
