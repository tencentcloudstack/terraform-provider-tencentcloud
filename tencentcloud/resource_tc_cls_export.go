/*
Provides a resource to create a cls export

Example Usage

```hcl
resource "tencentcloud_cls_export" "export" {
  topic_id  = "7e34a3a7-635e-4da8-9005-88106c1fde69"
  log_count = 2
  query     = "select count(*) as count"
  from      = 1607499107000
  to        = 1607499108000
  order     = "desc"
  format    = "json"
}

```

Import

cls export can be imported using the id, e.g.

```
terraform import tencentcloud_cls_export.export topic_id#export_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsExportCreate,
		Read:   resourceTencentCloudClsExportRead,
		Delete: resourceTencentCloudClsExportDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "topic id.",
			},

			"query": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "export query rules.",
			},

			"log_count": {
				Required:    true,
				Type:        schema.TypeInt,
				ForceNew:    true,
				Description: "export amount of log.",
			},

			"from": {
				Required:    true,
				Type:        schema.TypeInt,
				ForceNew:    true,
				Description: "export start time.",
			},

			"to": {
				Required:    true,
				Type:        schema.TypeInt,
				ForceNew:    true,
				Description: "export end time.",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "log export time sorting. desc or asc.",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "log export format.",
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

	if v, ok := d.GetOkExists("log_count"); ok {
		request.Count = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("query"); ok {
		request.Query = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("from"); ok {
		request.From = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("to"); ok {
		request.To = helper.IntInt64(v.(int))
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

	exportId = *response.Response.ExportId
	d.SetId(topicId + FILED_SP + exportId)

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

	if export.Count != nil {
		_ = d.Set("log_count", export.Count)
	}

	if export.Query != nil {
		_ = d.Set("query", export.Query)
	}

	if export.From != nil {
		_ = d.Set("from", export.From)
	}

	if export.To != nil {
		_ = d.Set("to", export.To)
	}

	if export.Order != nil {
		_ = d.Set("order", export.Order)
	}

	if export.Format != nil {
		_ = d.Set("format", export.Format)
	}

	return nil
}

func resourceTencentCloudClsExportDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_export.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)

	exportId := idSplit[1]

	if err := service.DeleteClsExportById(ctx, exportId); err != nil {
		return err
	}

	return nil
}
