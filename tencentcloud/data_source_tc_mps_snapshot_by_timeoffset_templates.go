/*
Use this data source to query detailed information of mps snapshot_by_timeoffset_templates

Example Usage

```hcl
data "tencentcloud_mps_snapshot_by_timeoffset_templates" "snapshot_by_timeoffset_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  snapshot_by_time_offset_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		resolution_adaptive = &lt;nil&gt;
		format = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		fill_type = &lt;nil&gt;

  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMpsSnapshotByTimeoffsetTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsSnapshotByTimeoffsetTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The snapshot by timeoffset template uniquely identifies the filter condition, and the array length limit: 100.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset, default: 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return the number of records, default value: 10, maximum value: 100.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template type filter condition, optional value:Preset: system preset template.Custom: user-defined template.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"snapshot_by_time_offset_template_set": {
				Type:        schema.TypeList,
				Description: "Snapshot by timeoffset template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Description: "The unique identifier of the snapshot by timeoffset template.",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Template type, optional value:Preset: system preset template.Custom: user-defined template.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Snapshot by timeoffset template name.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "The description information of snapshot by timeoffset template.",
						},
						"width": {
							Type:        schema.TypeInt,
							Description: "The maximum value of the snapshot width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
						},
						"height": {
							Type:        schema.TypeInt,
							Description: "The maximum value of the snapshot height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
						},
						"resolution_adaptive": {
							Type:        schema.TypeString,
							Description: "Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.",
						},
						"format": {
							Type:        schema.TypeString,
							Description: "Image format.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Template creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Template last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Description: "Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and use Gaussian blur for the rest of the edge.Default value: black.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMpsSnapshotByTimeoffsetTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_snapshot_by_timeoffset_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("definitions"); ok {
		definitionsSet := v.(*schema.Set).List()
		for i := range definitionsSet {
			definitions := definitionsSet[i].(int)
			paramMap["Definitions"] = append(paramMap["Definitions"], helper.IntUint64(definitions))
		}
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("snapshot_by_time_offset_template_set"); ok {
		snapshotByTimeOffsetTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.SnapshotByTimeOffsetTemplate, 0, len(snapshotByTimeOffsetTemplateSetSet))

		for _, item := range snapshotByTimeOffsetTemplateSetSet {
			snapshotByTimeOffsetTemplate := mps.SnapshotByTimeOffsetTemplate{}
			snapshotByTimeOffsetTemplateMap := item.(map[string]interface{})

			if v, ok := snapshotByTimeOffsetTemplateMap["definition"]; ok {
				snapshotByTimeOffsetTemplate.Definition = helper.IntUint64(v.(int))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["type"]; ok {
				snapshotByTimeOffsetTemplate.Type = helper.String(v.(string))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["name"]; ok {
				snapshotByTimeOffsetTemplate.Name = helper.String(v.(string))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["comment"]; ok {
				snapshotByTimeOffsetTemplate.Comment = helper.String(v.(string))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["width"]; ok {
				snapshotByTimeOffsetTemplate.Width = helper.IntUint64(v.(int))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["height"]; ok {
				snapshotByTimeOffsetTemplate.Height = helper.IntUint64(v.(int))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["resolution_adaptive"]; ok {
				snapshotByTimeOffsetTemplate.ResolutionAdaptive = helper.String(v.(string))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["format"]; ok {
				snapshotByTimeOffsetTemplate.Format = helper.String(v.(string))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["create_time"]; ok {
				snapshotByTimeOffsetTemplate.CreateTime = helper.String(v.(string))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["update_time"]; ok {
				snapshotByTimeOffsetTemplate.UpdateTime = helper.String(v.(string))
			}
			if v, ok := snapshotByTimeOffsetTemplateMap["fill_type"]; ok {
				snapshotByTimeOffsetTemplate.FillType = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &snapshotByTimeOffsetTemplate)
		}
		paramMap["snapshot_by_time_offset_template_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var snapshotByTimeOffsetTemplateSet []*mps.SnapshotByTimeOffsetTemplate

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsSnapshotByTimeoffsetTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		snapshotByTimeOffsetTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(snapshotByTimeOffsetTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(snapshotByTimeOffsetTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
