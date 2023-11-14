/*
Use this data source to query detailed information of mps word_samples

Example Usage

```hcl
data "tencentcloud_mps_word_samples" "word_samples" {
  keywords =
  usages =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMpsWordSamples() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsWordSamplesRead,
		Schema: map[string]*schema.Schema{
			"keywords": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Keyword filter. Array length limit: 100 words.",
			},

			"usages": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "&amp;amp;lt;b&amp;amp;gt;Keyword usage. Valid values:&amp;amp;lt;/b&amp;amp;gt;1. Recognition.Ocr: OCR-based content recognition2. Recognition.Asr: ASR-based content recognition3. Review.Ocr: OCR-based inappropriate information recognition4. Review.Asr: ASR-based inappropriate information recognition&amp;amp;lt;b&amp;amp;gt;Valid values can also be:&amp;amp;lt;/b&amp;amp;gt;5. Recognition: ASR- and OCR-based content recognition; equivalent to 1+26. Review: ASR- and OCR-based inappropriate information recognition; equivalent to 3+4You can select multiple elements, which are connected by OR logic. If a usage contains any element in this parameter, the keyword sample will be used.",
			},

			"word_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Keyword information.Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keyword": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Keyword.",
						},
						"tag_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Keyword tag.",
						},
						"usage_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Keyword use case.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in [ISO date format](https://intl.cloud.tencent.com/document/product/266/11732?from_cn_redirect=1#iso-.E6.97.A5.E6.9C.9F.E6.A0.BC.E5.BC.8F).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time in [ISO date format](https://intl.cloud.tencent.com/document/product/266/11732?from_cn_redirect=1#iso-.E6.97.A5.E6.9C.9F.E6.A0.BC.E5.BC.8F).",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMpsWordSamplesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_word_samples.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("keywords"); ok {
		keywordsSet := v.(*schema.Set).List()
		paramMap["Keywords"] = helper.InterfacesStringsPoint(keywordsSet)
	}

	if v, ok := d.GetOk("usages"); ok {
		usagesSet := v.(*schema.Set).List()
		paramMap["Usages"] = helper.InterfacesStringsPoint(usagesSet)
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var wordSet []*mps.AiSampleWord

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsWordSamplesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		wordSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(wordSet))
	tmpList := make([]map[string]interface{}, 0, len(wordSet))

	if wordSet != nil {
		for _, aiSampleWord := range wordSet {
			aiSampleWordMap := map[string]interface{}{}

			if aiSampleWord.Keyword != nil {
				aiSampleWordMap["keyword"] = aiSampleWord.Keyword
			}

			if aiSampleWord.TagSet != nil {
				aiSampleWordMap["tag_set"] = aiSampleWord.TagSet
			}

			if aiSampleWord.UsageSet != nil {
				aiSampleWordMap["usage_set"] = aiSampleWord.UsageSet
			}

			if aiSampleWord.CreateTime != nil {
				aiSampleWordMap["create_time"] = aiSampleWord.CreateTime
			}

			if aiSampleWord.UpdateTime != nil {
				aiSampleWordMap["update_time"] = aiSampleWord.UpdateTime
			}

			ids = append(ids, *aiSampleWord.Keyword)
			tmpList = append(tmpList, aiSampleWordMap)
		}

		_ = d.Set("word_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
