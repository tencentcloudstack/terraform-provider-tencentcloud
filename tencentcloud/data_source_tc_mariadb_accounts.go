package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbAccountsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "account list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "username.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host from which the user can log in (corresponding to the host field of MySQL users, UserName + Host uniquely identifies a user, in the form of IP, and the IP segment ends with %; supports filling in %; if it is empty, it defaults to %).",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User remarks.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"read_only": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Read-only flag, `0`: No, `1`: The SQL request of this account is preferentially executed on the standby machine, and the host machine is selected for execution when the standby machine is unavailable, `2`: The standby machine is preferentially selected for execution, and the operation fails when the standby machine is unavailable.",
						},
						"delay_thresh": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "This field is meaningful for read-only accounts, indicating that the standby machine with the active-standby delay less than this value is selected.",
						},
						"slave_const": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For read-only accounts, set whether the policy is to fix the standby machine, `0`: The standby machine is not fixed, that is, the standby machine does not meet the conditions and will not disconnect from the client, and the Proxy selects other available standby machines, `1`: The standby machine does not meet the conditions Disconnect, make sure one connection secures the standby.",
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

func dataSourceTencentCloudMariadbAccountsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_accounts.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var instanceId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["instance_id"] = helper.String(v.(string))
	}

	mariadbService := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var users []*mariadb.DBAccount
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := mariadbService.DescribeMariadbAccountsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		users = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Mariadb users failed, reason:%+v", logId, err)
		return err
	}

	userList := []interface{}{}
	if users != nil {
		for _, user := range users {
			userMap := map[string]interface{}{}
			if user.UserName != nil {
				userMap["user_name"] = user.UserName
			}
			if user.Host != nil {
				userMap["host"] = user.Host
			}
			if user.Description != nil {
				userMap["description"] = user.Description
			}
			if user.CreateTime != nil {
				userMap["create_time"] = user.CreateTime
			}
			if user.UpdateTime != nil {
				userMap["update_time"] = user.UpdateTime
			}
			if user.ReadOnly != nil {
				userMap["read_only"] = user.ReadOnly
			}
			if user.DelayThresh != nil {
				userMap["delay_thresh"] = user.DelayThresh
			}
			if user.SlaveConst != nil {
				userMap["slave_const"] = user.SlaveConst
			}

			userList = append(userList, userMap)
		}
		err = d.Set("list", userList)
		if err != nil {
			log.Printf("[CRITAL]%s provider set instances list fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}
	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), userList); e != nil {
			return e
		}
	}

	return nil
}
