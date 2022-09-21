/*
Provides a resource to create a teo application_proxy

Example Usage

```hcl
resource "tencentcloud_teo_application_proxy" "application_proxy" {
  accelerate_type      = 1
  plat_type            = "domain"
  proxy_name           = "applicationProxies-test-1"
  proxy_type           = "instance"
  security_type        = 1
  session_persist_time = 2400
  status               = "online"
  tags                 = {}
  zone_id              = "zone-2983wizgxqvm"

  ipv6 {
    switch = "off"
  }
}

```
Import

teo application_proxy can be imported using the zoneId#proxyId, e.g.
```
$ terraform import tencentcloud_teo_application_proxy.application_proxy zone-2983wizgxqvm#proxy-6972528a-373a-11ed-afca-52540044a456
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

			"proxy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Proxy ID.",
			},

			"proxy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "When `ProxyType` is hostname, `ProxyName` is the domain or subdomain name.When `ProxyType` is instance, `ProxyName` is the name of proxy application.",
			},

			"proxy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Layer 4 proxy mode. Valid values:- `hostname`: subdomain mode.- `instance`: instance mode.",
			},

			"plat_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Scheduling mode.- `ip`: Anycast IP.- `domain`: CNAME.",
			},

			"area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration area. Valid values: `mainland`, `overseas`.",
			},

			"security_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "- `0`: Disable security protection.- `1`: Enable security protection.",
			},

			"accelerate_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "- `0`: Disable acceleration.- `1`: Enable acceleration.",
			},

			"session_persist_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Session persistence duration. Value range: 30-3600 (in seconds), default value is 600.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status of this application proxy. Valid values to set is `online` and `offline`.- `online`: Enable.- `offline`: Disable.- `progress`: Deploying.- `stopping`: Deactivating.- `fail`: Deploy or deactivate failed.",
			},

			"ban_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Application proxy block status. Valid values: `banned`, `banning`, `recover`, `recovering`.",
			},

			"schedule_value": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Scheduling information.",
			},

			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When `ProxyType` is hostname, this field is the ID of the subdomain.",
			},

			"ipv6": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "IPv6 access configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "- `on`: Enable.- `off`: Disable.",
						},
					},
				},
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification date.",
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

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_name"); ok {
		request.ProxyName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_type"); ok {
		request.ProxyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plat_type"); ok {
		request.PlatType = helper.String(v.(string))
	}

	if v := d.Get("security_type"); v != nil {
		request.SecurityType = helper.IntInt64(v.(int))
	}

	if v := d.Get("accelerate_type"); v != nil {
		request.AccelerateType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("session_persist_time"); ok {
		request.SessionPersistTime = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ipv6"); ok {
		ipv6Access := teo.Ipv6{}
		if v, ok := dMap["switch"]; ok {
			ipv6Access.Switch = helper.String(v.(string))
		}
		request.Ipv6 = &ipv6Access
	}

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err := service.CheckZoneComplete(ctx, zoneId)
	if err != nil {
		log.Printf("[CRITAL]%s create teo dnsRecord failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
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

	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoApplicationProxy(ctx, zoneId, proxyId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.Status == "online" {
			return nil
		}
		if *instance.Status == "fail" {
			return resource.NonRetryableError(fmt.Errorf("applicationProxy status is %v, operate failed.",
				*instance.Status))
		}
		return resource.RetryableError(fmt.Errorf("applicationProxy status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

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

	if applicationProxy.ProxyId != nil {
		_ = d.Set("proxy_id", applicationProxy.ProxyId)
	}

	if applicationProxy.ProxyName != nil {
		_ = d.Set("proxy_name", applicationProxy.ProxyName)
	}

	if applicationProxy.ProxyType != nil {
		_ = d.Set("proxy_type", applicationProxy.ProxyType)
	}

	if applicationProxy.PlatType != nil {
		_ = d.Set("plat_type", applicationProxy.PlatType)
	}

	if applicationProxy.Area != nil {
		_ = d.Set("area", applicationProxy.Area)
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

	if applicationProxy.Status != nil {
		_ = d.Set("status", applicationProxy.Status)
	}

	if applicationProxy.BanStatus != nil {
		_ = d.Set("ban_status", applicationProxy.BanStatus)
	}

	if applicationProxy.ScheduleValue != nil {
		_ = d.Set("schedule_value", applicationProxy.ScheduleValue)
	}

	if applicationProxy.HostId != nil {
		_ = d.Set("host_id", applicationProxy.HostId)
	}

	if applicationProxy.Ipv6 != nil {
		iPv6Map := map[string]interface{}{}
		if applicationProxy.Ipv6.Switch != nil {
			iPv6Map["switch"] = applicationProxy.Ipv6.Switch
		}

		_ = d.Set("ipv6", []interface{}{iPv6Map})
	}

	if applicationProxy.UpdateTime != nil {
		_ = d.Set("update_time", applicationProxy.UpdateTime)
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

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if v, ok := d.GetOk("proxy_name"); ok {
		request.ProxyName = helper.String(v.(string))
	}

	if d.HasChange("proxy_type") {
		if v, ok := d.GetOk("proxy_type"); ok {
			request.ProxyType = helper.String(v.(string))
		}
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

	if d.HasChange("ipv6") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ipv6"); ok {
			ipv6Access := teo.Ipv6{}
			if v, ok := dMap["switch"]; ok {
				ipv6Access.Switch = helper.String(v.(string))
			}

			request.Ipv6 = &ipv6Access
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
		log.Printf("[CRITAL]%s create teo applicationProxy failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			statusRequest := teo.NewModifyApplicationProxyStatusRequest()

			statusRequest.ZoneId = &zoneId
			statusRequest.ProxyId = &proxyId
			statusRequest.Status = helper.String(v.(string))

			statusErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				statusResult, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyApplicationProxyStatus(statusRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, request.GetAction(), request.ToJsonString(), statusResult.ToJsonString())
				}
				return nil
			})

			if statusErr != nil {
				log.Printf("[CRITAL]%s create teo applicationProxy failed, reason:%+v", logId, statusErr)
				return statusErr
			}
			_ = d.Set("status", v.(string))
		}
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

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		e := resourceTencentCloudTeoApplicationProxyRead(d, meta)
		if e != nil {
			log.Printf("[CRITAL]%s get teo applicationProxy failed, reason:%+v", logId, e)
			return resource.RetryableError(e)
		}

		status, _ := d.Get("status").(string)
		if status == "offline" {
			return nil
		}
		if status == "stopping" {
			return resource.RetryableError(fmt.Errorf("applicationProxy stopping"))
		}

		statusRequest := teo.NewModifyApplicationProxyStatusRequest()
		statusRequest.ZoneId = &zoneId
		statusRequest.ProxyId = &proxyId
		statusRequest.Status = helper.String("offline")
		_, e = meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyApplicationProxyStatus(statusRequest)
		if e != nil {
			return resource.NonRetryableError(fmt.Errorf("setting applicationProxy `status` to offline failed, reason: %v", e))
		}
		return resource.RetryableError(fmt.Errorf("setting applicationProxy `status` to offline"))
	})
	if err != nil {
		return err
	}

	if err = service.DeleteTeoApplicationProxyById(ctx, zoneId, proxyId); err != nil {
		return err
	}

	return nil
}
