/*
Provides a resource to create a CAM role policy attachment.

Example Usage

```hcl
resource "tencentcloud_cam_role_policy_attachment" "foo" {
  role_id   = "${tencentcloud_cam_role.foo.id}"
  policy_id = "${tencentcloud_cam_policy.foo.id}"
}
```

Import

CAM role policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role_policy_attachment.foo 4611686018427922725#26800353
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
)

func resourceTencentCloudCamRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamRolePolicyAttachmentCreate,
		Read:   resourceTencentCloudCamRolePolicyAttachmentRead,
		Delete: resourceTencentCloudCamRolePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the attached CAM role.",
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
				Description: "Mode of Creation of the CAM role policy attachment. 1 means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the CAM role policy attachment.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the policy.",
			},
		},
	}
}

func resourceTencentCloudCamRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_policy_attachment.create")()

	logId := getLogId(contextNil)

	roleId := d.Get("role_id").(string)
	policyId, e := strconv.Atoi(d.Get("policy_id").(string))
	if e != nil {
		return e
	}
	policyId64 := uint64(policyId)
	request := cam.NewAttachRolePolicyRequest()
	request.AttachRoleId = &roleId
	request.PolicyId = &policyId64

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().AttachRolePolicy(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM role policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(roleId + "#" + strconv.Itoa(policyId))
	time.Sleep(3 * time.Second)

	return resourceTencentCloudCamRolePolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_policy_attachment.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	rolePolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cam.AttachedPolicyOfRole
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeRolePolicyAttachmentById(ctx, rolePolicyAttachmentId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	roleId, policyId, e := camService.decodeCamPolicyAttachmentId(rolePolicyAttachmentId)
	if e != nil {
		return e
	}
	d.Set("role_id", roleId)
	d.Set("policy_id", strconv.Itoa(int(policyId)))
	d.Set("policy_name", *instance.PolicyName)
	d.Set("create_time", *instance.AddTime)
	d.Set("create_mode", int(*instance.CreateMode))
	d.Set("policy_type", *instance.PolicyType)

	return nil
}

func resourceTencentCloudCamRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_policy_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	rolePolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteRolePolicyAttachmentById(ctx, rolePolicyAttachmentId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM role policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
