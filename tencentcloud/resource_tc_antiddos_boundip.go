/*
Provides a resource to create a antiddos boundip

Example Usage

```hcl
resource "tencentcloud_antiddos_boundip" "boundip" {
  business = "bgp-multip"
  id = "bgp-000000xe"
  bound_dev_list {
		ip = "1.1.1.1"
		biz_type = "public"
		instance_id = "ins-xxx"
		device_type = "cvm"
		isp_code = 5

  }
  un_bound_dev_list {
		ip = "1.1.1.1"
		biz_type = "public"
		instance_id = "ins-xxx"
		device_type = "cvm"
		isp_code = 5

  }
  copy_policy = ""
}
```

Import

antiddos boundip can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_boundip.boundip boundip_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudAntiddosBoundip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosBoundipCreate,
		Read:   resourceTencentCloudAntiddosBoundipRead,
		Delete: resourceTencentCloudAntiddosBoundipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"business": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Anti-DDoS service type. bgp: Anti-DDoS Pro (Single IP); bgp-multip: Anti-DDoS Pro (Multi-IP).",
			},

			"id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Anti-DDoS instance ID.",
			},

			"bound_dev_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Array of IPs to bind to the Anti-DDoS instance. For Anti-DDoS Pro Single IP instance, the array contains only one IP. If there are no IPs to bind, it is empty; however, either BoundDevList or UnBoundDevList must not be empty.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address.",
						},
						"biz_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Category of product that can be bound. Valid values: public (CVM and CLB), bm (BM), eni (ENI), vpngw (VPN gateway), natgw (NAT gateway), waf (WAF), fpc (financial products), gaap (GAAP), and other (hosted IP). This field is required when you perform binding.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Anti-DDoS instance ID of the IP. This field is required only when the instance is bound to an IP. For example, this field InstanceId will be eni-* if the instance ID is bound to an ENI IP; none if there is no instance to bind to a managed IP.",
						},
						"device_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sub-product category. Valid values: cvm (CVM), lb (Load balancer), eni (ENI), vpngw (VPN gateway), natgw (NAT gateway), waf (WAF), fpc (financial products), gaap (GAAP), eip (BM EIP) and other (managed IP). This field is required when you perform binding.",
						},
						"isp_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ISP. Valid values: 0 (China Telecom), 1 (China Unicom), 2 (China Mobile), and 5 (BGP). This field is required when you perform binding.",
						},
					},
				},
			},

			"un_bound_dev_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Array of IPs to unbind from the Anti-DDoS instance. For Anti-DDoS Pro Single IP instance, the array contains only one IP; if there are no IPs to unbind, it is empty; however, either BoundDevList or UnBoundDevList must not be empty.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address.",
						},
						"biz_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Category of product that can be bound. Valid values: public (CVM and CLB), bm (BM), eni (ENI), vpngw (VPN gateway), natgw (NAT gateway), waf (WAF), fpc (financial products), gaap (GAAP), and other (hosted IP). This field is required when you perform binding.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Anti-DDoS instance ID of the IP. This field is required only when the instance is bound to an IP. For example, this field InstanceId will be eni-* if the instance ID is bound to an ENI IP; none if there is no instance to bind to a managed IP.",
						},
						"device_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sub-product category. Valid values: cvm (CVM), lb (Load balancer), eni (ENI), vpngw (VPN gateway), natgw (NAT gateway), waf (WAF), fpc (financial products), gaap (GAAP), eip (BM EIP) and other (managed IP). This field is required when you perform binding.",
						},
						"isp_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ISP. Valid values: 0 (China Telecom), 1 (China Unicom), 2 (China Mobile), and 5 (BGP). This field is required when you perform binding.",
						},
					},
				},
			},

			"copy_policy": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Disused.",
			},
		},
	}
}

func resourceTencentCloudAntiddosBoundipCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_boundip.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = antiddos.NewCreateBoundIPRequest()
		response = antiddos.NewCreateBoundIPResponse()
		id       string
	)
	if v, ok := d.GetOk("business"); ok {
		request.Business = helper.String(v.(string))
	}

	if v, ok := d.GetOk("id"); ok {
		id = v.(string)
		request.Id = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bound_dev_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			boundIpInfo := antiddos.BoundIpInfo{}
			if v, ok := dMap["ip"]; ok {
				boundIpInfo.Ip = helper.String(v.(string))
			}
			if v, ok := dMap["biz_type"]; ok {
				boundIpInfo.BizType = helper.String(v.(string))
			}
			if v, ok := dMap["instance_id"]; ok {
				boundIpInfo.InstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["device_type"]; ok {
				boundIpInfo.DeviceType = helper.String(v.(string))
			}
			if v, ok := dMap["isp_code"]; ok {
				boundIpInfo.IspCode = helper.IntUint64(v.(int))
			}
			request.BoundDevList = append(request.BoundDevList, &boundIpInfo)
		}
	}

	if v, ok := d.GetOk("un_bound_dev_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			boundIpInfo := antiddos.BoundIpInfo{}
			if v, ok := dMap["ip"]; ok {
				boundIpInfo.Ip = helper.String(v.(string))
			}
			if v, ok := dMap["biz_type"]; ok {
				boundIpInfo.BizType = helper.String(v.(string))
			}
			if v, ok := dMap["instance_id"]; ok {
				boundIpInfo.InstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["device_type"]; ok {
				boundIpInfo.DeviceType = helper.String(v.(string))
			}
			if v, ok := dMap["isp_code"]; ok {
				boundIpInfo.IspCode = helper.IntUint64(v.(int))
			}
			request.UnBoundDevList = append(request.UnBoundDevList, &boundIpInfo)
		}
	}

	if v, ok := d.GetOk("copy_policy"); ok {
		request.CopyPolicy = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAntiddosClient().CreateBoundIP(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos boundip failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Id
	d.SetId(id)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"idle"}, 1*readRetryTimeout, time.Second, service.AntiddosBoundipStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudAntiddosBoundipRead(d, meta)
}

func resourceTencentCloudAntiddosBoundipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_boundip.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	boundipId := d.Id()

	boundip, err := service.DescribeAntiddosBoundipById(ctx, id)
	if err != nil {
		return err
	}

	if boundip == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosBoundip` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if boundip.Business != nil {
		_ = d.Set("business", boundip.Business)
	}

	if boundip.Id != nil {
		_ = d.Set("id", boundip.Id)
	}

	if boundip.BoundDevList != nil {
		boundDevListList := []interface{}{}
		for _, boundDevList := range boundip.BoundDevList {
			boundDevListMap := map[string]interface{}{}

			if boundip.BoundDevList.Ip != nil {
				boundDevListMap["ip"] = boundip.BoundDevList.Ip
			}

			if boundip.BoundDevList.BizType != nil {
				boundDevListMap["biz_type"] = boundip.BoundDevList.BizType
			}

			if boundip.BoundDevList.InstanceId != nil {
				boundDevListMap["instance_id"] = boundip.BoundDevList.InstanceId
			}

			if boundip.BoundDevList.DeviceType != nil {
				boundDevListMap["device_type"] = boundip.BoundDevList.DeviceType
			}

			if boundip.BoundDevList.IspCode != nil {
				boundDevListMap["isp_code"] = boundip.BoundDevList.IspCode
			}

			boundDevListList = append(boundDevListList, boundDevListMap)
		}

		_ = d.Set("bound_dev_list", boundDevListList)

	}

	if boundip.UnBoundDevList != nil {
		unBoundDevListList := []interface{}{}
		for _, unBoundDevList := range boundip.UnBoundDevList {
			unBoundDevListMap := map[string]interface{}{}

			if boundip.UnBoundDevList.Ip != nil {
				unBoundDevListMap["ip"] = boundip.UnBoundDevList.Ip
			}

			if boundip.UnBoundDevList.BizType != nil {
				unBoundDevListMap["biz_type"] = boundip.UnBoundDevList.BizType
			}

			if boundip.UnBoundDevList.InstanceId != nil {
				unBoundDevListMap["instance_id"] = boundip.UnBoundDevList.InstanceId
			}

			if boundip.UnBoundDevList.DeviceType != nil {
				unBoundDevListMap["device_type"] = boundip.UnBoundDevList.DeviceType
			}

			if boundip.UnBoundDevList.IspCode != nil {
				unBoundDevListMap["isp_code"] = boundip.UnBoundDevList.IspCode
			}

			unBoundDevListList = append(unBoundDevListList, unBoundDevListMap)
		}

		_ = d.Set("un_bound_dev_list", unBoundDevListList)

	}

	if boundip.CopyPolicy != nil {
		_ = d.Set("copy_policy", boundip.CopyPolicy)
	}

	return nil
}

func resourceTencentCloudAntiddosBoundipDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_boundip.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	boundipId := d.Id()

	if err := service.DeleteAntiddosBoundipById(ctx, id); err != nil {
		return err
	}

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"idle"}, 1*readRetryTimeout, time.Second, service.AntiddosBoundipStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
