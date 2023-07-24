/*
Use this resource to create TcaplusDB table group.

Example Usage

Create a tcaplusdb table group

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "your_pw_123111"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "example" {
  cluster_id      = tencentcloud_tcaplus_cluster.example.id
  tablegroup_name = "tf_example_group_name"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudTcaplusTableGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusTableGroupCreate,
		Read:   resourceTencentCloudTcaplusTableGroupRead,
		Update: resourceTencentCloudTcaplusTableGroupUpdate,
		Delete: resourceTencentCloudTcaplusTableGroupDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the TcaplusDB cluster to which the table group belongs.",
			},
			"tablegroup_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
				Description:  "Name of the TcaplusDB table group. Name length should be between 1 and 30.",
			},
			// Computed values.
			"table_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tables.",
			},
			"total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total storage size (MB).",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the TcaplusDB table group.",
			},
		},
	}
}

func resourceTencentCloudTcaplusTableGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_tablegroup.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		clusterId = d.Get("cluster_id").(string)
		groupName = d.Get("tablegroup_name").(string)
	)
	groupId, err := tcaplusService.CreateGroup(ctx, clusterId, groupName)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", clusterId, groupId))
	return resourceTencentCloudTcaplusTableGroupRead(d, meta)
}

func resourceTencentCloudTcaplusTableGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_tablegroup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Id()

	info, has, err := tcaplusService.DescribeGroup(ctx, clusterId, groupId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, err = tcaplusService.DescribeGroup(ctx, clusterId, groupId)
			if err != nil {
				return retryError(err)
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

	_ = d.Set("tablegroup_name", info.TableGroupName)
	_ = d.Set("table_count", int(*info.TableCount))
	_ = d.Set("total_size", int(*info.TotalSize))
	_ = d.Set("create_time", info.CreatedTime)

	return nil
}

func resourceTencentCloudTcaplusTableGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_tablegroup.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Id()

	if d.HasChange("tablegroup_name") {
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyGroupName(ctx, clusterId, groupId, d.Get("tablegroup_name").(string))
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return resourceTencentCloudTcaplusTableGroupRead(d, meta)
}

func resourceTencentCloudTcaplusTableGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_tablegroup.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Id()

	err := tcaplusService.DeleteGroup(ctx, clusterId, groupId)
	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err = tcaplusService.DeleteGroup(ctx, clusterId, groupId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeGroup(ctx, clusterId, groupId)
	if err != nil || has {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeGroup(ctx, clusterId, groupId)
			if err != nil {
				return retryError(err)
			}
			if has {
				err = fmt.Errorf("delete group fail, group still exist from sdk DescribeGroup")
				return resource.RetryableError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete group fail, group still exist from sdk DescribeGroup")
	}
}
