package tencentcloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

const (
	VOD_AUDIO_CHANNEL_MONO   = "mono"
	VOD_AUDIO_CHANNEL_DUAL   = "dual"
	VOD_AUDIO_CHANNEL_STEREO = "stereo"

	VOD_DEFAULT_OFFSET = 0
	VOD_MAX_LIMIT      = 100
)

var (
	VOD_AUDIO_CHANNEL_TYPE_TO_INT = map[string]int64{
		VOD_AUDIO_CHANNEL_MONO:   1,
		VOD_AUDIO_CHANNEL_DUAL:   2,
		VOD_AUDIO_CHANNEL_STEREO: 6,
	}
	VOD_AUDIO_CHANNEL_TYPE_TO_STRING = map[int64]string{
		1: VOD_AUDIO_CHANNEL_MONO,
		2: VOD_AUDIO_CHANNEL_DUAL,
		6: VOD_AUDIO_CHANNEL_STEREO,
	}
	DISABLE_HIGHER_VIDEO_BITRATE_TO_UNINT = map[bool]uint64{
		true:  1,
		false: 0,
	}
	DISABLE_HIGHER_VIDEO_RESOLUTION_TO_UNINT = map[bool]uint64{
		true:  1,
		false: 0,
	}
	RESOLUTION_ADAPTIVE_TO_STRING = map[bool]string{
		true:  "open",
		false: "close",
	}
	REMOVE_AUDIO_TO_UNINT = map[bool]uint64{
		true:  1,
		false: 0,
	}
	DRM_SWITCH_TO_STRING = map[bool]string{
		true:  "ON",
		false: "OFF",
	}
)

func VodWatermarkResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"definition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Watermarking template ID.",
			},
			"text_content": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 100),
				Description:  "Text content of up to `100` characters. This needs to be entered only when the watermark type is text. Note: this field may return null, indicating that no valid values can be obtained.",
			},
			"svg_content": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 2000000),
				Description:  "SVG content of up to `2000000` characters. This needs to be entered only when the watermark type is `SVG`. Note: this field may return null, indicating that no valid values can be obtained.",
			},
			"start_time_offset": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Start time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame. If this parameter is left blank or `0` is entered, the watermark will appear upon the first video frame; If this value is greater than `0` (e.g., n), the watermark will appear at second n after the first video frame; If this value is smaller than `0` (e.g., -n), the watermark will appear at second n before the last video frame.",
			},
			"end_time_offset": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "End time offset of a watermark in seconds. If this parameter is left blank or `0` is entered, the watermark will exist till the last video frame; If this value is greater than `0` (e.g., n), the watermark will exist till second n; If this value is smaller than `0` (e.g., -n), the watermark will exist till second n before the last video frame.",
			},
		},
	}
}
