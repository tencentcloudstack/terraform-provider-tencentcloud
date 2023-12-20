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

func ResourceTencentCloudCamUserPolicyAttachment() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user_policy_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	userId, _, err := getUserId(d)
	if err != nil {
		return err
	}
	policyId := d.Get("policy_id").(string)
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := camService.AddUserPolicyAttachment(ctx, userId, policyId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return tccommon.RetryError(e)
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
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeUserPolicyAttachmentById(ctx, userPolicyAttachmentId)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user_policy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	userPolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *cam.AttachPolicyInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeUserPolicyAttachmentById(ctx, userPolicyAttachmentId)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user_policy_attachment.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	userPolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteUserPolicyAttachmentById(ctx, userPolicyAttachmentId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return tccommon.RetryError(e)
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
