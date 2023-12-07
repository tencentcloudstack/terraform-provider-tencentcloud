package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeHostDeployRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostDeployRecordRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID to be deployed.",
			},

			"resource_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resource Type.",
			},

			"deploy_record_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Certificate deployment record listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Deployment record ID.",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment certificate ID.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deploy resource type.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Deployment state.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Recent update time.",
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

func dataSourceTencentCloudSslDescribeHostDeployRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_host_deploy_record.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("certificate_id"); ok {
		paramMap["CertificateId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		paramMap["ResourceType"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var deployRecordList []*ssl.DeployRecordInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeHostDeployRecordByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		deployRecordList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(deployRecordList))
	tmpList := make([]map[string]interface{}, 0, len(deployRecordList))

	if deployRecordList != nil {
		for _, deployRecordInfo := range deployRecordList {
			deployRecordInfoMap := map[string]interface{}{}

			if deployRecordInfo.Id != nil {
				deployRecordInfoMap["id"] = deployRecordInfo.Id
			}

			if deployRecordInfo.CertId != nil {
				deployRecordInfoMap["cert_id"] = deployRecordInfo.CertId
			}

			if deployRecordInfo.ResourceType != nil {
				deployRecordInfoMap["resource_type"] = deployRecordInfo.ResourceType
			}

			if deployRecordInfo.Region != nil {
				deployRecordInfoMap["region"] = deployRecordInfo.Region
			}

			if deployRecordInfo.Status != nil {
				deployRecordInfoMap["status"] = deployRecordInfo.Status
			}

			if deployRecordInfo.CreateTime != nil {
				deployRecordInfoMap["create_time"] = deployRecordInfo.CreateTime
			}

			if deployRecordInfo.UpdateTime != nil {
				deployRecordInfoMap["update_time"] = deployRecordInfo.UpdateTime
			}

			ids = append(ids, helper.UInt64ToStr(*deployRecordInfo.Id))
			tmpList = append(tmpList, deployRecordInfoMap)
		}

		_ = d.Set("deploy_record_list", tmpList)
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
