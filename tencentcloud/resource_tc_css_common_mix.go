/*
Provides a resource to create a css common_mix

Example Usage

```hcl
resource "tencentcloud_css_common_mix" "common_mix" {
  mix_stream_session_id = ""
  input_stream_list {
		input_stream_name = ""
		layout_params {
			image_layer =
			input_type =
			image_height =
			image_width =
			location_x =
			location_y =
			color = ""
			watermark_id =
		}
		crop_params {
			crop_width =
			crop_height =
			crop_start_location_x =
			crop_start_location_y =
		}

  }
  output_params {
		output_stream_name = ""
		output_stream_type =
		output_stream_bit_rate =
		output_stream_gop =
		output_stream_frame_rate =
		output_audio_bit_rate =
		output_audio_sample_rate =
		output_audio_channels =
		mix_sei = ""

  }
  mix_stream_template_id =
  control_params {
		use_mix_crop_center =
		allow_copy =
		pass_input_sei =

  }
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssCommonMix() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssCommonMixCreate,
		Read:   resourceTencentCloudCssCommonMixRead,
		Update: resourceTencentCloudCssCommonMixUpdate,
		Delete: resourceTencentCloudCssCommonMixDelete,
		Schema: map[string]*schema.Schema{
			"mix_stream_session_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of a stream mix session (from applying for the stream mix to cancelling it). This parameter can contain up to 80 bytes of letters, digits, and underscores.",
			},

			"input_stream_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Input stream list for stream mix.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_stream_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Input stream name, which can contain up to 80 bytes of letters, digits, and underscores.The value should be the name of an input stream for stream mix when `LayoutParams.InputType` is set to `0` (audio and video), `4` (pure audio), or `5` (pure video).The value can be a random name for identification, such as `Canvas1` or `Picture1`, when `LayoutParams.InputType` is set to `2` (image) or `3` (canvas).",
						},
						"layout_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Input stream layout parameter.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_layer": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Input layer. Value range: [1,16](1) For the background stream, i.e., the room owner’s image or the canvas, set this parameter to `1`.(2) This parameter is required for audio-only stream mixing as well.Note that two inputs cannot have the same `ImageLayer` value.",
									},
									"input_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Input type. Value range: [0,5].If this parameter is left empty, 0 will be used by default.0: the input stream is audio/video.2: the input stream is image.3: the input stream is canvas. 4: the input stream is audio.5: the input stream is pure video.",
									},
									"image_height": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Output height of input video image. Value range:Pixel: [0,2000]Percentage: [0.01,0.99]If this parameter is left empty, the input stream height will be used by default.If percentage is used, the expected output is (percentage * background height).",
									},
									"image_width": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Output width of input video image. Value range:Pixel: [0,2000]Percentage: [0.01,0.99]If this parameter is left empty, the input stream width will be used by default.If percentage is used, the expected output is (percentage * background width).",
									},
									"location_x": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "X-axis offset of input in output video image. Value range:Pixel: [0,2000]Percentage: [0.01,0.99]If this parameter is left empty, 0 will be used by default.Horizontal offset from the top-left corner of main host background video image. If percentage is used, the expected output is (percentage * background width).",
									},
									"location_y": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Y-axis offset of input in output video image. Value range:Pixel: [0,2000]Percentage: [0.01,0.99]If this parameter is left empty, 0 will be used by default.Vertical offset from the top-left corner of main host background video image. If percentage is used, the expected output is (percentage * background width).",
									},
									"color": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When `InputType` is 3 (canvas), this value indicates the canvas color.Commonly used colors include:Red: 0xcc0033.Yellow: 0xcc9900.Green: 0xcccc33.Blue: 0x99CCFF.Black: 0x000000.White: 0xFFFFFF.Gray: 0x999999.",
									},
									"watermark_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "When `InputType` is 2 (image), this value is the watermark ID.",
									},
								},
							},
						},
						"crop_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Input stream crop parameter.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crop_width": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Crop width. Value range: [0,2000].",
									},
									"crop_height": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Crop height. Value range: [0,2000].",
									},
									"crop_start_location_x": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Starting crop X coordinate. Value range: [0,2000].",
									},
									"crop_start_location_y": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Starting crop Y coordinate. Value range: [0,2000].",
									},
								},
							},
						},
					},
				},
			},

			"output_params": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Output stream parameter for stream mix.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"output_stream_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Output stream name.",
						},
						"output_stream_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream type. Valid values: [0,1].If this parameter is left empty, 0 will be used by default.If the output stream is a stream in the input stream list, enter 0.If you want the stream mix result to be a new stream, enter 1.If this value is 1, `output_stream_id` cannot appear in `input_stram_list`, and there cannot be a stream with the same ID on the LVB backend.",
						},
						"output_stream_bit_rate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The output bitrate. Value range: 1-10000.If you do not specify this, the system will select a bitrate automatically.",
						},
						"output_stream_gop": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream GOP size. Value range: [1,10].If this parameter is left empty, the system will automatically determine.",
						},
						"output_stream_frame_rate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream frame rate. Value range: [1,60].If this parameter is left empty, the system will automatically determine.",
						},
						"output_audio_bit_rate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream audio bitrate. Value range: [1,500]If this parameter is left empty, the system will automatically determine.",
						},
						"output_audio_sample_rate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream audio sample rate. Valid values: [96000, 88200, 64000, 48000, 44100, 32000,24000, 22050, 16000, 12000, 11025, 8000].If this parameter is left empty, the system will automatically determine.",
						},
						"output_audio_channels": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream audio sound channel. Valid values: [1,2].If this parameter is left empty, the system will automatically determine.",
						},
						"mix_sei": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SEI information in output stream. If there are no special needs, leave it empty.",
						},
					},
				},
			},

			"mix_stream_template_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Input template ID. If this parameter is set, the output will be generated according to the default template layout, and there is no need to enter the custom position parameters.If this parameter is left empty, 0 will be used by default.For two input sources, 10, 20, 30, 40, and 50 are supported.For three input sources, 310, 390, and 391 are supported.For four input sources, 410 is supported.For five input sources, 510 and 590 are supported.For six input sources, 610 is supported.",
			},

			"control_params": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Special control parameter for stream mix. If there are no special needs, leave it empty.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_mix_crop_center": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Value range: [0,1]. If 1 is entered, when the layer resolution in the parameter is different from the actual video resolution, the video will be automatically cropped according to the resolution set by the layer.",
						},
						"allow_copy": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Value range: [0,1].If this parameter is set to 1, when both `InputStreamList` and `OutputParams.OutputStreamType` are set to 1, you can copy a stream instead of canceling it.",
						},
						"pass_input_sei": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Valid values: 0, 1If you set this parameter to 1, SEI (Supplemental Enhanced Information) of the input streams will be passed through.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCssCommonMixCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_common_mix.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = css.NewCreateCommonMixStreamRequest()
		mixStreamSessionId string
	)
	if v, ok := d.GetOk("mix_stream_session_id"); ok {
		request.MixStreamSessionId = helper.String(v.(string))
		mixStreamSessionId = v.(string)
	}

	if v, ok := d.GetOk("input_stream_list"); ok {
		for ii, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			commonMixInputParam := css.CommonMixInputParam{}
			if v, ok := dMap["input_stream_name"]; ok {
				commonMixInputParam.InputStreamName = helper.String(v.(string))
			}

			indexLayoutParams :=fmt.Sprintf("input_stream_list.%v.layout_params.0", ii)
			// if layoutParamsMap, ok := helper.InterfaceToMap(dMap, "layout_params"); ok {
			if _, ok := d.GetOk(indexLayoutParams); ok {
				indexLayoutParams = indexLayoutParams + "."
				commonMixLayoutParams := css.CommonMixLayoutParams{}
				// if v, ok := layoutParamsMap["image_layer"]; ok {
				// 		commonMixLayoutParams.ImageLayer = helper.IntInt64(v.(int))
				// }
				if v, ok := d.GetOk(indexLayoutParams+"image_layer"); ok {
					commonMixLayoutParams.ImageLayer = helper.IntInt64(v.(int))
				}
				// if v, ok := layoutParamsMap["input_type"]; ok {
				// 	commonMixLayoutParams.InputType = helper.IntInt64(v.(int))
				// }
				if v, ok := d.GetOk(indexLayoutParams+"input_type"); ok {
					commonMixLayoutParams.InputType = helper.IntInt64(v.(int))
				}

				// if v, ok := layoutParamsMap["image_height"]; ok {
					if v, ok := d.GetOk(indexLayoutParams+"image_height"); ok {
					commonMixLayoutParams.ImageHeight = helper.Float64(v.(float64))
				}
				// if v, ok := layoutParamsMap["image_width"]; ok {
					if v, ok := d.GetOk(indexLayoutParams+"image_width"); ok {
					commonMixLayoutParams.ImageWidth = helper.Float64(v.(float64))
				}
				// if v, ok := layoutParamsMap["location_x"]; ok {
					if v, ok := d.GetOk(indexLayoutParams+"location_x"); ok {
					commonMixLayoutParams.LocationX = helper.Float64(v.(float64))
				}
				// if v, ok := layoutParamsMap["location_y"]; ok {
					if v, ok := d.GetOk(indexLayoutParams+"location_y"); ok {
					commonMixLayoutParams.LocationY = helper.Float64(v.(float64))
				}
				// if v, ok := layoutParamsMap["color"]; ok {
					if v, ok := d.GetOk(indexLayoutParams+"color"); ok {
					commonMixLayoutParams.Color = helper.String(v.(string))
				}
				// if v, ok := layoutParamsMap["watermark_id"]; ok {
					if v, ok := d.GetOk(indexLayoutParams+"watermark_id"); ok {
					commonMixLayoutParams.WatermarkId = helper.IntInt64(v.(int))
				}
				commonMixInputParam.LayoutParams = &commonMixLayoutParams
			}

			indexCropParams :=fmt.Sprintf("input_stream_list.%v.crop_params.0", ii)
			// if cropParamsMap, ok := helper.InterfaceToMap(dMap, "crop_params"); ok {
			if _, ok := d.GetOk(indexCropParams); ok {
				indexCropParams = indexCropParams + "."
				commonMixCropParams := css.CommonMixCropParams{}
				// if v, ok := cropParamsMap["crop_width"]; ok {
					if v, ok := d.GetOk(indexCropParams+"crop_width"); ok {
					commonMixCropParams.CropWidth = helper.Float64(v.(float64))
				}
				// if v, ok := cropParamsMap["crop_height"]; ok {
					if v, ok := d.GetOk(indexCropParams+"crop_height"); ok {
					commonMixCropParams.CropHeight = helper.Float64(v.(float64))
				}
				// if v, ok := cropParamsMap["crop_start_location_x"]; ok {
					if v, ok := d.GetOk(indexCropParams+"crop_start_location_x"); ok {
					commonMixCropParams.CropStartLocationX = helper.Float64(v.(float64))
				}
				// if v, ok := cropParamsMap["crop_start_location_y"]; ok {
					if v, ok := d.GetOk(indexCropParams+"crop_start_location_y"); ok {
					commonMixCropParams.CropStartLocationY = helper.Float64(v.(float64))
				}
				commonMixInputParam.CropParams = &commonMixCropParams
			}
			request.InputStreamList = append(request.InputStreamList, &commonMixInputParam)
		}
	}

	indexOutputParams :=fmt.Sprintf("output_params.%v", 0)
	// if dMap, ok := helper.InterfacesHeadMap(d, "output_params"); ok {
	if _, ok := d.GetOk(indexOutputParams); ok {
		indexOutputParams = indexOutputParams + "."
		commonMixOutputParams := css.CommonMixOutputParams{}
		// if v, ok := dMap["output_stream_name"]; ok {
			if v, ok := d.GetOk(indexOutputParams+"output_stream_name"); ok {
			commonMixOutputParams.OutputStreamName = helper.String(v.(string))
		}
		// if v, ok := dMap["output_stream_type"]; ok {
		if v, ok := d.GetOk(indexOutputParams+"output_stream_type"); ok {
			commonMixOutputParams.OutputStreamType = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["output_stream_bit_rate"]; ok {
			if v, ok := d.GetOk(indexOutputParams+"output_stream_bit_rate"); ok {
				commonMixOutputParams.OutputStreamBitRate = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["output_stream_gop"]; ok {
			if v, ok := d.GetOk(indexOutputParams+"output_stream_gop"); ok {
			commonMixOutputParams.OutputStreamGop = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["output_stream_frame_rate"]; ok {
			if v, ok := d.GetOk(indexOutputParams+"output_stream_frame_rate"); ok {
			commonMixOutputParams.OutputStreamFrameRate = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["output_audio_bit_rate"]; ok {
			if v, ok := d.GetOk(indexOutputParams+"output_audio_bit_rate"); ok {
			commonMixOutputParams.OutputAudioBitRate = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["output_audio_sample_rate"]; ok {
			if v, ok := d.GetOk(indexOutputParams+"output_audio_sample_rate"); ok {
			commonMixOutputParams.OutputAudioSampleRate = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["output_audio_channels"]; ok {
			if v, ok := d.GetOk(indexOutputParams+"output_audio_channels"); ok {
			commonMixOutputParams.OutputAudioChannels = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["mix_sei"]; ok {
			if v, ok := d.GetOk(indexOutputParams+"mix_sei"); ok {
			commonMixOutputParams.MixSei = helper.String(v.(string))
		}
		request.OutputParams = &commonMixOutputParams
	}

	if v, ok := d.GetOkExists("mix_stream_template_id"); ok {
		request.MixStreamTemplateId = helper.IntInt64(v.(int))
	}

	indexControlParams :=fmt.Sprintf("control_params.%v", 0)
	// if dMap, ok := helper.InterfacesHeadMap(d, "control_params"); ok {
	if _, ok := d.GetOk(indexControlParams); ok {
		indexControlParams = indexControlParams + "."
		commonMixControlParams := css.CommonMixControlParams{}
		// if v, ok := dMap["use_mix_crop_center"]; ok {
		if v, ok := d.GetOk(indexControlParams+"use_mix_crop_center"); ok {
			commonMixControlParams.UseMixCropCenter = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["allow_copy"]; ok {
		if v, ok := d.GetOk(indexControlParams+"allow_copy"); ok {
			commonMixControlParams.AllowCopy = helper.IntInt64(v.(int))
		}
		// if v, ok := dMap["pass_input_sei"]; ok {
		if v, ok := d.GetOk(indexControlParams+"pass_input_sei"); ok {
			commonMixControlParams.PassInputSei = helper.IntInt64(v.(int))
		}
		request.ControlParams = &commonMixControlParams
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateCommonMixStream(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css commonMix failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(mixStreamSessionId)

	return resourceTencentCloudCssCommonMixRead(d, meta)
}

func resourceTencentCloudCssCommonMixRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_common_mix.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCssCommonMixUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_common_mix.update")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCssCommonMixDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_common_mix.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	commonMixId := d.Id()

	if err := service.DeleteCssCommonMixById(ctx, commonMixId); err != nil {
		return err
	}

	return nil
}
