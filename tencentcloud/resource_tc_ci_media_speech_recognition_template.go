/*
Provides a resource to create a ci media_speech_recognition_template

Example Usage

```hcl
resource "tencentcloud_ci_media_speech_recognition_template" "media_speech_recognition_template" {
  bucket = "terraform-ci-1308919341"
  name = "speech_recognition_template"
  speech_recognition {
		engine_model_type = "16k_zh"
		channel_num = "1"
		res_text_format = "1"
		filter_dirty = "0"
		filter_modal = "1"
		convert_num_mode = "0"
		speaker_diarization = "1"
		speaker_number = "0"
		filter_punc = "0"
		output_file_type = "txt"
  }
}
```

Import

ci media_speech_recognition_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_speech_recognition_template.media_speech_recognition_template terraform-ci-xxxxxx#t1d794430f2f1f4350b11e905ce2c6167e
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

func resourceTencentCloudCiMediaSpeechRecognitionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaSpeechRecognitionTemplateCreate,
		Read:   resourceTencentCloudCiMediaSpeechRecognitionTemplateRead,
		Update: resourceTencentCloudCiMediaSpeechRecognitionTemplateUpdate,
		Delete: resourceTencentCloudCiMediaSpeechRecognitionTemplateDelete,
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

			"speech_recognition": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "audio configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_model_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Engine model type, divided into phone scene and non-phone scene, phone scene: 8k_zh: phone 8k Chinese Mandarin general (can be used for dual-channel audio), 8k_zh_s: phone 8k Chinese Mandarin speaker separation (only for monophonic audio), 8k_en: Telephone 8k English; non-telephone scene: 16k_zh: 16k Mandarin Chinese, 16k_zh_video: 16k audio and video field, 16k_en: 16k English, 16k_ca: 16k Cantonese, 16k_ja: 16k Japanese, 16k_zh_edu: Chinese education, 16k_en_edu: English education, 16k_zh_medical: medical, 16k_th: Thai, 16k_zh_dialect: multi-dialect, supports 23 dialects.",
						},
						"channel_num": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Number of voice channels: 1 means mono. EngineModelType supports only mono for non-telephone scenarios, and 2 means dual channels (only 8k_zh engine model supports dual channels, which should correspond to both sides of the call).",
						},
						"res_text_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Recognition result return form: 0 means the recognition result text (including segmented time stamps), 1 is the detailed recognition result at the word level granularity, without punctuation, and includes the speech rate value (a list of word time stamps, generally used to generate subtitle scenes), 2 Detailed recognition results at word-level granularity (including punctuation and speech rate values)..",
						},
						"filter_dirty": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to filter dirty words (currently supports Mandarin Chinese engine): 0 means not to filter dirty words, 1 means to filter dirty words, 2 means to replace dirty words with *, the default value is 0.",
						},
						"filter_modal": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to pass modal particles (currently supports Mandarin Chinese engine): 0 means not to filter modal particles, 1 means partial filtering, 2 means strict filtering, and the default value is 0.",
						},
						"convert_num_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to perform intelligent conversion of Arabic numerals (currently supports Mandarin Chinese engine): 0 means no conversion, directly output Chinese numbers, 1 means intelligently convert to Arabic numerals according to the scene, 3 means enable math-related digital conversion, the default value is 0.",
						},
						"speaker_diarization": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to enable speaker separation: 0 means not enabled, 1 means enabled (only supports 8k_zh, 16k_zh, 16k_zh_video, monophonic audio), the default value is 0, Note: 8K telephony scenarios suggest using dual-channel to distinguish between the two parties, set ChannelNum=2 is enough, no need to enable speaker separation.",
						},
						"speaker_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The number of speakers to be separated (need to be used in conjunction with enabling speaker separation), value range: 0-10, 0 means automatic separation (currently only supports <= 6 people), 1-10 represents the number of specified speakers to be separated. The default value is 0.",
						},
						"filter_punc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to filter punctuation (currently supports Mandarin Chinese engine): 0 means no filtering, 1 means filtering end-of-sentence punctuation, 2 means filtering all punctuation, the default value is 0.",
						},
						"output_file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Output file type, optional txt, srt. The default is txt.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaSpeechRecognitionTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_speech_recognition_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaSpeechRecognitionTemplateOptions{
			Tag: "SpeechRecognition",
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

	if dMap, ok := helper.InterfacesHeadMap(d, "speech_recognition"); ok {
		speechRecognition := cos.SpeechRecognition{}
		if v, ok := dMap["engine_model_type"]; ok {
			speechRecognition.EngineModelType = v.(string)
		}
		if v, ok := dMap["channel_num"]; ok {
			speechRecognition.ChannelNum = v.(string)
		}
		if v, ok := dMap["res_text_format"]; ok {
			speechRecognition.ResTextFormat = v.(string)
		}
		if v, ok := dMap["filter_dirty"]; ok {
			speechRecognition.FilterDirty = v.(string)
		}
		if v, ok := dMap["filter_modal"]; ok {
			speechRecognition.FilterModal = v.(string)
		}
		if v, ok := dMap["convert_num_mode"]; ok {
			speechRecognition.ConvertNumMode = v.(string)
		}
		if v, ok := dMap["speaker_diarization"]; ok {
			speechRecognition.SpeakerDiarization = v.(string)
		}
		if v, ok := dMap["speaker_number"]; ok {
			speechRecognition.SpeakerNumber = v.(string)
		}
		if v, ok := dMap["filter_punc"]; ok {
			speechRecognition.FilterPunc = v.(string)
		}
		if v, ok := dMap["output_file_type"]; ok {
			speechRecognition.OutputFileType = v.(string)
		}
		request.SpeechRecognition = &speechRecognition
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaSpeechRecognitionTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaSpeechRecognitionTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSpeechRecognitionTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaSpeechRecognitionTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSpeechRecognitionTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_speech_recognition_template.read")()
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

	mediaSpeechRecognitionTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaSpeechRecognitionTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if mediaSpeechRecognitionTemplate.Name != "" {
		_ = d.Set("name", mediaSpeechRecognitionTemplate.Name)
	}

	if mediaSpeechRecognitionTemplate.SpeechRecognition != nil {
		speechRecognitionMap := map[string]interface{}{}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.EngineModelType != "" {
			speechRecognitionMap["engine_model_type"] = mediaSpeechRecognitionTemplate.SpeechRecognition.EngineModelType
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.ChannelNum != "" {
			speechRecognitionMap["channel_num"] = mediaSpeechRecognitionTemplate.SpeechRecognition.ChannelNum
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.ResTextFormat != "" {
			speechRecognitionMap["res_text_format"] = mediaSpeechRecognitionTemplate.SpeechRecognition.ResTextFormat
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.FilterDirty != "" {
			speechRecognitionMap["filter_dirty"] = mediaSpeechRecognitionTemplate.SpeechRecognition.FilterDirty
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.FilterModal != "" {
			speechRecognitionMap["filter_modal"] = mediaSpeechRecognitionTemplate.SpeechRecognition.FilterModal
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.ConvertNumMode != "" {
			speechRecognitionMap["convert_num_mode"] = mediaSpeechRecognitionTemplate.SpeechRecognition.ConvertNumMode
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.SpeakerDiarization != "" {
			speechRecognitionMap["speaker_diarization"] = mediaSpeechRecognitionTemplate.SpeechRecognition.SpeakerDiarization
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.SpeakerNumber != "" {
			speechRecognitionMap["speaker_number"] = mediaSpeechRecognitionTemplate.SpeechRecognition.SpeakerNumber
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.FilterPunc != "" {
			speechRecognitionMap["filter_punc"] = mediaSpeechRecognitionTemplate.SpeechRecognition.FilterPunc
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.OutputFileType != "" {
			speechRecognitionMap["output_file_type"] = mediaSpeechRecognitionTemplate.SpeechRecognition.OutputFileType
		}

		_ = d.Set("speech_recognition", []interface{}{speechRecognitionMap})
	}

	return nil
}

func resourceTencentCloudCiMediaSpeechRecognitionTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_speech_recognition_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaSpeechRecognitionTemplateOptions{
		Tag: "SpeechRecognition",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("speech_recognition") {
		if dMap, ok := helper.InterfacesHeadMap(d, "speech_recognition"); ok {
			speechRecognition := cos.SpeechRecognition{}
			if v, ok := dMap["engine_model_type"]; ok {
				speechRecognition.EngineModelType = v.(string)
			}
			if v, ok := dMap["channel_num"]; ok {
				speechRecognition.ChannelNum = v.(string)
			}
			if v, ok := dMap["res_text_format"]; ok {
				speechRecognition.ResTextFormat = v.(string)
			}
			if v, ok := dMap["filter_dirty"]; ok {
				speechRecognition.FilterDirty = v.(string)
			}
			if v, ok := dMap["filter_modal"]; ok {
				speechRecognition.FilterModal = v.(string)
			}
			if v, ok := dMap["convert_num_mode"]; ok {
				speechRecognition.ConvertNumMode = v.(string)
			}
			if v, ok := dMap["speaker_diarization"]; ok {
				speechRecognition.SpeakerDiarization = v.(string)
			}
			if v, ok := dMap["speaker_number"]; ok {
				speechRecognition.SpeakerNumber = v.(string)
			}
			if v, ok := dMap["filter_punc"]; ok {
				speechRecognition.FilterPunc = v.(string)
			}
			if v, ok := dMap["output_file_type"]; ok {
				speechRecognition.OutputFileType = v.(string)
			}
			request.SpeechRecognition = &speechRecognition
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaSpeechRecognitionTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaSpeechRecognitionTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSpeechRecognitionTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaSpeechRecognitionTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSpeechRecognitionTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_speech_recognition_template.delete")()
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
