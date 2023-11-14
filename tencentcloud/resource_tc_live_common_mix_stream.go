/*
Provides a resource to create a live common_mix_stream

Example Usage

```hcl
resource "tencentcloud_live_common_mix_stream" "common_mix_stream" {
  mix_stream_session_id = "test_room"
  input_stream_list {
		input_stream_name = "demo"
		layout_params {
			image_layer = 1
			input_type = 1
			image_height =
			image_width =
			location_x =
			location_y =
			color = "0xcc0033"
			watermark_id = 123456
		}
		crop_params {
			crop_width =
			crop_height =
			crop_start_location_x =
			crop_start_location_y =
		}

  }
  output_params {
		output_stream_name = "demo"
		output_stream_type = 1
		output_stream_bit_rate = 20
		output_stream_gop = 5
		output_stream_frame_rate = 30
		output_audio_bit_rate = 20
		output_audio_sample_rate = 96000
		output_audio_channels = 1
		mix_sei = "demo_sei"

  }
  mix_stream_template_id = 123456
  control_params {
		use_mix_crop_center = 1
		allow_copy = 1
		pass_input_sei = 1

  }
}
```

Import

live common_mix_stream can be imported using the id, e.g.

```
terraform import tencentcloud_live_common_mix_stream.common_mix_stream common_mix_stream_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLiveCommonMixStream() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveCommonMixStreamCreate,
		Read:   resourceTencentCloudLiveCommonMixStreamRead,
		Update: resourceTencentCloudLiveCommonMixStreamUpdate,
		Delete: resourceTencentCloudLiveCommonMixStreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mix_stream_session_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of mixed stream session (from the beginning of applying for mixed stream to the end of canceling mixed stream). A string containing only letters, numbers, and underscores within 80 bytes.",
			},

			"input_stream_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Mixed stream input stream list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_stream_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the stream name. A string containing only letters, numbers, and underscores within 80 bytes.",
						},
						"layout_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Enter the flow layout parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_layer": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Import layers. Value range [1, 16]. 1) Image of background stream (i.e. big anchor screen or canvas)_ Fill in 1 for layer. 2) Pure audio mixed stream, this parameter also needs to be filled in.",
									},
									"input_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Enter the type. Value range [0, 5]. If it is not filled in, the default value is 0. 0 means the input stream is audio and video. 2 indicates that the input stream is a picture. 3 indicates that the input stream is a canvas. 4 indicates that the input stream is audio. 5 indicates that the input stream is pure video.",
									},
									"image_height": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Enter the height of the screen at the time of output. Value range: Pixel: [0,2000] Percentage: [0.01, 0.99] If it is left blank, the default is the height of the input stream. When using percentage, the expected output is (percentage * high background).",
									},
									"image_width": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "The width of the input screen at the time of output. Value range: Pixel: [0,2000] Percentage: [0.01, 0.99] If it is left blank, the default is the width of the input stream. When using percentage, the expected output is (percentage * background width).",
									},
									"location_x": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Enter the X offset on the output screen. Value Range: Pixel: [0,2000] Percentage: [0.01, 0.99] If not filled, the default value is 0. The horizontal offset relative to the upper left corner of the big anchor background screen. When using percentage, the expected output is (percentage * background width).",
									},
									"location_y": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Enter the Y offset on the output screen. Value Range: Pixel: [02000] Percentage: [0.01, 0.99] If not filled, the default value is 0. The vertical offset relative to the upper left corner of the big anchor background screen. When using percentage, the expected output is (percentage * background width).",
									},
									"color": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When InputType is 3 (Canvas), this value represents the color of the canvas.Common colors are:Red: 0xcc0033.Yellow: 0xcc9900.Green: 0xcccc33.Blue: 0x99CCFF.Black: 0x000000.White: 0xFFFFFF.Grey: 0x999999.",
									},
									"watermark_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "When InputType is 2 (picture), this value is the watermark ID.",
									},
								},
							},
						},
						"crop_params": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Enter flow clipping parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crop_width": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "The width of the crop. Value range [0,2000].",
									},
									"crop_height": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "The height of the crop. Value range [0,2000].",
									},
									"crop_start_location_x": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "The starting X coordinate of the clipping. Value range [0,2000].",
									},
									"crop_start_location_y": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "The starting Y coordinate of the clipping. Value range [0,2000].",
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
				Description: "Mixed stream output stream parameters.",
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
							Description: "Output stream type, value range [0,1].If it is not filled in, the default value is 0.When the output stream is one item in the input stream list, fill in 0.When the expected mixed flow result becomes a new flow, this value is filled as 1.When the value is 1, output_ stream_ ID cannot appear in input_ stram_ The stream with the same ID cannot exist in the list and the live broadcast background.",
						},
						"output_stream_bit_rate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream bit rate. The value range is [1，50000]. If it is not filled, the system will automatically judge.",
						},
						"output_stream_gop": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream GOP size. Value range [1,10]. If it is not filled, the system will automatically judge.",
						},
						"output_stream_frame_rate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream frame rate size. Value range [1,60]. If it is not filled, the system will automatically judge.",
						},
						"output_audio_bit_rate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream audio bit rate. Value range [1，500] If it is not filled, the system will automatically judge.",
						},
						"output_audio_sample_rate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Output stream audio sampling rate. Value range [96000, 88200, 64000, 48000, 44100, 3200024000, 22050, 16000, 12000, 11025, 8000]. If it is not filled, the system will automatically judge.",
						},
						"output_audio_channels": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of output stream audio channels. Value range [1,2]. If it is not filled, the system will automatically judge.",
						},
						"mix_sei": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sei information in the output stream. If there is no special need, do not fill in.",
						},
					},
				},
			},

			"mix_stream_template_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Enter the template ID. If this parameter is set, it will be output according to the default template layout. There is no need to fill in the user-defined location parameter.If it is not filled in, the default value is 0.Two input sources support 10, 20, 30, 40, 50.Three input sources support 310390391.Four input sources support 410.Five input sources support 510590.Six input sources support 610.",
			},

			"control_params": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Special control parameters of mixed flow. No need to fill in if there is no special demand.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_mix_crop_center": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Value range [0,1]. When filling in 1, when the layer resolution parameter in the parameter is inconsistent with the actual video resolution, it will automatically crop from the video according to the resolution ratio set by the layer.",
						},
						"allow_copy": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Value range [0,1] When filling in 1, when the number of InputStreamLists is 1 and the OutputParams. OutputStreamType is 1, do not cancel the operation, but copy the stream.",
						},
						"pass_input_sei": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Value range [0,1] When filling 1, the sei of the original flow is transmitted.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudLiveCommonMixStreamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_common_mix_stream.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = live.NewCreateCommonMixStreamRequest()
		response           = live.NewCreateCommonMixStreamResponse()
		mixStreamSessionId string
	)
	if v, ok := d.GetOk("mix_stream_session_id"); ok {
		mixStreamSessionId = v.(string)
		request.MixStreamSessionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("input_stream_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			commonMixInputParam := live.CommonMixInputParam{}
			if v, ok := dMap["input_stream_name"]; ok {
				commonMixInputParam.InputStreamName = helper.String(v.(string))
			}
			if layoutParamsMap, ok := helper.InterfaceToMap(dMap, "layout_params"); ok {
				commonMixLayoutParams := live.CommonMixLayoutParams{}
				if v, ok := layoutParamsMap["image_layer"]; ok {
					commonMixLayoutParams.ImageLayer = helper.IntInt64(v.(int))
				}
				if v, ok := layoutParamsMap["input_type"]; ok {
					commonMixLayoutParams.InputType = helper.IntInt64(v.(int))
				}
				if v, ok := layoutParamsMap["image_height"]; ok {
					commonMixLayoutParams.ImageHeight = helper.Float64(v.(float64))
				}
				if v, ok := layoutParamsMap["image_width"]; ok {
					commonMixLayoutParams.ImageWidth = helper.Float64(v.(float64))
				}
				if v, ok := layoutParamsMap["location_x"]; ok {
					commonMixLayoutParams.LocationX = helper.Float64(v.(float64))
				}
				if v, ok := layoutParamsMap["location_y"]; ok {
					commonMixLayoutParams.LocationY = helper.Float64(v.(float64))
				}
				if v, ok := layoutParamsMap["color"]; ok {
					commonMixLayoutParams.Color = helper.String(v.(string))
				}
				if v, ok := layoutParamsMap["watermark_id"]; ok {
					commonMixLayoutParams.WatermarkId = helper.IntInt64(v.(int))
				}
				commonMixInputParam.LayoutParams = &commonMixLayoutParams
			}
			if cropParamsMap, ok := helper.InterfaceToMap(dMap, "crop_params"); ok {
				commonMixCropParams := live.CommonMixCropParams{}
				if v, ok := cropParamsMap["crop_width"]; ok {
					commonMixCropParams.CropWidth = helper.Float64(v.(float64))
				}
				if v, ok := cropParamsMap["crop_height"]; ok {
					commonMixCropParams.CropHeight = helper.Float64(v.(float64))
				}
				if v, ok := cropParamsMap["crop_start_location_x"]; ok {
					commonMixCropParams.CropStartLocationX = helper.Float64(v.(float64))
				}
				if v, ok := cropParamsMap["crop_start_location_y"]; ok {
					commonMixCropParams.CropStartLocationY = helper.Float64(v.(float64))
				}
				commonMixInputParam.CropParams = &commonMixCropParams
			}
			request.InputStreamList = append(request.InputStreamList, &commonMixInputParam)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "output_params"); ok {
		commonMixOutputParams := live.CommonMixOutputParams{}
		if v, ok := dMap["output_stream_name"]; ok {
			commonMixOutputParams.OutputStreamName = helper.String(v.(string))
		}
		if v, ok := dMap["output_stream_type"]; ok {
			commonMixOutputParams.OutputStreamType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["output_stream_bit_rate"]; ok {
			commonMixOutputParams.OutputStreamBitRate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["output_stream_gop"]; ok {
			commonMixOutputParams.OutputStreamGop = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["output_stream_frame_rate"]; ok {
			commonMixOutputParams.OutputStreamFrameRate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["output_audio_bit_rate"]; ok {
			commonMixOutputParams.OutputAudioBitRate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["output_audio_sample_rate"]; ok {
			commonMixOutputParams.OutputAudioSampleRate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["output_audio_channels"]; ok {
			commonMixOutputParams.OutputAudioChannels = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["mix_sei"]; ok {
			commonMixOutputParams.MixSei = helper.String(v.(string))
		}
		request.OutputParams = &commonMixOutputParams
	}

	if v, ok := d.GetOkExists("mix_stream_template_id"); ok {
		request.MixStreamTemplateId = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "control_params"); ok {
		commonMixControlParams := live.CommonMixControlParams{}
		if v, ok := dMap["use_mix_crop_center"]; ok {
			commonMixControlParams.UseMixCropCenter = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["allow_copy"]; ok {
			commonMixControlParams.AllowCopy = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["pass_input_sei"]; ok {
			commonMixControlParams.PassInputSei = helper.IntInt64(v.(int))
		}
		request.ControlParams = &commonMixControlParams
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateCommonMixStream(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live commonMixStream failed, reason:%+v", logId, err)
		return err
	}

	mixStreamSessionId = *response.Response.MixStreamSessionId
	d.SetId(mixStreamSessionId)

	return resourceTencentCloudLiveCommonMixStreamRead(d, meta)
}

func resourceTencentCloudLiveCommonMixStreamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_common_mix_stream.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	commonMixStreamId := d.Id()

	commonMixStream, err := service.DescribeLiveCommonMixStreamById(ctx, mixStreamSessionId)
	if err != nil {
		return err
	}

	if commonMixStream == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveCommonMixStream` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if commonMixStream.MixStreamSessionId != nil {
		_ = d.Set("mix_stream_session_id", commonMixStream.MixStreamSessionId)
	}

	if commonMixStream.InputStreamList != nil {
		inputStreamListList := []interface{}{}
		for _, inputStreamList := range commonMixStream.InputStreamList {
			inputStreamListMap := map[string]interface{}{}

			if commonMixStream.InputStreamList.InputStreamName != nil {
				inputStreamListMap["input_stream_name"] = commonMixStream.InputStreamList.InputStreamName
			}

			if commonMixStream.InputStreamList.LayoutParams != nil {
				layoutParamsMap := map[string]interface{}{}

				if commonMixStream.InputStreamList.LayoutParams.ImageLayer != nil {
					layoutParamsMap["image_layer"] = commonMixStream.InputStreamList.LayoutParams.ImageLayer
				}

				if commonMixStream.InputStreamList.LayoutParams.InputType != nil {
					layoutParamsMap["input_type"] = commonMixStream.InputStreamList.LayoutParams.InputType
				}

				if commonMixStream.InputStreamList.LayoutParams.ImageHeight != nil {
					layoutParamsMap["image_height"] = commonMixStream.InputStreamList.LayoutParams.ImageHeight
				}

				if commonMixStream.InputStreamList.LayoutParams.ImageWidth != nil {
					layoutParamsMap["image_width"] = commonMixStream.InputStreamList.LayoutParams.ImageWidth
				}

				if commonMixStream.InputStreamList.LayoutParams.LocationX != nil {
					layoutParamsMap["location_x"] = commonMixStream.InputStreamList.LayoutParams.LocationX
				}

				if commonMixStream.InputStreamList.LayoutParams.LocationY != nil {
					layoutParamsMap["location_y"] = commonMixStream.InputStreamList.LayoutParams.LocationY
				}

				if commonMixStream.InputStreamList.LayoutParams.Color != nil {
					layoutParamsMap["color"] = commonMixStream.InputStreamList.LayoutParams.Color
				}

				if commonMixStream.InputStreamList.LayoutParams.WatermarkId != nil {
					layoutParamsMap["watermark_id"] = commonMixStream.InputStreamList.LayoutParams.WatermarkId
				}

				inputStreamListMap["layout_params"] = []interface{}{layoutParamsMap}
			}

			if commonMixStream.InputStreamList.CropParams != nil {
				cropParamsMap := map[string]interface{}{}

				if commonMixStream.InputStreamList.CropParams.CropWidth != nil {
					cropParamsMap["crop_width"] = commonMixStream.InputStreamList.CropParams.CropWidth
				}

				if commonMixStream.InputStreamList.CropParams.CropHeight != nil {
					cropParamsMap["crop_height"] = commonMixStream.InputStreamList.CropParams.CropHeight
				}

				if commonMixStream.InputStreamList.CropParams.CropStartLocationX != nil {
					cropParamsMap["crop_start_location_x"] = commonMixStream.InputStreamList.CropParams.CropStartLocationX
				}

				if commonMixStream.InputStreamList.CropParams.CropStartLocationY != nil {
					cropParamsMap["crop_start_location_y"] = commonMixStream.InputStreamList.CropParams.CropStartLocationY
				}

				inputStreamListMap["crop_params"] = []interface{}{cropParamsMap}
			}

			inputStreamListList = append(inputStreamListList, inputStreamListMap)
		}

		_ = d.Set("input_stream_list", inputStreamListList)

	}

	if commonMixStream.OutputParams != nil {
		outputParamsMap := map[string]interface{}{}

		if commonMixStream.OutputParams.OutputStreamName != nil {
			outputParamsMap["output_stream_name"] = commonMixStream.OutputParams.OutputStreamName
		}

		if commonMixStream.OutputParams.OutputStreamType != nil {
			outputParamsMap["output_stream_type"] = commonMixStream.OutputParams.OutputStreamType
		}

		if commonMixStream.OutputParams.OutputStreamBitRate != nil {
			outputParamsMap["output_stream_bit_rate"] = commonMixStream.OutputParams.OutputStreamBitRate
		}

		if commonMixStream.OutputParams.OutputStreamGop != nil {
			outputParamsMap["output_stream_gop"] = commonMixStream.OutputParams.OutputStreamGop
		}

		if commonMixStream.OutputParams.OutputStreamFrameRate != nil {
			outputParamsMap["output_stream_frame_rate"] = commonMixStream.OutputParams.OutputStreamFrameRate
		}

		if commonMixStream.OutputParams.OutputAudioBitRate != nil {
			outputParamsMap["output_audio_bit_rate"] = commonMixStream.OutputParams.OutputAudioBitRate
		}

		if commonMixStream.OutputParams.OutputAudioSampleRate != nil {
			outputParamsMap["output_audio_sample_rate"] = commonMixStream.OutputParams.OutputAudioSampleRate
		}

		if commonMixStream.OutputParams.OutputAudioChannels != nil {
			outputParamsMap["output_audio_channels"] = commonMixStream.OutputParams.OutputAudioChannels
		}

		if commonMixStream.OutputParams.MixSei != nil {
			outputParamsMap["mix_sei"] = commonMixStream.OutputParams.MixSei
		}

		_ = d.Set("output_params", []interface{}{outputParamsMap})
	}

	if commonMixStream.MixStreamTemplateId != nil {
		_ = d.Set("mix_stream_template_id", commonMixStream.MixStreamTemplateId)
	}

	if commonMixStream.ControlParams != nil {
		controlParamsMap := map[string]interface{}{}

		if commonMixStream.ControlParams.UseMixCropCenter != nil {
			controlParamsMap["use_mix_crop_center"] = commonMixStream.ControlParams.UseMixCropCenter
		}

		if commonMixStream.ControlParams.AllowCopy != nil {
			controlParamsMap["allow_copy"] = commonMixStream.ControlParams.AllowCopy
		}

		if commonMixStream.ControlParams.PassInputSei != nil {
			controlParamsMap["pass_input_sei"] = commonMixStream.ControlParams.PassInputSei
		}

		_ = d.Set("control_params", []interface{}{controlParamsMap})
	}

	return nil
}

func resourceTencentCloudLiveCommonMixStreamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_common_mix_stream.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"mix_stream_session_id", "input_stream_list", "output_params", "mix_stream_template_id", "control_params"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudLiveCommonMixStreamRead(d, meta)
}

func resourceTencentCloudLiveCommonMixStreamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_common_mix_stream.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	commonMixStreamId := d.Id()

	if err := service.DeleteLiveCommonMixStreamById(ctx, mixStreamSessionId); err != nil {
		return err
	}

	return nil
}
