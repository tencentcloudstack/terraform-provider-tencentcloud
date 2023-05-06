/*
Provides a resource to create a ci media_animation_template

Example Usage

```hcl
resource "tencentcloud_ci_media_animation_template" "media_animation_template" {
  bucket = "terraform-ci-1308919341"
  name = "animation_template-002"
  container {
		format = "gif"
  }
  video {
		codec = "gif"
		width = "1280"
		height = ""
		fps = "20"
		animate_only_keep_key_frame = "true"
		animate_time_interval_of_frame = ""
		animate_frames_per_second = ""
		quality = ""

  }
  time_interval {
		start = "0"
		duration = "60"

  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCiMediaAnimationTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaAnimationTemplateCreate,
		Read:   resourceTencentCloudCiMediaAnimationTemplateRead,
		Update: resourceTencentCloudCiMediaAnimationTemplateUpdate,
		Delete: resourceTencentCloudCiMediaAnimationTemplateDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
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

			"container": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "container format.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Package format.",
						},
					},
				},
			},

			"video": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "video information, do not upload Video, which is equivalent to deleting video information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Codec format `gif`, `webp`.",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "High, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video, must be even.",
						},
						"fps": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frame rate, value range: (0, 60], Unit: fps.",
						},
						"animate_only_keep_key_frame": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GIFs are kept only Keyframe, Priority: AnimateFramesPerSecond &gt; AnimateOnlyKeepKeyFrame &gt; AnimateTimeIntervalOfFrame.",
						},
						"animate_time_interval_of_frame": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Animation frame extraction every time, (0, video duration], Animation frame extraction time interval, If TimeInterval.Duration is set, it is less than this value.",
						},
						"animate_frames_per_second": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Animation per second frame number, Priority: AnimateFramesPerSecond &gt; AnimateOnlyKeepKeyFrame &gt; AnimateTimeIntervalOfFrame.",
						},
						"quality": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Set relative quality, [1, 100), webp image quality setting takes effect, gif has no quality parameter.",
						},
					},
				},
			},

			"time_interval": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "time interval.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Starting time, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
						},
						"duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "duration, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaAnimationTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_animation_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaAnimationTemplateOptions{
			Tag: "Animation",
		}
		templateId string
		bucket     string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		request.Container = &container
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
		video := cos.AnimationVideo{}
		if v, ok := dMap["codec"]; ok {
			video.Codec = v.(string)
		}
		if v, ok := dMap["width"]; ok {
			video.Width = v.(string)
		}
		if v, ok := dMap["height"]; ok {
			video.Height = v.(string)
		}
		if v, ok := dMap["fps"]; ok {
			video.Fps = v.(string)
		}
		if v, ok := dMap["animate_only_keep_key_frame"]; ok {
			video.AnimateOnlyKeepKeyFrame = v.(string)
		}
		if v, ok := dMap["animate_time_interval_of_frame"]; ok {
			video.AnimateTimeIntervalOfFrame = v.(string)
		}
		if v, ok := dMap["animate_frames_per_second"]; ok {
			video.AnimateFramesPerSecond = v.(string)
		}
		if v, ok := dMap["quality"]; ok {
			video.Quality = v.(string)
		}
		request.Video = &video
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "time_interval"); ok {
		timeInterval := cos.TimeInterval{}
		if v, ok := dMap["start"]; ok {
			timeInterval.Start = v.(string)
		}
		if v, ok := dMap["duration"]; ok {
			timeInterval.Duration = v.(string)
		}
		request.TimeInterval = &timeInterval
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaAnimationTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%+v], response body [%+v]\n", logId, "CreateMediaAnimationTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaAnimationTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaAnimationTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaAnimationTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_animation_template.read")()
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

	template, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if template.Name != "" {
		_ = d.Set("name", template.Name)
	}

	mediaAnimationTemplate := template.TransTpl
	if mediaAnimationTemplate != nil {
		if mediaAnimationTemplate.Container != nil {
			containerMap := map[string]interface{}{}

			if mediaAnimationTemplate.Container.Format != "" {
				containerMap["format"] = mediaAnimationTemplate.Container.Format
			}

			_ = d.Set("container", []interface{}{containerMap})
		}

		if mediaAnimationTemplate.Video != nil {
			videoMap := map[string]interface{}{}

			if mediaAnimationTemplate.Video.Codec != "" {
				videoMap["codec"] = mediaAnimationTemplate.Video.Codec
			}

			if mediaAnimationTemplate.Video.Width != "" {
				videoMap["width"] = mediaAnimationTemplate.Video.Width
			}

			if mediaAnimationTemplate.Video.Height != "" {
				videoMap["height"] = mediaAnimationTemplate.Video.Height
			}

			if mediaAnimationTemplate.Video.Fps != "" {
				videoMap["fps"] = mediaAnimationTemplate.Video.Fps
			}

			// if mediaAnimationTemplate.Video.AnimateOnlyKeepKeyFrame != "" {
			// 	videoMap["animate_only_keep_key_frame"] = mediaAnimationTemplate.Video.AnimateOnlyKeepKeyFrame
			// }

			// if mediaAnimationTemplate.Video.AnimateTimeIntervalOfFrame != "" {
			// 	videoMap["animate_time_interval_of_frame"] = mediaAnimationTemplate.Video.AnimateTimeIntervalOfFrame
			// }

			// if mediaAnimationTemplate.Video.AnimateFramesPerSecond != "" {
			// 	videoMap["animate_frames_per_second"] = mediaAnimationTemplate.Video.AnimateFramesPerSecond
			// }

			// if mediaAnimationTemplate.Video.Quality != "" {
			// 	videoMap["quality"] = mediaAnimationTemplate.Video.Quality
			// }

			err = d.Set("video", []interface{}{videoMap})
			if err != nil {
				return err
			}
		}

		if mediaAnimationTemplate.TimeInterval != nil {
			timeIntervalMap := map[string]interface{}{}

			if mediaAnimationTemplate.TimeInterval.Start != "" {
				timeIntervalMap["start"] = mediaAnimationTemplate.TimeInterval.Start
			}

			if mediaAnimationTemplate.TimeInterval.Duration != "" {
				timeIntervalMap["duration"] = mediaAnimationTemplate.TimeInterval.Duration
			}

			err = d.Set("time_interval", []interface{}{timeIntervalMap})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceTencentCloudCiMediaAnimationTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_animation_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaAnimationTemplateOptions{
		Tag: "Animation",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		request.Container = &container
	}

	if d.HasChange("video") {
		if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
			video := cos.AnimationVideo{}
			if v, ok := dMap["codec"]; ok {
				video.Codec = v.(string)
			}
			if v, ok := dMap["width"]; ok {
				video.Width = v.(string)
			}
			if v, ok := dMap["height"]; ok {
				video.Height = v.(string)
			}
			if v, ok := dMap["fps"]; ok {
				video.Fps = v.(string)
			}
			if v, ok := dMap["animate_only_keep_key_frame"]; ok {
				video.AnimateOnlyKeepKeyFrame = v.(string)
			}
			if v, ok := dMap["animate_time_interval_of_frame"]; ok {
				video.AnimateTimeIntervalOfFrame = v.(string)
			}
			if v, ok := dMap["animate_frames_per_second"]; ok {
				video.AnimateFramesPerSecond = v.(string)
			}
			if v, ok := dMap["quality"]; ok {
				video.Quality = v.(string)
			}
			request.Video = &video
		}
	}

	if d.HasChange("time_interval") {
		if dMap, ok := helper.InterfacesHeadMap(d, "time_interval"); ok {
			timeInterval := cos.TimeInterval{}
			if v, ok := dMap["start"]; ok {
				timeInterval.Start = v.(string)
			}
			if v, ok := dMap["duration"]; ok {
				timeInterval.Duration = v.(string)
			}
			request.TimeInterval = &timeInterval
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaAnimationTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaAnimationTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaAnimationTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaAnimationTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaAnimationTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_animation_template.delete")()
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
