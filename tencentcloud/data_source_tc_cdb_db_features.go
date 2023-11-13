/*
Use this data source to query detailed information of cdb db_features

Example Usage

```hcl
data "tencentcloud_cdb_db_features" "db_features" {
  instance_id = ""
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

func dataSourceTencentCloudCdbDbFeatures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbDbFeaturesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cdb-c1nl9rpv or cdbro-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.",
			},

			"is_support_audit": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to support the database audit function.",
			},

			"audit_need_upgrade": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable auditing needs to upgrade the kernel version.",
			},

			"is_support_encryption": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to support the database encryption function.",
			},

			"encryption_need_upgrade": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable encryption needs to upgrade the kernel version.",
			},

			"is_remote_ro": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether it is a remote read-only instance.",
			},

			"master_region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The region where the master instance is located.",
			},

			"is_support_update_sub_version": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to support minor version upgrades.",
			},

			"current_sub_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current kernel version.",
			},

			"target_sub_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Available kernel versions for upgrade.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbDbFeaturesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_db_features.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbDbFeaturesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		isSupportAudit = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(isSupportAudit))
	if isSupportAudit != nil {
		_ = d.Set("is_support_audit", isSupportAudit)
	}

	if auditNeedUpgrade != nil {
		_ = d.Set("audit_need_upgrade", auditNeedUpgrade)
	}

	if isSupportEncryption != nil {
		_ = d.Set("is_support_encryption", isSupportEncryption)
	}

	if encryptionNeedUpgrade != nil {
		_ = d.Set("encryption_need_upgrade", encryptionNeedUpgrade)
	}

	if isRemoteRo != nil {
		_ = d.Set("is_remote_ro", isRemoteRo)
	}

	if masterRegion != nil {
		_ = d.Set("master_region", masterRegion)
	}

	if isSupportUpdateSubVersion != nil {
		_ = d.Set("is_support_update_sub_version", isSupportUpdateSubVersion)
	}

	if currentSubVersion != nil {
		_ = d.Set("current_sub_version", currentSubVersion)
	}

	if targetSubVersion != nil {
		_ = d.Set("target_sub_version", targetSubVersion)
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
