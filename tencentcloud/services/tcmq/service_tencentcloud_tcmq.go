package tcmq

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	tcmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// basic information

func NewTcmqService(client *connectivity.TencentCloudClient) TcmqService {
	return TcmqService{client: client}
}

type TcmqService struct {
	client *connectivity.TencentCloudClient
}

func (me *TcmqService) DescribeTcmqTopicById(ctx context.Context, topicName string) (topic *tcmq.CmqTopic, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDescribeCmqTopicDetailRequest()
	request.TopicName = &topicName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqTopicDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.TopicDescribe == nil {
		return
	}

	topic = response.Response.TopicDescribe
	return
}

func (me *TcmqService) DeleteTcmqTopicById(ctx context.Context, topicName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDeleteCmqTopicRequest()
	request.TopicName = &topicName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteCmqTopic(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TcmqService) DescribeTcmqQueueById(ctx context.Context, queueName string) (queue *tcmq.CmqQueue, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDescribeCmqQueueDetailRequest()
	request.QueueName = &queueName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqQueueDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	queue = response.Response.QueueDescribe
	return
}

func (me *TcmqService) DescribeTcmqQueueByFilter(ctx context.Context, paramMap map[string]interface{}) (queueList []*tcmq.CmqQueue, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDescribeCmqQueuesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if v, ok := paramMap["offset"]; ok {
		request.Offset = helper.IntUint64(v.(int))
	}
	if v, ok := paramMap["limit"]; ok {
		request.Limit = helper.IntUint64(v.(int))
	}
	if v, ok := paramMap["queue_name"]; ok {
		request.QueueName = helper.String(v.(string))
	}
	if v, ok := paramMap["is_tag_filter"]; ok {
		request.IsTagFilter = helper.Bool(v.(bool))
	}
	if v, ok := paramMap["queue_name_list"]; ok {
		queueNameList := make([]*string, 0)
		for _, item := range v.([]interface{}) {
			queueName := helper.String(item.(string))
			queueNameList = append(queueNameList, queueName)
		}
		request.QueueNameList = queueNameList
	}
	if v, ok := paramMap["filters"]; ok {
		filters := make([]*tcmq.Filter, 0)
		for _, item := range v.([]interface{}) {
			itemMap := item.(map[string]interface{})
			name := helper.String(itemMap["name"].(string))
			values := make([]*string, 0)
			for _, item := range itemMap["values"].([]string) {
				values = append(values, helper.String(item))
			}
			filters = append(filters, &tcmq.Filter{
				Name:   name,
				Values: values,
			})
		}
		request.Filters = filters
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqQueues(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	queueList = response.Response.QueueList
	return
}

func (me *TcmqService) DeleteTcmqQueueById(ctx context.Context, queueName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDeleteCmqQueueRequest()
	request.QueueName = &queueName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteCmqQueue(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TcmqService) DescribeTcmqTopicByFilter(ctx context.Context, paramMap map[string]interface{}) (topicList []*tcmq.CmqTopic, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDescribeCmqTopicsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if v, ok := paramMap["offset"]; ok {
		request.Offset = helper.IntUint64(v.(int))
	}
	if v, ok := paramMap["limit"]; ok {
		request.Limit = helper.IntUint64(v.(int))
	}
	if v, ok := paramMap["topic_name"]; ok {
		request.TopicName = helper.String(v.(string))
	}
	if v, ok := paramMap["is_tag_filter"]; ok {
		request.IsTagFilter = helper.Bool(v.(bool))
	}
	if v, ok := paramMap["topic_name_list"]; ok {
		topicNameList := make([]*string, 0)
		for _, item := range v.([]interface{}) {
			topicName := helper.String(item.(string))
			topicNameList = append(topicNameList, topicName)
		}
		request.TopicNameList = topicNameList
	}
	if v, ok := paramMap["filters"]; ok {
		filters := make([]*tcmq.Filter, 0)
		for _, item := range v.([]interface{}) {
			itemMap := item.(map[string]interface{})
			name := helper.String(itemMap["name"].(string))
			values := make([]*string, 0)
			for _, item := range itemMap["values"].([]string) {
				values = append(values, helper.String(item))
			}
			filters = append(filters, &tcmq.Filter{
				Name:   name,
				Values: values,
			})
		}
		request.Filters = filters
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqTopics(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	topicList = response.Response.TopicList
	return
}

func (me *TcmqService) DescribeTcmqSubscribeById(ctx context.Context, topicName string, subscriptionName string) (subscribe *tcmq.CmqSubscription, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDescribeCmqSubscriptionDetailRequest()
	request.TopicName = &topicName
	request.SubscriptionName = &subscriptionName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqSubscriptionDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SubscriptionSet) < 1 {
		return
	}

	subscribe = response.Response.SubscriptionSet[0]
	return
}

func (me *TcmqService) DescribeTcmqSubscribeByFilter(ctx context.Context, paramMap map[string]interface{}) (subscriptionList []*tcmq.CmqSubscription, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDescribeCmqSubscriptionDetailRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if v, ok := paramMap["topic_name"]; ok {
		request.TopicName = helper.String(v.(string))
	}
	if v, ok := paramMap["offset"]; ok {
		request.Offset = helper.IntUint64(v.(int))
	}
	if v, ok := paramMap["limit"]; ok {
		request.Limit = helper.IntUint64(v.(int))
	}
	if v, ok := paramMap["subscription_name"]; ok {
		request.SubscriptionName = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqSubscriptionDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	subscriptionList = response.Response.SubscriptionSet
	return
}

func (me *TcmqService) DeleteTcmqSubscribeById(ctx context.Context, topicName string, subscriptionName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcmq.NewDeleteCmqSubscribeRequest()
	request.TopicName = &topicName
	request.SubscriptionName = &subscriptionName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteCmqSubscribe(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
