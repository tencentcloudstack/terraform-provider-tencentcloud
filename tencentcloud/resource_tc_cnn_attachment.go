package tencentcloud

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudCnnAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCnnAttachmentCreate,
		Read:   resourceTencentCloudCnnAttachmentRead,
		Delete: resourceTencentCloudCnnAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"cnn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{CNN_INSTANCE_TYPE_VPC, CNN_INSTANCE_TYPE_DIRECTCONNECT, CNN_INSTANCE_TYPE_BMVPC}),
				ForceNew:     true,
			},
			"instance_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Computed values
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attached_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
func resourceTencentCloudCnnAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_cnn_attachment.create")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		cnnId          = d.Get("cnn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)

	if len(cnnId) < 4 || len(instanceRegion) < 3 || len(instanceId) < 3 {
		return fmt.Errorf("param cnn_id or instance_region or instance_id  error")
	}

	_, has, err := service.DescribeCcn(ctx, cnnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return fmt.Errorf("cnn[%s] doesn't exist", cnnId)
	}

	if err := service.AttachCcnInstances(ctx, cnnId, instanceRegion, instanceType, instanceId); err != nil {
		return err
	}

	m := md5.New()
	m.Write([]byte(cnnId + instanceType + instanceRegion + instanceId))
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	return resourceTencentCloudCnnAttachmentRead(d, meta)
}

func resourceTencentCloudCnnAttachmentRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_cnn_attachment.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		cnnId          = d.Get("cnn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)

	_, has, err := service.DescribeCcn(ctx, cnnId)
	if err != nil {
		return err
	}

	if has == 0 {
		d.SetId("")
		return nil
	}

	info, has, err := service.DescribeCcnAttachedInstance(ctx, cnnId, instanceRegion, instanceType, instanceId)
	if err != nil {
		return err
	}
	if has == 0 {
		d.SetId("")
		return nil
	}

	d.Set("state", strings.ToUpper(info.state))
	d.Set("attached_time", info.attachedTime)
	d.Set("cidr_block", info.cidrBlock)

	return nil
}

func resourceTencentCloudCnnAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_cnn_attachment.delete")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var (
		cnnId          = d.Get("cnn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)

	_, has, err := service.DescribeCcn(ctx, cnnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}

	if err := service.DetachCcnInstances(ctx, cnnId, instanceRegion, instanceType, instanceId); err != nil {
		return err
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, has, err := service.DescribeCcnAttachedInstance(ctx, cnnId, instanceRegion, instanceType, instanceId)
		if err != nil {
			return resource.RetryableError(err)
		}
		if has == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete fail"))
	})

	return nil
}
