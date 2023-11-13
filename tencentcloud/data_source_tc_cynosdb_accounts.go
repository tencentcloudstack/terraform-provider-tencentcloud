/*
Use this data source to query detailed information of cynosdb accounts

Example Usage

```hcl
data "tencentcloud_cynosdb_accounts" "accounts" {
  cluster_id = "cynosdbmysql-on5xw0ni"
  account_names =
  db_type = "MYSQL"
  hosts =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbAccountsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"account_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of accounts that need to be filtered.",
			},

			"db_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database type, value range:&amp;amp;lt;li&amp;amp;gt;MYSQL&amp;amp;lt;/li&amp;amp;gt;This parameter has been deprecated.",
			},

			"hosts": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of accounts that need to be filtered.",
			},

			"account_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Database account list note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database account name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database account description.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Main engine.",
						},
						"max_user_connections": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of user connections.",
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

func dataSourceTencentCloudCynosdbAccountsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_accounts.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_names"); ok {
		accountNamesSet := v.(*schema.Set).List()
		paramMap["AccountNames"] = helper.InterfacesStringsPoint(accountNamesSet)
	}

	if v, ok := d.GetOk("db_type"); ok {
		paramMap["DbType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("hosts"); ok {
		hostsSet := v.(*schema.Set).List()
		paramMap["Hosts"] = helper.InterfacesStringsPoint(hostsSet)
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var accountSet []*cynosdb.Account

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbAccountsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		accountSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(accountSet))
	tmpList := make([]map[string]interface{}, 0, len(accountSet))

	if accountSet != nil {
		for _, account := range accountSet {
			accountMap := map[string]interface{}{}

			if account.AccountName != nil {
				accountMap["account_name"] = account.AccountName
			}

			if account.Description != nil {
				accountMap["description"] = account.Description
			}

			if account.CreateTime != nil {
				accountMap["create_time"] = account.CreateTime
			}

			if account.UpdateTime != nil {
				accountMap["update_time"] = account.UpdateTime
			}

			if account.Host != nil {
				accountMap["host"] = account.Host
			}

			if account.MaxUserConnections != nil {
				accountMap["max_user_connections"] = account.MaxUserConnections
			}

			ids = append(ids, *account.ClusterId)
			tmpList = append(tmpList, accountMap)
		}

		_ = d.Set("account_set", tmpList)
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
