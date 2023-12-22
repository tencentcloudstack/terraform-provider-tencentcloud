package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
)

func DataSourceTencentCloudMariadbDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDatabasesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"databases": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The database list of this instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
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

func dataSourceTencentCloudMariadbDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_databases.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var databases []*mariadb.Database

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDatabasesByFilter(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		databases = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(databases))
	if databases != nil {
		for _, database := range databases {
			databaseMap := map[string]interface{}{}

			if database.DbName != nil {
				databaseMap["db_name"] = database.DbName
			}

			tmpList = append(tmpList, databaseMap)
		}

		_ = d.Set("databases", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
