/*
Provides a resource to create a CAM user policy attachment.

Example Usage

```hcl
variable "cam_user_basic" {
  default = "keep-cam-user"
}

resource "tencentcloud_cam_policy" "policy_basic" {
  name        = "tf_cam_attach_user_policy"
  document    =jsonencode({
    "version":"2.0",
    "statement":[
      {
        "action":["cos:*"],
        "resource":["*"],
        "effect":"allow",
      },
      {
        "effect":"allow",
        "action":["monitor:*","cam:ListUsersForGroup","cam:ListGroups","cam:GetGroup"],
        "resource":["*"],
      }
    ]
  })
  description = "tf_test"
}

data "tencentcloud_cam_users" "users" {
  name = var.cam_user_basic
}

resource "tencentcloud_cam_user_policy_attachment" "user_policy_attachment_basic" {
  user_name = data.tencentcloud_cam_users.users.user_list.0.user_id
  policy_id = tencentcloud_cam_policy.policy_basic.id
}
```

Import

CAM user policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_user_policy_attachment.foo cam-test#26800353
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
)

func resourceTencentCloudCamUserPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserPolicyAttachmentCreate,
		Read:   resourceTencentCloudCamUserPolicyAttachmentRead,
		Delete: resourceTencentCloudCamUserPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"user_name", "user_id"},
				Deprecated:   "It has been deprecated from version 1.59.5. Use `user_name` instead.",
				Description:  "ID of the attached CAM user.",
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"user_name", "user_id"},
				Description:  "Name of the attached CAM user as uniq key.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the policy.",
			},
			"create_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Mode of Creation of the CAM user policy attachment. `1` means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the policy strategy. `User` means customer strategy and `QCS` means preset strategy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CAM user policy attachment.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the policy.",
			},
		},
	}
}

func resourceTencentCloudCamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_policy_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	userId, _, err := getUserId(d)
	if err != nil {
		return err
	}
	policyId := d.Get("policy_id").(string)
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := camService.AddUserPolicyAttachment(ctx, userId, policyId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM user policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(userId + "#" + policyId)

	//get really instance then read

	userPolicyAttachmentId := d.Id()
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeUserPolicyAttachmentById(ctx, userPolicyAttachmentId)
		if e != nil {
			return retryError(e)
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM user policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(10 * time.Second)
	return resourceTencentCloudCamUserPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_policy_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	userPolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cam.AttachPolicyInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeUserPolicyAttachmentById(ctx, userPolicyAttachmentId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM user policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	//split id
	userId, policyId, e := camService.decodeCamPolicyAttachmentId(userPolicyAttachmentId)
	if e != nil {
		return e
	}
	_ = d.Set("user_name", userId)
	_ = d.Set("policy_id", strconv.Itoa(int(policyId)))
	_ = d.Set("policy_name", *instance.PolicyName)
	_ = d.Set("create_time", *instance.AddTime)
	_ = d.Set("create_mode", int(*instance.CreateMode))
	_ = d.Set("policy_type", *instance.PolicyType)
	return nil
}

func resourceTencentCloudCamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_policy_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	userPolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteUserPolicyAttachmentById(ctx, userPolicyAttachmentId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM user policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}

func getUserId(d *schema.ResourceData) (value string, usingName bool, err error) {
	name, hasName := d.GetOk("user_name")
	id, hasId := d.GetOk("user_id")
	if hasName {
		return name.(string), true, nil
	} else if hasId {
		return id.(string), false, nil
	}
	return "", false, fmt.Errorf("no user name provided")
}
