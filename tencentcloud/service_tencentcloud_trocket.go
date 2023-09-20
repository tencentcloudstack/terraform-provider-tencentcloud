package tencentcloud

import (
	"context"
	"log"

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
	logId := getLogId(ctx)

	request := trocket.NewDescribeInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTrocketClient().DescribeInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	rocketmqInstance = response.Response
	return
}

func (me *TrocketService) DeleteTrocketRocketmqInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := trocket.NewDeleteInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTrocketClient().DeleteInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TrocketService) TrocketRocketmqInstanceStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

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
	logId := getLogId(ctx)

	request := trocket.NewDescribeTopicListRequest()
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

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	instances := make([]*trocket.TopicItem, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTrocketClient().DescribeTopicList(request)
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
	rocketmqTopic = instances[0]
	return
}

func (me *TrocketService) DeleteTrocketRocketmqTopicById(ctx context.Context, instanceId string, topic string) (errRet error) {
	logId := getLogId(ctx)

	request := trocket.NewDeleteTopicRequest()
	request.InstanceId = &instanceId
	request.Topic = &topic

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTrocketClient().DeleteTopic(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TrocketService) DescribeTrocketRocketmqConsumerGroupById(ctx context.Context, instanceId string, consumerGroup string) (rocketmqConsumerGroup *trocket.DescribeConsumerGroupResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := trocket.NewDescribeConsumerGroupRequest()
	request.InstanceId = &instanceId
	request.ConsumerGroup = &consumerGroup

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTrocketClient().DescribeConsumerGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	rocketmqConsumerGroup = response.Response
	return
}

func (me *TrocketService) DeleteTrocketRocketmqConsumerGroupById(ctx context.Context, instanceId string, consumerGroup string) (errRet error) {
	logId := getLogId(ctx)

	request := trocket.NewDeleteConsumerGroupRequest()
	request.InstanceId = &instanceId
	request.ConsumerGroup = &consumerGroup

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTrocketClient().DeleteConsumerGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
