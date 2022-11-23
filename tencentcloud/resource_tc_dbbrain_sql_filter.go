/*
Provides a resource to create a dbbrain sql_filter.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}
variable "region" {
  default = "ap-guangzhou"
}

data "tencentcloud_mysql_instance" "mysql" {
  instance_name = "instance_name"
}

locals {
  mysql_id = data.tencentcloud_mysql_instance.mysql.instance_list.0.mysql_id
}

resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = local.mysql_id
  session_token {
    user = "test"
	password = "===password==="
  }
  sql_type = "SELECT"
  filter_key = "filter_key"
  max_concurrency = 10
  duration = 3600
}

```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDbbrainSqlFilter() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDbbrainSqlFilterRead,
		Create: resourceTencentCloudDbbrainSqlFilterCreate,
		Update: resourceTencentCloudDbbrainSqlFilterUpdate,
		Delete: resourceTencentCloudDbbrainSqlFilterDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"session_token": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "session token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "user name.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "password.",
						},
					},
				},
			},

			"sql_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "sql type, optional value is SELECT, UPDATE, DELETE, INSERT, REPLACE.",
			},

			"filter_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "filter key.",
			},

			"max_concurrency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "maximum concurreny.",
			},

			"duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "filter duration.",
			},

			"product": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "product, optional value is &amp;#39;mysql&amp;#39;, &amp;#39;cynosdb&amp;#39;.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "filter status.",
			},
		},
	}
}
func resourceTencentCloudDbbrainSqlFilterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_sql_filter.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		request    = dbbrain.NewCreateSqlFilterRequest()
		response   *dbbrain.CreateSqlFilterResponse
		instanceId string
		filterId   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sql_type"); ok {
		request.SqlType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_key"); ok {
		request.FilterKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_concurrency"); ok {
		request.MaxConcurrency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("duration"); ok {
		request.Duration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	sessionToken, err := getSessionToken(d, meta, ctx)
	if err != nil {
		return err
	}
	if sessionToken != nil {
		request.SessionToken = sessionToken
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().CreateSqlFilter(request)
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
		log.Printf("[CRITAL]%s create dbbrain sqlFilter failed, reason:%+v", logId, err)
		return err
	}

	filterId = helper.Int64ToStr(*response.Response.FilterId)

	d.SetId(instanceId + FILED_SP + filterId)
	return resourceTencentCloudDbbrainSqlFilterRead(d, meta)
}

func getSessionToken(d *schema.ResourceData, meta interface{}, ctx context.Context) (sessionToken *string, errRet error) {
	var (
		logId      = getLogId(contextNil)
		instanceId *string
		product    *string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		product = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "session_token"); ok {
		var user *string
		var pw *string
		if v, ok := dMap["user"]; ok {
			user = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			pw = helper.String(v.(string))
		}

		service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}
		sessionToken, errRet = service.getSessionToken(ctx, instanceId, user, pw, product)

		if errRet != nil {
			return
		}

		log.Printf("[DEBUG]%s verify user account success, sessionToken [%s]\n", logId, *sessionToken)
	}
	return
}

func resourceTencentCloudDbbrainSqlFilterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_sql_filter.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}
		idSplit []string
	)

	idSplit = strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	filterId := idSplit[1]

	sqlFilter, err := service.DescribeDbbrainSqlFilter(ctx, helper.String(instanceId), helper.String(filterId))
	if err != nil {
		return err
	}

	if sqlFilter == nil {
		d.SetId("")
		return fmt.Errorf("resource `sqlFilter` %s does not exist", d.Id())
	}

	_ = d.Set("instance_id", instanceId)

	if sqlFilter.SqlType != nil {
		_ = d.Set("sql_type", sqlFilter.SqlType)
	}

	if sqlFilter.OriginKeys != nil {
		_ = d.Set("filter_key", sqlFilter.OriginKeys)
	}

	if sqlFilter.MaxConcurrency != nil {
		_ = d.Set("max_concurrency", sqlFilter.MaxConcurrency)
	}

	if sqlFilter.Status != nil {
		_ = d.Set("status", sqlFilter.Status)
	}

	return nil
}

func resourceTencentCloudDbbrainSqlFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_sql_filter.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		request = dbbrain.NewModifySqlFiltersRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := helper.String(idSplit[0])
	filterId := idSplit[1]
	request.InstanceId = instanceId
	request.FilterIds = []*int64{helper.StrToInt64Point(filterId)}

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("sql_type") {
		return fmt.Errorf("`sql_type` do not support change now.")
	}

	if d.HasChange("filter_key") {
		return fmt.Errorf("`filter_key` do not support change now.")
	}

	if d.HasChange("max_concurrency") {
		return fmt.Errorf("`max_concurrency` do not support change now.")
	}

	if d.HasChange("duration") {
		return fmt.Errorf("`duration` do not support change now.")
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	if d.HasChange("product") {
		if v, ok := d.GetOk("product"); ok {
			request.Product = helper.String(v.(string))
		}
	}

	sessionToken, err := getSessionToken(d, meta, ctx)
	if err != nil {
		return err
	}
	if sessionToken != nil {
		request.SessionToken = sessionToken
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().ModifySqlFilters(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify dbbrain sqlFilter failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDbbrainSqlFilterRead(d, meta)
}

func resourceTencentCloudDbbrainSqlFilterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_sql_filter.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := helper.String(idSplit[0])
	filterId := helper.String(idSplit[1])
	sessionToken, err := getSessionToken(d, meta, ctx)
	if err != nil {
		return err
	}

	if err := service.DeleteDbbrainSqlFilterById(ctx, instanceId, filterId, sessionToken); err != nil {
		return err
	}

	return nil
}
