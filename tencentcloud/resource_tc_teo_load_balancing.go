/*
Provides a resource to create a teo load_balancing

Example Usage

```hcl
resource "tencentcloud_teo_load_balancing" "load_balancing" {
  zone_id = &lt;nil&gt;
    host = &lt;nil&gt;
  type = &lt;nil&gt;
  origin_group_id = &lt;nil&gt;
  backup_origin_group_id = &lt;nil&gt;
  t_t_l = &lt;nil&gt;
  status = &lt;nil&gt;
    }
```

Import

teo load_balancing can be imported using the id, e.g.

```
terraform import tencentcloud_teo_load_balancing.load_balancing load_balancing_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudTeoLoadBalancing() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoLoadBalancingCreate,
		Read:   resourceTencentCloudTeoLoadBalancingRead,
		Update: resourceTencentCloudTeoLoadBalancingUpdate,
		Delete: resourceTencentCloudTeoLoadBalancingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"load_balancing_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Load balancer instance ID.",
			},

			"host": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subdomain name. You can use @ to represent the root domain.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Proxy mode.- `dns_only`: Only DNS.- `proxied`: Enable proxy.",
			},

			"origin_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the origin group to use.",
			},

			"backup_origin_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the backup origin group to use.",
			},

			"t_t_l": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Indicates DNS TTL time when `Type` is dns_only.",
			},

			"status": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Status of the task. Valid values to set: `online`, `offline`. During status change, the status is `process`.",
			},

			"cname": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Schedules domain names. Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Last modification date.",
			},
		},
	}
}

func resourceTencentCloudTeoLoadBalancingCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_load_balancing.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = teo.NewCreateLoadBalancingRequest()
		response        = teo.NewCreateLoadBalancingResponse()
		zoneId          string
		loadBalancingId string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_group_id"); ok {
		request.OriginGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_origin_group_id"); ok {
		request.BackupOriginGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("t_t_l"); ok {
		request.TTL = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateLoadBalancing(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo loadBalancing failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, loadBalancingId}, FILED_SP))

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"online"}, 60*readRetryTimeout, time.Second, service.TeoLoadBalancingStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTeoLoadBalancingRead(d, meta)
}

func resourceTencentCloudTeoLoadBalancingRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_load_balancing.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	loadBalancingId := idSplit[1]

	loadBalancing, err := service.DescribeTeoLoadBalancingById(ctx, zoneId, loadBalancingId)
	if err != nil {
		return err
	}

	if loadBalancing == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoLoadBalancing` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if loadBalancing.ZoneId != nil {
		_ = d.Set("zone_id", loadBalancing.ZoneId)
	}

	if loadBalancing.LoadBalancingId != nil {
		_ = d.Set("load_balancing_id", loadBalancing.LoadBalancingId)
	}

	if loadBalancing.Host != nil {
		_ = d.Set("host", loadBalancing.Host)
	}

	if loadBalancing.Type != nil {
		_ = d.Set("type", loadBalancing.Type)
	}

	if loadBalancing.OriginGroupId != nil {
		_ = d.Set("origin_group_id", loadBalancing.OriginGroupId)
	}

	if loadBalancing.BackupOriginGroupId != nil {
		_ = d.Set("backup_origin_group_id", loadBalancing.BackupOriginGroupId)
	}

	if loadBalancing.TTL != nil {
		_ = d.Set("t_t_l", loadBalancing.TTL)
	}

	if loadBalancing.Status != nil {
		_ = d.Set("status", loadBalancing.Status)
	}

	if loadBalancing.Cname != nil {
		_ = d.Set("cname", loadBalancing.Cname)
	}

	if loadBalancing.UpdateTime != nil {
		_ = d.Set("update_time", loadBalancing.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoLoadBalancingUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_load_balancing.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyLoadBalancingRequest  = teo.NewModifyLoadBalancingRequest()
		modifyLoadBalancingResponse = teo.NewModifyLoadBalancingResponse()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	loadBalancingId := idSplit[1]

	request.ZoneId = &zoneId
	request.LoadBalancingId = &loadBalancingId

	immutableArgs := []string{"zone_id", "load_balancing_id", "host", "type", "origin_group_id", "backup_origin_group_id", "t_t_l", "status", "cname", "update_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("origin_group_id") {
		if v, ok := d.GetOk("origin_group_id"); ok {
			request.OriginGroupId = helper.String(v.(string))
		}
	}

	if d.HasChange("t_t_l") {
		if v, ok := d.GetOkExists("t_t_l"); ok {
			request.TTL = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyLoadBalancing(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo loadBalancing failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoLoadBalancingRead(d, meta)
}

func resourceTencentCloudTeoLoadBalancingDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_load_balancing.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	loadBalancingId := idSplit[1]

	if err := service.DeleteTeoLoadBalancingById(ctx, zoneId, loadBalancingId); err != nil {
		return err
	}

	return nil
}
