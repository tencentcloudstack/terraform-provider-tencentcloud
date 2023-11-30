package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeHostUpdateRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostUpdateRecordRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "New certificate ID.",
			},

			"old_certificate_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Original certificate ID.",
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
							Description: "Record ID.",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "New certificate ID.",
						},
						"old_cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original certificate ID.",
						},
						"resource_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of resource types.",
						},
						"regions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of regional deploymentNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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
							Description: "Last update time.",
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

func dataSourceTencentCloudSslDescribeHostUpdateRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_host_update_record.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("certificate_id"); ok {
		paramMap["CertificateId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("old_certificate_id"); ok {
		paramMap["OldCertificateId"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var deployRecordList []*ssl.UpdateRecordInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeHostUpdateRecordByFilter(ctx, paramMap)
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
		for _, updateRecordInfo := range deployRecordList {
			updateRecordInfoMap := map[string]interface{}{}

			if updateRecordInfo.Id != nil {
				updateRecordInfoMap["id"] = updateRecordInfo.Id
			}

			if updateRecordInfo.CertId != nil {
				updateRecordInfoMap["cert_id"] = updateRecordInfo.CertId
			}

			if updateRecordInfo.OldCertId != nil {
				updateRecordInfoMap["old_cert_id"] = updateRecordInfo.OldCertId
			}

			if updateRecordInfo.ResourceTypes != nil {
				updateRecordInfoMap["resource_types"] = updateRecordInfo.ResourceTypes
			}

			if updateRecordInfo.Regions != nil {
				updateRecordInfoMap["regions"] = updateRecordInfo.Regions
			}

			if updateRecordInfo.Status != nil {
				updateRecordInfoMap["status"] = updateRecordInfo.Status
			}

			if updateRecordInfo.CreateTime != nil {
				updateRecordInfoMap["create_time"] = updateRecordInfo.CreateTime
			}

			if updateRecordInfo.UpdateTime != nil {
				updateRecordInfoMap["update_time"] = updateRecordInfo.UpdateTime
			}

			ids = append(ids, helper.UInt64ToStr(*updateRecordInfo.Id))
			tmpList = append(tmpList, updateRecordInfoMap)
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
