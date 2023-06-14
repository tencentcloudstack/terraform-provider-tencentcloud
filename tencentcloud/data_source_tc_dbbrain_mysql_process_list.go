/*
Use this data source to query detailed information of dbbrain mysql_process_list

Example Usage

```hcl
data "tencentcloud_dbbrain_mysql_process_list" "mysql_process_list" {
  instance_id = local.mysql_id
  product     = "mysql"
}
```
*/
package tencentcloud

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainMysqlProcessList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainMysqlProcessListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "thread ID, used to filter the thread list.",
			},

			"user": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The operating account name of the thread, used to filter the thread list.",
			},

			"host": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The operating host address of the thread, used to filter the thread list.",
			},

			"db": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The threads operations database, used to filter the thread list.",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The operational state of the thread, used to filter the thread list.",
			},

			"command": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The execution type of the thread, used to filter the thread list.",
			},

			"time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The minimum value of the operation duration of a thread, in seconds, used to filter the list of threads whose operation duration is longer than this value.",
			},

			"info": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The threads operation statement is used to filter the thread list.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values: `mysql` - cloud database MySQL; `cynosdb` - cloud database TDSQL-C for MySQL, the default is `mysql`.",
			},

			"process_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Live thread list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "thread ID.",
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operating account name of the thread.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operating host address of the thread.",
						},
						"db": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The thread that operates the database.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operational state of the thread.",
						},
						"command": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The execution type of the thread.",
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operation duration of the thread, in seconds.",
						},
						"info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operation statement for the thread.",
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

func dataSourceTencentCloudDbbrainMysqlProcessListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_mysql_process_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOkExists("id"); ok {
		paramMap["ID"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("user"); ok {
		paramMap["User"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		paramMap["Host"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db"); ok {
		paramMap["DB"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		paramMap["State"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command"); ok {
		paramMap["Command"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("time"); ok {
		paramMap["Time"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("info"); ok {
		paramMap["Info"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var processList []*dbbrain.MySqlProcess

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainMysqlProcessListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		processList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(processList))
	tmpList := make([]map[string]interface{}, 0, len(processList))

	if processList != nil {
		for _, mySqlProcess := range processList {
			mySqlProcessMap := map[string]interface{}{}

			if mySqlProcess.ID != nil {
				mySqlProcessMap["id"] = mySqlProcess.ID
			}

			if mySqlProcess.User != nil {
				mySqlProcessMap["user"] = mySqlProcess.User
			}

			if mySqlProcess.Host != nil {
				mySqlProcessMap["host"] = mySqlProcess.Host
			}

			if mySqlProcess.DB != nil {
				mySqlProcessMap["db"] = mySqlProcess.DB
			}

			if mySqlProcess.State != nil {
				mySqlProcessMap["state"] = mySqlProcess.State
			}

			if mySqlProcess.Command != nil {
				mySqlProcessMap["command"] = mySqlProcess.Command
			}

			if mySqlProcess.Time != nil {
				mySqlProcessMap["time"] = mySqlProcess.Time
			}

			if mySqlProcess.Info != nil {
				mySqlProcessMap["info"] = mySqlProcess.Info
			}

			ids = append(ids, strings.Join([]string{instanceId, *mySqlProcess.ID}, FILED_SP))
			tmpList = append(tmpList, mySqlProcessMap)
		}

		_ = d.Set("process_list", tmpList)
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
