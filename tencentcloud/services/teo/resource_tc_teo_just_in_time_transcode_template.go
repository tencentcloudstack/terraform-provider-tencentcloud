package teo

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoJustInTimeTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoJustInTimeTranscodeTemplateCreate,
		Read:   resourceTencentCloudTeoJustInTimeTranscodeTemplateRead,
		Delete: resourceTencentCloudTeoJustInTimeTranscodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Transcode template name. Max length: 64 characters.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Template description. Max length: 256 characters.",
			},
			"video_stream_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"on", "off"}),
				Description:  "Video stream switch. Valid values: on, off. Default: on.",
			},
			"audio_stream_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"on", "off"}),
				Description:  "Audio stream switch. Valid values: on, off. Default: on.",
			},
			"video_template": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "Video stream configuration parameters. Required when video_stream_switch is on.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"video_codec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video codec. Optional values: H.264, H.265.",
						},
						"fps": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Video frame rate. Range: [0, 30]. Default: 0.",
						},
						"bitrate": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Video bitrate in kbps. Range: 0 or [128, 10000]. Default: 0.",
							ValidateFunc: tccommon.ValidateIntegerInRange(128, 10000),
						},
						"resolution_adaptive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resolution adaptive mode. Optional values: open, close. Default: open.",
						},
						"width": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Video width/long-edge in pixels. Range: 0 or [128, 1920]. Default: 0.",
							ValidateFunc: tccommon.ValidateIntegerInRange(128, 1920),
						},
						"height": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Video height/short-edge in pixels. Range: 0 or [128, 1080]. Default: 0.",
							ValidateFunc: tccommon.ValidateIntegerInRange(128, 1080),
						},
						"fill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Fill type. Optional values: stretch, black, white, gauss. Default: black.",
						},
					},
				},
			},
			"audio_template": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "Audio stream configuration parameters. Required when audio_stream_switch is on.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Audio codec. Optional values: libfdk_aac.",
						},
						"audio_channel": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2,
							Description:  "Audio channel count. Optional values: 2. Default: 2.",
							ValidateFunc: tccommon.ValidateIntegerInRange(1, 2),
						},
					},
				},
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template ID returned after creation.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template type. Values: preset, custom.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template creation time in ISO 8601 format.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template last update time in ISO 8601 format.",
			},
		},
	}
}

func resourceTencentCloudTeoJustInTimeTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_just_in_time_transcode_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		zoneId string
	)

	zoneId = d.Get("zone_id").(string)
	request := teo.NewCreateJustInTimeTranscodeTemplateRequest()
	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("template_name"); ok {
		templateName := v.(string)
		if len(templateName) > 64 {
			return fmt.Errorf("template_name exceeds maximum length of 64 characters")
		}
		request.TemplateName = helper.String(templateName)
	}

	if v, ok := d.GetOk("comment"); ok {
		comment := v.(string)
		if len(comment) > 256 {
			return fmt.Errorf("comment exceeds maximum length of 256 characters")
		}
		request.Comment = helper.String(comment)
	}

	if v, ok := d.GetOk("video_stream_switch"); ok {
		request.VideoStreamSwitch = helper.String(v.(string))
	}

	if v, ok := d.GetOk("audio_stream_switch"); ok {
		request.AudioStreamSwitch = helper.String(v.(string))
	}

	if videoTemplateList, ok := d.GetOk("video_template"); ok && len(videoTemplateList.([]interface{})) > 0 {
		videoTemplateMap := videoTemplateList.([]interface{})[0].(map[string]interface{})
		request.VideoTemplate = buildVideoTemplateInfo(videoTemplateMap)
	}

	if audioTemplateList, ok := d.GetOk("audio_template"); ok && len(audioTemplateList.([]interface{})) > 0 {
		audioTemplateMap := audioTemplateList.([]interface{})[0].(map[string]interface{})
		request.AudioTemplate = buildAudioTemplateInfo(audioTemplateMap)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateJustInTimeTranscodeTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		d.SetId(fmt.Sprintf("%s#%s", zoneId, *result.Response.TemplateId))
		d.Set("template_id", result.Response.TemplateId)
		return nil
	})
	if err != nil {
		log.Printf("[CRITICAL]%s create teo just-in-time transcode template failed, reason: %s", logId, err.Error())
		return err
	}

	return resourceTencentCloudTeoJustInTimeTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudTeoJustInTimeTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_just_in_time_transcode_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("resource ID is missing")
	}

	zoneId, templateId, err := resourceTencentCloudTeoJustInTimeTranscodeTemplateParseId(id)
	if err != nil {
		return err
	}

	request := teo.NewDescribeJustInTimeTranscodeTemplatesRequest()
	request.ZoneId = helper.String(zoneId)
	request.Filters = []*teo.Filter{
		{
			Name:   helper.String("template-id"),
			Values: []*string{helper.String(templateId)},
		},
	}
	request.Limit = helper.Int64(10)

	var template *teo.JustInTimeTranscodeTemplate
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeJustInTimeTranscodeTemplates(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response.Response.TotalCount == nil || *response.Response.TotalCount == 0 {
			return resource.NonRetryableError(fmt.Errorf("template not found"))
		}

		if len(response.Response.TemplateSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("template list is empty"))
		}

		template = response.Response.TemplateSet[0]
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s read teo just-in-time transcode template failed, reason: %s", logId, err.Error())
		return err
	}

	if template == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("template_id", template.TemplateId)
	_ = d.Set("template_name", template.TemplateName)
	_ = d.Set("comment", template.Comment)
	_ = d.Set("video_stream_switch", template.VideoStreamSwitch)
	_ = d.Set("audio_stream_switch", template.AudioStreamSwitch)
	_ = d.Set("type", template.Type)
	_ = d.Set("create_time", template.CreateTime)
	_ = d.Set("update_time", template.UpdateTime)

	if template.VideoTemplate != nil {
		if err := d.Set("video_template", []interface{}{mapVideoTemplateInfoToSchema(template.VideoTemplate)}); err != nil {
			return fmt.Errorf("failed to set video_template: %s", err)
		}
	}

	if template.AudioTemplate != nil {
		if err := d.Set("audio_template", []interface{}{mapAudioTemplateInfoToSchema(template.AudioTemplate)}); err != nil {
			return fmt.Errorf("failed to set audio_template: %s", err)
		}
	}

	return nil
}

func resourceTencentCloudTeoJustInTimeTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_just_in_time_transcode_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("resource ID is missing")
	}

	zoneId, templateId, err := resourceTencentCloudTeoJustInTimeTranscodeTemplateParseId(id)
	if err != nil {
		return err
	}

	request := teo.NewDeleteJustInTimeTranscodeTemplatesRequest()
	request.ZoneId = helper.String(zoneId)
	request.TemplateIds = []*string{helper.String(templateId)}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteJustInTimeTranscodeTemplates(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, request.GetAction(), request.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s delete teo just-in-time transcode template failed, reason: %s", logId, err.Error())
		return err
	}

	d.SetId("")
	return nil
}

func resourceTencentCloudTeoJustInTimeTranscodeTemplateParseId(id string) (string, string, error) {
	items := strings.Split(id, "#")
	if len(items) != 2 {
		return "", "", fmt.Errorf("invalid resource ID format, expected: zone_id#template_id, got: %s", id)
	}
	return items[0], items[1], nil
}

func buildVideoTemplateInfo(videoTemplateMap map[string]interface{}) *teo.VideoTemplateInfo {
	info := &teo.VideoTemplateInfo{}

	if v, ok := videoTemplateMap["video_codec"]; ok {
		info.Codec = helper.String(v.(string))
	}
	if v, ok := videoTemplateMap["fps"]; ok {
		info.Fps = helper.Float64(v.(float64))
	}
	if v, ok := videoTemplateMap["bitrate"]; ok {
		info.Bitrate = helper.Uint64(uint64(v.(int)))
	}
	if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
		info.ResolutionAdaptive = helper.String(v.(string))
	}
	if v, ok := videoTemplateMap["width"]; ok {
		info.Width = helper.Uint64(uint64(v.(int)))
	}
	if v, ok := videoTemplateMap["height"]; ok {
		info.Height = helper.Uint64(uint64(v.(int)))
	}
	if v, ok := videoTemplateMap["fill_type"]; ok {
		info.FillType = helper.String(v.(string))
	}

	return info
}

func buildAudioTemplateInfo(audioTemplateMap map[string]interface{}) *teo.AudioTemplateInfo {
	info := &teo.AudioTemplateInfo{}

	if v, ok := audioTemplateMap["codec"]; ok {
		info.Codec = helper.String(v.(string))
	}
	if v, ok := audioTemplateMap["audio_channel"]; ok {
		info.AudioChannel = helper.Uint64(uint64(v.(int)))
	}

	return info
}

func mapVideoTemplateInfoToSchema(info *teo.VideoTemplateInfo) map[string]interface{} {
	result := make(map[string]interface{})

	if info.Codec != nil {
		result["video_codec"] = *info.Codec
	}
	if info.Fps != nil {
		result["fps"] = *info.Fps
	}
	if info.Bitrate != nil {
		result["bitrate"] = int(*info.Bitrate)
	}
	if info.ResolutionAdaptive != nil {
		result["resolution_adaptive"] = *info.ResolutionAdaptive
	}
	if info.Width != nil {
		result["width"] = int(*info.Width)
	}
	if info.Height != nil {
		result["height"] = int(*info.Height)
	}
	if info.FillType != nil {
		result["fill_type"] = *info.FillType
	}

	return result
}

func mapAudioTemplateInfoToSchema(info *teo.AudioTemplateInfo) map[string]interface{} {
	result := make(map[string]interface{})

	if info.Codec != nil {
		result["codec"] = *info.Codec
	}
	if info.AudioChannel != nil {
		result["audio_channel"] = int(*info.AudioChannel)
	}

	return result
}
