package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoExportZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoExportZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone ID.",
			},
			"types": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of configuration types to export. If not specified, all configuration types will be exported. Valid values: L7AccelerationConfig (L7 acceleration configuration).",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exported configuration content in JSON format.",
			},
		},
	}
}

func dataSourceTencentCloudTeoExportZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_export_zone_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId := d.Get("zone_id").(string)
	if zoneId == "" {
		return fmt.Errorf("zone_id is required")
	}

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var respContent *string
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.TeoExportZoneConfig(ctx, zoneId, d.Get("types").([]interface{}))
		if e != nil {
			return tccommon.RetryError(e)
		}
		respContent = result
		return nil
	})
	if reqErr != nil {
		return reqErr
	}

	if respContent != nil {
		_ = d.Set("content", *respContent)
	}

	d.SetId(helper.DataResourceIdsHash([]string{zoneId}))

	return nil
}
