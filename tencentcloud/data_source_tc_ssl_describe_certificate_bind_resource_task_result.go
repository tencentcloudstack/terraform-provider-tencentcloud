/*
Use this data source to query detailed information of ssl describe_certificate_bind_resource_task_result

Example Usage

```hcl
data "tencentcloud_ssl_describe_certificate_bind_resource_task_result" "describe_certificate_bind_resource_task_result" {
  task_ids =
  }
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeCertificateBindResourceTaskResult() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeCertificateBindResourceTaskResultRead,
		Schema: map[string]*schema.Schema{
			"task_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Task ID, query the results of binding cloud resources according to the task ID, support the maximum support of 100.",
			},

			"sync_task_bind_resource_result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of asynchronous tasks binding affiliated cloud resources resultsNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task ID.",
						},
						"bind_resource_result": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Related Cloud Resources ResultNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource types: CLB, CDN, Waf, LIVE, VOD, DDOS, TKE, Apigateway, TCB, Teo (Edgeone).",
									},
									"bind_resource_region_result": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Binding resource area results.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "areaNote: This field may return NULL, indicating that the valid value cannot be obtained.",
												},
												"total_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total number of related resources.",
												},
											},
										},
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Related Cloud Resources Inquiry results: 0 indicates that in the query, 1 means the query is successful.2 means the query is abnormal; if the status is 1, check the results of bindResourceResult; if the state is 2, check the reason for ERROR.",
						},
						"error": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated Cloud Resource Error InformationNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unusual error codeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unusual error messageNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
								},
							},
						},
						"cache_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current result cache time.",
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

func dataSourceTencentCloudSslDescribeCertificateBindResourceTaskResultRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_certificate_bind_resource_task_result.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("task_ids"); ok {
		taskIdsSet := v.(*schema.Set).List()
		paramMap["TaskIds"] = helper.InterfacesStringsPoint(taskIdsSet)
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var syncTaskBindResourceResult []*ssl.SyncTaskBindResourceResult

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeCertificateBindResourceTaskResultByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		syncTaskBindResourceResult = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(syncTaskBindResourceResult))
	tmpList := make([]map[string]interface{}, 0, len(syncTaskBindResourceResult))

	if syncTaskBindResourceResult != nil {
		for _, syncTaskBindResourceResultInc := range syncTaskBindResourceResult {
			syncTaskBindResourceResultMap := map[string]interface{}{}

			if syncTaskBindResourceResultInc.TaskId != nil {
				syncTaskBindResourceResultMap["task_id"] = syncTaskBindResourceResultInc.TaskId
			}

			if syncTaskBindResourceResultInc.BindResourceResult != nil {
				bindResourceResultList := []interface{}{}
				for _, bindResourceResult := range syncTaskBindResourceResultInc.BindResourceResult {
					bindResourceResultMap := map[string]interface{}{}

					if bindResourceResult.ResourceType != nil {
						bindResourceResultMap["resource_type"] = bindResourceResult.ResourceType
					}

					if bindResourceResult.BindResourceRegionResult != nil {
						bindResourceRegionResultList := []interface{}{}
						for _, bindResourceRegionResult := range bindResourceResult.BindResourceRegionResult {
							bindResourceRegionResultMap := map[string]interface{}{}

							if bindResourceRegionResult.Region != nil {
								bindResourceRegionResultMap["region"] = bindResourceRegionResult.Region
							}

							if bindResourceRegionResult.TotalCount != nil {
								bindResourceRegionResultMap["total_count"] = bindResourceRegionResult.TotalCount
							}

							bindResourceRegionResultList = append(bindResourceRegionResultList, bindResourceRegionResultMap)
						}

						bindResourceResultMap["bind_resource_region_result"] = []interface{}{bindResourceRegionResultList}
					}

					bindResourceResultList = append(bindResourceResultList, bindResourceResultMap)
				}

				syncTaskBindResourceResultMap["bind_resource_result"] = []interface{}{bindResourceResultList}
			}

			if syncTaskBindResourceResultInc.Status != nil {
				syncTaskBindResourceResultMap["status"] = syncTaskBindResourceResultInc.Status
			}

			if syncTaskBindResourceResultInc.Error != nil {
				errorMap := map[string]interface{}{}

				if syncTaskBindResourceResultInc.Error.Code != nil {
					errorMap["code"] = syncTaskBindResourceResultInc.Error.Code
				}

				if syncTaskBindResourceResultInc.Error.Message != nil {
					errorMap["message"] = syncTaskBindResourceResultInc.Error.Message
				}

				syncTaskBindResourceResultMap["error"] = []interface{}{errorMap}
			}

			if syncTaskBindResourceResultInc.CacheTime != nil {
				syncTaskBindResourceResultMap["cache_time"] = syncTaskBindResourceResultInc.CacheTime
			}

			ids = append(ids, *syncTaskBindResourceResultInc.TaskId)
			tmpList = append(tmpList, syncTaskBindResourceResultMap)
		}

		_ = d.Set("sync_task_bind_resource_result", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output3, ok := d.GetOk("result_output_file")
	if ok && output3.(string) != "" {
		if e := writeToFile(output3.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
