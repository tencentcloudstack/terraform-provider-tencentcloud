package cam

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
)

func ResourceTencentCloudCamRolePolicyAttachment() *schema.Resource {
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
				Description: "ID of the attached CAM role.",
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
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the policy.",
			},
		},
	}
}

func resourceTencentCloudCamRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_policy_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	roleId := d.Get("role_id").(string)
	policyId, e := strconv.Atoi(d.Get("policy_id").(string))
	if e != nil {
		return e
	}
	policyId64 := uint64(policyId)
	request := cam.NewAttachRolePolicyRequest()
	request.AttachRoleId = &roleId
	request.PolicyId = &policyId64

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

	d.SetId(roleId + "#" + strconv.Itoa(policyId))

	//get really instance then read
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	rolePolicyAttachmentId := d.Id()
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeRolePolicyAttachmentById(ctx, rolePolicyAttachmentId)
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
	return resourceTencentCloudCamRolePolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_policy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	rolePolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *cam.AttachedPolicyOfRole
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeRolePolicyAttachmentById(ctx, rolePolicyAttachmentId)
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

	roleId, policyId, e := camService.decodeCamPolicyAttachmentId(rolePolicyAttachmentId)
	if e != nil {
		return e
	}
	_ = d.Set("role_id", roleId)
	_ = d.Set("policy_id", strconv.Itoa(int(policyId)))
	_ = d.Set("policy_name", *instance.PolicyName)
	_ = d.Set("create_time", *instance.AddTime)
	_ = d.Set("create_mode", int(*instance.CreateMode))
	_ = d.Set("policy_type", *instance.PolicyType)

	return nil
}

func resourceTencentCloudCamRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role_policy_attachment.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	rolePolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteRolePolicyAttachmentById(ctx, rolePolicyAttachmentId)
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
