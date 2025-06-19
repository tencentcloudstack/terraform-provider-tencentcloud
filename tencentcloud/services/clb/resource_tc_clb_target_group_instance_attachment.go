package clb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func ResourceTencentCloudClbTGAttachmentInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTGAttachmentInstanceCreate,
		Read:   resourceTencentCloudClbTGAttachmentInstanceRead,
		Update: resourceTencentCloudClbTGAttachmentInstanceUpdate,
		Delete: resourceTencentCloudClbTGAttachmentInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "Target group ID.",
			},
			"bind_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "The Intranet IP of the target group instance.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The port of the target group instance, fully listening to the target group does not support passing this field.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The weight of the target group instance. Value range: 0-100.",
			},
		},
	}
}

func resourceTencentCloudClbTGAttachmentInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_instance_attachment.create")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = clb.NewRegisterTargetGroupInstancesRequest()
		targetGroupId string
		bindIp        string
		port          int
	)

	if v, ok := d.GetOk("target_group_id"); ok {
		request.TargetGroupId = helper.String(v.(string))
		targetGroupId = v.(string)
	}

	targetGroupInstances := clb.TargetGroupInstance{}
	if v, ok := d.GetOk("bind_ip"); ok {
		targetGroupInstances.BindIP = helper.String(v.(string))
		bindIp = v.(string)
	}

	if v, ok := d.GetOkExists("port"); ok {
		targetGroupInstances.Port = helper.IntUint64(v.(int))
		port = v.(int)
	}

	if v, ok := d.GetOkExists("weight"); ok {
		targetGroupInstances.Weight = helper.IntUint64(v.(int))
	}

	request.TargetGroupInstances = append(request.TargetGroupInstances, &targetGroupInstances)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().RegisterTargetGroupInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RequestId == nil {
			return resource.NonRetryableError(fmt.Errorf("Register target group instance failed, Response is nil."))
		}

		requestId := *result.Response.RequestId
		retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
		if retryErr != nil {
			return resource.NonRetryableError(errors.WithStack(retryErr))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s register target group instance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{targetGroupId, bindIp, strconv.Itoa(port)}, tccommon.FILED_SP))
	return resourceTencentCloudClbTGAttachmentInstanceRead(d, meta)
}

func resourceTencentCloudClbTGAttachmentInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_instance_attachment.read")()

	var (
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService           = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                   = d.Id()
		targetGroupInstances []*clb.TargetGroupBackend
	)

	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("target group instance attachment id is not set")
	}

	targetGroupId := idSplit[0]
	bindIp := idSplit[1]
	port, err := strconv.ParseUint(idSplit[2], 0, 64)
	if err != nil {
		return err
	}

	filters := make(map[string]string)
	filters["TargetGroupId"] = targetGroupId
	filters["BindIP"] = bindIp
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		targetGroupInstances, err = clbService.DescribeTargetGroupInstances(ctx, filters)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, tgInstance := range targetGroupInstances {
		if *tgInstance.Port == port {
			_ = d.Set("target_group_id", idSplit[0])
			_ = d.Set("bind_ip", idSplit[1])
			_ = d.Set("port", helper.StrToInt64(idSplit[2]))
			if tgInstance.Weight != nil {
				_ = d.Set("weight", *tgInstance.Weight)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceTencentCloudClbTGAttachmentInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_instance_attachment.update")()

	var (
		logId                 = tccommon.GetLogId(tccommon.ContextNil)
		ctx                   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService            = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                    = d.Id()
		port                  int
		bindIp, targetGroupId string
		err                   error
	)

	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("target group instance attachment id is not set")
	}

	targetGroupId = idSplit[0]
	bindIp = idSplit[1]
	port, err = strconv.Atoi(idSplit[2])
	if err != nil {
		return err
	}

	if d.HasChange("weight") {
		newWeight := d.Get("weight").(int)
		err := clbService.ModifyTargetGroupInstancesWeight(ctx, targetGroupId, bindIp, uint64(port), uint64(newWeight))
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudClbTGAttachmentInstanceRead(d, meta)
}

func resourceTencentCloudClbTGAttachmentInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_instance_attachment.delete")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id         = d.Id()
	)

	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("target group instance attachment id is not set")
	}

	targetGroupId := idSplit[0]
	bindIp := idSplit[1]
	port, err := strconv.ParseUint(idSplit[2], 0, 64)
	if err != nil {
		return err
	}

	return clbService.DeregisterTargetInstances(ctx, targetGroupId, bindIp, port)
}
