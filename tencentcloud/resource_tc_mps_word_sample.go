/*
Provides a resource to create a mps word_sample

Example Usage

```hcl
resource "tencentcloud_mps_word_sample" "word_sample" {
  usages =
  words {
		keyword = ""
		tags =

  }
}
```

Import

mps word_sample can be imported using the id, e.g.

```
terraform import tencentcloud_mps_word_sample.word_sample word_sample_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMpsWordSample() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsWordSampleCreate,
		Read:   resourceTencentCloudMpsWordSampleRead,
		Update: resourceTencentCloudMpsWordSampleUpdate,
		Delete: resourceTencentCloudMpsWordSampleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"usages": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "&amp;amp;lt;b&amp;amp;gt;Keyword usage. Valid values:&amp;amp;lt;/b&amp;amp;gt;1. Recognition.Ocr: OCR-based content recognition2. Recognition.Asr: ASR-based content recognition3. Review.Ocr: OCR-based inappropriate information recognition4. Review.Asr: ASR-based inappropriate information recognition&amp;amp;lt;b&amp;amp;gt;Valid values can also be:&amp;amp;lt;/b&amp;amp;gt;5. Recognition: ASR- and OCR-based content recognition; equivalent to 1+26. Review: ASR- and OCR-based inappropriate information recognition; equivalent to 3+47. All: ASR- and OCR-based content recognition and inappropriate information detection; equivalent to 1+2+3+4.",
			},

			"words": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Keyword. Array length limit: 100.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keyword": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Keyword. Length limit: 20 characters.",
						},
						"tags": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Keyword tag&amp;lt;li&amp;gt;Array length limit: 20 tags;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;Tag length limit: 128 characters.&amp;lt;/li&amp;gt;.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsWordSampleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_word_sample.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewCreateWordSamplesRequest()
		response = mps.NewCreateWordSamplesResponse()
		keyword  string
	)
	if v, ok := d.GetOk("usages"); ok {
		usagesSet := v.(*schema.Set).List()
		for i := range usagesSet {
			usages := usagesSet[i].(string)
			request.Usages = append(request.Usages, &usages)
		}
	}

	if v, ok := d.GetOk("words"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			aiSampleWordInfo := mps.AiSampleWordInfo{}
			if v, ok := dMap["keyword"]; ok {
				aiSampleWordInfo.Keyword = helper.String(v.(string))
			}
			if v, ok := dMap["tags"]; ok {
				tagsSet := v.(*schema.Set).List()
				for i := range tagsSet {
					tags := tagsSet[i].(string)
					aiSampleWordInfo.Tags = append(aiSampleWordInfo.Tags, &tags)
				}
			}
			request.Words = append(request.Words, &aiSampleWordInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateWordSamples(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps wordSample failed, reason:%+v", logId, err)
		return err
	}

	keyword = *response.Response.Keyword
	d.SetId(keyword)

	return resourceTencentCloudMpsWordSampleRead(d, meta)
}

func resourceTencentCloudMpsWordSampleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_word_sample.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	wordSampleId := d.Id()

	wordSample, err := service.DescribeMpsWordSampleById(ctx, keyword)
	if err != nil {
		return err
	}

	if wordSample == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsWordSample` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if wordSample.Usages != nil {
		_ = d.Set("usages", wordSample.Usages)
	}

	if wordSample.Words != nil {
		wordsList := []interface{}{}
		for _, words := range wordSample.Words {
			wordsMap := map[string]interface{}{}

			if wordSample.Words.Keyword != nil {
				wordsMap["keyword"] = wordSample.Words.Keyword
			}

			if wordSample.Words.Tags != nil {
				wordsMap["tags"] = wordSample.Words.Tags
			}

			wordsList = append(wordsList, wordsMap)
		}

		_ = d.Set("words", wordsList)

	}

	return nil
}

func resourceTencentCloudMpsWordSampleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_word_sample.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyWordSampleRequest()

	wordSampleId := d.Id()

	request.Keyword = &keyword

	immutableArgs := []string{"usages", "words"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("usages") {
		if v, ok := d.GetOk("usages"); ok {
			usagesSet := v.(*schema.Set).List()
			for i := range usagesSet {
				usages := usagesSet[i].(string)
				request.Usages = append(request.Usages, &usages)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyWordSample(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps wordSample failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsWordSampleRead(d, meta)
}

func resourceTencentCloudMpsWordSampleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_word_sample.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	wordSampleId := d.Id()

	if err := service.DeleteMpsWordSampleById(ctx, keyword); err != nil {
		return err
	}

	return nil
}
