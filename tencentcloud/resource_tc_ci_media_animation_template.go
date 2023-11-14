/*
Provides a resource to create a ci media_animation_template

Example Usage

```hcl
resource "tencentcloud_ci_media_animation_template" "media_animation_template" {
  name = &lt;nil&gt;
  container {
		format = &lt;nil&gt;

  }
  video {
		codec = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		fps = &lt;nil&gt;
		animate_only_keep_key_frame = &lt;nil&gt;
		animate_time_interval_of_frame = &lt;nil&gt;
		animate_frames_per_second = &lt;nil&gt;
		quality = &lt;nil&gt;

  }
  time_interval {
		start = &lt;nil&gt;
		duration = &lt;nil&gt;

  }
}
```

Import

ci media_animation_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_animation_template.media_animation_template media_animation_template_id
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

func resourceTencentCloudCiMediaAnimationTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaAnimationTemplateCreate,
		Read:   resourceTencentCloudCiMediaAnimationTemplateRead,
		Update: resourceTencentCloudCiMediaAnimationTemplateUpdate,
		Delete: resourceTencentCloudCiMediaAnimationTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"container": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Container format.",
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
				Description: "Video information, do not upload Video, which is equivalent to deleting video information.",
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
							Description: "Width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.",
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
							Description: "GIFs are kept only Keyframe, Priority: AnimateFramesPerSecond &amp;gt; AnimateOnlyKeepKeyFrame &amp;gt; AnimateTimeIntervalOfFrame.",
						},
						"animate_time_interval_of_frame": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Animation frame extraction every time, (0, video duration], Animation frame extraction time interval, If TimeInterval.Duration is set, it is less than this value.",
						},
						"animate_frames_per_second": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Animation per second frame number, Priority: AnimateFramesPerSecond &amp;gt; AnimateOnlyKeepKeyFrame &amp;gt; AnimateTimeIntervalOfFrame.",
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
				Description: "Time interval.",
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
							Description: "Duration, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
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

	var (
		request    = ci.NewCreateMediaAnimationTemplateRequest()
		response   = ci.NewCreateMediaAnimationTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := ci.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = helper.String(v.(string))
		}
		request.Container = &container
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
		video := ci.Video{}
		if v, ok := dMap["codec"]; ok {
			video.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["width"]; ok {
			video.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			video.Height = helper.String(v.(string))
		}
		if v, ok := dMap["fps"]; ok {
			video.Fps = helper.String(v.(string))
		}
		if v, ok := dMap["animate_only_keep_key_frame"]; ok {
			video.AnimateOnlyKeepKeyFrame = helper.String(v.(string))
		}
		if v, ok := dMap["animate_time_interval_of_frame"]; ok {
			video.AnimateTimeIntervalOfFrame = helper.String(v.(string))
		}
		if v, ok := dMap["animate_frames_per_second"]; ok {
			video.AnimateFramesPerSecond = helper.String(v.(string))
		}
		if v, ok := dMap["quality"]; ok {
			video.Quality = helper.String(v.(string))
		}
		request.Video = &video
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "time_interval"); ok {
		timeInterval := ci.TimeInterval{}
		if v, ok := dMap["start"]; ok {
			timeInterval.Start = helper.String(v.(string))
		}
		if v, ok := dMap["duration"]; ok {
			timeInterval.Duration = helper.String(v.(string))
		}
		request.TimeInterval = &timeInterval
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaAnimationTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaAnimationTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaAnimationTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaAnimationTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_animation_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaAnimationTemplateId := d.Id()

	mediaAnimationTemplate, err := service.DescribeCiMediaAnimationTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaAnimationTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaAnimationTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaAnimationTemplate.Name != nil {
		_ = d.Set("name", mediaAnimationTemplate.Name)
	}

	if mediaAnimationTemplate.Container != nil {
		containerMap := map[string]interface{}{}

		if mediaAnimationTemplate.Container.Format != nil {
			containerMap["format"] = mediaAnimationTemplate.Container.Format
		}

		_ = d.Set("container", []interface{}{containerMap})
	}

	if mediaAnimationTemplate.Video != nil {
		videoMap := map[string]interface{}{}

		if mediaAnimationTemplate.Video.Codec != nil {
			videoMap["codec"] = mediaAnimationTemplate.Video.Codec
		}

		if mediaAnimationTemplate.Video.Width != nil {
			videoMap["width"] = mediaAnimationTemplate.Video.Width
		}

		if mediaAnimationTemplate.Video.Height != nil {
			videoMap["height"] = mediaAnimationTemplate.Video.Height
		}

		if mediaAnimationTemplate.Video.Fps != nil {
			videoMap["fps"] = mediaAnimationTemplate.Video.Fps
		}

		if mediaAnimationTemplate.Video.AnimateOnlyKeepKeyFrame != nil {
			videoMap["animate_only_keep_key_frame"] = mediaAnimationTemplate.Video.AnimateOnlyKeepKeyFrame
		}

		if mediaAnimationTemplate.Video.AnimateTimeIntervalOfFrame != nil {
			videoMap["animate_time_interval_of_frame"] = mediaAnimationTemplate.Video.AnimateTimeIntervalOfFrame
		}

		if mediaAnimationTemplate.Video.AnimateFramesPerSecond != nil {
			videoMap["animate_frames_per_second"] = mediaAnimationTemplate.Video.AnimateFramesPerSecond
		}

		if mediaAnimationTemplate.Video.Quality != nil {
			videoMap["quality"] = mediaAnimationTemplate.Video.Quality
		}

		_ = d.Set("video", []interface{}{videoMap})
	}

	if mediaAnimationTemplate.TimeInterval != nil {
		timeIntervalMap := map[string]interface{}{}

		if mediaAnimationTemplate.TimeInterval.Start != nil {
			timeIntervalMap["start"] = mediaAnimationTemplate.TimeInterval.Start
		}

		if mediaAnimationTemplate.TimeInterval.Duration != nil {
			timeIntervalMap["duration"] = mediaAnimationTemplate.TimeInterval.Duration
		}

		_ = d.Set("time_interval", []interface{}{timeIntervalMap})
	}

	return nil
}

func resourceTencentCloudCiMediaAnimationTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_animation_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaAnimationTemplateRequest()

	mediaAnimationTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "container", "video", "time_interval"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaAnimationTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaAnimationTemplate failed, reason:%+v", logId, err)
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
	mediaAnimationTemplateId := d.Id()

	if err := service.DeleteCiMediaAnimationTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
