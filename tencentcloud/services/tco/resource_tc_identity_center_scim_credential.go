package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIdentityCenterScimCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterScimCredentialCreate,
		Read:   resourceTencentCloudIdentityCenterScimCredentialRead,
		Delete: resourceTencentCloudIdentityCenterScimCredentialDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Space ID. z-prefix starts with 12 random digits/lowercase letters.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCIM key status, Enabled-On, Disabled-Closed.",
			},

			"credential_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCIM key ID. scimcred-prefix and followed by 12 random digits/lowercase letters.",
			},

			"credential_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCIM credential type.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCIM create time.",
			},

			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCIM expire time.",
			},
		},
	}
}

func resourceTencentCloudIdentityCenterScimCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_credential.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId       string
		credentialId string
	)
	var (
		request  = organization.NewCreateSCIMCredentialRequest()
		response = organization.NewCreateSCIMCredentialResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	request.ZoneId = helper.String(zoneId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreateSCIMCredentialWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create identity center scim credential failed, reason:%+v", logId, err)
		return err
	}

	credentialId = *response.Response.CredentialId

	d.SetId(strings.Join([]string{zoneId, credentialId}, tccommon.FILED_SP))

	return resourceTencentCloudIdentityCenterScimCredentialRead(d, meta)
}

func resourceTencentCloudIdentityCenterScimCredentialRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_credential.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	credentialId := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeIdentityCenterScimCredentialById(ctx, zoneId, credentialId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_scim_credential` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.CredentialId != nil {
		_ = d.Set("credential_id", respData.CredentialId)
	}

	if respData.CredentialType != nil {
		_ = d.Set("credential_type", respData.CredentialType)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ExpireTime != nil {
		_ = d.Set("expire_time", respData.ExpireTime)
	}

	return nil
}

func resourceTencentCloudIdentityCenterScimCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_credential.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	credentialId := idSplit[1]

	var (
		request  = organization.NewDeleteSCIMCredentialRequest()
		response = organization.NewDeleteSCIMCredentialResponse()
	)

	request.ZoneId = helper.String(zoneId)

	request.CredentialId = helper.String(credentialId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().DeleteSCIMCredentialWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete identity center scim credential failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
