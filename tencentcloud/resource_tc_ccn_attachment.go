/*
Provides a CCN attaching resource.

Example Usage

```hcl
variable "region" {
  default = "ap-guangzhou"
}

variable "otheruin" {
  default = "123353"
}

variable "otherccn" {
  default = "ccn-151ssaga"
}

resource "tencentcloud_vpc" "vpc" {
  name         = "ci-temp-test-vpc"
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = false
}

resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.main.id
  instance_type   = "VPC"
  instance_id     = tencentcloud_vpc.vpc.id
  instance_region = var.region
}

resource "tencentcloud_ccn_attachment" "other_account" {
  ccn_id          = var.otherccn
  instance_type   = "VPC"
  instance_id     = tencentcloud_vpc.vpc.id
  instance_region = var.region
  ccn_uin         = var.otheruin
}
```
*/
package tencentcloud

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudCcnAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnAttachmentCreate,
		Read:   resourceTencentCloudCcnAttachmentRead,
		Update: resourceTencentCloudCcnAttachmentUpdate,
		Delete: resourceTencentCloudCcnAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CCN.",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{CNN_INSTANCE_TYPE_VPC, CNN_INSTANCE_TYPE_DIRECTCONNECT, CNN_INSTANCE_TYPE_BMVPC, CNN_INSTANCE_TYPE_VPNGW}),
				ForceNew:     true,
				Description:  "Type of attached instance network, and available values include `VPC`, `DIRECTCONNECT`, `BMVPC` and `VPNGW`. Note: `VPNGW` type is only for whitelist customer now.",
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark of attachment.",
			},
			"ccn_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Uin of the ccn attached. Default is ``, which means the uin of this account. This parameter is used with case when attaching ccn of other account to the instance of this account. For now only support instance type `VPC`.",
			},
			// Computed values
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "States of instance is attached. Valid values: `PENDING`, `ACTIVE`, `EXPIRED`, `REJECTED`, `DELETED`, `FAILED`, `ATTACHING`, `DETACHING` and `DETACHFAILED`. `FAILED` means asynchronous forced disassociation after 2 hours. `DETACHFAILED` means asynchronous forced disassociation after 2 hours.",
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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
		description    = ""
		ccnUin         = ""
	)

	if len(ccnId) < 4 || len(instanceRegion) < 3 || len(instanceId) < 3 {
		return fmt.Errorf("param ccn_id or instance_region or instance_id  error")
	}

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}

	if v, ok := d.GetOk("ccn_uin"); ok {
		ccnUin = v.(string)
		if ccnUin != "" && instanceType != CNN_INSTANCE_TYPE_VPC {
			return fmt.Errorf("Other ccn account attachment %s only support instance type of `VPC`.", ccnId)
		}
	} else {
		_, has, err := service.DescribeCcn(ctx, ccnId)
		if err != nil {
			return err
		}
		if has == 0 {
			return fmt.Errorf("ccn[%s] doesn't exist", ccnId)
		}
	}

	if err := service.AttachCcnInstances(ctx, ccnId, instanceRegion, instanceType, instanceId, ccnUin, description); err != nil {
		return err
	}

	m := md5.New()
	_, err := m.Write([]byte(ccnId + instanceType + instanceRegion + instanceId))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	return resourceTencentCloudCcnAttachmentRead(d, meta)
}

func resourceTencentCloudCcnAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	if v, ok := d.GetOk("ccn_uin"); ok {
		ccnUin := v.(string)
		ccnId := d.Get("ccn_id").(string)
		instanceType := d.Get("instance_type").(string)
		instanceRegion := d.Get("instance_region").(string)
		instanceId := d.Get("instance_id").(string)

		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			infos, e := service.DescribeCcnAttachmentsByInstance(ctx, instanceType, instanceId, instanceRegion)
			if e != nil {
				return retryError(e)
			}

			if len(infos) == 0 {
				d.SetId("")
				return nil
			}
			findFlag := false
			for _, info := range infos {
				if *info.CcnUin == ccnUin && *info.CcnId == ccnId {
					_ = d.Set("state", strings.ToUpper(*info.State))
					_ = d.Set("attached_time", info.AttachedTime)
					_ = d.Set("cidr_block", info.CidrBlock)
					findFlag = true
					break
				}
			}
			if !findFlag {
				d.SetId("")
				return nil
			}
			return nil
		})

		if err != nil {
			return err
		}
		return nil
	}

	var (
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
		onlineHas      = true
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, e := service.DescribeCcn(ctx, ccnId)
		if e != nil {
			return retryError(e)
		}

		if has == 0 {
			d.SetId("")
			onlineHas = false
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	if !onlineHas {
		return nil
	}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeCcnAttachedInstance(ctx, ccnId, instanceRegion, instanceType, instanceId)
		if e != nil {
			return retryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("description", info.description)
		_ = d.Set("state", strings.ToUpper(info.state))
		_ = d.Set("attached_time", info.attachedTime)
		_ = d.Set("cidr_block", info.cidrBlock)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudCcnAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_attachment.create")()
	logId := getLogId(contextNil)

	if d.HasChange("description") {
		var (
			ccnId          = d.Get("ccn_id").(string)
			instanceType   = d.Get("instance_type").(string)
			instanceRegion = d.Get("instance_region").(string)
			instanceId     = d.Get("instance_id").(string)
			description    = d.Get("description").(string)
		)

		request := vpc.NewModifyCcnAttachedInstancesAttributeRequest()
		request.CcnId = &ccnId

		var ccnInstance vpc.CcnInstance
		ccnInstance.InstanceId = &instanceId
		ccnInstance.InstanceRegion = &instanceRegion
		ccnInstance.InstanceType = &instanceType
		ccnInstance.Description = &description

		request.Instances = []*vpc.CcnInstance{&ccnInstance}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyCcnAttachedInstancesAttribute(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), err.Error())
				return retryError(err)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		})
		if err != nil {
			return err
		}
	}
	return resourceTencentCloudCcnAttachmentRead(d, meta)
}

func resourceTencentCloudCcnAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, e := service.DescribeCcn(ctx, ccnId)
		if e != nil {
			return retryError(e)
		}
		if has == 0 {
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := service.DetachCcnInstances(ctx, ccnId, instanceRegion, instanceType, instanceId); err != nil {
		return err
	}

	return resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
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
