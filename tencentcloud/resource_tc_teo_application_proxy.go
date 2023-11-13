/*
Provides a resource to create a teo application_proxy

Example Usage

```hcl
resource "tencentcloud_teo_application_proxy" "application_proxy" {
  zone_id = &lt;nil&gt;
    proxy_name = &lt;nil&gt;
  proxy_type = &lt;nil&gt;
  plat_type = &lt;nil&gt;
    security_type = &lt;nil&gt;
  accelerate_type = &lt;nil&gt;
  session_persist_time = &lt;nil&gt;
  status = &lt;nil&gt;
        i_pv6 {
		switch = &lt;nil&gt;

  }
  }
```

Import

teo application_proxy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_application_proxy.application_proxy application_proxy_id
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

func resourceTencentCloudTeoApplicationProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoApplicationProxyCreate,
		Read:   resourceTencentCloudTeoApplicationProxyRead,
		Update: resourceTencentCloudTeoApplicationProxyUpdate,
		Delete: resourceTencentCloudTeoApplicationProxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"proxy_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Proxy ID.",
			},

			"proxy_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "When `ProxyType` is hostname, `ProxyName` is the domain or subdomain name.When `ProxyType` is instance, `ProxyName` is the name of proxy application.",
			},

			"proxy_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Layer 4 proxy mode. Valid values:- `hostname`: subdomain mode.- `instance`: instance mode.",
			},

			"plat_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Scheduling mode.- `ip`: Anycast IP.- `domain`: CNAME.",
			},

			"area": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Acceleration area. Valid values: `mainland`, `overseas`.",
			},

			"security_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "- `0`: Disable security protection.- `1`: Enable security protection.",
			},

			"accelerate_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "- `0`: Disable acceleration.- `1`: Enable acceleration.",
			},

			"session_persist_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Session persistence duration. Value range: 30-3600 (in seconds), default value is 600.",
			},

			"status": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Status of this application proxy. Valid values to set is `online` and `offline`.- `online`: Enable.- `offline`: Disable.- `progress`: Deploying.- `stopping`: Deactivating.- `fail`: Deploy or deactivate failed.",
			},

			"ban_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application proxy block status. Valid values: `banned`, `banning`, `recover`, `recovering`.",
			},

			"schedule_value": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Scheduling information.",
			},

			"host_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "When `ProxyType` is hostname, this field is the ID of the subdomain.",
			},

			"i_pv6": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
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
				Computed:    true,
				Type:        schema.TypeString,
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
		response = teo.NewCreateApplicationProxyResponse()
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

	if v, ok := d.GetOkExists("security_type"); ok {
		request.SecurityType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("accelerate_type"); ok {
		request.AccelerateType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("session_persist_time"); ok {
		request.SessionPersistTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "i_pv6"); ok {
		ipv6Access := teo.Ipv6Access{}
		if v, ok := dMap["switch"]; ok {
			ipv6Access.Switch = helper.String(v.(string))
		}
		request.IPv6 = &ipv6Access
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateApplicationProxy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo applicationProxy failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, proxyId}, FILED_SP))

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"online"}, 60*readRetryTimeout, time.Second, service.TeoApplicationProxyStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

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

	applicationProxy, err := service.DescribeTeoApplicationProxyById(ctx, zoneId, proxyId)
	if err != nil {
		return err
	}

	if applicationProxy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoApplicationProxy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
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

	if applicationProxy.IPv6 != nil {
		iPv6Map := map[string]interface{}{}

		if applicationProxy.IPv6.Switch != nil {
			iPv6Map["switch"] = applicationProxy.IPv6.Switch
		}

		_ = d.Set("i_pv6", []interface{}{iPv6Map})
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

	var (
		modifyApplicationProxyRequest  = teo.NewModifyApplicationProxyRequest()
		modifyApplicationProxyResponse = teo.NewModifyApplicationProxyResponse()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	proxyId := idSplit[1]

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId

	immutableArgs := []string{"zone_id", "proxy_id", "proxy_name", "proxy_type", "plat_type", "area", "security_type", "accelerate_type", "session_persist_time", "status", "ban_status", "schedule_value", "host_id", "i_pv6", "update_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("proxy_name") {
		if v, ok := d.GetOk("proxy_name"); ok {
			request.ProxyName = helper.String(v.(string))
		}
	}

	if d.HasChange("proxy_type") {
		if v, ok := d.GetOk("proxy_type"); ok {
			request.ProxyType = helper.String(v.(string))
		}
	}

	if d.HasChange("session_persist_time") {
		if v, ok := d.GetOkExists("session_persist_time"); ok {
			request.SessionPersistTime = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	if d.HasChange("i_pv6") {
		if dMap, ok := helper.InterfacesHeadMap(d, "i_pv6"); ok {
			ipv6Access := teo.Ipv6Access{}
			if v, ok := dMap["switch"]; ok {
				ipv6Access.Switch = helper.String(v.(string))
			}
			request.IPv6 = &ipv6Access
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyApplicationProxy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo applicationProxy failed, reason:%+v", logId, err)
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
