/*
Provides a resource to create a teo application_proxy

Example Usage

```hcl
resource "tencentcloud_teo_application_proxy" "application_proxy" {
  zone_id   = tencentcloud_teo_zone.zone.id
  zone_name = "sfurnace.work"

  accelerate_type      = 1
  security_type        = 1
  plat_type            = "domain"
  proxy_name           = "www.sfurnace.work"
  proxy_type           = "hostname"
  session_persist_time = 2400
}

```
Import

teo application_proxy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_application_proxy.application_proxy zoneId#proxyId
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
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoApplicationProxy() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoApplicationProxyRead,
		Create: resourceTencentCloudTeoApplicationProxyCreate,
		Update: resourceTencentCloudTeoApplicationProxyUpdate,
		Delete: resourceTencentCloudTeoApplicationProxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site name.",
			},

			"proxy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Layer-4 proxy name.",
			},

			"plat_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Scheduling mode.- ip: Anycast IP.- domain: CNAME.",
			},

			"security_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "- 0: Disable security protection.- 1: Enable security protection.",
			},

			"accelerate_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "- 0: Disable acceleration.- 1: Enable acceleration.",
			},

			"session_persist_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Session persistence duration. Value range: 30-3600 (in seconds).",
			},

			"proxy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how a layer-4 proxy is created.- hostname: Subdomain name.- instance: Instance.",
			},

			"proxy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Proxy ID.",
			},

			"schedule_value": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Scheduling information.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification date.",
			},

			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the layer-7 domain name.",
			},
		},
	}
}

func resourceTencentCloudTeoApplicationProxyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_application_proxy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateApplicationProxyRequest()
		response *teo.CreateApplicationProxyResponse
		zoneId   string
		proxyId  string
	)

	request.ForwardClientIp = helper.String("")
	request.SessionPersist = helper.Bool(true)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_name"); ok {
		request.ZoneName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_name"); ok {
		request.ProxyName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plat_type"); ok {
		request.PlatType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_type"); ok {
		request.SecurityType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("accelerate_type"); ok {
		request.AccelerateType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("session_persist_time"); ok {
		request.SessionPersistTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("proxy_type"); ok {
		request.ProxyType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateApplicationProxy(request)
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
		log.Printf("[CRITAL]%s create teo applicationProxy failed, reason:%+v", logId, err)
		return err
	}

	proxyId = *response.Response.ProxyId

	d.SetId(zoneId + FILED_SP + proxyId)
	return resourceTencentCloudTeoApplicationProxyRead(d, meta)
}

func resourceTencentCloudTeoApplicationProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_application_proxy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	applicationProxy, err := service.DescribeTeoApplicationProxy(ctx, zoneId, proxyId)

	if err != nil {
		return err
	}

	if applicationProxy == nil {
		d.SetId("")
		return fmt.Errorf("resource `applicationProxy` %s does not exist", proxyId)
	}

	if applicationProxy.ZoneId != nil {
		_ = d.Set("zone_id", applicationProxy.ZoneId)
	}

	if applicationProxy.ZoneName != nil {
		_ = d.Set("zone_name", applicationProxy.ZoneName)
	}

	if applicationProxy.ProxyName != nil {
		_ = d.Set("proxy_name", applicationProxy.ProxyName)
	}

	if applicationProxy.PlatType != nil {
		_ = d.Set("plat_type", applicationProxy.PlatType)
	}

	if applicationProxy.SecurityType != nil {
		_ = d.Set("security_type", applicationProxy.SecurityType)
	}

	if applicationProxy.AccelerateType != nil {
		_ = d.Set("accelerate_type", applicationProxy.AccelerateType)
	}

	if applicationProxy.SessionPersistTime != nil {
		_ = d.Set("session_persist_time", applicationProxy.SessionPersistTime)
	}

	if applicationProxy.ProxyType != nil {
		_ = d.Set("proxy_type", applicationProxy.ProxyType)
	}

	if applicationProxy.ProxyId != nil {
		_ = d.Set("proxy_id", applicationProxy.ProxyId)
	}

	if applicationProxy.ScheduleValue != nil {
		_ = d.Set("schedule_value", applicationProxy.ScheduleValue)
	}

	if applicationProxy.UpdateTime != nil {
		_ = d.Set("update_time", applicationProxy.UpdateTime)
	}

	if applicationProxy.HostId != nil {
		_ = d.Set("host_id", applicationProxy.HostId)
	}

	return nil
}

func resourceTencentCloudTeoApplicationProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_application_proxy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyApplicationProxyRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId
	request.ForwardClientIp = helper.String("")
	request.SessionPersist = helper.Bool(true)

	if v, ok := d.GetOk("proxy_name"); ok {
		request.ProxyName = helper.String(v.(string))
	}

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if d.HasChange("zone_name") {
		return fmt.Errorf("`zone_name` do not support change now.")
	}

	if d.HasChange("plat_type") {
		return fmt.Errorf("`plat_type` do not support change now.")
	}

	if d.HasChange("security_type") {
		return fmt.Errorf("`security_type` do not support change now.")
	}

	if d.HasChange("accelerate_type") {
		return fmt.Errorf("`accelerate_type` do not support change now.")
	}

	if d.HasChange("session_persist_time") {
		if v, ok := d.GetOk("session_persist_time"); ok {
			request.SessionPersistTime = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("proxy_type") {
		if v, ok := d.GetOk("proxy_type"); ok {
			request.ProxyType = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyApplicationProxy(request)
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

	return resourceTencentCloudTeoApplicationProxyRead(d, meta)
}

func resourceTencentCloudTeoApplicationProxyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_application_proxy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	if err := service.DeleteTeoApplicationProxyById(ctx, zoneId, proxyId); err != nil {
		return err
	}

	return nil
}
