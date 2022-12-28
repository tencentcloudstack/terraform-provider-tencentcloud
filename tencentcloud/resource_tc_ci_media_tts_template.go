/*
Provides a resource to create a ci media_tts_template

Example Usage

```hcl
resource "tencentcloud_ci_media_tts_template" "media_tts_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "tts_template"
  mode = "Asyc"
  codec = "pcm"
  voice_type = "ruxue"
  volume = "0"
  speed = "100"
}
```

Import

ci media_tts_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_tts_template.media_tts_template terraform-ci-xxxxxx#t1ed421df8bd2140b6b73474f70f99b0f8
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
	"github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCiMediaTtsTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaTtsTemplateCreate,
		Read:   resourceTencentCloudCiMediaTtsTemplateRead,
		Update: resourceTencentCloudCiMediaTtsTemplateUpdate,
		Delete: resourceTencentCloudCiMediaTtsTemplateDelete,
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

			"mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Processing mode, default value Asyc, Asyc (asynchronous composition), Sync (synchronous composition), When Asyc is selected, the codec only supports pcm.",
			},

			"codec": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Audio format, default wav (synchronous)/pcm (asynchronous, wav, mp3, pcm.",
			},

			"voice_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Timbre, the default value is ruxue.",
			},

			"volume": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Volume, default value 0, [-10,10].",
			},

			"speed": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Speech rate, the default value is 100, [50,200].",
			},
		},
	}
}

func resourceTencentCloudCiMediaTtsTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_tts_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaTtsTemplateOptions{
			Tag: "Tts",
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

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = v.(string)
	}

	if v, ok := d.GetOk("codec"); ok {
		request.Codec = v.(string)
	}

	if v, ok := d.GetOk("voice_type"); ok {
		request.VoiceType = v.(string)
	}

	if v, ok := d.GetOk("volume"); ok {
		request.Volume = v.(string)
	}

	if v, ok := d.GetOk("speed"); ok {
		request.Speed = v.(string)
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaTtsTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaTtsTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTtsTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaTtsTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTtsTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_tts_template.read")()
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

	if template.Name != "" {
		_ = d.Set("name", template.Name)
	}

	mediaTtsTemplate := template.TtsTpl
	if mediaTtsTemplate != nil {
		if mediaTtsTemplate.Mode != "" {
			_ = d.Set("mode", mediaTtsTemplate.Mode)
		}

		if mediaTtsTemplate.Codec != "" {
			_ = d.Set("codec", mediaTtsTemplate.Codec)
		}

		if mediaTtsTemplate.VoiceType != "" {
			_ = d.Set("voice_type", mediaTtsTemplate.VoiceType)
		}

		if mediaTtsTemplate.Volume != "" {
			_ = d.Set("volume", mediaTtsTemplate.Volume)
		}

		if mediaTtsTemplate.Speed != "" {
			_ = d.Set("speed", mediaTtsTemplate.Speed)
		}
	}
	return nil
}

func resourceTencentCloudCiMediaTtsTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_tts_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaTtsTemplateOptions{
		Tag: "Tts",
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("mode") {
		if v, ok := d.GetOk("mode"); ok {
			request.Mode = v.(string)
		}
	}

	if d.HasChange("codec") {
		if v, ok := d.GetOk("codec"); ok {
			request.Codec = v.(string)
		}
	}

	if d.HasChange("voice_type") {
		if v, ok := d.GetOk("voice_type"); ok {
			request.VoiceType = v.(string)
		}
	}

	if d.HasChange("volume") {
		if v, ok := d.GetOk("volume"); ok {
			request.Volume = v.(string)
		}
	}

	if d.HasChange("speed") {
		if v, ok := d.GetOk("speed"); ok {
			request.Speed = v.(string)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaTtsTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaTtsTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTtsTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaTtsTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTtsTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_tts_template.delete")()
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
