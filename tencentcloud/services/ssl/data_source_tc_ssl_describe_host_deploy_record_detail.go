package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeHostDeployRecordDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostDeployRecordDetailRead,
		Schema: map[string]*schema.Schema{
			"deploy_record_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deployment record ID.",
			},

			"deploy_record_detail_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Certificate deployment record listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Deployment record details ID.",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment certificate ID.",
						},
						"old_cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original binding certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment example name.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment monitor IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"domains": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of deployment domain.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment monitoring protocolNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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
							Description: "Deployment record details Create time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment record details last update time.",
						},
						"listener_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Delicate monitor name.",
						},
						"sni_switch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to turn on SNI.",
						},
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "COS storage barrel nameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Named space nameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"secret_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secret nameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "portNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"env_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TCB environment IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"tcb_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployed TCB typeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployed TCB regionNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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

func dataSourceTencentCloudSslDescribeHostDeployRecordDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_host_deploy_record_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("deploy_record_id"); ok {
		paramMap["DeployRecordId"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var deployRecordDetailList []*ssl.DeployRecordDetail
	var successTotalCount, failedTotalCount, runningTotalCount *int64

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, successTotal, failedTotal, runningTotal, e := service.DescribeSslDescribeHostDeployRecordDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		deployRecordDetailList = result
		successTotalCount, failedTotalCount, runningTotalCount = successTotal, failedTotal, runningTotal
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(deployRecordDetailList))
	tmpList := make([]map[string]interface{}, 0, len(deployRecordDetailList))

	if deployRecordDetailList != nil {
		for _, deployRecordDetail := range deployRecordDetailList {
			deployRecordDetailMap := map[string]interface{}{}

			if deployRecordDetail.Id != nil {
				deployRecordDetailMap["id"] = deployRecordDetail.Id
			}

			if deployRecordDetail.CertId != nil {
				deployRecordDetailMap["cert_id"] = deployRecordDetail.CertId
			}

			if deployRecordDetail.OldCertId != nil {
				deployRecordDetailMap["old_cert_id"] = deployRecordDetail.OldCertId
			}

			if deployRecordDetail.InstanceId != nil {
				deployRecordDetailMap["instance_id"] = deployRecordDetail.InstanceId
			}

			if deployRecordDetail.InstanceName != nil {
				deployRecordDetailMap["instance_name"] = deployRecordDetail.InstanceName
			}

			if deployRecordDetail.ListenerId != nil {
				deployRecordDetailMap["listener_id"] = deployRecordDetail.ListenerId
			}

			if deployRecordDetail.Domains != nil {
				deployRecordDetailMap["domains"] = deployRecordDetail.Domains
			}

			if deployRecordDetail.Protocol != nil {
				deployRecordDetailMap["protocol"] = deployRecordDetail.Protocol
			}

			if deployRecordDetail.Status != nil {
				deployRecordDetailMap["status"] = deployRecordDetail.Status
			}

			if deployRecordDetail.ErrorMsg != nil {
				deployRecordDetailMap["error_msg"] = deployRecordDetail.ErrorMsg
			}

			if deployRecordDetail.CreateTime != nil {
				deployRecordDetailMap["create_time"] = deployRecordDetail.CreateTime
			}

			if deployRecordDetail.UpdateTime != nil {
				deployRecordDetailMap["update_time"] = deployRecordDetail.UpdateTime
			}

			if deployRecordDetail.ListenerName != nil {
				deployRecordDetailMap["listener_name"] = deployRecordDetail.ListenerName
			}

			if deployRecordDetail.SniSwitch != nil {
				deployRecordDetailMap["sni_switch"] = deployRecordDetail.SniSwitch
			}

			if deployRecordDetail.Bucket != nil {
				deployRecordDetailMap["bucket"] = deployRecordDetail.Bucket
			}

			if deployRecordDetail.Namespace != nil {
				deployRecordDetailMap["namespace"] = deployRecordDetail.Namespace
			}

			if deployRecordDetail.SecretName != nil {
				deployRecordDetailMap["secret_name"] = deployRecordDetail.SecretName
			}

			if deployRecordDetail.Port != nil {
				deployRecordDetailMap["port"] = deployRecordDetail.Port
			}

			if deployRecordDetail.EnvId != nil {
				deployRecordDetailMap["env_id"] = deployRecordDetail.EnvId
			}

			if deployRecordDetail.TCBType != nil {
				deployRecordDetailMap["tcb_type"] = deployRecordDetail.TCBType
			}

			if deployRecordDetail.Region != nil {
				deployRecordDetailMap["region"] = deployRecordDetail.Region
			}

			ids = append(ids, *deployRecordDetail.InstanceId)
			tmpList = append(tmpList, deployRecordDetailMap)
		}

		_ = d.Set("deploy_record_detail_list", tmpList)
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

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
