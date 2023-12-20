package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func DataSourceTencentCloudMysqlInstanceInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"instance_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "instance name.",
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

func dataSourceTencentCloudMysqlInstanceInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_instance_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var instanceInfo *cdb.DescribeDBInstanceInfoResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlInstanceInfoById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	if instanceInfo.InstanceName != nil {
		_ = d.Set("instance_name", instanceInfo.InstanceName)
	}

	if instanceInfo.Encryption != nil {
		_ = d.Set("encryption", instanceInfo.Encryption)
	}

	if instanceInfo.KeyId != nil {
		_ = d.Set("key_id", instanceInfo.KeyId)
	}

	if instanceInfo.KeyRegion != nil {
		_ = d.Set("key_region", instanceInfo.KeyRegion)
	}

	if instanceInfo.DefaultKmsRegion != nil {
		_ = d.Set("default_kms_region", instanceInfo.DefaultKmsRegion)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
