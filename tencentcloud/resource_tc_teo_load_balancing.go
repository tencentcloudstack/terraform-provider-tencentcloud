/*
Provides a resource to create a teo load_balancing

Example Usage

```hcl
resource "tencentcloud_teo_load_balancing" "load_balancing" {
#  backup_origin_group_id = "origin-a499ca4b-3721-11ed-b9c1-5254005a52aa"
  host                   = "www.toutiao2.com"
  origin_group_id        = "origin-4f8a30b2-3720-11ed-b66b-525400dceb86"
  status                 = "online"
  tags                   = {}
  ttl                    = 600
  type                   = "proxied"
  zone_id                = "zone-297z8rf93cfw"
}

```
Import

teo load_balancing can be imported using the zone_id#loadBalancing_id, e.g.
```
$ terraform import tencentcloud_teo_load_balancing.load_balancing zone-297z8rf93cfw#lb-2a93c649-3719-11ed-b9c1-5254005a52aa
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
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoLoadBalancing() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoLoadBalancingRead,
		Create: resourceTencentCloudTeoLoadBalancingCreate,
		Update: resourceTencentCloudTeoLoadBalancingUpdate,
		Delete: resourceTencentCloudTeoLoadBalancingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"load_balancing_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer instance ID.",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subdomain name. You can use @ to represent the root domain.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Proxy mode.- `dns_only`: Only DNS.- `proxied`: Enable proxy.",
			},

			"origin_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the origin group to use.",
			},

			"backup_origin_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the backup origin group to use.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Indicates DNS TTL time when `Type` is dns_only.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status of the task. Valid values to set: `online`, `offline`. During status change, the status is `process`.",
			},

			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Schedules domain names. Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
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
		response        *teo.CreateLoadBalancingResponse
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
	} else {
		request.BackupOriginGroupId = helper.String("")
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateLoadBalancing(request)
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
		log.Printf("[CRITAL]%s create teo loadBalancing failed, reason:%+v", logId, err)
		return err
	}

	loadBalancingId = *response.Response.LoadBalancingId

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoLoadBalancing(ctx, zoneId, loadBalancingId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.Status == "online" || *instance.Status == "init" {
			return nil
		}
		if *instance.Status == "process" {
			return resource.RetryableError(fmt.Errorf("loadBalancing status is %v, operate failed.", *instance.Status))
		}
		return resource.RetryableError(fmt.Errorf("loadBalancing status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

	d.SetId(zoneId + FILED_SP + loadBalancingId)
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

	loadBalancing, err := service.DescribeTeoLoadBalancing(ctx, zoneId, loadBalancingId)

	if err != nil {
		return err
	}

	if loadBalancing == nil {
		d.SetId("")
		return fmt.Errorf("resource `loadBalancing` %s does not exist", loadBalancingId)
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
		_ = d.Set("ttl", loadBalancing.TTL)
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

	var (
		logId         = getLogId(contextNil)
		request       = teo.NewModifyLoadBalancingRequest()
		statusRequest = teo.NewModifyLoadBalancingStatusRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	loadBalancingId := idSplit[1]

	request.ZoneId = &zoneId
	request.LoadBalancingId = &loadBalancingId

	if d.HasChange("zone_id") {

		return fmt.Errorf("`zone_id` do not support change now.")

	}

	if d.HasChange("host") {

		return fmt.Errorf("`host` do not support change now.")

	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_group_id"); ok {
		request.OriginGroupId = helper.String(v.(string))
	}

	if d.HasChange("backup_origin_group_id") {
		if v, ok := d.GetOk("backup_origin_group_id"); ok {
			request.BackupOriginGroupId = helper.String(v.(string))
		}
	} else {
		request.BackupOriginGroupId = helper.String("")
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			request.TTL = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyLoadBalancing(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo loadBalancing failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("status") {
		statusRequest.ZoneId = &zoneId
		statusRequest.LoadBalancingId = &loadBalancingId
		if v, ok := d.GetOk("status"); ok {
			statusRequest.Status = helper.String(v.(string))
		}

		statusErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			statusResult, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyLoadBalancingStatus(statusRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, statusRequest.GetAction(), statusRequest.ToJsonString(), statusResult.ToJsonString())
			}
			return nil
		})

		if statusErr != nil {
			log.Printf("[CRITAL]%s create teo loadBalancing failed, reason:%+v", logId, statusErr)
			return statusErr
		}
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
	statusRequest := teo.NewModifyLoadBalancingStatusRequest()
	statusRequest.ZoneId = &zoneId
	statusRequest.LoadBalancingId = &loadBalancingId
	statusRequest.Status = helper.String("offline")
	statusErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		statusResult, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyLoadBalancingStatus(statusRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, statusRequest.GetAction(), statusRequest.ToJsonString(), statusResult.ToJsonString())
		}
		return nil
	})

	if statusErr != nil {
		log.Printf("[CRITAL]%s offline teo loadBalancing failed, reason:%+v", logId, statusErr)
		return statusErr
	}

	if err := service.DeleteTeoLoadBalancingById(ctx, zoneId, loadBalancingId); err != nil {
		return err
	}

	return nil
}
