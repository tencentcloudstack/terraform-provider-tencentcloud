/*
Use this data source to query detailed information of mysql db_features

Example Usage

```hcl
data "tencentcloud_mysql_db_features" "db_features" {
  instance_id = "cdb-fitq5t9h"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlDbFeatures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlDbFeaturesRead,
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

func dataSourceTencentCloudMysqlDbFeaturesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_db_features.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var dbFeatures *cdb.DescribeDBFeaturesResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlDbFeaturesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		dbFeatures = result
		return nil
	})
	if err != nil {
		return err
	}

	if dbFeatures.IsSupportAudit != nil {
		_ = d.Set("is_support_audit", dbFeatures.IsSupportAudit)
	}

	if dbFeatures.AuditNeedUpgrade != nil {
		_ = d.Set("audit_need_upgrade", dbFeatures.AuditNeedUpgrade)
	}

	if dbFeatures.IsSupportEncryption != nil {
		_ = d.Set("is_support_encryption", dbFeatures.IsSupportEncryption)
	}

	if dbFeatures.EncryptionNeedUpgrade != nil {
		_ = d.Set("encryption_need_upgrade", dbFeatures.EncryptionNeedUpgrade)
	}

	if dbFeatures.IsRemoteRo != nil {
		_ = d.Set("is_remote_ro", dbFeatures.IsRemoteRo)
	}

	if dbFeatures.MasterRegion != nil {
		_ = d.Set("master_region", dbFeatures.MasterRegion)
	}

	if dbFeatures.IsSupportUpdateSubVersion != nil {
		_ = d.Set("is_support_update_sub_version", dbFeatures.IsSupportUpdateSubVersion)
	}

	if dbFeatures.CurrentSubVersion != nil {
		_ = d.Set("current_sub_version", dbFeatures.CurrentSubVersion)
	}

	if dbFeatures.TargetSubVersion != nil {
		_ = d.Set("target_sub_version", dbFeatures.TargetSubVersion)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
