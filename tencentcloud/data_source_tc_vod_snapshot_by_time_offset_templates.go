/*
Use this data source to query detailed information of VOD snapshot by time offset templates.

Example Usage

```hcl
resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  width               = 128
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}

data "tencentcloud_vod_snapshot_by_time_offset_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVodSnapshotByTimeOffsetTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodSnapshotByTimeOffsetTemplatesRead,

		Schema: map[string]*schema.Schema{
			"definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique ID filter of snapshot by time offset template.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template type filter. Valid values: `Preset`: preset template; `Custom`: custom template.",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"template_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of snapshot by time offset templates. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of snapshot by time offset template.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template type filter. Valid values: `Preset`: preset template; `Custom`: custom template.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of a time point screen capturing template.",
						},
						"width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum value of the `width` (or long side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, width will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used.",
						},
						"height": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum value of the `height` (or short side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used.",
						},
						"resolution_adaptive": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Resolution adaption. Valid values: `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height.",
						},
						"format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image format. Valid values: `jpg`, `png`.",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template description.",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot `shorter` or `longer`; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. `white`: fill with white. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with white color blocks. `gauss`: fill with Gaussian blur. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with Gaussian blur.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of template in ISO date format.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time of template in ISO date format.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVodSnapshotByTimeOffsetTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vod_snapshot_by_time_offset_templates.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	filter := make(map[string]interface{})
	if v, ok := d.GetOk("definition"); ok {
		filter["definitions"] = []string{v.(string)}
	}
	if v, ok := d.GetOk("type"); ok {
		filter["type"] = v.(string)
	}
	if v, ok := d.GetOk("sub_app_id"); ok {
		filter["sub_appid"] = v.(int)
	}

	vodService := VodService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	templates, err := vodService.DescribeSnapshotByTimeOffsetTemplatesByFilter(ctx, filter)
	if err != nil {
		return err
	}

	templatesList := make([]map[string]interface{}, 0, len(templates))
	ids := make([]string, 0, len(templates))
	for _, item := range templates {
		definitionStr := strconv.FormatUint(*item.Definition, 10)
		templatesList = append(templatesList, map[string]interface{}{
			"definition":          definitionStr,
			"type":                item.Type,
			"name":                item.Name,
			"width":               item.Width,
			"height":              item.Height,
			"resolution_adaptive": *item.ResolutionAdaptive == "open",
			"format":              item.Format,
			"comment":             item.Comment,
			"fill_type":           item.FillType,
			"create_time":         item.CreateTime,
			"update_time":         item.UpdateTime,
		})
		ids = append(ids, definitionStr)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("template_list", templatesList); e != nil {
		log.Printf("[CRITAL]%s provider set vod snapshot by time offset template list fail, reason:%s ", logId, e.Error())
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), templatesList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
