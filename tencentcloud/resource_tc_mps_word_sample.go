/*
Provides a resource to create a mps word_sample

Example Usage

```hcl
resource "tencentcloud_mps_word_sample" "word_sample" {
  usages = ["Recognition.Ocr","Review.Ocr","Review.Asr"]
  keyword = "tf_test_kw_1"
  tags = ["tags_1", "tags_2"]
}
```

Import

mps word_sample can be imported using the id, e.g.

```
terraform import tencentcloud_mps_word_sample.word_sample keyword
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "Keyword usage. Valid values: 1.`Recognition.Ocr`: OCR-based content recognition. 2.`Recognition.Asr`: ASR-based content recognition. 3.`Review.Ocr`: OCR-based inappropriate information recognition. 4.`Review.Asr`: ASR-based inappropriate information recognition.",
			},

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
				Description: "Keyword tag. Array length limit: 20 tags. Each tag length limit: 128 characters.",
			},
		},
	}
}

func resourceTencentCloudMpsWordSampleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_word_sample.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = mps.NewCreateWordSamplesRequest()
		// response = mps.NewCreateWordSamplesResponse()
		keyword string
	)
	if v, ok := d.GetOk("usages"); ok {
		usagesSet := v.(*schema.Set).List()
		for i := range usagesSet {
			if usagesSet[i] != nil {
				usages := usagesSet[i].(string)
				request.Usages = append(request.Usages, &usages)
			}
		}
	}

	aiSampleWordInfo := mps.AiSampleWordInfo{}
	if v, ok := d.GetOk("keyword"); ok {
		aiSampleWordInfo.Keyword = helper.String(v.(string))
		keyword = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.(*schema.Set).List()
		for i := range tagsSet {
			if tagsSet[i] != nil {
				tags := tagsSet[i].(string)
				aiSampleWordInfo.Tags = append(aiSampleWordInfo.Tags, &tags)
			}
		}
	}
	request.Words = append(request.Words, &aiSampleWordInfo)

	//if v, ok := d.GetOk("words"); ok {
	//	for _, item := range v.([]interface{}) {
	//		dMap := item.(map[string]interface{})
	//		aiSampleWordInfo := mps.AiSampleWordInfo{}
	//		if v, ok := dMap["keyword"]; ok {
	//			aiSampleWordInfo.Keyword = helper.String(v.(string))
	//			keywords = append(keywords, v.(string))
	//		}
	//		if v, ok := dMap["tags"]; ok {
	//			tagsSet := v.(*schema.Set).List()
	//			for i := range tagsSet {
	//				if tagsSet[i] != nil {
	//					tags := tagsSet[i].(string)
	//					aiSampleWordInfo.Tags = append(aiSampleWordInfo.Tags, &tags)
	//				}
	//			}
	//		}
	//		request.Words = append(request.Words, &aiSampleWordInfo)
	//	}
	//}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateWordSamples(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps wordSample failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(keyword)

	return resourceTencentCloudMpsWordSampleRead(d, meta)
}

func resourceTencentCloudMpsWordSampleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_word_sample.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	//keywords := strings.Split(d.Id(), FILED_SP)
	//if len(keywords) == 0 {
	//	return fmt.Errorf("id is broken,%s", d.Id())
	//}
	keyword := d.Id()

	wordSample, err := service.DescribeMpsWordSampleById(ctx, keyword)
	if err != nil {
		return err
	}

	if wordSample == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsWordSample` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if len(wordSample.UsageSet) > 0 {
		_ = d.Set("usages", helper.StringsInterfaces(wordSample.UsageSet))
	}

	if wordSample.Keyword != nil {
		_ = d.Set("keyword", wordSample.Keyword)
	}

	if len(wordSample.TagSet) > 0 {
		_ = d.Set("tags", helper.StringsInterfaces(wordSample.TagSet))
	}

	//wordsList := []interface{}{}
	//for _, kw := range keywords{
	//	for _, ws :=range wordSamples{
	//		wordsMap := map[string]interface{}{}
	//		if ws.Keyword!=nil && kw == *ws.Keyword{
	//			wordsMap["keyword"] = ws.Keyword
	//			if len(ws.TagSet)>0 {
	//				wordsMap["tags"] = ws.TagSet
	//			}
	//			wordsList = append(wordsList, wordsMap)
	//			break
	//		}
	//	}
	//}
	//_ = d.Set("words", wordsList)

	return nil
}

func resourceTencentCloudMpsWordSampleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_word_sample.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyWordSampleRequest()

	//keywords := strings.Split(d.Id(), FILED_SP)
	//if len(keywords) == 0 {
	//	return fmt.Errorf("id is broken,%s", d.Id())
	//}
	keyword := d.Id()

	request.Keyword = &keyword

	immutableArgs := []string{"keyword"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("usages") {
		if v, ok := d.GetOk("usages"); ok {
			usagesSet := v.(*schema.Set).List()
			for i := range usagesSet {
				if usagesSet[i] != nil {
					usages := usagesSet[i].(string)
					request.Usages = append(request.Usages, &usages)
				}
			}
		}
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tagSet := v.(*schema.Set).List()

			request.TagOperationInfo = &mps.AiSampleTagOperation{
				Type: helper.String("reset"),
				Tags: helper.InterfacesStringsPoint(tagSet),
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
	keyword := d.Id()

	if err := service.DeleteMpsWordSamplesById(ctx, []string{keyword}); err != nil {
		return err
	}

	return nil
}
