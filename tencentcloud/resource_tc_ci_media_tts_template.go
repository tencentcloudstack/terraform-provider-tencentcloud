/*
Provides a resource to create a ci media_tts_template

Example Usage

```hcl
resource "tencentcloud_ci_media_tts_template" "media_tts_template" {
  name = &lt;nil&gt;
  mode = &lt;nil&gt;
  codec = &lt;nil&gt;
  voice_type = &lt;nil&gt;
  volume = &lt;nil&gt;
  speed = &lt;nil&gt;
}
```

Import

ci media_tts_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_tts_template.media_tts_template media_tts_template_id
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

	var (
		request    = ci.NewCreateMediaTtsTemplateRequest()
		response   = ci.NewCreateMediaTtsTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("codec"); ok {
		request.Codec = helper.String(v.(string))
	}

	if v, ok := d.GetOk("voice_type"); ok {
		request.VoiceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("volume"); ok {
		request.Volume = helper.String(v.(string))
	}

	if v, ok := d.GetOk("speed"); ok {
		request.Speed = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaTtsTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTtsTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaTtsTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTtsTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_tts_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaTtsTemplateId := d.Id()

	mediaTtsTemplate, err := service.DescribeCiMediaTtsTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaTtsTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaTtsTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaTtsTemplate.Name != nil {
		_ = d.Set("name", mediaTtsTemplate.Name)
	}

	if mediaTtsTemplate.Mode != nil {
		_ = d.Set("mode", mediaTtsTemplate.Mode)
	}

	if mediaTtsTemplate.Codec != nil {
		_ = d.Set("codec", mediaTtsTemplate.Codec)
	}

	if mediaTtsTemplate.VoiceType != nil {
		_ = d.Set("voice_type", mediaTtsTemplate.VoiceType)
	}

	if mediaTtsTemplate.Volume != nil {
		_ = d.Set("volume", mediaTtsTemplate.Volume)
	}

	if mediaTtsTemplate.Speed != nil {
		_ = d.Set("speed", mediaTtsTemplate.Speed)
	}

	return nil
}

func resourceTencentCloudCiMediaTtsTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_tts_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaTtsTemplateRequest()

	mediaTtsTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "mode", "codec", "voice_type", "volume", "speed"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaTtsTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaTtsTemplate failed, reason:%+v", logId, err)
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
	mediaTtsTemplateId := d.Id()

	if err := service.DeleteCiMediaTtsTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
