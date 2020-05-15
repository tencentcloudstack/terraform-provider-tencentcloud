/*
Provides a resource to create a CAM group policy attachment.

Example Usage

```hcl
resource "tencentcloud_cam_group_policy_attachment" "foo" {
  group_id  = tencentcloud_cam_group.foo.id
  policy_id = tencentcloud_cam_policy.foo.id
}
```

Import

CAM group policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_policy_attachment.foo 12515263#26800353
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
)

func resourceTencentCloudCamGroupPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamGroupPolicyAttachmentCreate,
		Read:   resourceTencentCloudCamGroupPolicyAttachmentRead,
		Delete: resourceTencentCloudCamGroupPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the attached CAM group.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the policy.",
			},
			"create_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Mode of Creation of the CAM group policy attachment. 1 means the cam policy attachment is created by production, and the others indicate syntax strategy ways.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the policy strategy. 'Group' means customer strategy and 'QCS' means preset strategy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CAM group policy attachment.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the policy.",
			},
		},
	}
}

func resourceTencentCloudCamGroupPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group_policy_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	groupId := d.Get("group_id").(string)
	policyId := d.Get("policy_id").(string)
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := camService.AddGroupPolicyAttachment(ctx, groupId, policyId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM group policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(groupId + "#" + policyId)

	//get really instance then read
	groupPolicyAttachmentId := d.Id()
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeGroupPolicyAttachmentById(ctx, groupPolicyAttachmentId)
		if e != nil {
			return retryError(e)
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(10 * time.Second)
	return resourceTencentCloudCamGroupPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCamGroupPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group_policy_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	groupPolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cam.AttachPolicyInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeGroupPolicyAttachmentById(ctx, groupPolicyAttachmentId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	//split id
	groupId, policyId, e := camService.decodeCamPolicyAttachmentId(groupPolicyAttachmentId)
	if e != nil {
		return e
	}
	_ = d.Set("group_id", groupId)
	_ = d.Set("policy_id", strconv.Itoa(int(policyId)))
	_ = d.Set("policy_name", *instance.PolicyName)
	_ = d.Set("create_time", *instance.AddTime)
	_ = d.Set("create_mode", int(*instance.CreateMode))
	_ = d.Set("policy_type", *instance.PolicyType)
	return nil
}

func resourceTencentCloudCamGroupPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group_policy_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	groupPolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteGroupPolicyAttachmentById(ctx, groupPolicyAttachmentId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM group policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
