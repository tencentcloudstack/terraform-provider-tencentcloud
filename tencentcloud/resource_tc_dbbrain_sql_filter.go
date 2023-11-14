/*
Provides a resource to create a dbbrain sql_filter

Example Usage

```hcl
resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = &lt;nil&gt;
  session_token {
		user = &lt;nil&gt;
		password = &lt;nil&gt;

  }
  sql_type = &lt;nil&gt;
  filter_key = &lt;nil&gt;
  max_concurrency = &lt;nil&gt;
  duration = &lt;nil&gt;
  product = &lt;nil&gt;
  status = &lt;nil&gt;
}
```

Import

dbbrain sql_filter can be imported using the id, e.g.

```
terraform import tencentcloud_dbbrain_sql_filter.sql_filter sql_filter_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudDbbrainSqlFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbbrainSqlFilterCreate,
		Read:   resourceTencentCloudDbbrainSqlFilterRead,
		Update: resourceTencentCloudDbbrainSqlFilterUpdate,
		Delete: resourceTencentCloudDbbrainSqlFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"session_token": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Session token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User name.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Password.",
						},
					},
				},
			},

			"sql_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sql type, optional value is SELECT, UPDATE, DELETE, INSERT, REPLACE.",
			},

			"filter_key": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Filter key.",
			},

			"max_concurrency": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Maxmum concurreny.",
			},

			"duration": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Filter duration.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Product, optional value is &amp;amp;#39;mysql&amp;amp;#39;, &amp;amp;#39;cynosdb&amp;amp;#39;.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter status.",
			},
		},
	}
}

func resourceTencentCloudDbbrainSqlFilterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_sql_filter.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dbbrain.NewCreateSqlFilterRequest()
		response   = dbbrain.NewCreateSqlFilterResponse()
		instanceId string
		filterId   int
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("session_token"); ok {
		request.SessionToken = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sql_type"); ok {
		request.SqlType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_key"); ok {
		request.FilterKey = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_concurrency"); ok {
		request.MaxConcurrency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("duration"); ok {
		request.Duration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().CreateSqlFilter(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dbbrain sqlFilter failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, helper.Int64ToStr(filterId)}, FILED_SP))

	return resourceTencentCloudDbbrainSqlFilterRead(d, meta)
}

func resourceTencentCloudDbbrainSqlFilterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_sql_filter.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	filterId := idSplit[1]

	sqlFilter, err := service.DescribeDbbrainSqlFilterById(ctx, instanceId, filterId)
	if err != nil {
		return err
	}

	if sqlFilter == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DbbrainSqlFilter` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if sqlFilter.InstanceId != nil {
		_ = d.Set("instance_id", sqlFilter.InstanceId)
	}

	if sqlFilter.SessionToken != nil {
	}

	if sqlFilter.SqlType != nil {
		_ = d.Set("sql_type", sqlFilter.SqlType)
	}

	if sqlFilter.FilterKey != nil {
		_ = d.Set("filter_key", sqlFilter.FilterKey)
	}

	if sqlFilter.MaxConcurrency != nil {
		_ = d.Set("max_concurrency", sqlFilter.MaxConcurrency)
	}

	if sqlFilter.Duration != nil {
		_ = d.Set("duration", sqlFilter.Duration)
	}

	if sqlFilter.Product != nil {
		_ = d.Set("product", sqlFilter.Product)
	}

	if sqlFilter.Status != nil {
		_ = d.Set("status", sqlFilter.Status)
	}

	return nil
}

func resourceTencentCloudDbbrainSqlFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_sql_filter.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dbbrain.NewModifySqlFiltersRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	filterId := idSplit[1]

	request.InstanceId = &instanceId
	request.FilterId = &filterId

	immutableArgs := []string{"instance_id", "session_token", "sql_type", "filter_key", "max_concurrency", "duration", "product", "status"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().ModifySqlFilters(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dbbrain sqlFilter failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDbbrainSqlFilterRead(d, meta)
}

func resourceTencentCloudDbbrainSqlFilterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_sql_filter.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	filterId := idSplit[1]

	if err := service.DeleteDbbrainSqlFilterById(ctx, instanceId, filterId); err != nil {
		return err
	}

	return nil
}
