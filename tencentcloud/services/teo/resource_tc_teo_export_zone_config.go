package teo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoExportZoneConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoExportZoneConfigCreate,
		Read:   resourceTencentCloudTeoExportZoneConfigRead,
		Update: resourceTencentCloudTeoExportZoneConfigUpdate,
		Delete: resourceTencentCloudTeoExportZoneConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of configuration types to export. If left blank, all configuration types are exported. Supported values include: `L7AccelerationConfig`: Export Layer 7 acceleration configuration.",
			},

			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exported configuration content in JSON format, encoded in UTF-8.",
			},
		},
	}
}

func resourceTencentCloudTeoExportZoneConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId := d.Get("zone_id").(string)

	var types []*string
	if v, ok := d.GetOk("types"); ok {
		typeList := v.([]interface{})
		types = make([]*string, 0, len(typeList))
		for _, item := range typeList {
			types = append(types, helper.String(item.(string)))
		}
	}

	// Call ExportZoneConfig API
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		content, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoService().ExportZoneConfig(ctx, zoneId, types)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, zoneId [%s], types [%v], content length [%d]\n",
				logId, "ExportZoneConfig", zoneId, types, len(*content))
			if err := d.Set("content", content); err != nil {
				return resource.NonRetryableError(fmt.Errorf("set content failed: %s", err))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s export zone config failed, reason:%+v", logId, err)
		return err
	}

	// Set Resource ID as zoneId
	d.SetId(zoneId)

	return resourceTencentCloudTeoExportZoneConfigRead(d, meta)
}

func resourceTencentCloudTeoExportZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId := d.Id()

	// Verify the zone exists
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoService().DescribeTeoZone(ctx, zoneId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read zone config failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func resourceTencentCloudTeoExportZoneConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	// ExportZoneConfig is a read-only resource, changes require recreation
	// This function should never be called due to ForceNew attributes
	return fmt.Errorf("resource.tencentcloud_teo_export_zone_config does not support update, changes require recreation")
}

func resourceTencentCloudTeoExportZoneConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	log.Printf("[DEBUG]%s delete zone config export, resource id [%s]\n", logId, d.Id())

	// ExportZoneConfig is a read-only resource, perform logical deletion
	d.SetId("")

	return nil
}
