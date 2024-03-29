package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func DataSourceTencentCloudMysqlSupportedPrivileges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlSupportedPrivilegesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"global_supported_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Global permissions supported by the instance.",
			},

			"database_supported_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database permissions supported by the instance.",
			},

			"table_supported_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database table permissions supported by the instance.",
			},

			"column_supported_privileges": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The database column permissions supported by the instance.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMysqlSupportedPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_supported_privileges.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var supportedPrivileges *cdb.DescribeSupportedPrivilegesResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlSupportedPrivilegesById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		supportedPrivileges = result
		return nil
	})
	if err != nil {
		return err
	}

	if supportedPrivileges.GlobalSupportedPrivileges != nil {
		_ = d.Set("global_supported_privileges", supportedPrivileges.GlobalSupportedPrivileges)
	}

	if supportedPrivileges.DatabaseSupportedPrivileges != nil {
		_ = d.Set("database_supported_privileges", supportedPrivileges.DatabaseSupportedPrivileges)
	}

	if supportedPrivileges.TableSupportedPrivileges != nil {
		_ = d.Set("table_supported_privileges", supportedPrivileges.TableSupportedPrivileges)
	}

	if supportedPrivileges.ColumnSupportedPrivileges != nil {
		_ = d.Set("column_supported_privileges", supportedPrivileges.ColumnSupportedPrivileges)
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
