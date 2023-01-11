/*
Provides a resource to create a antiddos ip. Only support for bgp-multip.

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_ip_attachment_v2" "boundip" {
  id = "bgp-xxxxxx"
  bound_ip_list {
	ip = "1.1.1.1"
	biz_type = "public"
	instance_id = "ins-xxx"
	device_type = "cvm"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDayuDDosIpAttachmentV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDDosIpAttachmentCreateV2,
		Read:   resourceTencentCloudDayuDDosIpAttachmentReadV2,
		Delete: resourceTencentCloudDayuDDosIpAttachmentDeleteV2,
		Schema: map[string]*schema.Schema{
			"bgp_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Anti-DDoS instance ID.",
			},

			"bound_ip_list": {
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
					},
				},
			},
		},
	}
}

func resourceTencentCloudDayuDDosIpAttachmentCreateV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_ip_attachment_v2.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	var bgpInstanceId string
	boundIps := make([]string, 0)
	request := antiddos.NewCreateBoundIPRequest()

	request.Business = helper.String(BUSINESS_BGP_MULTIP)
	if v, ok := d.GetOk("bgp_instance_id"); ok {
		bgpInstanceId = v.(string)
		request.Id = helper.String(bgpInstanceId)
	}

	if v, ok := d.GetOk("bound_ip_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			boundIpInfo := antiddos.BoundIpInfo{}
			if v, ok := dMap["ip"]; ok {
				boundIps = append(boundIps, v.(string))
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
			boundIpInfo.IspCode = helper.IntUint64(ISP_CODE_BGP)
			request.BoundDevList = append(request.BoundDevList, &boundIpInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAntiddosClient().CreateBoundIP(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if *result.Response.Success.Code != RESPONSE_SUCCESS_CODE {
			return resource.RetryableError(errors.New("request failed"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos boundip failed, reason:%+v", logId, err)
		return err
	}
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		boundip, e := service.DescribeAntiddosBoundipById(ctx, bgpInstanceId)
		if e != nil {
			return retryError(e)
		}
		if *boundip.BoundStatus == DAYU_BOUNDSTATUS_IDLE {
			return nil
		}
		return resource.RetryableError(errors.New("still building."))

	})
	if err != nil {
		return err
	}

	d.SetId(bgpInstanceId + FILED_SP + strings.Join(boundIps, COMMA_SP))

	return resourceTencentCloudDayuDDosIpAttachmentReadV2(d, meta)
}

func resourceTencentCloudDayuDDosIpAttachmentReadV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_ip_attachment_v2.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bgpInstanceId := idSplit[0]
	boundIps := idSplit[1]
	boundIpMap := make(map[string]bool)
	for _, boundIp := range strings.Split(boundIps, COMMA_SP) {
		boundIpMap[boundIp] = true
	}

	boundip, err := service.DescribeAntiddosBoundipById(ctx, bgpInstanceId)
	if err != nil {
		return err
	}

	if boundip == nil {
		d.SetId("")
		log.Printf("resource `AntiddosIp` %s does not exist", d.Id())
		return nil
	}
	_ = d.Set("bgp_instance_id", bgpInstanceId)
	boundIpList := make([]map[string]interface{}, 0)
	if boundip.EipProductInfos != nil {
		boundIpListItem := make(map[string]interface{})
		for _, item := range boundip.EipProductInfos {
			if _, ok := boundIpMap[*item.Ip]; !ok {
				continue
			}
			boundIpListItem["ip"] = *item.Ip
			boundIpListItem["biz_type"] = *item.BizType
			boundIpListItem["instance_id"] = *item.InstanceId
			boundIpListItem["device_type"] = *item.DeviceType
			boundIpList = append(boundIpList, boundIpListItem)
		}
	}
	_ = d.Set("bound_ip_list", boundIpList)

	return nil
}

func resourceTencentCloudDayuDDosIpAttachmentDeleteV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_ip_attachment_v2.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bgpInstanceId := idSplit[0]
	boundIps := idSplit[1]

	request := antiddos.NewCreateBoundIPRequest()
	request.Business = helper.String(BUSINESS_BGP_MULTIP)
	request.Id = helper.String(bgpInstanceId)
	ubBoundDevList := make([]*antiddos.BoundIpInfo, 0)
	for _, boundIp := range strings.Split(boundIps, COMMA_SP) {
		ubBoundDevList = append(
			ubBoundDevList,
			&antiddos.BoundIpInfo{
				Ip: &boundIp,
			},
		)
	}
	request.UnBoundDevList = ubBoundDevList

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAntiddosClient().CreateBoundIP(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if *result.Response.Success.Code != RESPONSE_SUCCESS_CODE {
			return resource.RetryableError(errors.New("request failed"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos boundip failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		boundip, e := service.DescribeAntiddosBoundipById(ctx, bgpInstanceId)
		if e != nil {
			return retryError(e)
		}
		if *boundip.BoundStatus == DAYU_BOUNDSTATUS_IDLE {
			return nil
		}
		return resource.RetryableError(errors.New("still building."))

	})
	if err != nil {
		return err
	}

	return nil
}
