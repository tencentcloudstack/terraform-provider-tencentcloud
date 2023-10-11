/*
Use this data source to query detailed information of ssl describe_host_update_record_detail

Example Usage

```hcl
data "tencentcloud_ssl_describe_host_update_record_detail" "describe_host_update_record_detail" {
  deploy_record_id = "35364"
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

func dataSourceTencentCloudSslDescribeHostUpdateRecordDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostUpdateRecordDetailRead,
		Schema: map[string]*schema.Schema{
			"deploy_record_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "One -click update record ID.",
			},

			"record_detail_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Certificate deployment record listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deploy resource type.",
						},
						"list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of deployment resources details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Detailed record ID.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "New certificate ID.",
									},
									"old_cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Old certificate ID.",
									},
									"domains": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "List of deployment domainNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy resource type.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DeploymentNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Deployment state.",
									},
									"error_msg": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deployment error messageNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deployment time.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last update time.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deployment instance IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deployment example nameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy listener ID (CLB for CLB)Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy listener name (CLB for CLB)Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "protocolNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"sni_switch": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether to turn on SNI (CLB dedicated)Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"bucket": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "BUCKET name (COS dedicated)Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "portNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Naming Space (TKE)Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"secret_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Secret Name (TKE for TKE)Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"env_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Environment IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"t_c_b_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TCB deployment typeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of deployment resources.",
						},
					},
				},
			},

			"success_total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total successNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"failed_total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of failuresNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"running_total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of deploymentNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeHostUpdateRecordDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_host_update_record_detail.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var id string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("deploy_record_id"); ok {
		id = v.(string)
		paramMap["DeployRecordId"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var recordDetailList []*ssl.UpdateRecordDetails
	var successTotalCount, failedTotalCount, runningTotalCount *int64
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, successTotal, failedTotal, runningTotal, e := service.DescribeSslDescribeHostUpdateRecordDetailByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		recordDetailList = result
		successTotalCount, failedTotalCount, runningTotalCount = successTotal, failedTotal, runningTotal
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(recordDetailList))

	if recordDetailList != nil {
		for _, updateRecordDetails := range recordDetailList {
			updateRecordDetailsMap := map[string]interface{}{}

			if updateRecordDetails.ResourceType != nil {
				updateRecordDetailsMap["resource_type"] = updateRecordDetails.ResourceType
			}

			if updateRecordDetails.List != nil {
				var listList []interface{}
				for _, list := range updateRecordDetails.List {
					listMap := map[string]interface{}{}

					if list.Id != nil {
						listMap["id"] = list.Id
					}

					if list.CertId != nil {
						listMap["cert_id"] = list.CertId
					}

					if list.OldCertId != nil {
						listMap["old_cert_id"] = list.OldCertId
					}

					if list.Domains != nil {
						listMap["domains"] = list.Domains
					}

					if list.ResourceType != nil {
						listMap["resource_type"] = list.ResourceType
					}

					if list.Region != nil {
						listMap["region"] = list.Region
					}

					if list.Status != nil {
						listMap["status"] = list.Status
					}

					if list.ErrorMsg != nil {
						listMap["error_msg"] = list.ErrorMsg
					}

					if list.CreateTime != nil {
						listMap["create_time"] = list.CreateTime
					}

					if list.UpdateTime != nil {
						listMap["update_time"] = list.UpdateTime
					}

					if list.InstanceId != nil {
						listMap["instance_id"] = list.InstanceId
					}

					if list.InstanceName != nil {
						listMap["instance_name"] = list.InstanceName
					}

					if list.ListenerId != nil {
						listMap["listener_id"] = list.ListenerId
					}

					if list.ListenerName != nil {
						listMap["listener_name"] = list.ListenerName
					}

					if list.Protocol != nil {
						listMap["protocol"] = list.Protocol
					}

					if list.SniSwitch != nil {
						listMap["sni_switch"] = list.SniSwitch
					}

					if list.Bucket != nil {
						listMap["bucket"] = list.Bucket
					}

					if list.Port != nil {
						listMap["port"] = list.Port
					}

					if list.Namespace != nil {
						listMap["namespace"] = list.Namespace
					}

					if list.SecretName != nil {
						listMap["secret_name"] = list.SecretName
					}

					if list.EnvId != nil {
						listMap["env_id"] = list.EnvId
					}

					if list.TCBType != nil {
						listMap["t_c_b_type"] = list.TCBType
					}

					listList = append(listList, listMap)
				}

				updateRecordDetailsMap["list"] = listList
			}

			if updateRecordDetails.TotalCount != nil {
				updateRecordDetailsMap["total_count"] = updateRecordDetails.TotalCount
			}

			tmpList = append(tmpList, updateRecordDetailsMap)
		}

		_ = d.Set("record_detail_list", tmpList)
	}

	if successTotalCount != nil {
		_ = d.Set("success_total_count", successTotalCount)
	}

	if failedTotalCount != nil {
		_ = d.Set("failed_total_count", failedTotalCount)
	}

	if runningTotalCount != nil {
		_ = d.Set("running_total_count", runningTotalCount)
	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
