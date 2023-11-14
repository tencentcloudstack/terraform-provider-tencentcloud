/*
Provides a resource to create a cls export

Example Usage

```hcl
resource "tencentcloud_cls_export" "export" {
  topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  query = "* | select count(*) as count"
  from = 1607499107000
  order = "desc"
  format = "json"
}
```

Import

cls export can be imported using the id, e.g.

```
terraform import tencentcloud_cls_export.export export_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudClsExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsExportCreate,
		Read:   resourceTencentCloudClsExportRead,
		Update: resourceTencentCloudClsExportUpdate,
		Delete: resourceTencentCloudClsExportDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic id.",
			},

			"query": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Export query rules.",
			},

			"from": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Export start time.",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log export time sorting. desc or asc.",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log export format.",
			},
		},
	}
}

func resourceTencentCloudClsExportCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_export.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateExportRequest()
		response = cls.NewCreateExportResponse()
		topicId  string
		exportId string
	)
	if v, ok := d.GetOk("topic_id"); ok {
		topicId = v.(string)
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("query"); ok {
		request.Query = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("from"); ok {
		request.From = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("order"); ok {
		request.Order = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateExport(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls export failed, reason:%+v", logId, err)
		return err
	}

	topicId = *response.Response.TopicId
	d.SetId(strings.Join([]string{topicId, exportId}, FILED_SP))

	return resourceTencentCloudClsExportRead(d, meta)
}

func resourceTencentCloudClsExportRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_export.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	exportId := idSplit[1]

	export, err := service.DescribeClsExportById(ctx, topicId, exportId)
	if err != nil {
		return err
	}

	if export == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsExport` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if export.TopicId != nil {
		_ = d.Set("topic_id", export.TopicId)
	}

	if export.Query != nil {
		_ = d.Set("query", export.Query)
	}

	if export.From != nil {
		_ = d.Set("from", export.From)
	}

	if export.Order != nil {
		_ = d.Set("order", export.Order)
	}

	if export.Format != nil {
		_ = d.Set("format", export.Format)
	}

	return nil
}

func resourceTencentCloudClsExportUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_export.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"topic_id", "query", "from", "order", "format"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudClsExportRead(d, meta)
}

func resourceTencentCloudClsExportDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_export.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	exportId := idSplit[1]

	if err := service.DeleteClsExportById(ctx, topicId, exportId); err != nil {
		return err
	}

	return nil
}
