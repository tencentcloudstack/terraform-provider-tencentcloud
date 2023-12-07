package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbAccountsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cloud database account information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User Name.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "From which host the user can log in (corresponding to the host field of MySQL users, UserName + Host uniquely identifies a user, in the form of IP, the IP segment ends with %; supports filling in %; if it is empty, it defaults to %).",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User remarks info.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time.",
						},
						"read_only": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Read-only flag, 0: No, 1: The SQL request of this account is preferentially executed on the standby machine, and the host is selected for execution when the standby machine is unavailable. 2: The standby machine is preferentially selected for execution, and the operation fails when the standby machine is unavailable.",
						},
						"delay_thresh": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If the standby machine delay exceeds the setting value of this parameter, the system will consider that the standby machine is faulty and recommend that the parameter value be greater than 10. This parameter takes effect when ReadOnly selects 1 and 2.",
						},
						"slave_const": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For read-only accounts, set the policy whether to fix the standby machine, 0: not fix the standby machine, that is, the standby machine will not disconnect from the client if it does not meet the conditions, the Proxy selects other available standby machines, 1: the standby machine will be disconnected if the conditions are not met, Make sure a connection is secured to the standby machine.",
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

func dataSourceTencentCloudDcdbAccountsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_accounts.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	dcdbService := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var dbAccountList []*dcdb.DBAccount
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := dcdbService.DescribeDcdbAccountsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		dbAccountList = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dcdb list failed, reason:%+v", logId, err)
		return err
	}

	retList := []interface{}{}
	if dbAccountList != nil {
		ids := make([]string, 0, len(dbAccountList))
		for _, dbA := range dbAccountList {
			listMap := map[string]interface{}{}
			if dbA.UserName != nil {
				listMap["user_name"] = dbA.UserName
			}
			if dbA.Host != nil {
				listMap["host"] = dbA.Host
			}
			if dbA.Description != nil {
				listMap["description"] = dbA.Description
			}
			if dbA.CreateTime != nil {
				listMap["create_time"] = dbA.CreateTime
			}
			if dbA.UpdateTime != nil {
				listMap["update_time"] = dbA.UpdateTime
			}
			if dbA.ReadOnly != nil {
				listMap["read_only"] = dbA.ReadOnly
			}
			if dbA.DelayThresh != nil {
				listMap["delay_thresh"] = dbA.DelayThresh
			}
			if dbA.SlaveConst != nil {
				listMap["slave_const"] = dbA.SlaveConst
			}
			ids = append(ids, *dbA.UserName+FILED_SP+*dbA.Host)
			retList = append(retList, listMap)
		}

		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", retList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), retList); e != nil {
			return e
		}
	}

	return nil
}
