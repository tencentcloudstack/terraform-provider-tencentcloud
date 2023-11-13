/*
Provides a resource to create a css live_transcode_template

Example Usage

```hcl
resource "tencentcloud_css_live_transcode_template" "live_transcode_template" {
  template_name = &lt;nil&gt;
  video_bitrate = &lt;nil&gt;
  acodec = &lt;nil&gt;
  audio_bitrate = &lt;nil&gt;
  vcodec = &lt;nil&gt;
  description = &lt;nil&gt;
  need_video = &lt;nil&gt;
  width = &lt;nil&gt;
  need_audio = &lt;nil&gt;
  height = &lt;nil&gt;
  fps = &lt;nil&gt;
  gop = &lt;nil&gt;
  rotate = &lt;nil&gt;
  profile = &lt;nil&gt;
  bitrate_to_orig = &lt;nil&gt;
  height_to_orig = &lt;nil&gt;
  fps_to_orig = &lt;nil&gt;
  ai_trans_code = &lt;nil&gt;
  adapt_bitrate_percent = &lt;nil&gt;
  short_edge_as_height = &lt;nil&gt;
  d_r_m_type = &lt;nil&gt;
  d_r_m_tracks = &lt;nil&gt;
}
```

Import

css live_transcode_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_live_transcode_template.live_transcode_template live_transcode_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCssLiveTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssLiveTranscodeTemplateCreate,
		Read:   resourceTencentCloudCssLiveTranscodeTemplateRead,
		Update: resourceTencentCloudCssLiveTranscodeTemplateUpdate,
		Delete: resourceTencentCloudCssLiveTranscodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name, only support 0-9 and a-z.",
			},

			"video_bitrate": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Video bitrate, 0 for origin, range 0kbps - 8000kbps.",
			},

			"acodec": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Default acc, not support now.",
			},

			"audio_bitrate": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Default 0, range 0 - 500.",
			},

			"vcodec": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Video codec, default origin, support h264/h265/origin.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template desc.",
			},

			"need_video": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Keep video or not, default 1 for yes, 0 for no.",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Template width, default 0, range 0 - 3000, must be pow of 2.",
			},

			"need_audio": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Keep audio or not, default 1 for yes, 0 for no.",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Template height, default 0, range 0 - 3000, must be pow of 2, needed while AiTransCode = 1.",
			},

			"fps": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Video fps, default 0, range 0 - 60.",
			},

			"gop": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Gop of the video, second, default origin of the video, range 2 - 6.",
			},

			"rotate": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Roate degree, default 0, support 0/90/180/270.",
			},

			"profile": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Quality of the video, default baseline, support baseline/main/high.",
			},

			"bitrate_to_orig": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Base on origin bitrate if origin bitrate is lower than the setting bitrate. default 0, 1 for yes, 0 for no.",
			},

			"height_to_orig": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Base on origin height if origin height is lower than the setting height. default 0, 1 for yes, 0 for no.",
			},

			"fps_to_orig": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Base on origin fps if origin fps is lower than the setting fps. default 0, 1 for yes, 0 for no.",
			},

			"ai_trans_code": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Enable high speed mode, default 0, 1 for enable, 0 for no.",
			},

			"adapt_bitrate_percent": {
				Optional:    true,
				Type:        schema.TypeFloat,
				Description: "High speed mode adapt bitrate, support 0 - 0.5.",
			},

			"short_edge_as_height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Let the short edge as the height.",
			},

			"d_r_m_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DRM type, support fairplay/normalaes/widevine.",
			},

			"d_r_m_tracks": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DRM tracks, support AUDIO/SD/HD/UHD1/UHD2.",
			},
		},
	}
}

func resourceTencentCloudCssLiveTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewCreateLiveTranscodeTemplateRequest()
		response   = css.NewCreateLiveTranscodeTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("video_bitrate"); ok {
		request.VideoBitrate = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("acodec"); ok {
		request.Acodec = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("audio_bitrate"); ok {
		request.AudioBitrate = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vcodec"); ok {
		request.Vcodec = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("need_video"); ok {
		request.NeedVideo = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("need_audio"); ok {
		request.NeedAudio = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("height"); ok {
		request.Height = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("fps"); ok {
		request.Fps = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("gop"); ok {
		request.Gop = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("rotate"); ok {
		request.Rotate = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("profile"); ok {
		request.Profile = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("bitrate_to_orig"); ok {
		request.BitrateToOrig = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("height_to_orig"); ok {
		request.HeightToOrig = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("fps_to_orig"); ok {
		request.FpsToOrig = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("ai_trans_code"); ok {
		request.AiTransCode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("adapt_bitrate_percent"); ok {
		request.AdaptBitratePercent = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOkExists("short_edge_as_height"); ok {
		request.ShortEdgeAsHeight = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("d_r_m_type"); ok {
		request.DRMType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_r_m_tracks"); ok {
		request.DRMTracks = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveTranscodeTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css liveTranscodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudCssLiveTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudCssLiveTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	liveTranscodeTemplateId := d.Id()

	liveTranscodeTemplate, err := service.DescribeCssLiveTranscodeTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if liveTranscodeTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssLiveTranscodeTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if liveTranscodeTemplate.TemplateName != nil {
		_ = d.Set("template_name", liveTranscodeTemplate.TemplateName)
	}

	if liveTranscodeTemplate.VideoBitrate != nil {
		_ = d.Set("video_bitrate", liveTranscodeTemplate.VideoBitrate)
	}

	if liveTranscodeTemplate.Acodec != nil {
		_ = d.Set("acodec", liveTranscodeTemplate.Acodec)
	}

	if liveTranscodeTemplate.AudioBitrate != nil {
		_ = d.Set("audio_bitrate", liveTranscodeTemplate.AudioBitrate)
	}

	if liveTranscodeTemplate.Vcodec != nil {
		_ = d.Set("vcodec", liveTranscodeTemplate.Vcodec)
	}

	if liveTranscodeTemplate.Description != nil {
		_ = d.Set("description", liveTranscodeTemplate.Description)
	}

	if liveTranscodeTemplate.NeedVideo != nil {
		_ = d.Set("need_video", liveTranscodeTemplate.NeedVideo)
	}

	if liveTranscodeTemplate.Width != nil {
		_ = d.Set("width", liveTranscodeTemplate.Width)
	}

	if liveTranscodeTemplate.NeedAudio != nil {
		_ = d.Set("need_audio", liveTranscodeTemplate.NeedAudio)
	}

	if liveTranscodeTemplate.Height != nil {
		_ = d.Set("height", liveTranscodeTemplate.Height)
	}

	if liveTranscodeTemplate.Fps != nil {
		_ = d.Set("fps", liveTranscodeTemplate.Fps)
	}

	if liveTranscodeTemplate.Gop != nil {
		_ = d.Set("gop", liveTranscodeTemplate.Gop)
	}

	if liveTranscodeTemplate.Rotate != nil {
		_ = d.Set("rotate", liveTranscodeTemplate.Rotate)
	}

	if liveTranscodeTemplate.Profile != nil {
		_ = d.Set("profile", liveTranscodeTemplate.Profile)
	}

	if liveTranscodeTemplate.BitrateToOrig != nil {
		_ = d.Set("bitrate_to_orig", liveTranscodeTemplate.BitrateToOrig)
	}

	if liveTranscodeTemplate.HeightToOrig != nil {
		_ = d.Set("height_to_orig", liveTranscodeTemplate.HeightToOrig)
	}

	if liveTranscodeTemplate.FpsToOrig != nil {
		_ = d.Set("fps_to_orig", liveTranscodeTemplate.FpsToOrig)
	}

	if liveTranscodeTemplate.AiTransCode != nil {
		_ = d.Set("ai_trans_code", liveTranscodeTemplate.AiTransCode)
	}

	if liveTranscodeTemplate.AdaptBitratePercent != nil {
		_ = d.Set("adapt_bitrate_percent", liveTranscodeTemplate.AdaptBitratePercent)
	}

	if liveTranscodeTemplate.ShortEdgeAsHeight != nil {
		_ = d.Set("short_edge_as_height", liveTranscodeTemplate.ShortEdgeAsHeight)
	}

	if liveTranscodeTemplate.DRMType != nil {
		_ = d.Set("d_r_m_type", liveTranscodeTemplate.DRMType)
	}

	if liveTranscodeTemplate.DRMTracks != nil {
		_ = d.Set("d_r_m_tracks", liveTranscodeTemplate.DRMTracks)
	}

	return nil
}

func resourceTencentCloudCssLiveTranscodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewModifyLiveTranscodeTemplateRequest()

	liveTranscodeTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "video_bitrate", "acodec", "audio_bitrate", "vcodec", "description", "need_video", "width", "need_audio", "height", "fps", "gop", "rotate", "profile", "bitrate_to_orig", "height_to_orig", "fps_to_orig", "ai_trans_code", "adapt_bitrate_percent", "short_edge_as_height", "d_r_m_type", "d_r_m_tracks"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("video_bitrate") {
		if v, ok := d.GetOkExists("video_bitrate"); ok {
			request.VideoBitrate = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("acodec") {
		if v, ok := d.GetOk("acodec"); ok {
			request.Acodec = helper.String(v.(string))
		}
	}

	if d.HasChange("vcodec") {
		if v, ok := d.GetOk("vcodec"); ok {
			request.Vcodec = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("need_video") {
		if v, ok := d.GetOkExists("need_video"); ok {
			request.NeedVideo = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("width") {
		if v, ok := d.GetOkExists("width"); ok {
			request.Width = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("need_audio") {
		if v, ok := d.GetOkExists("need_audio"); ok {
			request.NeedAudio = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("height") {
		if v, ok := d.GetOkExists("height"); ok {
			request.Height = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("fps") {
		if v, ok := d.GetOkExists("fps"); ok {
			request.Fps = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("gop") {
		if v, ok := d.GetOkExists("gop"); ok {
			request.Gop = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("rotate") {
		if v, ok := d.GetOkExists("rotate"); ok {
			request.Rotate = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("profile") {
		if v, ok := d.GetOk("profile"); ok {
			request.Profile = helper.String(v.(string))
		}
	}

	if d.HasChange("bitrate_to_orig") {
		if v, ok := d.GetOkExists("bitrate_to_orig"); ok {
			request.BitrateToOrig = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("height_to_orig") {
		if v, ok := d.GetOkExists("height_to_orig"); ok {
			request.HeightToOrig = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("fps_to_orig") {
		if v, ok := d.GetOkExists("fps_to_orig"); ok {
			request.FpsToOrig = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("adapt_bitrate_percent") {
		if v, ok := d.GetOkExists("adapt_bitrate_percent"); ok {
			request.AdaptBitratePercent = helper.Float64(v.(float64))
		}
	}

	if d.HasChange("short_edge_as_height") {
		if v, ok := d.GetOkExists("short_edge_as_height"); ok {
			request.ShortEdgeAsHeight = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("d_r_m_type") {
		if v, ok := d.GetOk("d_r_m_type"); ok {
			request.DRMType = helper.String(v.(string))
		}
	}

	if d.HasChange("d_r_m_tracks") {
		if v, ok := d.GetOk("d_r_m_tracks"); ok {
			request.DRMTracks = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLiveTranscodeTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css liveTranscodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssLiveTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudCssLiveTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_live_transcode_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	liveTranscodeTemplateId := d.Id()

	if err := service.DeleteCssLiveTranscodeTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
