package cam

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCamRolePolicyAttachmentByName() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamRolePolicyAttachmentByNameCreate,
		Read:   resourceTencentCloudCamRolePolicyAttachmentByNameRead,
		Delete: resourceTencentCloudCamRolePolicyAttachmentByNameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the attached CAM role.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the policy.",
			},
			"create_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Mode of Creation of the CAM role policy attachment. `1` means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the policy strategy. `User` means customer strategy and `QCS` means preset strategy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the CAM role policy attachment.",
			},
		},
	}
}

func resourceTencentCloudCamRolePolicyAttachmentByNameCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_policy_attachment_by_name.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	roleName := d.Get("role_name").(string)
	policyName := d.Get("policy_name").(string)

	request := cam.NewAttachRolePolicyRequest()
	request.PolicyName = helper.String(policyName)
	request.AttachRoleName = helper.String(roleName)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().AttachRolePolicy(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
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

	d.SetId(roleName + "#" + policyName)

	//get really instance then read
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	params := make(map[string]interface{})
	params["policy_name"] = policyName
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeRolePolicyAttachmentByName(ctx, roleName, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(10 * time.Second)
	return resourceTencentCloudCamRolePolicyAttachmentByNameRead(d, meta)
}

func resourceTencentCloudCamRolePolicyAttachmentByNameRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_policy_attachment_by_name.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	rolePolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *cam.AttachedPolicyOfRole
	items := strings.Split(rolePolicyAttachmentId, "#")
	if len(items) < 2 {
		return fmt.Errorf("RolePolicyAttachmentId is invalid!")
	}
	roleName, policyName := items[0], items[1]
	params := make(map[string]interface{})
	params["policy_name"] = policyName
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeRolePolicyAttachmentByName(ctx, roleName, params)
		if e != nil {
			return tccommon.RetryError(e)
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
	_ = d.Set("role_name", roleName)
	_ = d.Set("policy_name", policyName)
	_ = d.Set("create_time", *instance.AddTime)
	_ = d.Set("create_mode", int(*instance.CreateMode))
	_ = d.Set("policy_type", *instance.PolicyType)

	return nil
}

func resourceTencentCloudCamRolePolicyAttachmentByNameDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_policy_attachment_by_name.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	rolePolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	items := strings.Split(rolePolicyAttachmentId, "#")
	if len(items) < 2 {
		return fmt.Errorf("RolePolicyAttachmentId is invalid!")
	}
	roleName, policyName := items[0], items[1]
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteRolePolicyAttachmentByName(ctx, roleName, policyName)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM role policy attachment failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
