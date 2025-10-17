package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsCloudProductLogTaskReadPreRequest0(ctx context.Context, req *clsv20201016.DescribeCloudProductLogTasksRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	assumerName := idSplit[1]
	logType := idSplit[2]

	req.Filters = []*clsv20201016.Filter{
		{
			Key:    common.StringPtr("instanceId"),
			Values: common.StringPtrs([]string{instanceId}),
		},
		{
			Key:    common.StringPtr("assumerName"),
			Values: common.StringPtrs([]string{assumerName}),
		},
		{
			Key:    common.StringPtr("logType"),
			Values: common.StringPtrs([]string{logType}),
		},
	}

	return nil
}

func resourceTencentCloudClsCloudProductLogTaskReadPreHandleResponse0(ctx context.Context, resp *clsv20201016.DescribeCloudProductLogTasksResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	if resp.Tasks == nil || len(resp.Tasks) < 1 {
		return fmt.Errorf("resource `cls_cloud_product_log_task` not found")
	}

	var (
		logId    = tccommon.GetLogId(ctx)
		meta     = tccommon.ProviderMetaFromContext(ctx)
		logsetId string
		topicId  string
	)

	for _, item := range resp.Tasks {
		if item.ClsRegion != nil {
			_ = d.Set("cls_region", item.ClsRegion)
		}

		if item.LogsetId != nil {
			_ = d.Set("logset_id", item.LogsetId)
			logsetId = *item.LogsetId
		}

		if item.TopicId != nil {
			_ = d.Set("topic_id", item.TopicId)
			topicId = *item.TopicId
		}

		if item.Extend != nil {
			_ = d.Set("extend", item.Extend)
		}
	}

	// get logset name
	logsetRequest := clsv20201016.NewDescribeLogsetsRequest()
	logsetResponse := clsv20201016.NewDescribeLogsetsResponse()
	logsetRequest.Filters = []*clsv20201016.Filter{
		{
			Key:    common.StringPtr("logsetId"),
			Values: common.StringPtrs([]string{logsetId}),
		},
	}
	logsetRequest.Offset = common.Int64Ptr(0)
	logsetRequest.Limit = common.Int64Ptr(20)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DescribeLogsetsWithContext(ctx, logsetRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, logsetRequest.GetAction(), logsetRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("get logset name error")
			return resource.NonRetryableError(e)
		}

		logsetResponse = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s get logset name failed, reason:%+v", logId, err)
		return err
	}

	if logsetResponse.Response != nil && logsetResponse.Response.Logsets == nil && len(logsetResponse.Response.Logsets) != 1 {
		return fmt.Errorf("get logset name failed")
	}

	if logsetResponse.Response.Logsets[0].LogsetName != nil {
		_ = d.Set("logset_name", logsetResponse.Response.Logsets[0].LogsetName)
	}

	// get topic name
	topicRequest := clsv20201016.NewDescribeTopicsRequest()
	topicResponse := clsv20201016.NewDescribeTopicsResponse()
	topicRequest.Filters = []*clsv20201016.Filter{
		{
			Key:    common.StringPtr("topicId"),
			Values: common.StringPtrs([]string{topicId}),
		},
	}
	topicRequest.Offset = common.Int64Ptr(0)
	topicRequest.Limit = common.Int64Ptr(20)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DescribeTopicsWithContext(ctx, topicRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, topicRequest.GetAction(), topicRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("get topic name error")
			return resource.NonRetryableError(e)
		}

		topicResponse = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s get topic name failed, reason:%+v", logId, err)
		return err
	}

	if topicResponse.Response != nil && topicResponse.Response.Topics == nil && len(topicResponse.Response.Topics) != 1 {
		return fmt.Errorf("get topic name failed")
	}

	if topicResponse.Response.Topics[0].TopicName != nil {
		_ = d.Set("topic_name", topicResponse.Response.Topics[0].TopicName)
	}

	return nil
}

func resourceTencentCloudClsCloudProductLogTaskDeletePostFillRequest1(ctx context.Context, req *clsv20201016.DeleteTopicRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	if v, ok := d.GetOk("topic_id"); ok {
		req.TopicId = helper.String(v.(string))
	}

	return nil
}

func resourceTencentCloudClsCloudProductLogTaskDeletePostFillRequest2(ctx context.Context, req *clsv20201016.DeleteLogsetRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	if v, ok := d.GetOk("logset_id"); ok {
		req.LogsetId = helper.String(v.(string))
	}

	return nil
}
