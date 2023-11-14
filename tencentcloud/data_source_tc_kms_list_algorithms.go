/*
Use this data source to query detailed information of kms list_algorithms

Example Usage

```hcl
data "tencentcloud_kms_list_algorithms" "list_algorithms" {
      }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKmsListAlgorithms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsListAlgorithmsRead,
		Schema: map[string]*schema.Schema{
			"symmetric_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Symmetric encryption algorithms supported in this region.",
			},

			"asymmetric_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Asymmetric encryption algorithms supported in this region.",
			},

			"asymmetric_sign_verify_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Asymmetric signature verification algorithms supported in this region.",
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

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var symmetricAlgorithms []*kms.AlgorithmInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsListAlgorithmsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		symmetricAlgorithms = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(symmetricAlgorithms))
	tmpList := make([]map[string]interface{}, 0, len(symmetricAlgorithms))

	if symmetricAlgorithms != nil {
		_ = d.Set("symmetric_algorithms", symmetricAlgorithms)
	}

	if asymmetricAlgorithms != nil {
		_ = d.Set("asymmetric_algorithms", asymmetricAlgorithms)
	}

	if asymmetricSignVerifyAlgorithms != nil {
		_ = d.Set("asymmetric_sign_verify_algorithms", asymmetricSignVerifyAlgorithms)
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
