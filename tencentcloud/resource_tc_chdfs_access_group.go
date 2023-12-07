package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudChdfsAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudChdfsAccessGroupCreate,
		Read:   resourceTencentCloudChdfsAccessGroupRead,
		Update: resourceTencentCloudChdfsAccessGroupUpdate,
		Delete: resourceTencentCloudChdfsAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Permission group name.",
			},

			"vpc_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "vpc network type(1:CVM, 2:BM 1.0).",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "VPC ID.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Permission group description, default empty.",
			},
		},
	}
}

func resourceTencentCloudChdfsAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_access_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = chdfs.NewCreateAccessGroupRequest()
		response      = chdfs.NewCreateAccessGroupResponse()
		accessGroupId string
	)
	if v, ok := d.GetOk("access_group_name"); ok {
		request.AccessGroupName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("vpc_type"); v != nil {
		request.VpcType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().CreateAccessGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create chdfs accessGroup failed, reason:%+v", logId, err)
		return err
	}

	accessGroupId = *response.Response.AccessGroup.AccessGroupId
	d.SetId(accessGroupId)

	return resourceTencentCloudChdfsAccessGroupRead(d, meta)
}

func resourceTencentCloudChdfsAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_access_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	accessGroupId := d.Id()

	accessGroup, err := service.DescribeChdfsAccessGroupById(ctx, accessGroupId)
	if err != nil {
		return err
	}

	if accessGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ChdfsAccessGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if accessGroup.AccessGroupName != nil {
		_ = d.Set("access_group_name", accessGroup.AccessGroupName)
	}

	if accessGroup.VpcType != nil {
		_ = d.Set("vpc_type", accessGroup.VpcType)
	}

	if accessGroup.VpcId != nil {
		_ = d.Set("vpc_id", accessGroup.VpcId)
	}

	if accessGroup.Description != nil {
		_ = d.Set("description", accessGroup.Description)
	}

	return nil
}

func resourceTencentCloudChdfsAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_access_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := chdfs.NewModifyAccessGroupRequest()

	accessGroupId := d.Id()

	request.AccessGroupId = &accessGroupId

	immutableArgs := []string{"vpc_type", "vpc_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("access_group_name") {
		if v, ok := d.GetOk("access_group_name"); ok {
			request.AccessGroupName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().ModifyAccessGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update chdfs accessGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudChdfsAccessGroupRead(d, meta)
}

func resourceTencentCloudChdfsAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_access_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}
	accessGroupId := d.Id()

	if err := service.DeleteChdfsAccessGroupById(ctx, accessGroupId); err != nil {
		return err
	}

	return nil
}
