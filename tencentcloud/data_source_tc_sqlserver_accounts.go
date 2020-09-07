/*
Use this data source to query the list of SQL Server accounts.

Example Usage

```hcl
data "tencentcloud_sqlserver_accounts" "name" {
  instance_id           = "mssql-3cdq7kx5"
  name = "myaccount"
}

data "tencentcloud_sqlserver_accounts" "foo" {
  instance_id           = "mssql-3cdq7kx5"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverAccountsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SQL server instance ID that the account belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SQL server account to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of SQL Server account. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL server instance ID that the account belongs to.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SQL server account.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark of the SQL Server account.",
						},
						//computed
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the SQL Server account. 1 for creating, 2 for running, 3 for modifying, 4 for resetting password, -1 for deleting.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the SQL Server account.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last updated time of the SQL Server account.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverAccountsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_accounts.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Get("instance_id").(string)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	accounts, err := sqlserverService.DescribeSqlserverAccounts(ctx, instanceId)

	if err != nil {
		return fmt.Errorf("api[DescribeAccounts]fail, return %s", err.Error())
	}

	var list []map[string]interface{}
	var ids = make([]string, len(accounts))

	for _, item := range accounts {
		mapping := map[string]interface{}{
			"remark":      item.Remark,
			"name":        item.Name,
			"status":      item.Status,
			"create_time": item.CreateTime,
			"update_time": item.UpdateTime,
			"instance_id": instanceId,
		}
		if v, ok := d.GetOk("name"); ok && v.(string) != *item.Name {
			continue
		}
		list = append(list, mapping)
		ids = append(ids, fmt.Sprintf("%s%s%s", instanceId, FILED_SP, *item.Name))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}

	return nil
}
