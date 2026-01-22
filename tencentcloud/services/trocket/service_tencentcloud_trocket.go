package trocket

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TrocketService struct {
	client *connectivity.TencentCloudClient
}

func (me *TrocketService) DescribeTrocketRocketmqInstanceById(ctx context.Context, instanceId string) (rocketmqInstance *trocket.DescribeInstanceResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := trocket.NewDescribeInstanceRequest()
	response := trocket.NewDescribeInstanceResponse()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTrocketClient().DescribeInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe trocket rocketmqInstance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		return nil, err
	}

	rocketmqInstance = response.Response
	return
}

func (me *TrocketService) DeleteTrocketRocketmqInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := trocket.NewDeleteInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTrocketClient().DeleteInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete trocket rocketmqInstance failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return
}

func (me *TrocketService) TrocketRocketmqInstanceStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeTrocketRocketmqInstanceById(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}
		if *object.InstanceStatus == "RUNNING" {
			for _, endpoint := range object.EndpointList {
				if *endpoint.Status != "OPEN" {
					return object, "PROCESSING", nil
				}
			}
		}
		return object, helper.PString(object.InstanceStatus), nil
	}
}

func (me *TrocketService) DescribeTrocketRocketmqTopicById(ctx context.Context, instanceId string, topic string) (rocketmqTopic *trocket.TopicItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := trocket.NewDescribeTopicListRequest()
	response := trocket.NewDescribeTopicListResponse()
	request.InstanceId = &instanceId
	filter := &trocket.Filter{
		Name:   helper.String("TopicName"),
		Values: []*string{&topic},
	}
	request.Filters = []*trocket.Filter{filter}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset    int64 = 0
		limit     int64 = 100
		instances       = make([]*trocket.TopicItem, 0)
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTrocketClient().DescribeTopicList(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe topicList failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.Data) < 1 {
			break
		}

		instances = append(instances, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	for _, item := range instances {
		if *item.Topic == topic {
			rocketmqTopic = item
			break
		}
	}

	return
}

func (me *TrocketService) DeleteTrocketRocketmqTopicById(ctx context.Context, instanceId string, topic string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := trocket.NewDeleteTopicRequest()
	request.InstanceId = &instanceId
	request.Topic = &topic

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTrocketClient().DeleteTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *TrocketService) DescribeTrocketRocketmqConsumerGroupById(ctx context.Context, instanceId string, consumerGroup string) (rocketmqConsumerGroup *trocket.DescribeConsumerGroupResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := trocket.NewDescribeConsumerGroupRequest()
	response := trocket.NewDescribeConsumerGroupResponse()
	request.InstanceId = &instanceId
	request.ConsumerGroup = &consumerGroup

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTrocketClient().DescribeConsumerGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe consumer group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	rocketmqConsumerGroup = response.Response
	return
}

func (me *TrocketService) DeleteTrocketRocketmqConsumerGroupById(ctx context.Context, instanceId string, consumerGroup string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := trocket.NewDeleteConsumerGroupRequest()
	request.InstanceId = &instanceId
	request.ConsumerGroup = &consumerGroup

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTrocketClient().DeleteConsumerGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *TrocketService) DescribeTrocketRocketmqRoleById(ctx context.Context, instanceId string, role string) (rocketmqRole *trocket.RoleItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := trocket.NewDescribeRoleListRequest()
	request.InstanceId = &instanceId
	filter := &trocket.Filter{
		Name:   helper.String("RoleName"),
		Values: []*string{&role},
	}
	request.Filters = []*trocket.Filter{filter}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	instances := make([]*trocket.RoleItem, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTrocketClient().DescribeRoleList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}
		instances = append(instances, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	rocketmqRole = instances[0]
	return
}

func (me *TrocketService) DeleteTrocketRocketmqRoleById(ctx context.Context, instanceId string, role string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := trocket.NewDeleteRoleRequest()
	request.InstanceId = &instanceId
	request.Role = &role

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTrocketClient().DeleteRole(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TrocketService) DescribeTrocketRocketmqInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret []*trocket.InstanceItem, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = trocket.NewDescribeInstanceListRequest()
		response = trocket.NewDescribeInstanceListResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*trocket.Filter)
		}

		if k == "TagFilters" {
			request.TagFilters = v.([]*trocket.TagFilter)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseTrocketClient().DescribeInstanceList(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe instance list failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.Data) < 1 {
			break
		}

		ret = append(ret, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
