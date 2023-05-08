/*
Use this data source to query detailed information of cat probe data

Example Usage

```hcl
data "tencentcloud_cat_probe_data" "probe_data" {
  begin_time = 1667923200000
  end_time = 1667996208428
  task_type = "AnalyzeTaskType_Network"
  sort_field = "ProbeTime"
  ascending = true
  selected_fields = ["terraform"]
  offset = 0
  limit = 20
  task_id = ["task-knare1mk"]
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCatProbeData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatProbedataRead,
		Schema: map[string]*schema.Schema{
			"begin_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Start timestamp (in milliseconds).",
			},

			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "End timestamp (in milliseconds).",
			},

			"task_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task Type in AnalyzeTaskType_Network,AnalyzeTaskType_Browse,AnalyzeTaskType_UploadDownload,AnalyzeTaskType_Transport,AnalyzeTaskType_MediaStream.",
			},

			"sort_field": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Fields to be sorted ProbeTime dial test time sorting can be filled in You can also fill in the selected fields in SelectedFields.",
			},

			"ascending": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "true is Ascending.",
			},

			"selected_fields": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Selected Fields.",
			},

			"offset": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Offset.",
			},

			"limit": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Limit.",
			},

			"task_id": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "TaskID list.",
			},

			"operators": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Operators list.",
			},

			"districts": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Districts list.",
			},

			"error_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "ErrorTypes list.",
			},

			"city": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "City list.",
			},

			"code": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Code list.",
			},

			"detailed_single_data_define": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Probe node list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"probe_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Probe time.",
						},
						"labels": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Labels.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom Field Name/Description.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value.",
									},
								},
							},
						},
						"fields": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Fields.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom Field Name/Description.",
									},
									"value": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Value.",
									},
								},
							},
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

func dataSourceTencentCloudCatProbedataRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cat_probedata.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("begin_time"); v != nil {
		paramMap["begin_time"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["end_time"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("task_type"); v != nil {
		paramMap["task_type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_field"); ok {
		paramMap["sort_field"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ascending"); v != nil {
		paramMap["ascending"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("selected_fields"); ok {
		selectedFieldsSet := v.(*schema.Set).List()
		selectedFields := make([]*string, 0, len(selectedFieldsSet))
		for _, vv := range selectedFieldsSet {
			selectedFields = append(selectedFields, helper.String(vv.(string)))
		}
		paramMap["selected_fields"] = selectedFields
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("task_id"); ok {
		taskIdSet := v.(*schema.Set).List()
		taskId := make([]*string, 0, len(taskIdSet))
		for _, vv := range taskIdSet {
			taskId = append(taskId, helper.String(vv.(string)))
		}
		paramMap["task_id"] = taskId
	}

	if v, ok := d.GetOk("operators"); ok {
		operatorsSet := v.(*schema.Set).List()
		operators := make([]*string, 0, len(operatorsSet))
		for _, vv := range operatorsSet {
			operators = append(operators, helper.String(vv.(string)))
		}
		paramMap["operators"] = operators
	}

	if v, ok := d.GetOk("districts"); ok {
		districtsSet := v.(*schema.Set).List()
		districts := make([]*string, 0, len(districtsSet))
		for _, vv := range districtsSet {
			districts = append(districts, helper.String(vv.(string)))
		}
		paramMap["districts"] = districts
	}

	if v, ok := d.GetOk("error_types"); ok {
		errorTypesSet := v.(*schema.Set).List()
		errorTypes := make([]*string, 0, len(errorTypesSet))
		for _, vv := range errorTypesSet {
			errorTypes = append(errorTypes, helper.String(vv.(string)))
		}
		paramMap["error_types"] = errorTypes
	}

	if v, ok := d.GetOk("city"); ok {
		citySet := v.(*schema.Set).List()
		city := make([]*string, 0, len(citySet))
		for _, vv := range citySet {
			city = append(city, helper.String(vv.(string)))
		}
		paramMap["city"] = city
	}

	catService := CatService{client: meta.(*TencentCloudClient).apiV3Conn}

	var dataSets []*cat.DetailedSingleDataDefine
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := catService.DescribeCatProbeDataByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		dataSets = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Cat dataSet failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(dataSets))
	dataSetList := make([]map[string]interface{}, 0, len(dataSets))

	if dataSets != nil {
		for _, dataSet := range dataSets {
			dataSetMap := map[string]interface{}{}
			if dataSet.ProbeTime != nil {
				dataSetMap["probe_time"] = dataSet.ProbeTime
			}
			if dataSet.Labels != nil {
				labelsList := []interface{}{}
				for _, labels := range dataSet.Labels {
					labelsMap := map[string]interface{}{}
					if labels.ID != nil {
						labelsMap["id"] = labels.ID
					}
					if labels.Name != nil {
						labelsMap["name"] = labels.Name
					}
					if labels.Value != nil {
						labelsMap["value"] = labels.Value
					}

					labelsList = append(labelsList, labelsMap)
				}
				dataSetMap["labels"] = labelsList
			}
			if dataSet.Fields != nil {
				fieldsList := []interface{}{}
				for _, fields := range dataSet.Fields {
					fieldsMap := map[string]interface{}{}
					if fields.ID != nil {
						fieldsMap["id"] = fields.ID
					}
					if fields.Name != nil {
						fieldsMap["name"] = fields.Name
					}
					if fields.Value != nil {
						fieldsMap["value"] = fields.Value
					}

					fieldsList = append(fieldsList, fieldsMap)
				}
				dataSetMap["fields"] = fieldsList
			}
			ids = append(ids, helper.UInt64ToStr(*dataSet.ProbeTime))
			dataSetList = append(dataSetList, dataSetMap)
		}

		_ = d.Set("detailed_single_data_define", dataSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), dataSetList); e != nil {
			return e
		}
	}

	return nil
}
