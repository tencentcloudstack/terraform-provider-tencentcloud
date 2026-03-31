package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID. Example: zone-2zpqp7qztest",
			},
			"types": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of configuration types to export. If not specified, all configuration types will be exported. Currently supported types include: <li>L7AccelerationConfig: Export L7 acceleration configuration, corresponding to \"Site Acceleration - Global Acceleration Configuration\" and \"Site Acceleration - Rule Engine\" in the console.</li>Note: The supported export types will increase with iteration. When exporting all types, please pay attention to the size of the exported file. It is recommended to specify the configuration types to be exported to control the size of the request response payload. Example: L7AccelerationConfig",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The specific content of the exported configuration. Returned in JSON format and encoded in UTF-8. For configuration content, refer to the example below.",
			},
		},
	}
}

func resourceTencentCloudTeoExportZoneConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	zoneId := d.Get("zone_id").(string)

	d.SetId(zoneId)

	return resourceTencentCloudTeoExportZoneConfigRead(d, meta)
}

func resourceTencentCloudTeoExportZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	zoneId := d.Id()

	var types []*string
	if v, ok := d.Get("types").([]interface{}); ok && len(v) > 0 {
		types = make([]*string, 0, len(v))
		for _, item := range v {
			types = append(types, helper.String(item.(string)))
		}
	}

	respData, err := service.ExportTeoZoneConfigById(ctx, zoneId, types)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_export_zone_config` [%s] not found.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	return nil
}

func resourceTencentCloudTeoExportZoneConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudTeoExportZoneConfigRead(d, meta)
}

func resourceTencentCloudTeoExportZoneConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	log.Printf("[INFO]%s resource `teo_export_zone_config` [%s] deleted.\n", logId, d.Id())

	d.SetId("")
	return nil
}
