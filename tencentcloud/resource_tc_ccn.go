/*
Provides a resource to create a CCN instance.

Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}
```

Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn.test ccn-id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCcn() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnCreate,
		Read:   resourceTencentCloudCcnRead,
		Update: resourceTencentCloudCcnUpdate,
		Delete: resourceTencentCloudCcnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the CCN to be queried, and maximum length does not exceed 60 bytes.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 100),
				Description:  "Description of CCN, and maximum length does not exceed 100 bytes.",
			},
			"qos": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      CNN_QOS_AU,
				ValidateFunc: validateAllowedStringValue([]string{CNN_QOS_PT, CNN_QOS_AU, CNN_QOS_AG}),
				Description:  "Service quality of CCN, and the available value include 'PT', 'AU', 'AG'. The default is 'AU'.",
			},
			// Computed values
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "States of instance. The available value include 'ISOLATED'(arrears) and 'AVAILABLE'.",
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of attached instances.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of resource.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tag.",
			},
		},
	}
}

func resourceTencentCloudCcnCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name        = d.Get("name").(string)
		description = ""
		qos         = d.Get("qos").(string)
	)
	if temp, ok := d.GetOk("description"); ok {
		description = temp.(string)
	}
	info, err := service.CreateCcn(ctx, name, description, qos)
	if err != nil {
		return err
	}
	d.SetId(info.ccnId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("vpc", "ccn", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCcnRead(d, meta)
}

func resourceTencentCloudCcnRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeCcn(ctx, d.Id())
		if e != nil {
			return retryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("name", info.name)
		_ = d.Set("description", info.description)
		_ = d.Set("qos", strings.ToUpper(info.qos))
		_ = d.Set("state", strings.ToUpper(info.state))
		_ = d.Set("instance_count", info.instanceCount)
		_ = d.Set("create_time", info.createTime)

		return nil
	})
	if err != nil {
		return err
	}
	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "ccn", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudCcnUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name        = ""
		description = ""
		change      = false
		changeList  = []string{}
	)
	if d.HasChange("name") {
		name = d.Get("name").(string)
		changeList = append(changeList, "name")
		change = true
	}

	if d.HasChange("description") {
		if temp, ok := d.GetOk("description"); ok {
			description = temp.(string)
		}
		if description == "" {
			return fmt.Errorf("can not set description='' ")
		}
		changeList = append(changeList, "description")
		change = true
	}

	d.Partial(true)
	if change {
		if err := service.ModifyCcnAttribute(ctx, d.Id(), name, description); err != nil {
			return err
		}
		for _, val := range changeList {
			d.SetPartial(val)
		}
	}

	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("vpc", "ccn", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceTencentCloudCcnRead(d, meta)
}

func resourceTencentCloudCcnDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, e := service.DescribeCcn(ctx, d.Id())
		if e != nil {
			return retryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err = service.DeleteCcn(ctx, d.Id()); err != nil {
		return err
	}

	return resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		_, has, err := service.DescribeCcn(ctx, d.Id())
		if err != nil {
			return resource.RetryableError(err)
		}
		if has == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}
