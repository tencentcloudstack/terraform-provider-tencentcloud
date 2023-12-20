package clb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewClsService(client *connectivity.TencentCloudClient) ClsService {
	return ClsService{client: client}
}

type ClsService struct {
	client *connectivity.TencentCloudClient
}

// cls logset
func (me *ClsService) DescribeClsLogset(ctx context.Context, logsetId string) (logset *cls.LogsetInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cls.NewDescribeLogsetsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&cls.Filter{
			Key:    helper.String("logsetId"),
			Values: []*string{&logsetId},
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

func (me *ClsService) DeleteClsLogsetById(ctx context.Context, logsetId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteLogsetRequest()
	request.LogsetId = &logsetId

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

func (me *ClsService) DescribeClsLogsetByFilter(ctx context.Context, filters map[string]string) (instances []*cls.LogsetInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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

// cls topic
func (me *ClsService) DescribeClsTopicByFilter(ctx context.Context, filters map[string]string) (instances []*cls.TopicInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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

func (me *ClsService) DescribeClsMachineGroupByConfigId(ctx context.Context, configId, groupId string) (machineGroup *cls.MachineGroupInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cls.NewDescribeConfigMachineGroupsRequest()
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ConfigId = &configId

	response, err := me.client.UseClsClient().DescribeConfigMachineGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instances := response.Response.MachineGroups

	for _, instance := range instances {
		if *instance.GroupId == groupId {
			machineGroup = instance
			break
		}
	}

	return
}

// cls cos shipper
func (me *ClsService) DescribeClsCosShippersByFilter(ctx context.Context, filters map[string]string) (instances []*cls.ShipperInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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

// cls config
func (me *ClsService) DescribeClsConfigById(ctx context.Context, configId string) (config *cls.ConfigInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeConfigsRequest()

	filter := &cls.Filter{
		Key:    helper.String("configId"),
		Values: []*string{&configId},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*cls.ConfigInfo, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseClsClient().DescribeConfigs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Configs) < 1 {
			break
		}
		instances = append(instances, response.Response.Configs...)
		if len(response.Response.Configs) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	config = instances[0]
	return
}

func (me *ClsService) DeleteClsConfig(ctx context.Context, id string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteConfigRequest()
	request.ConfigId = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DeleteConfig(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// cls config extra
func (me *ClsService) DescribeClsConfigExtraById(ctx context.Context, configExtraId string) (config *cls.ConfigExtraInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeConfigExtrasRequest()

	filter := &cls.Filter{
		Key:    helper.String("configExtraId"),
		Values: []*string{&configExtraId},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*cls.ConfigExtraInfo, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseClsClient().DescribeConfigExtras(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Configs) < 1 {
			break
		}
		instances = append(instances, response.Response.Configs...)
		if len(response.Response.Configs) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	config = instances[0]
	return
}

func (me *ClsService) DeleteClsConfigExtra(ctx context.Context, id string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteConfigExtraRequest()
	request.ConfigExtraId = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DeleteConfigExtra(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// cls index
func (me *ClsService) DeleteClsIndex(ctx context.Context, id string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteIndexRequest()
	request.TopicId = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DeleteIndex(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ClsService) DescribeClsAlarmById(ctx context.Context, alarmId string) (alarm *cls.AlarmInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeAlarmsRequest()
	filter := &cls.Filter{
		Key:    helper.String("alarmId"),
		Values: []*string{&alarmId},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*cls.AlarmInfo, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseClsClient().DescribeAlarms(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Alarms) < 1 {
			break
		}
		instances = append(instances, response.Response.Alarms...)
		if len(response.Response.Alarms) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	alarm = instances[0]
	return
}

func (me *ClsService) DeleteClsAlarmById(ctx context.Context, alarmId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteAlarmRequest()
	request.AlarmId = &alarmId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DeleteAlarm(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ClsService) DescribeClsAlarmNoticeById(ctx context.Context, alarmNoticeId string) (alarmNotice *cls.AlarmNotice, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeAlarmNoticesRequest()
	filter := &cls.Filter{
		Key:    helper.String("alarmNoticeId"),
		Values: []*string{&alarmNoticeId},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*cls.AlarmNotice, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseClsClient().DescribeAlarmNotices(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AlarmNotices) < 1 {
			break
		}
		instances = append(instances, response.Response.AlarmNotices...)
		if len(response.Response.AlarmNotices) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	alarmNotice = instances[0]
	return
}

func (me *ClsService) DeleteClsAlarmNoticeById(ctx context.Context, alarmNoticeId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteAlarmNoticeRequest()
	request.AlarmNoticeId = &alarmNoticeId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DeleteAlarmNotice(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ClsService) DescribeClsCkafkaConsumerById(ctx context.Context, topicId string) (ckafkaConsumer *cls.DescribeConsumerResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeConsumerRequest()
	request.TopicId = &topicId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeConsumer(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ckafkaConsumer = response.Response
	return
}

func (me *ClsService) DeleteClsCkafkaConsumerById(ctx context.Context, topicId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteConsumerRequest()
	request.TopicId = &topicId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DeleteConsumer(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ClsService) DescribeClsCosRechargeById(ctx context.Context, topicId string, rechargeId string) (cosRecharge *cls.CosRechargeInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeCosRechargesRequest()
	request.TopicId = &topicId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeCosRecharges(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, info := range response.Response.Data {
		if *info.Id == rechargeId {
			cosRecharge = info
			break
		}
	}
	return
}

func (me *ClsService) DescribeClsExportById(ctx context.Context, topicId string, exportId string) (export *cls.ExportInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeExportsRequest()
	request.TopicId = &topicId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeExports(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, info := range response.Response.Exports {
		if *info.ExportId == exportId {
			export = info
			break
		}
	}

	return
}

func (me *ClsService) DeleteClsExportById(ctx context.Context, exportId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteExportRequest()
	request.ExportId = &exportId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DeleteExport(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ClsService) DescribeClsShipperTasksByFilter(ctx context.Context, param map[string]interface{}) (shipperTasks []*cls.ShipperTaskInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cls.NewDescribeShipperTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ShipperId" {
			request.ShipperId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeShipperTasks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	shipperTasks = response.Response.Tasks

	return
}

func (me *ClsService) DescribeClsMachinesByFilter(ctx context.Context, param map[string]interface{}) (machines []*cls.MachineInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cls.NewDescribeMachinesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GroupId" {
			request.GroupId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeMachines(request)
	if err != nil {
		errRet = err
		return
	}

	machines = response.Response.Machines

	return
}

func (me *ClsService) DescribeClsMachineGroupConfigsByFilter(ctx context.Context, param map[string]interface{}) (machineGroupConfigs []*cls.ConfigInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cls.NewDescribeMachineGroupConfigsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GroupId" {
			request.GroupId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeMachineGroupConfigs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	machineGroupConfigs = response.Response.Configs
	return
}

func (me *ClsService) DescribeClsDataTransformById(ctx context.Context, taskId string) (dataTransform *cls.DataTransformTaskInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeDataTransformInfoRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeDataTransformInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DataTransformTaskInfos) < 1 {
		return
	}

	dataTransform = response.Response.DataTransformTaskInfos[0]
	return
}

func (me *ClsService) DeleteClsDataTransformById(ctx context.Context, taskId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteDataTransformRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DeleteDataTransform(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *ClsService) DescribeClsKafkaRechargeById(ctx context.Context, id string, topic string) (kafkaRecharge *cls.KafkaRechargeInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeKafkaRechargesRequest()
	request.Id = &id
	request.TopicId = &topic
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeKafkaRecharges(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Infos) < 1 {
		return
	}

	kafkaRecharge = response.Response.Infos[0]
	return
}

func (me *ClsService) DeleteClsKafkaRechargeById(ctx context.Context, id string, topic string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteKafkaRechargeRequest()
	request.Id = &id
	request.TopicId = &topic
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DeleteKafkaRecharge(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ClsService) DescribeClsScheduledSqlById(ctx context.Context, taskId string) (scheduledSql *cls.ScheduledSqlTaskInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDescribeScheduledSqlInfoRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DescribeScheduledSqlInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ScheduledSqlTaskInfos) < 1 {
		return
	}

	scheduledSql = response.Response.ScheduledSqlTaskInfos[0]
	return
}

func (me *ClsService) DeleteClsScheduledSqlById(ctx context.Context, taskId string, srcTopicId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cls.NewDeleteScheduledSqlRequest()
	request.TaskId = &taskId
	request.SrcTopicId = &srcTopicId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseClsClient().DeleteScheduledSql(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
