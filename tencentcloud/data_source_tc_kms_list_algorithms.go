package tencentcloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

func dataSourceTencentCloudKmsListAlgorithms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsListAlgorithmsRead,
		Schema: map[string]*schema.Schema{
			"symmetric_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Symmetric encryption algorithms supported in this region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key usage.",
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm.",
						},
					},
				},
			},
			"asymmetric_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Asymmetric encryption algorithms supported in this region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key usage.",
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm.",
						},
					},
				},
			},
			"asymmetric_sign_verify_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Asymmetric signature verification algorithms supported in this region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key usage.",
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm.",
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

func dataSourceTencentCloudKmsListAlgorithmsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_list_algorithms.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		service        = KmsService{client: meta.(*TencentCloudClient).apiV3Conn}
		listAlgorithms *kms.ListAlgorithmsResponseParams
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsListAlgorithmsByFilter(ctx)
		if e != nil {
			return retryError(e)
		}

		listAlgorithms = result
		return nil
	})

	if err != nil {
		return err
	}

	if listAlgorithms.SymmetricAlgorithms != nil {
		tmpList := make([]map[string]interface{}, 0, len(listAlgorithms.SymmetricAlgorithms))
		for _, item := range listAlgorithms.SymmetricAlgorithms {
			itemMap := map[string]interface{}{}
			if item.KeyUsage != nil {
				itemMap["key_usage"] = item.KeyUsage
			}

			if item.Algorithm != nil {
				itemMap["algorithm"] = item.Algorithm
			}

			tmpList = append(tmpList, itemMap)
		}

		_ = d.Set("symmetric_algorithms", tmpList)
	}

	if listAlgorithms.AsymmetricAlgorithms != nil {
		tmpList := make([]map[string]interface{}, 0, len(listAlgorithms.AsymmetricAlgorithms))
		for _, item := range listAlgorithms.AsymmetricAlgorithms {
			itemMap := map[string]interface{}{}
			if item.KeyUsage != nil {
				itemMap["key_usage"] = item.KeyUsage
			}

			if item.Algorithm != nil {
				itemMap["algorithm"] = item.Algorithm
			}

			tmpList = append(tmpList, itemMap)
		}

		_ = d.Set("asymmetric_algorithms", tmpList)
	}

	if listAlgorithms.AsymmetricSignVerifyAlgorithms != nil {
		tmpList := make([]map[string]interface{}, 0, len(listAlgorithms.AsymmetricSignVerifyAlgorithms))
		for _, item := range listAlgorithms.AsymmetricSignVerifyAlgorithms {
			itemMap := map[string]interface{}{}
			if item.KeyUsage != nil {
				itemMap["key_usage"] = item.KeyUsage
			}

			if item.Algorithm != nil {
				itemMap["algorithm"] = item.Algorithm
			}

			tmpList = append(tmpList, itemMap)
		}

		_ = d.Set("asymmetric_sign_verify_algorithms", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
