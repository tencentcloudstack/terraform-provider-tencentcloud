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

func ResourceTencentCloudIdentityCenterScimCredentialStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterScimCredentialStatusCreate,
		Read:   resourceTencentCloudIdentityCenterScimCredentialStatusRead,
		Update: resourceTencentCloudIdentityCenterScimCredentialStatusUpdate,
		Delete: resourceTencentCloudIdentityCenterScimCredentialStatusDelete,
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

			"credential_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SCIM key ID. scimcred-prefix and followed by 12 random digits/lowercase letters.",
			},

			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SCIM key status. Enabled-enabled. Disabled-disabled.",
			},
		},
	}
}

func resourceTencentCloudIdentityCenterScimCredentialStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_credential_status.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		zoneId       string
		credentialId string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("credential_id"); ok {
		credentialId = v.(string)
	}

	d.SetId(strings.Join([]string{zoneId, credentialId}, tccommon.FILED_SP))

	return resourceTencentCloudIdentityCenterScimCredentialStatusUpdate(d, meta)
}

func resourceTencentCloudIdentityCenterScimCredentialStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_credential_status.read")()
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
	_ = d.Set("credential_id", credentialId)

	respData, err := service.DescribeIdentityCenterScimCredentialStatusById(ctx, zoneId, credentialId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_scim_credential_status` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	return nil
}

func resourceTencentCloudIdentityCenterScimCredentialStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_credential_status.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	credentialId := idSplit[1]

	needChange := false
	mutableArgs := []string{"status"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := organization.NewUpdateSCIMCredentialStatusRequest()

		request.ZoneId = helper.String(zoneId)

		request.CredentialId = helper.String(credentialId)

		if v, ok := d.GetOk("status"); ok {
			request.NewStatus = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateSCIMCredentialStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update identity center scim credential status failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudIdentityCenterScimCredentialStatusRead(d, meta)
}

func resourceTencentCloudIdentityCenterScimCredentialStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_credential_status.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
