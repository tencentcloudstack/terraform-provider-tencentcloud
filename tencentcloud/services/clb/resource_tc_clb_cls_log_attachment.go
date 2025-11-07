package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clbv20180317 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbClsLogAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbClsLogAttachmentCreate,
		Read:   resourceTencentCloudClbClsLogAttachmentRead,
		Delete: resourceTencentCloudClbClsLogAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CLB instance ID.",
			},

			"log_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Logset ID of the Cloud Log Service (CLS).<li>When adding or updating a log topic, call the [DescribeLogsets](https://intl.cloud.tencent.com/document/product/614/58624?from_cn_redirect=1) API to obtain the logset ID.</li><li>When deleting a log topic, set this parameter to null.</li>.",
			},

			"log_topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log topic ID of the CLS.<li>When adding or updating a log topic, call the [DescribeTopics](https://intl.cloud.tencent.com/document/product/614/56454?from_cn_redirect=1) API to obtain the log topic ID.</li><li>When deleting a log topic, set this parameter to null.</li>.",
			},

			// "log_type": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	ForceNew:    true,
			// 	Computed:    true,
			// 	Description: "Log type:\n<li>`ACCESS`: access logs</li>\n<li>`HEALTH`: health check logs</li>\nDefault: `ACCESS`.",
			// },
		},
	}
}

func resourceTencentCloudClbClsLogAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_cls_log_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request        = clbv20180317.NewSetLoadBalancerClsLogRequest()
		loadBalancerId string
		logSetId       string
		logTopicId     string
	)

	if v, ok := d.GetOk("load_balancer_id"); ok {
		request.LoadBalancerId = helper.String(v.(string))
		loadBalancerId = v.(string)
	}

	if v, ok := d.GetOk("log_set_id"); ok {
		request.LogSetId = helper.String(v.(string))
		logSetId = v.(string)
	}

	if v, ok := d.GetOk("log_topic_id"); ok {
		request.LogTopicId = helper.String(v.(string))
		logTopicId = v.(string)
	}

	// if v, ok := d.GetOk("log_type"); ok {
	// 	request.LogType = helper.String(v.(string))
	// }

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerClsLogWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create clb cls log attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{loadBalancerId, logSetId, logTopicId}, tccommon.FILED_SP))
	return resourceTencentCloudClbClsLogAttachmentRead(d, meta)
}

func resourceTencentCloudClbClsLogAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_cls_log_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	loadBalancerId := idSplit[0]
	logSetId := idSplit[1]
	logTopicId := idSplit[2]

	respData, err := service.DescribeClbClsLogAttachmentById(ctx, loadBalancerId, logSetId, logTopicId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_clb_cls_log_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.LoadBalancerId != nil {
		_ = d.Set("load_balancer_id", respData.LoadBalancerId)
	}

	if respData.LogSetId != nil && respData.LogTopicId != nil {
		if *respData.LogSetId == logSetId && *respData.LogTopicId == logTopicId {
			_ = d.Set("log_set_id", respData.LogSetId)
			_ = d.Set("log_topic_id", respData.LogTopicId)
			// _ = d.Set("log_type", "ACCESS")
		}
	}

	// if respData.HealthLogSetId != nil && respData.HealthLogTopicId != nil {
	// 	if *respData.HealthLogSetId == logSetId && *respData.HealthLogTopicId == logTopicId {
	// 		_ = d.Set("log_set_id", respData.HealthLogSetId)
	// 		_ = d.Set("log_topic_id", respData.HealthLogTopicId)
	// 		_ = d.Set("log_type", "HEALTH")
	// 	}
	// }

	return nil
}

func resourceTencentCloudClbClsLogAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_cls_log_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = clbv20180317.NewSetLoadBalancerClsLogRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	loadBalancerId := idSplit[0]

	request.LoadBalancerId = &loadBalancerId
	request.LogSetId = helper.String("")
	request.LogTopicId = helper.String("")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerClsLogWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete clb cls log attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
