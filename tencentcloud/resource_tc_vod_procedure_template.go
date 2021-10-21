/*
Provide a resource to create a VOD procedure template.

Example Usage

```hcl
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec               = "libx265"
      fps                 = 4
      bitrate             = 129
      resolution_adaptive = false
      width               = 128
      height              = 128
      fill_type           = "stretch"
    }
    audio {
      codec         = "libmp3lame"
      bitrate       = 129
      sample_rate   = 44100
      audio_channel = "dual"
    }
    remove_audio = false
  }
  stream_info {
    video {
      codec   = "libx264"
      fps     = 4
      bitrate = 256
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 256
      sample_rate = 44100
    }
    remove_audio = true
  }
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  width               = 130
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}

resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure"
  comment = "test"
  media_process_task {
    adaptive_dynamic_streaming_task_list {
      definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
    }
    snapshot_by_time_offset_task_list {
      definition           = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
      ext_time_offset_list = [
        "3.5s"
      ]
    }
    image_sprite_task_list {
      definition = tencentcloud_vod_image_sprite_template.foo.id
    }
  }
}
```

Import

VOD procedure template can be imported using the name, e.g.

```
$ terraform import tencentcloud_vod_procedure_template.foo tf-procedure
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudVodProcedureTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodProcedureTemplateCreate,
		Read:   resourceTencentCloudVodProcedureTemplateRead,
		Update: resourceTencentCloudVodProcedureTemplateUpdate,
		Delete: resourceTencentCloudVodProcedureTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 20),
				Description:  "Task flow name (up to 20 characters).",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 256),
				Description:  "Template description. Length limit: 256 characters.",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.",
			},
			"media_process_task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				MinItems:    1,
				Description: "Parameter of video processing task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transcode_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of transcoding tasks. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Video transcoding template ID.",
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem:        VodWatermarkResource(),
									},
									"mosaic_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "List of blurs. Up to 10 ones can be supported.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"coordinate_origin": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "TopLeft",
													Description: "Origin position, which currently can only be: `TopLeft`: the origin of coordinates is in the top-left corner of the video, and the origin of the blur is in the top-left corner of the image or text. Default value: TopLeft.",
												},
												"x_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "0px",
													Description: "The horizontal position of the origin of the blur relative to the origin of coordinates of the video. `%` and `px` formats are supported: If the string ends in `%`, the XPos of the blur will be the specified percentage of the video width; for example, 10% means that XPos is 10% of the video width; If the string ends in `px`, the XPos of the blur will be the specified px; for example, 100px means that XPos is 100 px. Default value: `0px`.",
												},
												"y_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "0px",
													Description: "Vertical position of the origin of blur relative to the origin of coordinates of video. `%` and `px` formats are supported: If the string ends in `%`, the YPos of the blur will be the specified percentage of the video height; for example, 10% means that YPos is 10% of the video height; If the string ends in `px`, the YPos of the blur will be the specified px; for example, 100px means that YPos is 100 px. Default value: `0px`.",
												},
												"width": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "10%",
													Description: "Blur width. `%` and `px` formats are supported: If the string ends in `%`, the `width` of the blur will be the specified percentage of the video width; for example, 10% means that `width` is 10% of the video width; If the string ends in `px`, the `width` of the blur will be in px; for example, 100px means that Width is 100 px. Default value: `10%`.",
												},
												"height": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "10%",
													Description: "Blur height. `%` and `px` formats are supported: If the string ends in `%`, the `height` of the blur will be the specified percentage of the video height; for example, 10% means that Height is 10% of the video height; If the string ends in `px`, the `height` of the blur will be in px; for example, 100px means that Height is 100 px. Default value: `10%`.",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Start time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame. If this parameter is left empty or `0` is entered, the blur will appear upon the first video frame; If this value is greater than `0` (e.g., n), the blur will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the blur will appear at second n before the last video frame.",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "End time offset of blur in seconds. If this parameter is left empty or `0` is entered, the blur will exist till the last video frame; If this value is greater than `0` (e.g., n), the blur will exist till second n; If this value is smaller than `0` (e.g., -n), the blur will exist till second n before the last video frame.",
												},
											},
										},
									},
								},
							},
						},
						"animated_graphic_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of animated image generating tasks. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Animated image generating template ID.",
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Start time of animated image in video in seconds.",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "End time of animated image in video in seconds.",
									},
								},
							},
						},
						"snapshot_by_time_offset_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of time point screen capturing tasks. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Time point screen capturing template ID.",
									},
									"ext_time_offset_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of screenshot time points. `s` and `%` formats are supported: When a time point string ends with `s`, its unit is second. For example, `3.5s` means the 3.5th second of the video; When a time point string ends with `%`, it is marked with corresponding percentage of the video duration. For example, `10%` means that the time point is at the 10% of the video entire duration.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem:        VodWatermarkResource(),
									},
								},
							},
						},
						"sample_snapshot_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of sampled screen capturing tasks. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Sampled screen capturing template ID.",
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem:        VodWatermarkResource(),
									},
								},
							},
						},
						"image_sprite_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of image sprite generating tasks. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Image sprite generating template ID.",
									},
								},
							},
						},
						"cover_by_snapshot_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of cover generating tasks. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Time point screen capturing template ID.",
									},
									"position_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateAllowedStringValue([]string{"Time", "Percent"}),
										Description:  "Screen capturing mode. Valid values: `Time`, `Percent`. `Time`: screen captures by time point, `Percent`: screen captures by percentage.",
									},
									"position_value": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Screenshot position: For time point screen capturing, this means to take a screenshot at a specified time point (in seconds) and use it as the cover. For percentage screen capturing, this value means to take a screenshot at a specified percentage of the video duration and use it as the cover.",
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem:        VodWatermarkResource(),
									},
								},
							},
						},
						"adaptive_dynamic_streaming_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of adaptive bitrate streaming tasks. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Adaptive bitrate streaming template ID.",
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "List of up to `10` image or text watermarks. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem:        VodWatermarkResource(),
									},
								},
							},
						},
					},
				},
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
	}
}

func genWatermarkList(item map[string]interface{}) (list []*vod.WatermarkInput) {
	if _, ok := item["watermark_list"]; !ok {
		return nil
	}
	waterV := item["watermark_list"].([]interface{})
	list = make([]*vod.WatermarkInput, 0, len(waterV))
	for _, water := range waterV {
		waterVV := water.(map[string]interface{})
		list = append(list, &vod.WatermarkInput{
			Definition: func(str string) *uint64 {
				idUint, _ := strconv.ParseUint(str, 0, 64)
				return &idUint
			}(waterVV["definition"].(string)),
			TextContent: func() *string {
				if content, ok := waterVV["text_content"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			SvgContent: func() *string {
				if content, ok := waterVV["svg_content"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			StartTimeOffset: func() *float64 {
				if content, ok := waterVV["start_time_offset"]; !ok {
					return nil
				} else {
					return helper.Float64(content.(float64))
				}
			}(),
			EndTimeOffset: func() *float64 {
				if content, ok := waterVV["end_time_offset"]; !ok {
					return nil
				} else {
					return helper.Float64(content.(float64))
				}
			}(),
		})
	}
	return list
}

func genMosaicList(item map[string]interface{}) (list []*vod.MosaicInput) {
	if _, ok := item["mosaic_list"]; !ok {
		return nil
	}
	mosaicV := item["mosaic_list"].([]interface{})
	list = make([]*vod.MosaicInput, 0, len(mosaicV))
	for _, mosaic := range mosaicV {
		mosaicVV := mosaic.(map[string]interface{})
		list = append(list, &vod.MosaicInput{
			CoordinateOrigin: func() *string {
				if content, ok := mosaicVV["coordinate_origin"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			XPos: func() *string {
				if content, ok := mosaicVV["x_pos"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			YPos: func() *string {
				if content, ok := mosaicVV["y_pos"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			Width: func() *string {
				if content, ok := mosaicVV["width"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			Height: func() *string {
				if content, ok := mosaicVV["height"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			StartTimeOffset: func() *float64 {
				if content, ok := mosaicVV["start_time_offset"]; !ok {
					return nil
				} else {
					return helper.Float64(content.(float64))
				}
			}(),
			EndTimeOffset: func() *float64 {
				if content, ok := mosaicVV["end_time_offset"]; !ok {
					return nil
				} else {
					return helper.Float64(content.(float64))
				}
			}(),
		})
	}
	return list
}

func generateMediaProcessTask(d *schema.ResourceData) (mediaReq *vod.MediaProcessTaskInput) {
	mediaReq = &vod.MediaProcessTaskInput{}
	mediaProcessTask := d.Get("media_process_task").([]interface{})[0].(map[string]interface{})
	// transcode_task_list
	if transcodeTask, ok := mediaProcessTask["transcode_task_list"]; ok {
		transcodeTaskList := transcodeTask.([]interface{})
		transcodeReq := make([]*vod.TranscodeTaskInput, 0, len(transcodeTaskList))
		for _, itemV := range transcodeTaskList {
			item := itemV.(map[string]interface{})
			transcodeReq = append(transcodeReq, &vod.TranscodeTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet: genWatermarkList(item),
				MosaicSet:    genMosaicList(item),
			})
		}
		mediaReq.TranscodeTaskSet = transcodeReq
	}
	// animated_graphic_task_list
	if animateTask, ok := mediaProcessTask["animated_graphic_task_list"]; ok {
		animateTaskList := animateTask.([]interface{})
		animateReq := make([]*vod.AnimatedGraphicTaskInput, 0, len(animateTaskList))
		for _, itemV := range animateTaskList {
			item := itemV.(map[string]interface{})
			animateReq = append(animateReq, &vod.AnimatedGraphicTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				StartTimeOffset: helper.Float64(item["start_time_offset"].(float64)),
				EndTimeOffset:   helper.Float64(item["end_time_offset"].(float64)),
			})
		}
		mediaReq.AnimatedGraphicTaskSet = animateReq
	}
	// snapshot_by_time_offset_task_list
	if snapshotTask, ok := mediaProcessTask["snapshot_by_time_offset_task_list"]; ok {
		snapshotTaskList := snapshotTask.([]interface{})
		snapshotReq := make([]*vod.SnapshotByTimeOffsetTaskInput, 0, len(snapshotTaskList))
		for _, itemV := range snapshotTaskList {
			item := itemV.(map[string]interface{})
			snapshotReq = append(snapshotReq, &vod.SnapshotByTimeOffsetTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet: genWatermarkList(item),
				ExtTimeOffsetSet: func() (list []*string) {
					if _, ok := item["ext_time_offset_list"]; !ok {
						return nil
					}
					extTimeV := item["ext_time_offset_list"].([]interface{})
					list = make([]*string, 0, len(extTimeV))
					for _, extTimeVV := range extTimeV {
						list = append(list, helper.String(extTimeVV.(string)))
					}
					return list
				}(),
			})
		}
		mediaReq.SnapshotByTimeOffsetTaskSet = snapshotReq
	}
	// sample_snapshot_task_list
	if sampleTask, ok := mediaProcessTask["sample_snapshot_task_list"]; ok {
		sampleTaskList := sampleTask.([]interface{})
		sampleReq := make([]*vod.SampleSnapshotTaskInput, 0, len(sampleTaskList))
		for _, itemV := range sampleTaskList {
			item := itemV.(map[string]interface{})
			sampleReq = append(sampleReq, &vod.SampleSnapshotTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet: genWatermarkList(item),
			})
		}
		mediaReq.SampleSnapshotTaskSet = sampleReq
	}
	// image_sprite_task_list
	if imageTask, ok := mediaProcessTask["image_sprite_task_list"]; ok {
		imageTaskList := imageTask.([]interface{})
		imageReq := make([]*vod.ImageSpriteTaskInput, 0, len(imageTaskList))
		for _, itemV := range imageTaskList {
			item := itemV.(map[string]interface{})
			imageReq = append(imageReq, &vod.ImageSpriteTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
			})
		}
		mediaReq.ImageSpriteTaskSet = imageReq
	}
	// cover_by_snapshot_task_list
	if coverTask, ok := mediaProcessTask["cover_by_snapshot_task_list"]; ok {
		coverTaskList := coverTask.([]interface{})
		coverReq := make([]*vod.CoverBySnapshotTaskInput, 0, len(coverTaskList))
		for _, itemV := range coverTaskList {
			item := itemV.(map[string]interface{})
			coverReq = append(coverReq, &vod.CoverBySnapshotTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet:  genWatermarkList(item),
				PositionType:  helper.String(item["iposition_type"].(string)),
				PositionValue: helper.Float64(item["position_value"].(float64)),
			})
		}
		mediaReq.CoverBySnapshotTaskSet = coverReq
	}
	// adaptive_dynamic_streaming_task_list
	if adaptiveTask, ok := mediaProcessTask["adaptive_dynamic_streaming_task_list"]; ok {
		adaptiveTaskList := adaptiveTask.([]interface{})
		adaptiveReq := make([]*vod.AdaptiveDynamicStreamingTaskInput, 0, len(adaptiveTaskList))
		for _, itemV := range adaptiveTaskList {
			item := itemV.(map[string]interface{})
			adaptiveReq = append(adaptiveReq, &vod.AdaptiveDynamicStreamingTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet: genWatermarkList(item),
			})
		}
		mediaReq.AdaptiveDynamicStreamingTaskSet = adaptiveReq
	}

	return mediaReq
}

func resourceTencentCloudVodProcedureTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_procedure_template.create")()

	var (
		logId   = getLogId(contextNil)
		request = vod.NewCreateProcedureTemplateRequest()
	)

	request.Name = helper.String(d.Get("name").(string))
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sub_app_id"); ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}
	if _, ok := d.GetOk("media_process_task"); ok {
		mediaReq := generateMediaProcessTask(d)
		request.MediaProcessTask = mediaReq
	}

	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().CreateProcedureTemplate(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(d.Get("name").(string))

	return resourceTencentCloudVodProcedureTemplateRead(d, meta)
}

func resourceTencentCloudVodProcedureTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_procedure_template.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		id         = d.Id()
		subAppId   = d.Get("sub_app_id").(int)
		client     = meta.(*TencentCloudClient).apiV3Conn
		vodService = VodService{client: client}
	)
	template, has, err := vodService.DescribeProcedureTemplatesById(ctx, id, subAppId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", template.Name)
	_ = d.Set("comment", template.Comment)
	_ = d.Set("create_time", template.CreateTime)
	_ = d.Set("update_time", template.UpdateTime)

	mediaProcessTaskElem := make(map[string]interface{})
	if template.MediaProcessTask != nil {
		// transcode_task_list
		if template.MediaProcessTask.TranscodeTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.TranscodeTaskSet))
			for _, item := range template.MediaProcessTask.TranscodeTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
					"mosaic_list": func() interface{} {
						if item.MosaicSet == nil {
							return nil
						}
						mosaicList := make([]map[string]interface{}, 0, len(item.MosaicSet))
						for _, mosaicV := range item.MosaicSet {
							mosaicList = append(mosaicList, map[string]interface{}{
								"coordinate_origin": mosaicV.CoordinateOrigin,
								"x_pos":             mosaicV.XPos,
								"y_pos":             mosaicV.YPos,
								"width":             mosaicV.Width,
								"height":            mosaicV.Height,
								"start_time_offset": mosaicV.StartTimeOffset,
								"end_time_offset":   mosaicV.EndTimeOffset,
							})
						}
						return mosaicList
					}(),
				})
			}
			mediaProcessTaskElem["transcode_task_list"] = list
		}
		// animated_graphic_task_list
		if template.MediaProcessTask.AnimatedGraphicTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.AnimatedGraphicTaskSet))
			for _, item := range template.MediaProcessTask.AnimatedGraphicTaskSet {
				list = append(list, map[string]interface{}{
					"definition":        strconv.FormatUint(*item.Definition, 10),
					"start_time_offset": item.StartTimeOffset,
					"end_time_offset":   item.EndTimeOffset,
				})
			}
			mediaProcessTaskElem["animated_graphic_task_list"] = list
		}
		// snapshot_by_time_offset_task_list
		if template.MediaProcessTask.SnapshotByTimeOffsetTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.SnapshotByTimeOffsetTaskSet))
			for _, item := range template.MediaProcessTask.SnapshotByTimeOffsetTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
					"ext_time_offset_list": func() interface{} {
						if item.ExtTimeOffsetSet == nil {
							return nil
						}
						extList := make([]interface{}, 0, len(item.ExtTimeOffsetSet))
						for _, extV := range item.ExtTimeOffsetSet {
							extList = append(extList, extV)
						}
						return extList
					}(),
				})
			}
			mediaProcessTaskElem["snapshot_by_time_offset_task_list"] = list
		}
		// sample_snapshot_task_list
		if template.MediaProcessTask.SampleSnapshotTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.SampleSnapshotTaskSet))
			for _, item := range template.MediaProcessTask.SampleSnapshotTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
				})
			}
			mediaProcessTaskElem["sample_snapshot_task_list"] = list
		}
		// image_sprite_task_list
		if template.MediaProcessTask.ImageSpriteTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.ImageSpriteTaskSet))
			for _, item := range template.MediaProcessTask.ImageSpriteTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
				})
			}
			mediaProcessTaskElem["image_sprite_task_list"] = list
		}
		// cover_by_snapshot_task_list
		if template.MediaProcessTask.CoverBySnapshotTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.CoverBySnapshotTaskSet))
			for _, item := range template.MediaProcessTask.CoverBySnapshotTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
					"position_type":  item.PositionType,
					"position_value": item.PositionValue,
				})
			}
			mediaProcessTaskElem["cover_by_snapshot_task_list"] = list
		}
		// adaptive_dynamic_streaming_task_list
		if template.MediaProcessTask.AdaptiveDynamicStreamingTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.AdaptiveDynamicStreamingTaskSet))
			for _, item := range template.MediaProcessTask.AdaptiveDynamicStreamingTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
				})
			}
			mediaProcessTaskElem["adaptive_dynamic_streaming_task_list"] = list
		}

		_ = d.Set("media_process_task", []interface{}{mediaProcessTaskElem})
	}

	return nil
}

func resourceTencentCloudVodProcedureTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_procedure_template.update")()

	var (
		logId      = getLogId(contextNil)
		request    = vod.NewResetProcedureTemplateRequest()
		id         = d.Id()
		changeFlag = false
	)

	request.Name = &id
	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("sub_app_id") {
		changeFlag = true
		request.SubAppId = helper.IntUint64(d.Get("sub_app_id").(int))
	}
	if d.HasChange("media_process_task") {
		changeFlag = true
		mediaReq := generateMediaProcessTask(d)
		request.MediaProcessTask = mediaReq
	}

	if changeFlag {
		var err error
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ResetProcedureTemplate(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodProcedureTemplateRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodProcedureTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_procedure_template.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	vodService := VodService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if err := vodService.DeleteProcedureTemplate(ctx, id, uint64(d.Get("sub_app_id").(int))); err != nil {
		return err
	}

	return nil
}
