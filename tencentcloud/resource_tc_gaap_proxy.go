package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
)

func resourceTencentCloudGaapProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapProxyCreate,
		Read:   resourceTencentCloudGaapProxyRead,
		Update: resourceTencentCloudGaapProxyUpdate,
		Delete: resourceTencentCloudGaapProxyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAllowedIntValue([]int{10, 20, 50, 100, 200, 500, 1000}),
			},
			"concurrent": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAllowedIntValue([]int{2, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}),
			},
			"access_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"realserver_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			// computed
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scalable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_protocols": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"forward_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudGaapProxyCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_proxy.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	name := d.Get("name").(string)
	projectId := d.Get("project_id").(int)
	bandwidth := d.Get("bandwidth").(int)
	concurrent := d.Get("concurrent").(int)
	accessRegion := d.Get("access_region").(string)
	realserverRegion := d.Get("realserver_region").(string)
	enable := d.Get("enable").(bool)
	tags := getTags(d, "tags")

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateProxy(ctx, name, accessRegion, realserverRegion, bandwidth, concurrent, projectId, tags)
	if err != nil {
		return err
	}

	if !enable {
		if err := service.DisableProxy(ctx, id); err != nil {
			return err
		}
	}

	d.SetId(id)

	return resourceTencentCloudGaapProxyRead(d, m)
}

func resourceTencentCloudGaapProxyRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_proxy.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	proxies, err := service.DescribeProxies(ctx, []string{id}, nil, nil, nil, nil)
	if err != nil {
		return err
	}

	var proxy *gaap.ProxyInfo
	for _, p := range proxies {
		if p.InstanceId == nil {
			return errors.New("proxy id is nil")
		}
		if *p.InstanceId == id {
			proxy = p
			break
		}
	}

	if proxy == nil {
		d.SetId("")
		return nil
	}

	if proxy.ProxyName == nil {
		return fmt.Errorf("proxy %s name is nil", id)
	}
	d.Set("name", proxy.ProxyName)

	if proxy.ProjectId == nil {
		return fmt.Errorf("proxy %s project id is nil", id)
	}
	d.Set("project_id", proxy.ProjectId)

	if proxy.Bandwidth == nil {
		return fmt.Errorf("proxy %s bandwidth is nil", id)
	}
	d.Set("bandwidth", proxy.Bandwidth)

	if proxy.Concurrent == nil {
		return fmt.Errorf("proxy %s concurrent is nil", id)
	}
	d.Set("concurrent", *proxy.Concurrent)

	if proxy.AccessRegion == nil {
		return fmt.Errorf("proxy %s access region is nil", id)
	}
	d.Set("access_region", proxy.AccessRegion)

	if proxy.RealServerRegion == nil {
		return fmt.Errorf("proxy %s realserver region is nil", id)
	}
	d.Set("realserver_region", proxy.RealServerRegion)

	if proxy.Status == nil {
		return fmt.Errorf("proxy %s status is nil", id)
	}
	d.Set("enable", *proxy.Status == "RUNNING")
	d.Set("status", proxy.Status)

	if len(proxy.TagSet) > 0 {
		tags := make(map[string]string, len(proxy.TagSet))
		for _, tag := range proxy.TagSet {
			tags[*tag.TagKey] = *tag.TagValue
		}
		d.Set("tags", tags)
	}

	if proxy.CreateTime == nil {
		return fmt.Errorf("proxy %s create time is nil", id)
	}
	d.Set("create_time", proxy.CreateTime)

	if proxy.Domain == nil {
		return fmt.Errorf("proxy %s access domain is nil", id)
	}
	d.Set("access_domain", proxy.Domain)

	if proxy.IP == nil {
		return fmt.Errorf("proxy %s access IP is nil", id)
	}
	d.Set("access_ip", proxy.IP)

	if proxy.Scalarable == nil {
		return fmt.Errorf("proxy %s scalable is nil", id)
	}
	d.Set("scalable", *proxy.Scalarable == 1)

	if len(proxy.SupportProtocols) == 0 {
		return fmt.Errorf("proxy %s no suuport protocols", id)
	}
	supportProtocols := make([]string, 0, len(proxy.SupportProtocols))
	for _, sp := range proxy.SupportProtocols {
		supportProtocols = append(supportProtocols, *sp)
	}
	d.Set("support_protocols", supportProtocols)

	if proxy.ForwardIP == nil {
		return fmt.Errorf("proxy %s forward IPs is nil", id)
	}
	d.Set("forward_ip", proxy.ForwardIP)

	return nil
}

func resourceTencentCloudGaapProxyUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_proxy.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	d.Partial(true)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		if err := service.ModifyProxyName(ctx, id, name); err != nil {
			return err
		}
		d.SetPartial("name")
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		if err := service.ModifyProxyProjectId(ctx, id, projectId); err != nil {
			return err
		}
		d.SetPartial("project_id")
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
		if err := service.ModifyProxyConfiguration(ctx, id, bandwidth, concurrent); err != nil {
			return err
		}
		if d.HasChange("bandwidth") {
			d.SetPartial("bandwidth")
		}
		if d.HasChange("concurrent") {
			d.SetPartial("concurrent")
		}
	}

	if d.HasChange("enable") {
		enable := d.Get("enable").(bool)
		if enable {
			if err := service.EnableProxy(ctx, id); err != nil {
				return err
			}
		} else {
			if err := service.DisableProxy(ctx, id); err != nil {
				return err
			}
		}
		d.SetPartial("enable")
	}

	d.Partial(false)

	return resourceTencentCloudGaapProxyRead(d, m)
}

func resourceTencentCloudGaapProxyDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_proxy.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	createTime := d.Get("create_time").(int)

	if time.Now().Unix()-int64(createTime) < 120 {
		log.Printf("[DEBUG]%s proxy can't be deleted unless it has lived 2 minutes", logId)
		sleepTime := 2*time.Minute - time.Duration(time.Now().UnixNano()-int64(createTime)*1000000000)
		time.Sleep(sleepTime)
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteProxy(ctx, id)
}
