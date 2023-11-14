/*
Use this data source to query detailed information of cdb instance_info

Example Usage

```hcl
data "tencentcloud_cdb_instance_info" "instance_info" {
          }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbInstanceInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbInstanceInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance name.",
			},

			"encryption": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable encryption, YES is enabled, NO is not enabled.",
			},

			"key_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The key ID used for encryption.",
			},

			"key_region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The region where the key is located.",
			},

			"default_kms_region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The default region of the KMS service used by the current CDB backend service.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbInstanceInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_instance_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbInstanceInfoByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceId = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceId))
	if instanceName != nil {
		_ = d.Set("instance_name", instanceName)
	}

	if encryption != nil {
		_ = d.Set("encryption", encryption)
	}

	if keyId != nil {
		_ = d.Set("key_id", keyId)
	}

	if keyRegion != nil {
		_ = d.Set("key_region", keyRegion)
	}

	if defaultKmsRegion != nil {
		_ = d.Set("default_kms_region", defaultKmsRegion)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
