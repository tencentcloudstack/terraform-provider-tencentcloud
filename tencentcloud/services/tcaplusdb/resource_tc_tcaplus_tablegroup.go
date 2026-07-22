package tcaplusdb

import (
	"context"
	"errors"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTcaplusTableGroup() *schema.Resource {
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
			ValidateFunc: tccommon.ValidateStringLengthInRange(1, 30),
			Description:  "Name of the TcaplusDB table group. Name length should be between 1 and 30.",
		},
		"table_group_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ID of the TcaplusDB table group, can be user-specified (must be unique within the cluster) or auto-incremented by the API when not set. Immutable after creation.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_tablegroup.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		clusterId    = d.Get("cluster_id").(string)
		groupName    = d.Get("tablegroup_name").(string)
		tableGroupId = d.Get("table_group_id").(string)
	)
	groupId, err := tcaplusService.CreateGroup(ctx, clusterId, groupName, tableGroupId)
	if err != nil {
		return err
	}
	log.Printf("[CRUD] tcaplus_tablegroup create success, logId=%s, d.Id()=%s", logId, fmt.Sprintf("%s:%s", clusterId, groupId))
	d.SetId(fmt.Sprintf("%s:%s", clusterId, groupId))
	return resourceTencentCloudTcaplusTableGroupRead(d, meta)
}

func resourceTencentCloudTcaplusTableGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_tablegroup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Id()

	info, has, err := tcaplusService.DescribeGroup(ctx, clusterId, groupId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, err = tcaplusService.DescribeGroup(ctx, clusterId, groupId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		log.Printf("[CRUD] tcaplus_tablegroup id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if info.TableGroupId != nil {
		_ = d.Set("table_group_id", info.TableGroupId)
	}
	_ = d.Set("tablegroup_name", info.TableGroupName)
	_ = d.Set("table_count", int(*info.TableCount))
	_ = d.Set("total_size", int(*info.TotalSize))
	_ = d.Set("create_time", info.CreatedTime)

	return nil
}

func resourceTencentCloudTcaplusTableGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_tablegroup.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Id()

	immutableArgs := []string{"table_group_id"}
	for _, arg := range immutableArgs {
		if d.HasChange(arg) {
			return fmt.Errorf("tcaplus_tablegroup argument `%s` cannot be changed, it is immutable after creation. Please recreate the resource if you need a different value.", arg)
		}
	}

	if d.HasChange("tablegroup_name") {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyGroupName(ctx, clusterId, groupId, d.Get("tablegroup_name").(string))
			if err != nil {
				return tccommon.RetryError(err)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_tablegroup.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Id()

	err := tcaplusService.DeleteGroup(ctx, clusterId, groupId)
	if err != nil {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err = tcaplusService.DeleteGroup(ctx, clusterId, groupId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeGroup(ctx, clusterId, groupId)
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeGroup(ctx, clusterId, groupId)
			if err != nil {
				return tccommon.RetryError(err)
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
