/*
Use this data source to query detailed information of mps person_samples

Example Usage

```hcl
data "tencentcloud_mps_person_samples" "person_samples" {
  type = &lt;nil&gt;
  person_ids = &lt;nil&gt;
  names = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  total_count = &lt;nil&gt;
  person_set {
		person_id = &lt;nil&gt;
		name = &lt;nil&gt;
		description = &lt;nil&gt;
		face_info_set {
			face_id = &lt;nil&gt;
			url = &lt;nil&gt;
		}
		tag_set = &lt;nil&gt;
		usage_set = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;

  }
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

func dataSourceTencentCloudMpsPersonSamples() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsPersonSamplesRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The material type to pull, optional value:UserDefine: User-defined material library.Default: system default library.Default value: UserDefine, pull user-defined material library material.Note: If you are pulling the default material library of the system, you can only use the material name or material ID + material name to pull, and only one face image will be returned.",
			},

			"person_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Material ids, array length limit: 100.",
			},

			"names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Material names, array length limit: 20.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset, default: 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return the number of records, default value: 100, maximum value: 100.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"person_set": {
				Type:        schema.TypeList,
				Description: "Material information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"person_id": {
							Type:        schema.TypeString,
							Description: "Person id.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Person name.",
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Person description.",
						},
						"face_info_set": {
							Type:        schema.TypeList,
							Description: "Face information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"face_id": {
										Type:        schema.TypeString,
										Description: "Face id.",
									},
									"url": {
										Type:        schema.TypeString,
										Description: "Face url.",
									},
								},
							},
						},
						"tag_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Person tag.",
						},
						"usage_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Application Scenario.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
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

func dataSourceTencentCloudMpsPersonSamplesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_person_samples.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("person_ids"); ok {
		personIdsSet := v.(*schema.Set).List()
		paramMap["PersonIds"] = helper.InterfacesStringsPoint(personIdsSet)
	}

	if v, ok := d.GetOk("names"); ok {
		namesSet := v.(*schema.Set).List()
		paramMap["Names"] = helper.InterfacesStringsPoint(namesSet)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("person_set"); ok {
		personSetSet := v.([]interface{})
		tmpSet := make([]*mps.AiSamplePerson, 0, len(personSetSet))

		for _, item := range personSetSet {
			aiSamplePerson := mps.AiSamplePerson{}
			aiSamplePersonMap := item.(map[string]interface{})

			if v, ok := aiSamplePersonMap["person_id"]; ok {
				aiSamplePerson.PersonId = helper.String(v.(string))
			}
			if v, ok := aiSamplePersonMap["name"]; ok {
				aiSamplePerson.Name = helper.String(v.(string))
			}
			if v, ok := aiSamplePersonMap["description"]; ok {
				aiSamplePerson.Description = helper.String(v.(string))
			}
			if v, ok := aiSamplePersonMap["face_info_set"]; ok {
				for _, item := range v.([]interface{}) {
					faceInfoSetMap := item.(map[string]interface{})
					aiSampleFaceInfo := mps.AiSampleFaceInfo{}
					if v, ok := faceInfoSetMap["face_id"]; ok {
						aiSampleFaceInfo.FaceId = helper.String(v.(string))
					}
					if v, ok := faceInfoSetMap["url"]; ok {
						aiSampleFaceInfo.Url = helper.String(v.(string))
					}
					aiSamplePerson.FaceInfoSet = append(aiSamplePerson.FaceInfoSet, &aiSampleFaceInfo)
				}
			}
			if v, ok := aiSamplePersonMap["tag_set"]; ok {
				tagSetSet := v.(*schema.Set).List()
				aiSamplePerson.TagSet = helper.InterfacesStringsPoint(tagSetSet)
			}
			if v, ok := aiSamplePersonMap["usage_set"]; ok {
				usageSetSet := v.(*schema.Set).List()
				aiSamplePerson.UsageSet = helper.InterfacesStringsPoint(usageSetSet)
			}
			if v, ok := aiSamplePersonMap["create_time"]; ok {
				aiSamplePerson.CreateTime = helper.String(v.(string))
			}
			if v, ok := aiSamplePersonMap["update_time"]; ok {
				aiSamplePerson.UpdateTime = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &aiSamplePerson)
		}
		paramMap["person_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var personSet []*mps.AiSamplePerson

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsPersonSamplesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		personSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(personSet))
	tmpList := make([]map[string]interface{}, 0, len(personSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
