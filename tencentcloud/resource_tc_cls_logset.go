/*
Provides a resource to create a cls logset

Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "logset"
  tags        = {
    "test" = "test"
  }
}

```

Import

cls logset can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_logset.logset 5cd3a17e-fb0b-418c-afd7-77b365397426
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsLogset() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsLogsetCreate,
		Read:   resourceTencentCloudClsLogsetRead,
		Delete: resourceTencentCloudClsLogsetDelete,
		Update: resourceTencentCloudClsLogsetUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"logset_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Logset name, which must be unique.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list. Up to 10 tag key-value pairs are supported and must be unique.",
			},

			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"topic_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of log topics in logset.",
			},
			"role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If AssumerUin is not empty, it indicates the service provider who creates the logset.",
			},
		},
	}
}

func resourceTencentCloudClsLogsetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_logset.create")()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateLogsetRequest()
		response *cls.CreateLogsetResponse
	)

	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		request.Tags = make([]*cls.Tag, 0, len(tags))
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &key,
				Value: &value,
			})
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateLogset(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls logset failed, reason:%+v", logId, err)
		return err
	}

	id := *response.Response.LogsetId
	d.SetId(id)
	return resourceTencentCloudClsLogsetRead(d, meta)
}

func resourceTencentCloudClsLogsetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_logset.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	logSet, err := service.DescribeClsLogsetById(ctx, id)

	if err != nil {
		return err
	}

	if logSet == nil {
		d.SetId("")
		return fmt.Errorf("resource `Logset` %s does not exist", id)
	}

	_ = d.Set("logset_name", logSet.LogsetName)

	tags := make(map[string]string, len(logSet.Tags))
	for _, tag := range logSet.Tags {
		tags[*tag.Key] = *tag.Value
	}
	_ = d.Set("tags", tags)

	_ = d.Set("create_time", logSet.CreateTime)
	_ = d.Set("topic_count", logSet.TopicCount)
	_ = d.Set("role_name", logSet.RoleName)

	return nil
}

func resourceTencentCloudClsLogsetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_logset.update")()
	logId := getLogId(contextNil)
	request := cls.NewModifyLogsetRequest()

	request.LogsetId = helper.String(d.Id())

	if d.HasChange("logset_name") || d.HasChange("tags") {
		request.LogsetName = helper.String(d.Get("logset_name").(string))

		tags := d.Get("tags").(map[string]interface{})
		request.Tags = make([]*cls.Tag, 0, len(tags))
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &key,
				Value: helper.String(value.(string)),
			})
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyLogset(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsLogsetRead(d, meta)
}

func resourceTencentCloudClsLogsetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_logset.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteClsLogset(ctx, id); err != nil {
		return err
	}

	return nil
}
