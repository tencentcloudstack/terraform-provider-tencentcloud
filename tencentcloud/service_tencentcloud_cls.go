package tencentcloud

import (
	"context"
	"log"

	"github.com/pkg/errors"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ClsService struct {
	client *connectivity.TencentCloudClient
}

func (me *ClsService) DescribeTopicsByTopicName(ctx context.Context, topicName string) (clbInstance *cls.TopicInfo, errRet error) {
	logId := getLogId(ctx)
	request := cls.NewDescribeTopicsRequest()
	request.Filters = []*cls.Filter{
		{
			Key:    common.StringPtr("topicName"),
			Values: []*string{&topicName},
		},
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClsClient().DescribeTopics(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Topics) < 1 {
		return
	}
	clbInstance = response.Response.Topics[0]
	return
}

func (me *ClsService) DeleteTopicsByTopicName(ctx context.Context, topicName string) (topicInfo *cls.TopicInfo, errRet error) {
	logId := getLogId(ctx)
	request := cls.NewDescribeTopicsRequest()
	request.Filters = []*cls.Filter{
		{
			Key:    common.StringPtr("topicName"),
			Values: []*string{&topicName},
		},
	}
	ratelimit.Check(request.GetAction())
	client := me.client.UseClsClient()
	response, err := client.DescribeTopics(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Topics) < 1 {
		return
	}
	clbInstance := response.Response.Topics[0]

	delRequest := cls.NewDeleteTopicRequest()
	delRequest.TopicId = clbInstance.TopicId
	_, err = client.DeleteTopic(delRequest)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	topicInfo = clbInstance
	return
}
