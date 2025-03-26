package mqtt

import (
	"context"
	"log"

	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewMqttService(client *connectivity.TencentCloudClient) MqttService {
	return MqttService{client: client}
}

type MqttService struct {
	client *connectivity.TencentCloudClient
}

func (me *MqttService) DescribeMqttById(ctx context.Context, instanceId string) (ret *mqttv20240516.DescribeInstanceResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeInstanceRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMqttV20240516Client().DescribeInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttInstancePublicEndpointById(ctx context.Context, instanceId string) (ret *mqttv20240516.DescribeInsPublicEndpointsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeInsPublicEndpointsRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMqttV20240516Client().DescribeInsPublicEndpoints(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttTopicById(ctx context.Context, instanceId string, topic string) (ret *mqttv20240516.DescribeTopicResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeTopicRequest()
	request.InstanceId = helper.String(instanceId)
	request.Topic = helper.String(topic)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMqttV20240516Client().DescribeTopic(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}
