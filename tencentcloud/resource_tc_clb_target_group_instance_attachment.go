package tencentcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbTGAttachmentInstance() *schema.Resource {
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
				ValidateFunc: validateNotEmpty,
				ForceNew:     true,
				Description:  "Target group ID.",
			},
			"bind_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				ForceNew:     true,
				Description:  "The Intranet IP of the target group instance.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Port of the target group instance.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The weight of the target group instance.",
			},
		},
	}
}

func resourceTencentCloudClbTGAttachmentInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_instance_attachment.create")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		clbService    = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		targetGroupId = d.Get("target_group_id").(string)
		bindIp        = d.Get("bind_ip").(string)
		port          = d.Get("port").(int)
		weight        = d.Get("weight").(int)
		err           error
	)

	err = clbService.RegisterTargetInstances(ctx, targetGroupId, bindIp, uint64(port), uint64(weight))

	if err != nil {
		return err
	}
	time.Sleep(time.Duration(3) * time.Second)

	d.SetId(strings.Join([]string{targetGroupId, bindIp, strconv.Itoa(port)}, FILED_SP))

	return resourceTencentCloudClbTGAttachmentInstanceRead(d, meta)
}

func resourceTencentCloudClbTGAttachmentInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_instance_attachment.read")()

	var (
		logId                = getLogId(contextNil)
		ctx                  = context.WithValue(context.TODO(), logIdKey, logId)
		clbService           = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                   = d.Id()
		targetGroupInstances []*clb.TargetGroupBackend
	)
	idSplit := strings.Split(id, FILED_SP)
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
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		targetGroupInstances, err = clbService.DescribeTargetGroupInstances(ctx, filters)
		if err != nil {
			return retryError(err, InternalError)
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
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceTencentCloudClbTGAttachmentInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_instance_attachment.update")()

	var (
		logId                 = getLogId(contextNil)
		ctx                   = context.WithValue(context.TODO(), logIdKey, logId)
		clbService            = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                    = d.Id()
		port                  int
		bindIp, targetGroupId string
		err                   error
	)
	idSplit := strings.Split(id, FILED_SP)
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
	defer logElapsed("resource.tencentcloud_clb_target_group_instance_attachment.delete")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		clbService = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		id         = d.Id()
	)
	idSplit := strings.Split(id, FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("target group instance attachment id is not set")
	}
	targetGroupId := idSplit[0]
	bindIp := idSplit[1]
	port, err := strconv.ParseUint(idSplit[2], 0, 64)
	if err != nil {
		return err
	}

	err = clbService.DeregisterTargetInstances(ctx, targetGroupId, bindIp, port)

	if err != nil {
		return err
	}
	return nil
}
