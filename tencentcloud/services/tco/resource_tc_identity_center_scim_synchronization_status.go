package tco

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIdentityCenterScimSynchronizationStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterScimSynchronizationStatusCreate,
		Read:   resourceTencentCloudIdentityCenterScimSynchronizationStatusRead,
		Update: resourceTencentCloudIdentityCenterScimSynchronizationStatusUpdate,
		Delete: resourceTencentCloudIdentityCenterScimSynchronizationStatusDelete,
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

			"scim_synchronization_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SCIM synchronization status. Enabled-enabled. Disabled-disables.",
			},
		},
	}
}

func resourceTencentCloudIdentityCenterScimSynchronizationStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_synchronization_status.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		zoneId string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)

	return resourceTencentCloudIdentityCenterScimSynchronizationStatusUpdate(d, meta)
}

func resourceTencentCloudIdentityCenterScimSynchronizationStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_synchronization_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	zoneId := d.Id()

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeIdentityCenterScimSynchronizationStatusById(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_scim_synchronization_status` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.SCIMSynchronizationStatus != nil {
		_ = d.Set("scim_synchronization_status", respData.SCIMSynchronizationStatus)
	}

	return nil
}

func resourceTencentCloudIdentityCenterScimSynchronizationStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_synchronization_status.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId := d.Id()

	needChange := false
	mutableArgs := []string{"scim_synchronization_status"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := organization.NewUpdateSCIMSynchronizationStatusRequest()

		request.ZoneId = helper.String(zoneId)

		if v, ok := d.GetOk("scim_synchronization_status"); ok {
			request.SCIMSynchronizationStatus = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateSCIMSynchronizationStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update identity center scim synchronization status failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudIdentityCenterScimSynchronizationStatusRead(d, meta)
}

func resourceTencentCloudIdentityCenterScimSynchronizationStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_scim_synchronization_status.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
