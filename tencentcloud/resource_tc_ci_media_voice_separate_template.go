/*
Provides a resource to create a ci media_voice_separate_template

Example Usage

```hcl
resource "tencentcloud_ci_media_voice_separate_template" "media_voice_separate_template" {
  name = &lt;nil&gt;
  audio_mode = &lt;nil&gt;
  audio_config {
		codec = &lt;nil&gt;
		samplerate = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		channels = &lt;nil&gt;

  }
}
```

Import

ci media_voice_separate_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_voice_separate_template.media_voice_separate_template media_voice_separate_template_id
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

func resourceTencentCloudCiMediaVoiceSeparateTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaVoiceSeparateTemplateCreate,
		Read:   resourceTencentCloudCiMediaVoiceSeparateTemplateRead,
		Update: resourceTencentCloudCiMediaVoiceSeparateTemplateUpdate,
		Delete: resourceTencentCloudCiMediaVoiceSeparateTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"audio_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Output audio IsAudio: output human voice, IsBackground: output background sound, AudioAndBackground: output vocal and background sound.",
			},

			"audio_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Audio configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Codec format, value aac, mp3, flac, amr.",
						},
						"samplerate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sampling Rate- 1: Unit: Hz- 2: Optional 8000, 11025, 22050, 32000, 44100, 48000, 96000- 3: When Codec is set to aac/flac, 8000 is not supported- 4: When Codec is set to mp3, 8000 and 96000 are not supported- 5: When Codec is set to amr, only 8000 is supported.",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Original audio bit rate, unit: Kbps, Value range: [8, 1000].",
						},
						"channels": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Number of channels- When Codec is set to aac/flac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3, support 1, 2- When Codec is set to amr, only 1 is supported.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaVoiceSeparateTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_voice_separate_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ci.NewCreateMediaVoiceSeparateTemplateRequest()
		response   = ci.NewCreateMediaVoiceSeparateTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("audio_mode"); ok {
		request.AudioMode = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio_config"); ok {
		audioConfig := ci.AudioConfig{}
		if v, ok := dMap["codec"]; ok {
			audioConfig.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["samplerate"]; ok {
			audioConfig.Samplerate = helper.String(v.(string))
		}
		if v, ok := dMap["bitrate"]; ok {
			audioConfig.Bitrate = helper.String(v.(string))
		}
		if v, ok := dMap["channels"]; ok {
			audioConfig.Channels = helper.String(v.(string))
		}
		request.AudioConfig = &audioConfig
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaVoiceSeparateTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaVoiceSeparateTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaVoiceSeparateTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVoiceSeparateTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_voice_separate_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaVoiceSeparateTemplateId := d.Id()

	mediaVoiceSeparateTemplate, err := service.DescribeCiMediaVoiceSeparateTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaVoiceSeparateTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaVoiceSeparateTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaVoiceSeparateTemplate.Name != nil {
		_ = d.Set("name", mediaVoiceSeparateTemplate.Name)
	}

	if mediaVoiceSeparateTemplate.AudioMode != nil {
		_ = d.Set("audio_mode", mediaVoiceSeparateTemplate.AudioMode)
	}

	if mediaVoiceSeparateTemplate.AudioConfig != nil {
		audioConfigMap := map[string]interface{}{}

		if mediaVoiceSeparateTemplate.AudioConfig.Codec != nil {
			audioConfigMap["codec"] = mediaVoiceSeparateTemplate.AudioConfig.Codec
		}

		if mediaVoiceSeparateTemplate.AudioConfig.Samplerate != nil {
			audioConfigMap["samplerate"] = mediaVoiceSeparateTemplate.AudioConfig.Samplerate
		}

		if mediaVoiceSeparateTemplate.AudioConfig.Bitrate != nil {
			audioConfigMap["bitrate"] = mediaVoiceSeparateTemplate.AudioConfig.Bitrate
		}

		if mediaVoiceSeparateTemplate.AudioConfig.Channels != nil {
			audioConfigMap["channels"] = mediaVoiceSeparateTemplate.AudioConfig.Channels
		}

		_ = d.Set("audio_config", []interface{}{audioConfigMap})
	}

	return nil
}

func resourceTencentCloudCiMediaVoiceSeparateTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_voice_separate_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaVoiceSeparateTemplateRequest()

	mediaVoiceSeparateTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "audio_mode", "audio_config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaVoiceSeparateTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaVoiceSeparateTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaVoiceSeparateTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVoiceSeparateTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_voice_separate_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	mediaVoiceSeparateTemplateId := d.Id()

	if err := service.DeleteCiMediaVoiceSeparateTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
