/*
Use this resource to create tcaplus table

Example Usage

```hcl
resource "tencentcloud_tcaplus_application" "test" {
  idl_type                 = "PROTO"
  app_name                 = "tf_tcaplus_app_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_zone" "zone" {
  app_id         = tencentcloud_tcaplus_application.test.id
  zone_name      = "tf_test_zone_name"
}

resource "tencentcloud_tcaplus_idl" "main" {
  app_id         = tencentcloud_tcaplus_application.test.id
  file_name      = "tf_idl_test_2"
  file_type      = "PROTO"
  file_ext_type  = "proto"
  file_content   = <<EOF
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
  app_id     = tencentcloud_tcaplus_application.test.id
  zone_id            = tencentcloud_tcaplus_zone.zone.id
  table_name         = "tb_online"
  table_type         = "GENERIC"
  description        = "test"
  idl_id             = tencentcloud_tcaplus_idl.main.id
  table_idl_type     = "PROTO"
  reserved_read_qps  = 1000
  reserved_write_qps = 20
  reserved_volume    = 1
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application of this table belongs.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone of this table belongs.",
			},
			"table_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of this table.",
			},
			"table_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_TABLE_TYPES),
				Description:  "Type of this table, Valid values are " + strings.Join(TCAPLUS_TABLE_TYPES, ",") + ".",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of this table.",
			},
			"idl_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Idl file for this table.",
			},
			"table_idl_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_TABLE_IDL_TYPES),
				Description:  "Type of this table idl, Valid values are " + strings.Join(TCAPLUS_TABLE_IDL_TYPES, ",") + ".",
			},
			"reserved_read_qps": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Table reserved read QPS.",
			},
			"reserved_write_qps": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Table reserved write QPS.",
			},
			"reserved_volume": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Table reserved capacity(GB).",
			},
			// Computed values.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the tcapplus table.",
			},
			"error": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Show if this table  create error.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of this table.",
			},
			"table_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of this table.",
			},
		},
	}
}

func resourceTencentCloudTcaplusTableCreate(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_tcaplus_table.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var tcaplusIdlId TcaplusIdlId

	if err := json.Unmarshal([]byte(d.Get("idl_id").(string)), &tcaplusIdlId); err != nil {
		return fmt.Errorf("field `idl_id` is illegal,%s", err.Error())
	}
	applicationId := d.Get("app_id").(string)
	zoneId := d.Get("zone_id").(string)
	tableName := d.Get("table_name").(string)
	tableType := d.Get("table_type").(string)
	description := d.Get("description").(string)
	tableIdlType := d.Get("table_idl_type").(string)
	reservedReadQps := int64(d.Get("reserved_read_qps").(int))
	reservedWriteQps := int64(d.Get("reserved_write_qps").(int))
	reservedVolume := int64(d.Get("reserved_volume").(int))

	taskId, tableInstanceId, err := tcaplusService.CreateTables(ctx,
		tcaplusIdlId,
		applicationId,
		zoneId,
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
		info, has, err := tcaplusService.DescribeTask(ctx, applicationId, taskId)
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
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationId := d.Get("app_id").(string)

	zoneId := d.Get("zone_id").(string)
	tableName := d.Get("table_name").(string)
	tableId := d.Id()

	d.Partial(true)

	//description
	if d.HasChange("description") {
		err := tcaplusService.ModifyTableMemo(ctx, applicationId, zoneId, tableId, tableName, d.Get("description").(string))
		if err != nil {

			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				err = tcaplusService.ModifyTableMemo(ctx, applicationId, zoneId, tableId, tableName, d.Get("description").(string))
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
		taskId, err := tcaplusService.ModifyTables(ctx, tcaplusIdlId, applicationId, zoneId, tableId, tableName, d.Get("table_idl_type").(string))
		if err != nil {
			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				taskId, err = tcaplusService.ModifyTables(ctx, tcaplusIdlId, applicationId, zoneId, tableId, tableName, d.Get("table_idl_type").(string))
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
			info, has, err := tcaplusService.DescribeTask(ctx, applicationId, taskId)
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
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationId := d.Get("app_id").(string)

	tableInfo, has, err := tcaplusService.DescribeTable(ctx, applicationId, d.Id())

	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			tableInfo, has, err = tcaplusService.DescribeTable(ctx, applicationId, d.Id())
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
	_ = d.Set("app_id", tableInfo.ApplicationId)
	_ = d.Set("zone_id", fmt.Sprintf("%s:%s", *tableInfo.ApplicationId, *tableInfo.LogicZoneId))
	_ = d.Set("table_name", tableInfo.TableName)
	_ = d.Set("table_type", tableInfo.TableType)
	_ = d.Set("description", tableInfo.Memo)
	_ = d.Set("table_idl_type", tableInfo.TableIdlType)
	_ = d.Set("reserved_read_qps", tableInfo.ReservedReadQps)
	_ = d.Set("reserved_write_qps", tableInfo.ReservedWriteQps)
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
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}
	applicationId := d.Get("app_id").(string)
	zoneId := d.Get("zone_id").(string)
	tableName := d.Get("table_name").(string)
	instanceTableId := d.Id()

	_, err := tcaplusService.DeleteTable(ctx, applicationId, zoneId, instanceTableId, tableName)

	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err = tcaplusService.DeleteTable(ctx, applicationId, zoneId, instanceTableId, tableName)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeTable(ctx, applicationId, instanceTableId)

	if err != nil || has {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeTable(ctx, applicationId, instanceTableId)
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

	taskId, err := tcaplusService.DeleteTable(ctx, applicationId, zoneId, instanceTableId, tableName)
	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			taskId, err = tcaplusService.DeleteTable(ctx, applicationId, zoneId, instanceTableId, tableName)
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
		info, has, err := tcaplusService.DescribeTask(ctx, applicationId, taskId)
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
