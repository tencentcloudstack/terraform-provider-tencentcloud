/*
Use this resource to create TcaplusDB table.

Example Usage

```hcl
resource "tencentcloud_tcaplus_cluster" "test" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_cluster_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "tablegroup" {
  cluster_id      = tencentcloud_tcaplus_cluster.test.id
  tablegroup_name = "tf_test_group_name"
}

resource "tencentcloud_tcaplus_idl" "main" {
  cluster_id    = tencentcloud_tcaplus_cluster.test.id
  tablegroup_id = tencentcloud_tcaplus_tablegroup.tablegroup.id
  file_name     = "tf_idl_test_2"
  file_type     = "PROTO"
  file_ext_type = "proto"
  file_content  = <<EOF
    syntax = "proto2";
    package myTcaplusTable;
    import "tcaplusservice.optionv1.proto";
    message tb_online {
       option(tcaplusservice.tcaplus_primary_key) = "uin,name,region";
        required int64 uin = 1;
        required string name = 2;
        required int32 region = 3;
        required int32 gamesvrid = 4;
        optional int32 logintime = 5 [default = 1];
        repeated int64 lockid = 6 [packed = true];
        optional bool is_available = 7 [default = false];
        optional pay_info pay = 8;
    }

    message pay_info {
        required int64 pay_id = 1;
        optional uint64 total_money = 2;
        optional uint64 pay_times = 3;
        optional pay_auth_info auth = 4;
        message pay_auth_info {
            required string pay_keys = 1;
            optional int64 update_time = 2;
        }
    }
    EOF
}

resource "tencentcloud_tcaplus_table" "table" {
  cluster_id        = tencentcloud_tcaplus_cluster.test.id
  tablegroup_id     = tencentcloud_tcaplus_tablegroup.tablegroup.id
  table_name        = "tb_online"
  table_type        = "GENERIC"
  description       = "test"
  idl_id            = tencentcloud_tcaplus_idl.main.id
  table_idl_type    = "PROTO"
  reserved_read_cu  = 1000
  reserved_write_cu = 20
  reserved_volume   = 1
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudTcaplusTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusTableCreate,
		Read:   resourceTencentCloudTcaplusTableRead,
		Update: resourceTencentCloudTcaplusTableUpdate,
		Delete: resourceTencentCloudTcaplusTableDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the TcaplusDB cluster to which the table belongs.",
			},
			"tablegroup_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the table group to which the table belongs.",
			},
			"table_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TcaplusDB table.",
			},
			"table_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_TABLE_TYPES),
				Description:  "Type of the TcaplusDB table. Valid values are `GENERIC` and `LIST`.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the TcaplusDB table.",
			},
			"idl_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the IDL File.",
			},
			"table_idl_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_TABLE_IDL_TYPES),
				Description:  "IDL type of the TcaplusDB table. Valid values: `PROTO` and `TDR`.",
			},
			"reserved_read_cu": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Reserved read capacity units of the TcaplusDB table.",
			},
			"reserved_write_cu": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Reserved write capacity units of the TcaplusDB table.",
			},
			"reserved_volume": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Reserved storage capacity of the TcaplusDB table (unit: GB).",
			},
			// Computed values.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the TcaplusDB table.",
			},
			"error": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Error messages for creating TcaplusDB table.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the TcaplusDB table.",
			},
			"table_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the TcaplusDB table.",
			},
		},
	}
}

func resourceTencentCloudTcaplusTableCreate(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_tcaplus_table.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var tcaplusIdlId TcaplusIdlId

	if err := json.Unmarshal([]byte(d.Get("idl_id").(string)), &tcaplusIdlId); err != nil {
		return fmt.Errorf("field `idl_id` is illegal,%s", err.Error())
	}
	clusterId := d.Get("cluster_id").(string)
	groupId := d.Get("tablegroup_id").(string)
	tableName := d.Get("table_name").(string)
	tableType := d.Get("table_type").(string)
	description := d.Get("description").(string)
	tableIdlType := d.Get("table_idl_type").(string)
	reservedReadQps := int64(d.Get("reserved_read_cu").(int))
	reservedWriteQps := int64(d.Get("reserved_write_cu").(int))
	reservedVolume := int64(d.Get("reserved_volume").(int))

	taskId, tableInstanceId, err := tcaplusService.CreateTables(ctx,
		tcaplusIdlId,
		clusterId,
		groupId,
		tableName,
		tableType,
		description,
		tableIdlType,
		reservedReadQps,
		reservedWriteQps,
		reservedVolume)

	if err != nil {
		return err
	}

	d.SetId(tableInstanceId)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		info, has, err := tcaplusService.DescribeTask(ctx, clusterId, taskId)
		if err != nil {
			return retryError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("create table task has been deleted"))
		}

		if *info.Progress == 100 {
			return nil
		}

		if *info.Progress >= 0 {
			return resource.RetryableError(fmt.Errorf("the table creation is in progress, and our wait has timed out"))
		}
		if *info.Progress < 0 {
			return resource.NonRetryableError(fmt.Errorf("TencentCloud SDK return %d task status,create table task failed", *info.Progress))
		}

		return nil
	})

	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	return resourceTencentCloudTcaplusTableRead(d, meta)
}

func resourceTencentCloudTcaplusTableUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_table.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)

	groupId := d.Get("tablegroup_id").(string)
	tableName := d.Get("table_name").(string)
	tableId := d.Id()

	d.Partial(true)

	//description
	if d.HasChange("description") {
		err := tcaplusService.ModifyTableMemo(ctx, clusterId, groupId, tableId, tableName, d.Get("description").(string))
		if err != nil {

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				err = tcaplusService.ModifyTableMemo(ctx, clusterId, groupId, tableId, tableName, d.Get("description").(string))
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		d.SetPartial("description")
	}

	//idl_id
	if d.HasChange("idl_id") || d.HasChange("table_name") || d.HasChange("table_idl_type") {
		var tcaplusIdlId TcaplusIdlId
		if err := json.Unmarshal([]byte(d.Get("idl_id").(string)), &tcaplusIdlId); err != nil {
			return fmt.Errorf("field `idl_id` is illegal,%s", err.Error())
		}
		taskId, err := tcaplusService.ModifyTables(ctx, tcaplusIdlId, clusterId, groupId, tableId, tableName, d.Get("table_idl_type").(string))
		if err != nil {
			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				taskId, err = tcaplusService.ModifyTables(ctx, tcaplusIdlId, clusterId, groupId, tableId, tableName, d.Get("table_idl_type").(string))
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			info, has, err := tcaplusService.DescribeTask(ctx, clusterId, taskId)
			if err != nil {
				return retryError(err)
			}
			if !has {
				return resource.NonRetryableError(fmt.Errorf("modify table idl task has been deleted"))
			}
			if *info.Progress == 100 {
				return nil
			}
			if *info.Progress >= 0 {
				return resource.RetryableError(fmt.Errorf("modify table idl is in progress, and our wait has timed out"))
			}
			if *info.Progress < 0 {
				return resource.NonRetryableError(fmt.Errorf("TencentCloud SDK return %d task status,modify table idl failed", *info.Progress))
			}
			return nil
		})

		if err != nil {
			return err
		}

		for _, key := range []string{"idl_id", "table_name", "table_idl_type"} {
			if d.HasChange(key) {
				d.SetPartial(key)
			}
		}
	}

	d.Partial(false)

	time.Sleep(time.Second)

	return resourceTencentCloudTcaplusTableRead(d, meta)
}

func resourceTencentCloudTcaplusTableRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_table.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)

	tableInfo, has, err := tcaplusService.DescribeTable(ctx, clusterId, d.Id())

	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			tableInfo, has, err = tcaplusService.DescribeTable(ctx, clusterId, d.Id())
			if err != nil {
				return retryError(err)
			}
			if !has {
				return nil
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("cluster_id", tableInfo.ClusterId)
	_ = d.Set("tablegroup_id", fmt.Sprintf("%s:%s", *tableInfo.ClusterId, *tableInfo.TableGroupId))
	_ = d.Set("table_name", tableInfo.TableName)
	_ = d.Set("table_type", tableInfo.TableType)
	_ = d.Set("description", tableInfo.Memo)
	_ = d.Set("table_idl_type", tableInfo.TableIdlType)
	_ = d.Set("reserved_read_cu", tableInfo.ReservedReadQps)
	_ = d.Set("reserved_write_cu", tableInfo.ReservedWriteQps)
	_ = d.Set("reserved_volume", tableInfo.ReservedVolume)
	_ = d.Set("create_time", tableInfo.CreatedTime)
	if tableInfo.Error != nil && tableInfo.Error.Message != nil {
		_ = d.Set("error", tableInfo.Error.Message)
	} else {
		_ = d.Set("error", "")
	}
	_ = d.Set("status", tableInfo.Status)
	_ = d.Set("table_size", tableInfo.TableSize)
	return nil
}

func resourceTencentCloudTcaplusTableDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_table.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}
	clusterId := d.Get("cluster_id").(string)
	groupId := d.Get("tablegroup_id").(string)
	tableName := d.Get("table_name").(string)
	instanceTableId := d.Id()

	_, err := tcaplusService.DeleteTable(ctx, clusterId, groupId, instanceTableId, tableName)

	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err = tcaplusService.DeleteTable(ctx, clusterId, groupId, instanceTableId, tableName)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeTable(ctx, clusterId, instanceTableId)

	if err != nil || has {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeTable(ctx, clusterId, instanceTableId)
			if err != nil {
				return retryError(err)
			}
			if has {
				err = fmt.Errorf("delete table fail, table still exist from sdk DescribeTable")
				return resource.RetryableError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	if has {
		return fmt.Errorf("delete table fail, table still exist from sdk DescribeTable")
	}

	taskId, err := tcaplusService.DeleteTable(ctx, clusterId, groupId, instanceTableId, tableName)
	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			taskId, err = tcaplusService.DeleteTable(ctx, clusterId, groupId, instanceTableId, tableName)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		info, has, err := tcaplusService.DescribeTask(ctx, clusterId, taskId)
		if err != nil {
			return retryError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("delete table task has been deleted"))
		}

		if *info.Progress == 100 {
			return nil
		}

		if *info.Progress >= 0 {
			return resource.RetryableError(fmt.Errorf("the table delete is in progress, and our wait has timed out"))
		}
		if *info.Progress < 0 {
			return resource.NonRetryableError(fmt.Errorf("TencentCloud SDK return %d task status,delete table task failed", *info.Progress))
		}

		return nil
	})

	return err
}
