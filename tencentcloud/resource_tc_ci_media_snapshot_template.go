/*
Provides a resource to create a ci media_snapshot_template

Example Usage

```hcl
resource "tencentcloud_ci_media_snapshot_template" "media_snapshot_template" {
  name = &lt;nil&gt;
  snapshot {
		mode = &lt;nil&gt;
		start = &lt;nil&gt;
		time_interval = &lt;nil&gt;
		count = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		c_i_param = &lt;nil&gt;
		is_check_count = &lt;nil&gt;
		is_check_black = &lt;nil&gt;
		black_level = &lt;nil&gt;
		pixel_black_threshold = &lt;nil&gt;
		snapshot_out_mode = &lt;nil&gt;
		sprite_snapshot_config {
			cell_width = &lt;nil&gt;
			cell_height = &lt;nil&gt;
			padding = &lt;nil&gt;
			margin = &lt;nil&gt;
			color = &lt;nil&gt;
			columns = &lt;nil&gt;
			lines = &lt;nil&gt;
		}

  }
      }
```

Import

ci media_snapshot_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_snapshot_template.media_snapshot_template media_snapshot_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	ci "github.com/tencentyun/cos-go-sdk-v5"
	"log"
)

func resourceTencentCloudCiMediaSnapshotTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaSnapshotTemplateCreate,
		Read:   resourceTencentCloudCiMediaSnapshotTemplateRead,
		Update: resourceTencentCloudCiMediaSnapshotTemplateUpdate,
		Delete: resourceTencentCloudCiMediaSnapshotTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"snapshot": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Screenshot.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Screenshot mode, value range: {Interval, Average, KeyFrame}- Interval means interval mode Average means average mode- KeyFrame represents the key frame mode- Interval mode: Start, TimeInterval, The Count parameter takes effect. When Count is set and TimeInterval is not set, Indicates to capture all frames, a total of Count pictures- Average mode: Start, the Count parameter takes effect. express.",
						},
						"start": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Starting time, [0 video duration] in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
						},
						"time_interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Screenshot time interval, (0 3600], in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
						},
						"count": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Number of screenshots, range (0 10000].",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Wide, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "High, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video.",
						},
						"c_i_param": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Screenshot image processing parameters, for example: imageMogr2/format/png.",
						},
						"is_check_count": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to check the number of screenshots forcibly, when using custom interval mode to take screenshots, the video time is not long enough to capture Count screenshots, you can switch to average screenshot mode to capture Count screenshots.",
						},
						"is_check_black": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to enable black screen detection true/false.",
						},
						"black_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Screenshot black screen detection parameters, Valid when IsCheckBlack=true, Value reference range [30, 100], indicating the proportion of black pixels, the smaller the value, the smaller the proportion of black pixels, Start&amp;gt;0, the parameter setting is invalid, no filter black screen, Start =0 parameter is valid, the start time of the frame capture is the first frame non-black screen start.",
						},
						"pixel_black_threshold": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Screenshot black screen detection parameters, Valid when IsCheckBlack=true, The threshold for judging whether a pixel is a black point, value range: [0, 255].",
						},
						"snapshot_out_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Screenshot output mode parameters, Value range: {OnlySnapshot, OnlySprite, SnapshotAndSprite}, OnlySnapshot means output only screenshot mode OnlySprite means only output sprite mode SnapshotAndSprite means output screenshot and sprite mode.",
						},
						"sprite_snapshot_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Screenshot output configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cell_width": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Single image width Value range: [8, 4096], Unit: px.",
									},
									"cell_height": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Single image height Value range: [8, 4096], Unit: px.",
									},
									"padding": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Screenshot padding size, Value range: [8, 4096], Unit: px.",
									},
									"margin": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Screenshot margin size, Value range: [8, 4096], Unit: px.",
									},
									"color": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "See `https://www.ffmpeg.org/ffmpeg-utils.html#color-syntax` for details on supported colors.",
									},
									"columns": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Number of screenshot columns, value range: [1, 10000].",
									},
									"lines": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Number of screenshot lines, value range: [1, 10000].",
									},
								},
							},
						},
					},
				},
			},

			"template_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Template ID.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},
		},
	}
}

func resourceTencentCloudCiMediaSnapshotTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_snapshot_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ci.NewCreateMediaSnapshotTemplateRequest()
		response   = ci.NewCreateMediaSnapshotTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			snapshot := ci.Snapshot{}
			if v, ok := dMap["mode"]; ok {
				snapshot.Mode = helper.String(v.(string))
			}
			if v, ok := dMap["start"]; ok {
				snapshot.Start = helper.String(v.(string))
			}
			if v, ok := dMap["time_interval"]; ok {
				snapshot.TimeInterval = helper.String(v.(string))
			}
			if v, ok := dMap["count"]; ok {
				snapshot.Count = helper.String(v.(string))
			}
			if v, ok := dMap["width"]; ok {
				snapshot.Width = helper.String(v.(string))
			}
			if v, ok := dMap["height"]; ok {
				snapshot.Height = helper.String(v.(string))
			}
			if v, ok := dMap["c_i_param"]; ok {
				snapshot.CIParam = helper.String(v.(string))
			}
			if v, ok := dMap["is_check_count"]; ok {
				snapshot.IsCheckCount = helper.String(v.(string))
			}
			if v, ok := dMap["is_check_black"]; ok {
				snapshot.IsCheckBlack = helper.String(v.(string))
			}
			if v, ok := dMap["black_level"]; ok {
				snapshot.BlackLevel = helper.String(v.(string))
			}
			if v, ok := dMap["pixel_black_threshold"]; ok {
				snapshot.PixelBlackThreshold = helper.String(v.(string))
			}
			if v, ok := dMap["snapshot_out_mode"]; ok {
				snapshot.SnapshotOutMode = helper.String(v.(string))
			}
			if v, ok := dMap["sprite_snapshot_config"]; ok {
				for _, item := range v.([]interface{}) {
					spriteSnapshotConfigMap := item.(map[string]interface{})
					snapshot := ci.Snapshot{}
					if v, ok := spriteSnapshotConfigMap["cell_width"]; ok {
						snapshot.CellWidth = helper.String(v.(string))
					}
					if v, ok := spriteSnapshotConfigMap["cell_height"]; ok {
						snapshot.CellHeight = helper.String(v.(string))
					}
					if v, ok := spriteSnapshotConfigMap["padding"]; ok {
						snapshot.Padding = helper.String(v.(string))
					}
					if v, ok := spriteSnapshotConfigMap["margin"]; ok {
						snapshot.Margin = helper.String(v.(string))
					}
					if v, ok := spriteSnapshotConfigMap["color"]; ok {
						snapshot.Color = helper.String(v.(string))
					}
					if v, ok := spriteSnapshotConfigMap["columns"]; ok {
						snapshot.Columns = helper.String(v.(string))
					}
					if v, ok := spriteSnapshotConfigMap["lines"]; ok {
						snapshot.Lines = helper.String(v.(string))
					}
					snapshot.SpriteSnapshotConfig = append(snapshot.SpriteSnapshotConfig, &snapshot)
				}
			}
			request.Snapshot = append(request.Snapshot, &snapshot)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaSnapshotTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSnapshotTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_snapshot_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaSnapshotTemplateId := d.Id()

	mediaSnapshotTemplate, err := service.DescribeCiMediaSnapshotTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaSnapshotTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaSnapshotTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaSnapshotTemplate.Name != nil {
		_ = d.Set("name", mediaSnapshotTemplate.Name)
	}

	if mediaSnapshotTemplate.Snapshot != nil {
		snapshotList := []interface{}{}
		for _, snapshot := range mediaSnapshotTemplate.Snapshot {
			snapshotMap := map[string]interface{}{}

			if mediaSnapshotTemplate.Snapshot.Mode != nil {
				snapshotMap["mode"] = mediaSnapshotTemplate.Snapshot.Mode
			}

			if mediaSnapshotTemplate.Snapshot.Start != nil {
				snapshotMap["start"] = mediaSnapshotTemplate.Snapshot.Start
			}

			if mediaSnapshotTemplate.Snapshot.TimeInterval != nil {
				snapshotMap["time_interval"] = mediaSnapshotTemplate.Snapshot.TimeInterval
			}

			if mediaSnapshotTemplate.Snapshot.Count != nil {
				snapshotMap["count"] = mediaSnapshotTemplate.Snapshot.Count
			}

			if mediaSnapshotTemplate.Snapshot.Width != nil {
				snapshotMap["width"] = mediaSnapshotTemplate.Snapshot.Width
			}

			if mediaSnapshotTemplate.Snapshot.Height != nil {
				snapshotMap["height"] = mediaSnapshotTemplate.Snapshot.Height
			}

			if mediaSnapshotTemplate.Snapshot.CIParam != nil {
				snapshotMap["c_i_param"] = mediaSnapshotTemplate.Snapshot.CIParam
			}

			if mediaSnapshotTemplate.Snapshot.IsCheckCount != nil {
				snapshotMap["is_check_count"] = mediaSnapshotTemplate.Snapshot.IsCheckCount
			}

			if mediaSnapshotTemplate.Snapshot.IsCheckBlack != nil {
				snapshotMap["is_check_black"] = mediaSnapshotTemplate.Snapshot.IsCheckBlack
			}

			if mediaSnapshotTemplate.Snapshot.BlackLevel != nil {
				snapshotMap["black_level"] = mediaSnapshotTemplate.Snapshot.BlackLevel
			}

			if mediaSnapshotTemplate.Snapshot.PixelBlackThreshold != nil {
				snapshotMap["pixel_black_threshold"] = mediaSnapshotTemplate.Snapshot.PixelBlackThreshold
			}

			if mediaSnapshotTemplate.Snapshot.SnapshotOutMode != nil {
				snapshotMap["snapshot_out_mode"] = mediaSnapshotTemplate.Snapshot.SnapshotOutMode
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig != nil {
				spriteSnapshotConfigList := []interface{}{}
				for _, spriteSnapshotConfig := range mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig {
					spriteSnapshotConfigMap := map[string]interface{}{}

					if spriteSnapshotConfig.CellWidth != nil {
						spriteSnapshotConfigMap["cell_width"] = spriteSnapshotConfig.CellWidth
					}

					if spriteSnapshotConfig.CellHeight != nil {
						spriteSnapshotConfigMap["cell_height"] = spriteSnapshotConfig.CellHeight
					}

					if spriteSnapshotConfig.Padding != nil {
						spriteSnapshotConfigMap["padding"] = spriteSnapshotConfig.Padding
					}

					if spriteSnapshotConfig.Margin != nil {
						spriteSnapshotConfigMap["margin"] = spriteSnapshotConfig.Margin
					}

					if spriteSnapshotConfig.Color != nil {
						spriteSnapshotConfigMap["color"] = spriteSnapshotConfig.Color
					}

					if spriteSnapshotConfig.Columns != nil {
						spriteSnapshotConfigMap["columns"] = spriteSnapshotConfig.Columns
					}

					if spriteSnapshotConfig.Lines != nil {
						spriteSnapshotConfigMap["lines"] = spriteSnapshotConfig.Lines
					}

					spriteSnapshotConfigList = append(spriteSnapshotConfigList, spriteSnapshotConfigMap)
				}

				snapshotMap["sprite_snapshot_config"] = []interface{}{spriteSnapshotConfigList}
			}

			snapshotList = append(snapshotList, snapshotMap)
		}

		_ = d.Set("snapshot", snapshotList)

	}

	if mediaSnapshotTemplate.TemplateId != nil {
		_ = d.Set("template_id", mediaSnapshotTemplate.TemplateId)
	}

	if mediaSnapshotTemplate.UpdateTime != nil {
		_ = d.Set("update_time", mediaSnapshotTemplate.UpdateTime)
	}

	if mediaSnapshotTemplate.CreateTime != nil {
		_ = d.Set("create_time", mediaSnapshotTemplate.CreateTime)
	}

	return nil
}

func resourceTencentCloudCiMediaSnapshotTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_snapshot_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaSnapshotTemplateRequest()

	mediaSnapshotTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "snapshot", "template_id", "update_time", "create_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaSnapshotTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSnapshotTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_snapshot_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	mediaSnapshotTemplateId := d.Id()

	if err := service.DeleteCiMediaSnapshotTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
