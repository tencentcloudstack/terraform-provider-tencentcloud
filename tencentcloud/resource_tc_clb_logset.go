/*
Provides a resource to create an exclusive CLB Logset.

Example Usage

```hcl
resource "tencentcloud_clb_logset" "foo" {
  name = "clb_logset"
  perioid = 7
}
```

Import

CLB attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_logset.foo 4eb9e3a8-9c42-4b32-9ddf-e215e9c92764
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClbLogSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbLogSetCreate,
		Read:   resourceTencentCloudClbLogSetRead,
		Delete: resourceTencentCloudClbLogSetDelete,
		//Update: resourceTencentCloudClbLogSetUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Logset retention period in days. Maximun value is `90`.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Logset name, which unique and fixed `clb_logset` among all CLS logsets.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Logset creation time.",
			},
			"topic_count": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Number of log topics in logset.",
			},
		},
	}
}

func resourceTencentCloudClbLogSetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_logset.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	info, err := service.DescribeClsLogSetById(ctx, id)

	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return fmt.Errorf("resource `Logset` %s does not exist", id)
	}

	_ = d.Set("name", info.LogsetName)

	//
	//if _, ok := d.GetOk("period"); !ok {
	//	_ = d.Set("period", info)
	//}
	_ = d.Set("create_time", info.CreateTime)
	_ = d.Set("topic_count", info.TopicCount)

	return nil
}

func resourceTencentCloudClbLogSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_logset.create")()
	defer clbActionMu.Unlock()
	clbActionMu.Lock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		period = d.Get("period").(int)
	)

	// We're not support specify name and health logs for now
	id, err := service.CreateClbLogSet(ctx, "clb_logset", "", period)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudClbLogSetRead(d, meta)
}

// All fields are now Computed/ForceNew, means it does not support update
func resourceTencentCloudClbLogSetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_logset.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	request := cls.NewModifyLogsetRequest()

	request.LogsetId = helper.String(d.Id())

	if d.HasChange("name") {
		request.LogsetName = helper.String(d.Get("name").(string))
	}

	err := service.UpdateClsLogSet(ctx, request)

	if err != nil {
		return err
	}

	return resourceTencentCloudClbLogSetCreate(d, meta)
}

func resourceTencentCloudClbLogSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_logset.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteClsLogSet(ctx, id); err != nil {
		return err
	}

	return nil
}
