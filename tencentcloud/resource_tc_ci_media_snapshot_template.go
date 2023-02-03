/*
Provides a resource to create a ci media_snapshot_template

Example Usage

```hcl
resource "tencentcloud_ci_media_snapshot_template" "media_snapshot_template" {
    bucket = "terraform-ci-xxxxxx"
  	name = "snapshot_template_test"
  	snapshot {
      count = "10"
      snapshot_out_mode = "SnapshotAndSprite"
      sprite_snapshot_config {
        color = "White"
        columns = "10"
        lines = "10"
        margin = "10"
        padding = "10"
      }
  	}
}
```

Import

ci media_snapshot_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_snapshot_template.media_snapshot_template terraform-ci-xxxxxx#t18210645f96564eaf80e86b1f58c20152
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentyun/cos-go-sdk-v5"
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
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "bucket name.",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"snapshot": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "screenshot.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot mode, value range: {Interval, Average, KeyFrame}- Interval means interval mode Average means average mode- KeyFrame represents the key frame mode- Interval mode: Start, TimeInterval, The Count parameter takes effect. When Count is set and TimeInterval is not set, Indicates to capture all frames, a total of Count pictures- Average mode: Start, the Count parameter takes effect. express.",
						},
						"start": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Starting time, [0 video duration] in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
						},
						"time_interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
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
							Computed:    true,
							Description: "wide, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "high, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video.",
						},
						"ci_param": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot image processing parameters, for example: imageMogr2/format/png.",
						},
						"is_check_count": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Whether to check the number of screenshots forcibly, when using custom interval mode to take screenshots, the video time is not long enough to capture Count screenshots, you can switch to average screenshot mode to capture Count screenshots.",
						},
						"is_check_black": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Whether to enable black screen detection true/false.",
						},
						"black_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot black screen detection parameters, Valid when IsCheckBlack=true, Value reference range [30, 100], indicating the proportion of black pixels, the smaller the value, the smaller the proportion of black pixels, Start&gt;0, the parameter setting is invalid, no filter black screen, Start =0 parameter is valid, the start time of the frame capture is the first frame non-black screen start.",
						},
						"pixel_black_threshold": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot black screen detection parameters, Valid when IsCheckBlack=true, The threshold for judging whether a pixel is a black point, value range: [0, 255].",
						},
						"snapshot_out_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot output mode parameters, Value range: {OnlySnapshot, OnlySprite, SnapshotAndSprite}, OnlySnapshot means output only screenshot mode OnlySprite means only output sprite mode SnapshotAndSprite means output screenshot and sprite mode.",
						},
						"sprite_snapshot_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Screenshot output configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cell_width": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Single image width Value range: [8, 4096], Unit: px.",
									},
									"cell_height": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Single image height Value range: [8, 4096], Unit: px.",
									},
									"padding": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "screenshot padding size, Value range: [8, 4096], Unit: px.",
									},
									"margin": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "screenshot margin size, Value range: [8, 4096], Unit: px.",
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
				Description: "update time.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "creation time.",
			},
		},
	}
}

func resourceTencentCloudCiMediaSnapshotTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_snapshot_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var templateId string
	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}
	request := cos.CreateMediaSnapshotTemplateOptions{
		Tag: "Snapshot",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "snapshot"); ok {
		snapshot := cos.Snapshot{}
		if v, ok := dMap["mode"]; ok {
			snapshot.Mode = v.(string)
		}
		if v, ok := dMap["start"]; ok {
			snapshot.Start = v.(string)
		}
		if v, ok := dMap["time_interval"]; ok {
			snapshot.TimeInterval = v.(string)
		}
		if v, ok := dMap["count"]; ok {
			snapshot.Count = v.(string)
		}
		if v, ok := dMap["width"]; ok {
			snapshot.Width = v.(string)
		}
		if v, ok := dMap["height"]; ok {
			snapshot.Height = v.(string)
		}
		if v, ok := dMap["ci_param"]; ok {
			snapshot.CIParam = v.(string)
		}
		if v, ok := dMap["is_check_count"]; ok {
			if v.(string) == "true" {
				snapshot.IsCheckCount = true
			} else {
				snapshot.IsCheckCount = false
			}
		}
		if v, ok := dMap["is_check_black"]; ok {
			if v.(string) == "true" {
				snapshot.IsCheckBlack = true
			} else {
				snapshot.IsCheckBlack = false
			}
		}
		if v, ok := dMap["black_level"]; ok {
			snapshot.BlackLevel = v.(string)
		}
		if v, ok := dMap["pixel_black_threshold"]; ok {
			snapshot.PixelBlackThreshold = v.(string)
		}
		if v, ok := dMap["snapshot_out_mode"]; ok {
			snapshot.SnapshotOutMode = v.(string)
		}
		if spriteSnapshotConfigMap, ok := helper.InterfaceToMap(dMap, "sprite_snapshot_config"); ok {
			spriteSnapshotConfig := cos.SpriteSnapshotConfig{}
			if v, ok := spriteSnapshotConfigMap["cell_width"]; ok {
				spriteSnapshotConfig.CellWidth = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["cell_height"]; ok {
				spriteSnapshotConfig.CellHeight = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["padding"]; ok {
				spriteSnapshotConfig.Padding = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["margin"]; ok {
				spriteSnapshotConfig.Margin = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["color"]; ok {
				spriteSnapshotConfig.Color = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["columns"]; ok {
				spriteSnapshotConfig.Columns = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["lines"]; ok {
				spriteSnapshotConfig.Lines = v.(string)
			}
			snapshot.SpriteSnapshotConfig = &spriteSnapshotConfig
		}
		request.Snapshot = &snapshot
	}

	var response *cos.CreateMediaTemplateResult
	ciClient := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := ciClient.CI.CreateMediaSnapshotTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%+v]\n", logId, "CreateMediaSnapshotTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSnapshotTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_snapshot_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	mediaSnapshotTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaSnapshotTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)
	if mediaSnapshotTemplate.Name != "" {
		_ = d.Set("name", mediaSnapshotTemplate.Name)
	}

	log.Printf("[DEBUG]Snapshot api[%+v]", mediaSnapshotTemplate.Snapshot)
	if mediaSnapshotTemplate.Snapshot != nil {
		snapshotMap := map[string]interface{}{}

		if mediaSnapshotTemplate.Snapshot.Mode != "" {
			snapshotMap["mode"] = mediaSnapshotTemplate.Snapshot.Mode
		}

		if mediaSnapshotTemplate.Snapshot.Start != "" {
			snapshotMap["start"] = mediaSnapshotTemplate.Snapshot.Start
		}

		if mediaSnapshotTemplate.Snapshot.TimeInterval != "" {
			snapshotMap["time_interval"] = mediaSnapshotTemplate.Snapshot.TimeInterval
		}

		if mediaSnapshotTemplate.Snapshot.Count != "" {
			snapshotMap["count"] = mediaSnapshotTemplate.Snapshot.Count
		}

		if mediaSnapshotTemplate.Snapshot.Width != "" {
			snapshotMap["width"] = mediaSnapshotTemplate.Snapshot.Width
		}

		if mediaSnapshotTemplate.Snapshot.Height != "" {
			snapshotMap["height"] = mediaSnapshotTemplate.Snapshot.Height
		}

		if mediaSnapshotTemplate.Snapshot.CIParam != "" {
			snapshotMap["ci_param"] = mediaSnapshotTemplate.Snapshot.CIParam
		}

		snapshotMap["is_check_count"] = fmt.Sprintf("%t", mediaSnapshotTemplate.Snapshot.IsCheckCount)
		snapshotMap["is_check_black"] = fmt.Sprintf("%t", mediaSnapshotTemplate.Snapshot.IsCheckBlack)

		if mediaSnapshotTemplate.Snapshot.BlackLevel != "" {
			snapshotMap["black_level"] = mediaSnapshotTemplate.Snapshot.BlackLevel
		}

		if mediaSnapshotTemplate.Snapshot.PixelBlackThreshold != "" {
			snapshotMap["pixel_black_threshold"] = mediaSnapshotTemplate.Snapshot.PixelBlackThreshold
		}

		if mediaSnapshotTemplate.Snapshot.SnapshotOutMode != "" {
			snapshotMap["snapshot_out_mode"] = mediaSnapshotTemplate.Snapshot.SnapshotOutMode
		}

		if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig != nil {
			spriteSnapshotConfigMap := map[string]interface{}{}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.CellWidth != "" {
				spriteSnapshotConfigMap["cell_width"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.CellWidth
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.CellHeight != "" {
				spriteSnapshotConfigMap["cell_height"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.CellHeight
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Padding != "" {
				spriteSnapshotConfigMap["padding"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Padding
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Margin != "" {
				spriteSnapshotConfigMap["margin"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Margin
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Color != "" {
				spriteSnapshotConfigMap["color"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Color
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Columns != "" {
				spriteSnapshotConfigMap["columns"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Columns
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Lines != "" {
				spriteSnapshotConfigMap["lines"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Lines
			}

			snapshotMap["sprite_snapshot_config"] = []interface{}{spriteSnapshotConfigMap}
		}

		err = d.Set("snapshot", []interface{}{snapshotMap})
		if err != nil {
			return err
		}
	}

	if mediaSnapshotTemplate.TemplateId != "" {
		_ = d.Set("template_id", mediaSnapshotTemplate.TemplateId)
	}

	if mediaSnapshotTemplate.UpdateTime != "" {
		_ = d.Set("update_time", mediaSnapshotTemplate.UpdateTime)
	}

	if mediaSnapshotTemplate.CreateTime != "" {
		_ = d.Set("create_time", mediaSnapshotTemplate.CreateTime)
	}

	return nil
}

func resourceTencentCloudCiMediaSnapshotTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_snapshot_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	request := cos.CreateMediaSnapshotTemplateOptions{
		Tag: "Snapshot",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}
	if d.HasChange("snapshot") {
		if dMap, ok := helper.InterfacesHeadMap(d, "snapshot"); ok {
			snapshot := cos.Snapshot{}
			if v, ok := dMap["mode"]; ok {
				snapshot.Mode = v.(string)
			}
			if v, ok := dMap["start"]; ok {
				snapshot.Start = v.(string)
			}
			if v, ok := dMap["time_interval"]; ok {
				snapshot.TimeInterval = v.(string)
			}
			if v, ok := dMap["count"]; ok {
				snapshot.Count = v.(string)
			}
			if v, ok := dMap["width"]; ok {
				snapshot.Width = v.(string)
			}
			if v, ok := dMap["height"]; ok {
				snapshot.Height = v.(string)
			}
			if v, ok := dMap["ci_param"]; ok {
				snapshot.CIParam = v.(string)
			}
			if v, ok := dMap["is_check_count"]; ok {
				snapshot.IsCheckCount = v.(bool)
			}
			if v, ok := dMap["is_check_black"]; ok {
				snapshot.IsCheckBlack = v.(bool)
			}
			if v, ok := dMap["black_level"]; ok {
				snapshot.BlackLevel = v.(string)
			}
			if v, ok := dMap["pixel_black_threshold"]; ok {
				snapshot.PixelBlackThreshold = v.(string)
			}
			if v, ok := dMap["snapshot_out_mode"]; ok {
				snapshot.SnapshotOutMode = v.(string)
			}
			if spriteSnapshotConfigMap, ok := helper.InterfacesHeadMap(d, "sprite_snapshot_config"); ok {
				spriteSnapshotConfig := cos.SpriteSnapshotConfig{}
				if v, ok := spriteSnapshotConfigMap["cell_width"]; ok {
					spriteSnapshotConfig.CellWidth = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["cell_height"]; ok {
					spriteSnapshotConfig.CellHeight = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["padding"]; ok {
					spriteSnapshotConfig.Padding = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["margin"]; ok {
					spriteSnapshotConfig.Margin = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["color"]; ok {
					spriteSnapshotConfig.Color = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["columns"]; ok {
					spriteSnapshotConfig.Columns = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["lines"]; ok {
					spriteSnapshotConfig.Lines = v.(string)
				}
				snapshot.SpriteSnapshotConfig = &spriteSnapshotConfig
			}
			request.Snapshot = &snapshot
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaSnapshotTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaSnapshotTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSnapshotTemplate failed, reason:%+v", logId, err)
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
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteCiMediaTemplateById(ctx, bucket, templateId); err != nil {
		return err
	}

	return nil
}
