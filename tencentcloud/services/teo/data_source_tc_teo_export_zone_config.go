package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

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
				Description: "Specifies the site ID.",
			},

			"types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Types of configuration to export. If not specified, all types of configuration will be exported. Valid values: `L7AccelerationConfig`, `WebSecurity`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exported zone configuration content in JSON format.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoExportZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_export_zone_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("types"); ok {
		typesList := v.([]interface{})
		types := make([]*string, 0, len(typesList))
		for _, item := range typesList {
			types = append(types, helper.String(item.(string)))
		}
		paramMap["Types"] = types
	}

	var respData interface{}
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.ExportZoneConfigByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	exportResult := respData.(*teo.ExportZoneConfigResponseParams)
	if exportResult.Content != nil {
		_ = d.Set("content", exportResult.Content)
	}

	d.SetId(zoneId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		resultMap := map[string]interface{}{}
		if exportResult.Content != nil {
			resultMap["content"] = exportResult.Content
		}
		if e := tccommon.WriteToFile(output.(string), resultMap); e != nil {
			return e
		}
	}

	return nil
}
