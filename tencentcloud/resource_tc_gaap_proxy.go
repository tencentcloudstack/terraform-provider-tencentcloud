/*
Provides a resource to create a GAAP proxy.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"

  tags = {
    test = "test"
  }
}
```

Import

GAAP proxy can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_proxy.foo link-11112222
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var gaapActionMu = &sync.Mutex{}

func resourceTencentCloudGaapProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapProxyCreate,
		Read:   resourceTencentCloudGaapProxyRead,
		Update: resourceTencentCloudGaapProxyUpdate,
		Delete: resourceTencentCloudGaapProxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
				Description:  "Name of the GAAP proxy, the maximum length is 30.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the project within the GAAP proxy, `0` means is default project.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum bandwidth of the GAAP proxy, unit is Mbps. Valid value: `10`, `20`, `50`, `100`, `200`, `500`, `1000`, `2000`, `5000` and `10000`. To set `2000`, `5000` or `10000`, you need to apply for a whitelist from Tencent Cloud.",
			},
			"concurrent": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum concurrency of the GAAP proxy, unit is 10k. Valid value: `2`, `5`, `10`, `20`, `30`, `40`, `50`, `60`, `70`, `80`, `90`, `100`, `150`, `200`, `250` and `300`. To set `150`, `200`, `250` or `300`, you need to apply for a whitelist from Tencent Cloud.",
			},
			"access_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Access region of the GAAP proxy. Valid value: `NorthChina`, `EastChina`, `SouthChina`, `SouthwestChina`, `Hongkong`, `SL_TAIWAN`, `SoutheastAsia`, `Korea`, `SL_India`, `SL_Australia`, `Europe`, `SL_UK`, `SL_SouthAmerica`, `NorthAmerica`, `SL_MiddleUSA`, `Canada`, `SL_VIET`, `WestIndia`, `Thailand`, `Virginia`, `Russia`, `Japan` and `SL_Indonesia`.",
			},
			"realserver_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the GAAP realserver. Valid value: `NorthChina`, `EastChina`, `SouthChina`, `SouthwestChina`, `Hongkong`, `SL_TAIWAN`, `SoutheastAsia`, `Korea`, `SL_India`, `SL_Australia`, `Europe`, `SL_UK`, `SL_SouthAmerica`, `NorthAmerica`, `SL_MiddleUSA`, `Canada`, `SL_VIET`, `WestIndia`, `Thailand`, `Virginia`, `Russia`, `Japan` and `SL_Indonesia`.",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether GAAP proxy is enabled, default value is `true`.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the GAAP proxy. Tags that do not exist are not created automatically.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(PROXY_NETWORK_TYPE),
				Description:  "Network type. `normal`: regular BGP, `cn2`: boutique BGP, `triple`: triple play.",
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the GAAP proxy.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the GAAP proxy.",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access domain of the GAAP proxy.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access IP of the GAAP proxy.",
			},
			"scalable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether GAAP proxy can scalable.",
			},
			"support_protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Supported protocols of the GAAP proxy.",
			},
			"forward_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Forwarding IP of the GAAP proxy.",
			},
		},
	}
}

func resourceTencentCloudGaapProxyCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_proxy.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]interface{})
	name := d.Get("name").(string)
	projectId := d.Get("project_id").(int)
	bandwidth := d.Get("bandwidth").(int)
	concurrent := d.Get("concurrent").(int)
	accessRegion := d.Get("access_region").(string)
	realserverRegion := d.Get("realserver_region").(string)
	enable := d.Get("enable").(bool)
	tags := helper.GetTags(d, "tags")

	if v, ok := d.GetOk("network_type"); ok {
		params["network_type"] = v.(string)
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateProxy(ctx, name, accessRegion, realserverRegion, bandwidth, concurrent, projectId, tags, params)
	if err != nil {
		return err
	}

	d.SetId(id)

	if !enable {
		if err := service.DisableProxy(ctx, id); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapProxyRead(d, m)
}

func resourceTencentCloudGaapProxyRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_proxy.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	proxies, err := service.DescribeProxies(ctx, []string{id}, nil, nil, nil, nil)
	if err != nil {
		return err
	}

	var proxy *gaap.ProxyInfo
	for _, p := range proxies {
		if p.ProxyId == nil {
			return errors.New("proxy id is nil")
		}
		if *p.ProxyId == id {
			proxy = p
			break
		}
	}

	if proxy == nil {
		d.SetId("")
		return nil
	}

	if proxy.ProxyName == nil {
		return errors.New("proxy name is nil")
	}
	_ = d.Set("name", proxy.ProxyName)

	if proxy.ProjectId == nil {
		return errors.New("proxy project id is nil")
	}
	_ = d.Set("project_id", proxy.ProjectId)

	if proxy.Bandwidth == nil {
		return errors.New("proxy bandwidth is nil")
	}
	_ = d.Set("bandwidth", proxy.Bandwidth)

	if proxy.Concurrent == nil {
		return errors.New("proxy concurrent is nil")
	}
	_ = d.Set("concurrent", proxy.Concurrent)

	if proxy.AccessRegion == nil {
		return errors.New("proxy access region is nil")
	}
	_ = d.Set("access_region", proxy.AccessRegion)

	if proxy.RealServerRegion == nil {
		return errors.New("proxy realserver region is nil")
	}
	_ = d.Set("realserver_region", proxy.RealServerRegion)

	if proxy.Status == nil {
		return errors.New("proxy status is nil")
	}
	_ = d.Set("enable", *proxy.Status == GAAP_PROXY_RUNNING)
	_ = d.Set("status", proxy.Status)

	if len(proxy.TagSet) > 0 {
		tags := make(map[string]string, len(proxy.TagSet))
		for _, tag := range proxy.TagSet {
			tags[*tag.TagKey] = *tag.TagValue
		}
		_ = d.Set("tags", tags)
	}

	if proxy.CreateTime == nil {
		return errors.New("proxy create time is nil")
	}
	_ = d.Set("create_time", helper.FormatUnixTime(*proxy.CreateTime))

	if proxy.Domain == nil {
		return errors.New("proxy access domain is nil")
	}
	_ = d.Set("domain", proxy.Domain)

	if proxy.IP == nil {
		return errors.New("proxy access IP is nil")
	}
	_ = d.Set("ip", proxy.IP)

	if proxy.Scalarable == nil {
		return errors.New("proxy scalable is nil")
	}
	_ = d.Set("scalable", *proxy.Scalarable == 1)

	if len(proxy.SupportProtocols) == 0 {
		return errors.New("proxy support protocols is empty")
	}
	supportProtocols := make([]string, 0, len(proxy.SupportProtocols))
	for _, sp := range proxy.SupportProtocols {
		supportProtocols = append(supportProtocols, *sp)
	}
	_ = d.Set("support_protocols", supportProtocols)

	if proxy.ForwardIP == nil {
		return errors.New("proxy forward ip is nil")
	}
	_ = d.Set("forward_ip", proxy.ForwardIP)

	if proxy.NetworkType != nil {
		_ = d.Set("network_type", proxy.NetworkType)
	}
	return nil
}

func resourceTencentCloudGaapProxyUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_proxy.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	gaapService := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	d.Partial(true)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		if err := gaapService.ModifyProxyName(ctx, id, name); err != nil {
			return err
		}

	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		if err := gaapService.ModifyProxyProjectId(ctx, id, projectId); err != nil {
			return err
		}

	}

	if d.HasChange("bandwidth") || d.HasChange("concurrent") {
		var (
			bandwidth  *int
			concurrent *int
		)
		if d.HasChange("bandwidth") {
			bandwidth = common.IntPtr(d.Get("bandwidth").(int))
		}
		if d.HasChange("concurrent") {
			concurrent = common.IntPtr(d.Get("concurrent").(int))
		}
		if err := gaapService.ModifyProxyConfiguration(ctx, id, bandwidth, concurrent); err != nil {
			return err
		}
		//deal with sync delay
		time.Sleep(time.Duration(10) * time.Second)
	}

	if d.HasChange("enable") {
		enable := d.Get("enable").(bool)
		if enable {
			if err := gaapService.EnableProxy(ctx, id); err != nil {
				return err
			}
		} else {
			if err := gaapService.DisableProxy(ctx, id); err != nil {
				return err
			}
		}

	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := TagService{client: m.(*TencentCloudClient).apiV3Conn}

		region := m.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::gaap:%s:uin/:proxy/%s", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudGaapProxyRead(d, m)
}

func resourceTencentCloudGaapProxyDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_proxy.update")()
	//gaapActionMu.Lock()
	//defer gaapActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	createTimeStr := d.Get("create_time").(string)

	if createTime, err := helper.ParseTime(createTimeStr); err == nil {
		if !time.Now().After(createTime.Add(2 * time.Minute)) {
			log.Printf("[DEBUG]%s proxy can't be deleted unless it has lived 2 minutes", logId)
			time.Sleep(time.Until(createTime.Add(2 * time.Minute)))
		}
	} else {
		log.Printf("[WARN]%s parse create time failed, delete immediately", logId)
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteProxy(ctx, id)
}
