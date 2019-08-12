/*
Provides a CCN attaching resource.

Example Usage

```hcl
variable "region" {
    default = "ap-guangzhou"
}

resource  "tencentcloud_vpc"   "vpc"  {
    name = "ci-temp-test-vpc"
    cidr_block = "10.0.0.0/16"
    dns_servers=["119.29.29.29","8.8.8.8"]
    is_multicast=false
}

resource "tencentcloud_ccn" "main"{
	name ="ci-temp-test-ccn"
	description="ci-temp-test-ccn-des"
	qos ="AG"
}

resource "tencentcloud_ccn_attachment" "attachment"{
	ccn_id = "${tencentcloud_ccn.main.id}"
	instance_type ="VPC"
	instance_id ="${tencentcloud_vpc.vpc.id}"
	instance_region="${var.region}"
}
```
*/
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

func resourceTencentCloudCcnAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnAttachmentCreate,
		Read:   resourceTencentCloudCcnAttachmentRead,
		Delete: resourceTencentCloudCcnAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CCN",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{CNN_INSTANCE_TYPE_VPC, CNN_INSTANCE_TYPE_DIRECTCONNECT, CNN_INSTANCE_TYPE_BMVPC}),
				ForceNew:     true,
				Description:  "Type of attached instance network, and available values include VPC, DIRECTCONNECT and BMVPC.",
			},
			"instance_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region that the instance locates at.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of instance is attached.",
			},

			// Computed values
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "States of instance is attached, and available values include PENDING, ACTIVE, EXPIRED, REJECTED, DELETED, FAILED(asynchronous forced disassociation after 2 hours), ATTACHING, DETACHING and DETACHFAILED(asynchronous forced disassociation after 2 hours).",
			},
			"attached_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time of attaching.",
			},
			"cidr_block": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A network address block of the instance that is attached.",
			},
		},
	}
}

func resourceTencentCloudCcnAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_attachment.create")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)

	if len(ccnId) < 4 || len(instanceRegion) < 3 || len(instanceId) < 3 {
		return fmt.Errorf("param ccn_id or instance_region or instance_id  error")
	}

	_, has, err := service.DescribeCcn(ctx, ccnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return fmt.Errorf("ccn[%s] doesn't exist", ccnId)
	}

	if err := service.AttachCcnInstances(ctx, ccnId, instanceRegion, instanceType, instanceId); err != nil {
		return err
	}

	m := md5.New()
	m.Write([]byte(ccnId + instanceType + instanceRegion + instanceId))
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	return resourceTencentCloudCcnAttachmentRead(d, meta)
}

func resourceTencentCloudCcnAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_attachment.read")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)

	_, has, err := service.DescribeCcn(ctx, ccnId)
	if err != nil {
		return err
	}

	if has == 0 {
		d.SetId("")
		return nil
	}

	info, has, err := service.DescribeCcnAttachedInstance(ctx, ccnId, instanceRegion, instanceType, instanceId)
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

func resourceTencentCloudCcnAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_attachment.delete")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)

	_, has, err := service.DescribeCcn(ctx, ccnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}

	if err := service.DetachCcnInstances(ctx, ccnId, instanceRegion, instanceType, instanceId); err != nil {
		return err
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, has, err := service.DescribeCcnAttachedInstance(ctx, ccnId, instanceRegion, instanceType, instanceId)
		if err != nil {
			return resource.RetryableError(err)
		}
		if has == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}
