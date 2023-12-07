package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlDatabasesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of single requests, the default value is 20, the minimum value is 1, and the maximum value is 100.",
			},

			"database_regexp": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Regular expression to match database library names.",
			},

			"items": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Returned instance information.",
			},

			"database_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Database name and character set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "The name of database.",
						},
						"character_set": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "character set type.",
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

func dataSourceTencentCloudMysqlDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_databases.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database_regexp"); ok {
		paramMap["DatabaseRegexp"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var databases *cdb.DescribeDatabasesResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlDatabasesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		databases = result
		return nil
	})
	if err != nil {
		return err
	}

	if databases.Items != nil {
		_ = d.Set("items", databases.Items)
	}

	ids := make([]string, 0, len(databases.DatabaseList))
	tmpList := make([]map[string]interface{}, 0, len(databases.DatabaseList))
	if databases.DatabaseList != nil {
		for _, databasesWithCharacterLists := range databases.DatabaseList {
			databasesWithCharacterListsMap := map[string]interface{}{}

			if databasesWithCharacterLists.DatabaseName != nil {
				databasesWithCharacterListsMap["database_name"] = databasesWithCharacterLists.DatabaseName
			}

			if databasesWithCharacterLists.CharacterSet != nil {
				databasesWithCharacterListsMap["character_set"] = databasesWithCharacterLists.CharacterSet
			}

			ids = append(ids, *databasesWithCharacterLists.DatabaseName)
			tmpList = append(tmpList, databasesWithCharacterListsMap)
		}

		_ = d.Set("database_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"items":         databases.Items,
			"database_list": tmpList,
		}); e != nil {
			return e
		}
	}
	return nil
}
