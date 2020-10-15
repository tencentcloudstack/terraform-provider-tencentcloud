/*
Provides a resource to create a CLB target group.

Example Usage

```hcl
resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    port              = 33
}
```

Import

CLB target group can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group.test lbtg-3k3io0i0
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbTargetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTargetCreate,
		Read:   resourceTencentCloudClbTargetRead,
		Update: resourceTencentCloudClbTargetUpdate,
		Delete: resourceTencentCloudClbTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TF_target_group",
				Description: "Target group name.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				ForceNew:    true,
				Description: "VPC ID, default is based on the network.",
			},
			//computed
			"target_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target group ID.",
			},
		},
	}
}

func resourceTencentCloudClbTargetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group.create")()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		clbService      = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		vpcId           = d.Get("vpc_id").(string)
		targetGroupName = d.Get("target_group_name").(string)
		insAttachments  = make([]*clb.TargetGroupInstance, 0)
		targetGroupId   string
		err             error
	)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		targetGroupId, err = clbService.CreateTargetGroup(ctx, targetGroupName, vpcId, insAttachments)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(targetGroupId)

	return resourceTencentCloudClbTargetRead(d, meta)

}

func resourceTencentCloudClbTargetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		clbService = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		id         = d.Id()
	)
	filters := make(map[string]string)
	targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, id, filters)
	if err != nil {
		return err
	}
	if len(targetGroupInfos) < 1 {
		d.SetId("")
		return nil
	}
	_ = d.Set("target_group_name", targetGroupInfos[0].TargetGroupName)
	_ = d.Set("vpc_id", targetGroupInfos[0].VpcId)

	return nil
}

func resourceTencentCloudClbTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group.update")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		clbService    = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		targetGroupId = d.Id()
	)

	if d.HasChange("target_group_name") {
		targetGroupName := d.Get("target_group_name").(string)
		err := clbService.ModifyTargetGroup(ctx, targetGroupId, targetGroupName)
		if err != nil {
			return nil
		}
	}
	return resourceTencentCloudClbTargetRead(d, meta)
}

func resourceTencentCloudClbTargetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group.delete")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		clbService    = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		targetGroupId = d.Id()
	)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteTarget(ctx, targetGroupId)
		if e != nil {
			return retryError(e, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
