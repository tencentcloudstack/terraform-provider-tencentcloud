package sts

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudStsAssumeRoleOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudStsAssumeRoleOperationCreate,
		Read:   resourceTencentCloudStsAssumeRoleOperationRead,
		Delete: resourceTencentCloudStsAssumeRoleOperationDelete,
		Schema: map[string]*schema.Schema{
			"role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource description of the role, which can be obtained by clicking the role name in [Access Management](https://console.cloud.tencent.com/cam/role).",
			},
			"role_session_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User-defined temporary session name. Length is between 2 and 128, can contain uppercase and lowercase characters, numbers, and special characters: =,.@_-.",
			},
			"duration_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the validity period of the temporary access credential in seconds. Default is 7200 seconds, maximum is 43200 seconds.",
			},
			"policy": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Policy description. The policy syntax refers to [CAM Policy Syntax](https://cloud.tencent.com/document/product/598/10603). The policy cannot contain the principal element.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Role external ID, which can be obtained by clicking the role name in [Access Management](https://console.cloud.tencent.com/cam/role). Length is between 2 and 128.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Session tag list. A maximum of 50 session tags can be passed, and duplicate tag keys are not supported.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key, up to 128 characters, case-sensitive.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value, up to 256 characters, case-sensitive.",
						},
					},
				},
			},
			"source_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Caller identity uin.",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "MFA serial number associated with the CAM user making the call. Format: qcs::cam:uin/${ownerUin}::mfa/${mfaType}. mfaType supports softToken.",
			},
			"token_code": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "MFA authentication code.",
			},
			"credentials": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Temporary access credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Token. The token length is up to 4096 bytes.",
						},
						"tmp_secret_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Temporary certificate secret ID. Maximum length is 1024 bytes.",
						},
						"tmp_secret_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Temporary certificate secret key. Maximum length is 1024 bytes.",
						},
					},
				},
			},
			"expired_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Expiration time of the temporary access credential, returned as a Unix timestamp accurate to seconds.",
			},
			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of the temporary access credential in ISO8601 format UTC time.",
			},
		},
	}
}

func resourceTencentCloudStsAssumeRoleOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sts_assume_role_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.Background()
		request = sts.NewAssumeRoleRequest()
	)

	request.RoleArn = helper.String(d.Get("role_arn").(string))
	request.RoleSessionName = helper.String(d.Get("role_session_name").(string))

	if v, ok := d.GetOk("duration_seconds"); ok {
		request.DurationSeconds = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("policy"); ok {
		request.Policy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("external_id"); ok {
		request.ExternalId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagList := v.([]interface{})
		tags := make([]*sts.Tag, 0, len(tagList))
		for _, item := range tagList {
			tagMap := item.(map[string]interface{})
			tag := &sts.Tag{
				Key:   helper.String(tagMap["key"].(string)),
				Value: helper.String(tagMap["value"].(string)),
			}
			tags = append(tags, tag)
		}
		request.Tags = tags
	}

	if v, ok := d.GetOk("source_identity"); ok {
		request.SourceIdentity = helper.String(v.(string))
	}

	if v, ok := d.GetOk("serial_number"); ok {
		request.SerialNumber = helper.String(v.(string))
	}

	if v, ok := d.GetOk("token_code"); ok {
		request.TokenCode = helper.String(v.(string))
	}

	var response *sts.AssumeRoleResponse
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseStsClient().AssumeRoleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s api[%s] fail, reason:%+v", logId, request.GetAction(), reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil {
		return fmt.Errorf("STS AssumeRole response is nil")
	}

	d.SetId(helper.BuildToken())

	if response.Response.Credentials != nil {
		credentialsList := []map[string]interface{}{
			{
				"token":          response.Response.Credentials.Token,
				"tmp_secret_id":  response.Response.Credentials.TmpSecretId,
				"tmp_secret_key": response.Response.Credentials.TmpSecretKey,
			},
		}
		_ = d.Set("credentials", credentialsList)
	}

	if response.Response.ExpiredTime != nil {
		_ = d.Set("expired_time", *response.Response.ExpiredTime)
	}

	if response.Response.Expiration != nil {
		_ = d.Set("expiration", *response.Response.Expiration)
	}

	return nil
}

func resourceTencentCloudStsAssumeRoleOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sts_assume_role_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudStsAssumeRoleOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sts_assume_role_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
