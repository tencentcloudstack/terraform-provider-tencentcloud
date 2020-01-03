/*
Use this resource to create dayu DDoS policy case

~> **NOTE:** when a dayu DDoS policy case is created, there will be a dayu DDoS policy created with the same prefix name in the same time. This resource only supports Anti-DDoS of type `bgp`, `bgp-multip` and `bgpip`. One Anti-DDoS resource can only has one DDoS policy case resource. When there is only one Anti-DDoS resource and one policy case, those two resource will be bind automatically.

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy_case" "foo" {
  resource_type         = "bgpip"
  name                  = "tf_test_policy_case"
  platform_types        = ["PC", "MOBILE"]
  app_type              = "WEB"
  app_protocols         = ["tcp", "udp"]
  tcp_start_port		= "1000"
  tcp_end_port          = "2000"
  udp_start_port		= "3000"
  udp_end_port			= "4000"
  has_abroad			= "yes"
  has_initiate_tcp		= "yes"
  has_initiate_udp		= "yes"
  peer_tcp_port			= "1111"
  peer_udp_port			= "3333"
  tcp_footprint		= "511"
  udp_footprint		= "500"
  web_api_urls			= ["abc.com", "test.cn/aaa.png"]
  min_tcp_package_len	= "1000"
  max_tcp_package_len	= "1200"
  min_udp_package_len	= "1000"
  max_udp_package_len	= "1200"
  has_vpn				= "yes"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDayuDdosPolicyCase() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDdosPolicyCaseCreate,
		Read:   resourceTencentCloudDayuDdosPolicyCaseRead,
		Update: resourceTencentCloudDayuDdosPolicyCaseUpdate,
		Delete: resourceTencentCloudDayuDdosPolicyCaseDelete,

		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE_CASE),
				ForceNew:     true,
				Description:  "Type of the resource that the DDoS policy case works for, valid values are `bgpip`, `bgp` and `bgp-multip`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 64),
				Description:  "Name of the DDoS policy case. Length should between 1 and 64.",
			},
			"platform_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAllowedStringValue(DAYU_APP_PLATFORM),
					Description:  "Platform of the DDoS policy case, and valid values are `PC`, `MOBILE`, `TV` and `SERVER`.",
				},
				Required:    true,
				Description: "Platform set of the DDoS policy case.",
			},
			"app_protocols": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAllowedStringValue(DAYU_PROTOCOL),
					Description:  "App protocol of the DDoS policy case, and valid values are `tcp`, `udp`, `icmp` and `all`.",
				},
				Required:    true,
				Description: "App protocol set of the DDoS policy case.",
			},
			"app_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_APP_TYPE), //to see the max
				Description:  "App type of the DDoS policy case, and valid values are `WEB`, `GAME`, `APP` and `OTHER`.",
			},
			"tcp_start_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePort,
				Description:  "Start port of the TCP service, valid value is range from 0 to 65535.",
			},
			"tcp_end_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePort,
				Description:  "End port of the TCP service, valid value is range from 0 to 65535. It must be greater than `tcp_start_port`.",
			},
			"udp_start_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePort,
				Description:  "Start port of the UDP service, valid value is range from 0 to 65535.",
			},
			"udp_end_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePort,
				Description:  "End port of the UDP service, valid value is range from 0 to 65535. It must be greater than `udp_start_port`.",
			},
			"has_abroad": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_BOOL_FLAG),
				Description:  "Indicate whether the service involves overseas or not, valid values are `no` and `yes`.",
			},
			"has_initiate_tcp": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_BOOL_FLAG),
				Description:  "Indicate whether the service actively initiates TCP requests or not, valid values are `no` and `yes`.",
			},
			"has_initiate_udp": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_BOOL_FLAG),
				Description:  "Indicate whether the actively initiate UDP requests or not, valid values are `no` and `yes`.",
			},
			"has_vpn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_BOOL_FLAG),
				Description:  "Indicate whether the service involves VPN service or not, valid values are `no` and `yes`.",
			},
			"peer_tcp_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePort,
				Description:  "The port that actively initiates TCP requests, valid value is range from 1 to 65535.",
			},
			"peer_udp_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePort,
				Description:  "The port that actively initiates UDP requests, valid value is range from 1 to 65535.",
			},
			"tcp_footprint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 512),
				Description:  "The fixed signature of TCP protocol load, valid value length is range from 1 to 512.",
			},
			"udp_footprint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 512),
				Description:  "The fixed signature of TCP protocol load, valid value length is range from 1 to 512.",
			},
			"web_api_urls": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "Web API url.",
				},
				Required:    true,
				Description: "Web API url set.",
			},
			"min_tcp_package_len": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 1499),
				Description:  "The minimum length of TCP message package, valid value length should be greater than 0 and less than 1500.",
			},
			"max_tcp_package_len": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 1499),
				Description:  "The max length of TCP message package, valid value length should be greater than 0 and less than 1500. It should be greater than `min_tcp_package_len`.",
			},
			"min_udp_package_len": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 1499),
				Description:  "The minimum length of UDP message package, valid value length should be greater than 0 and less than 1500.",
			},
			"max_udp_package_len": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 1499),
				Description:  "The max length of UDP message package, valid value length should be greater than 0 and less than 1500. It should be greater than `min_udp_package_len`.",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the DDos policy case.",
			},
			"scene_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the DDos policy case.",
			},
		},
	}
}

func resourceTencentCloudDayuDdosPolicyCaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_case.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	request := dayu.NewCreateDDoSPolicyCaseRequest()
	resourceType := d.Get("resource_type").(string)
	request.Business = &resourceType
	request.CaseName = helper.String(d.Get("name").(string))
	platforms := d.Get("platform_types").(*schema.Set).List()
	for _, plat := range platforms {
		request.PlatformTypes = append(request.PlatformTypes, helper.String(plat.(string)))
	}
	protocols := d.Get("app_protocols").(*schema.Set).List()
	for _, protocol := range protocols {
		request.AppProtocols = append(request.AppProtocols, helper.String(protocol.(string)))
	}
	urls := d.Get("web_api_urls").(*schema.Set).List()
	for _, url := range urls {
		request.WebApiUrl = append(request.WebApiUrl, helper.String(url.(string)))
	}
	request.AppType = helper.String(d.Get("app_type").(string))
	request.HasAbroad = helper.String(d.Get("has_abroad").(string))
	request.HasInitiateTcp = helper.String(d.Get("has_initiate_tcp").(string))
	request.HasInitiateUdp = helper.String(d.Get("has_initiate_udp").(string))
	request.HasVPN = helper.String(d.Get("has_vpn").(string))
	request.PeerTcpPort = helper.String(d.Get("peer_tcp_port").(string))
	request.PeerUdpPort = helper.String(d.Get("peer_udp_port").(string))
	request.TcpFootprint = helper.String(d.Get("tcp_footprint").(string))
	request.UdpFootprint = helper.String(d.Get("udp_footprint").(string))

	tcpPortStart := d.Get("tcp_start_port").(string)
	tcpPortEnd := d.Get("tcp_end_port").(string)
	startInt, sErr := strconv.Atoi(tcpPortStart)
	endInt, eErr := strconv.Atoi(tcpPortEnd)
	if sErr != nil {
		return sErr
	}
	if eErr != nil {
		return eErr
	}
	if endInt < startInt {
		return fmt.Errorf("`tcp_start_port`:%s should not be greater than `tcp_end_port`:%s.", tcpPortStart, tcpPortEnd)
	}
	udpPortStart := d.Get("udp_start_port").(string)
	udpPortEnd := d.Get("udp_end_port").(string)
	startInt, sErr = strconv.Atoi(udpPortStart)
	endInt, eErr = strconv.Atoi(udpPortEnd)
	if sErr != nil {
		return sErr
	}
	if eErr != nil {
		return eErr
	}
	if endInt < startInt {
		return fmt.Errorf("`udp_start_port`:%s should not be greater than `udp_end_port`:%s.", udpPortStart, udpPortEnd)
	}
	request.TcpSportStart = &tcpPortStart
	request.TcpSportEnd = &tcpPortEnd
	request.UdpSportStart = &udpPortStart
	request.UdpSportEnd = &udpPortEnd

	minTcpPackageLen := d.Get("min_tcp_package_len").(string)
	maxTcpPackageLen := d.Get("max_tcp_package_len").(string)
	minTcpPackageLenInt, _ := strconv.Atoi(minTcpPackageLen)
	maxTcpPackageLenInt, _ := strconv.Atoi(maxTcpPackageLen)
	if maxTcpPackageLenInt < minTcpPackageLenInt {
		return fmt.Errorf("`min_tcp_package_len`:%s should not be greater than `max_tcp_package_len`:%s.", minTcpPackageLen, maxTcpPackageLen)
	}
	minUdpPackageLen := d.Get("min_udp_package_len").(string)
	maxUdpPackageLen := d.Get("max_udp_package_len").(string)
	minUdpPackageLenInt, _ := strconv.Atoi(minUdpPackageLen)
	maxUdpPackageLenInt, _ := strconv.Atoi(maxUdpPackageLen)
	if maxUdpPackageLenInt < minUdpPackageLenInt {
		return fmt.Errorf("`min_udp_package_len`:%s should not be greater than `max_udp_package_len`:%s.", minUdpPackageLen, maxUdpPackageLen)
	}
	request.MinTcpPackageLen = &minTcpPackageLen
	request.MaxTcpPackageLen = &maxTcpPackageLen
	request.MinUdpPackageLen = &minUdpPackageLen
	request.MaxUdpPackageLen = &maxTcpPackageLen

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	sceneId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateDdosPolicyCase(ctx, request)
		if e != nil {
			return retryError(e)
		}
		sceneId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + FILED_SP + sceneId)

	return resourceTencentCloudDayuDdosPolicyCaseRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyCaseRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_case.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDos policy case")
	}
	resourceType := items[0]
	sceneId := items[1]

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	ddosPolicyCase, has, err := dayuService.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ddosPolicyCase, has, err = dayuService.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	for _, record := range ddosPolicyCase.Record {
		key := *record.Key
		if key == "CaseName" {
			_ = d.Set("name", record.Value)
		}
		if key == "HasInitiateTcp" {
			_ = d.Set("has_initiate_tcp", record.Value)
		}
		if key == "HasInitiateUdp" {
			_ = d.Set("has_initiate_udp", record.Value)
		}
		if key == "HasVPN" {
			_ = d.Set("has_vpn", record.Value)
		}
		if key == "PeerTcpPort" {
			_ = d.Set("peer_tcp_port", record.Value)
		}
		if key == "PeerUdpPort" {
			_ = d.Set("peer_udp_port", record.Value)
		}
		if key == "TcpFootprint" {
			_ = d.Set("tcp_footprint", record.Value)
		}
		if key == "UdpFootprint" {
			_ = d.Set("udp_footprint", record.Value)
		}
		if key == "HasAbroad" {
			_ = d.Set("has_abroad", record.Value)
		}
		if key == "TcpSportStart" {
			_ = d.Set("tcp_start_port", record.Value)
		}
		if key == "TcpSportEnd" {
			_ = d.Set("tcp_end_port", record.Value)
		}
		if key == "UdpSportStart" {
			_ = d.Set("udp_start_port", record.Value)
		}
		if key == "UdpSportEnd" {
			_ = d.Set("udp_end_port", record.Value)
		}
		if key == "MaxUdpPackageLen" {
			_ = d.Set("max_udp_package_len", record.Value)
		}
		if key == "MinUdpPackageLen" {
			_ = d.Set("min_udp_package_len", record.Value)
		}
		if key == "MaxTcpPackageLen" {
			_ = d.Set("max_tcp_package_len", record.Value)
		}
		if key == "MinTcpPackageLen" {
			_ = d.Set("min_tcp_package_len", record.Value)
		}
		if key == "AppType" {
			_ = d.Set("app_type", record.Value)
		}
		if key == "AppProtocols" {
			_ = d.Set("app_protocols", strings.Split(*record.Value, ";"))
		}
		if key == "WebApiUrl" {
			_ = d.Set("web_api_urls", strings.Split(*record.Value, ";"))
		}
		if key == "PlatformTypes" {
			_ = d.Set("platform_types", strings.Split(*record.Value, ";"))
		}
		if key == "Id" {
			_ = d.Set("scene_id", record.Value)
		}
		if key == "CreateTime" {
			_ = d.Set("create_time", record.Value)
		}
	}
	return nil
}

func resourceTencentCloudDayuDdosPolicyCaseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_case.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDos policy case")
	}
	resourceType := items[0]
	sceneId := items[1]

	request := dayu.NewModifyDDoSPolicyCaseRequest()
	request.Business = &resourceType
	request.SceneId = &sceneId

	platforms := d.Get("platform_types").(*schema.Set).List()
	for _, plat := range platforms {
		request.PlatformTypes = append(request.PlatformTypes, helper.String(plat.(string)))
	}
	protocols := d.Get("app_protocols").(*schema.Set).List()
	for _, protocol := range protocols {
		request.AppProtocols = append(request.AppProtocols, helper.String(protocol.(string)))
	}
	urls := d.Get("web_api_urls").(*schema.Set).List()
	for _, url := range urls {
		request.WebApiUrl = append(request.WebApiUrl, helper.String(url.(string)))
	}
	request.AppType = helper.String(d.Get("app_type").(string))
	request.HasAbroad = helper.String(d.Get("has_abroad").(string))
	request.HasInitiateTcp = helper.String(d.Get("has_initiate_tcp").(string))
	request.HasInitiateUdp = helper.String(d.Get("has_initiate_udp").(string))
	request.HasVPN = helper.String(d.Get("has_vpn").(string))
	request.PeerTcpPort = helper.String(d.Get("peer_tcp_port").(string))
	request.PeerUdpPort = helper.String(d.Get("peer_udp_port").(string))
	request.TcpFootprint = helper.String(d.Get("tcp_footprint").(string))
	request.UdpFootprint = helper.String(d.Get("udp_footprint").(string))

	tcpPortStart := d.Get("tcp_start_port").(string)
	tcpPortEnd := d.Get("tcp_end_port").(string)
	startInt, sErr := strconv.Atoi(tcpPortStart)
	endInt, eErr := strconv.Atoi(tcpPortEnd)
	if sErr != nil {
		return sErr
	}
	if eErr != nil {
		return eErr
	}

	if endInt < startInt {
		return fmt.Errorf("`tcp_start_port`:%s should not be greater than `tcp_end_port`:%s.", tcpPortStart, tcpPortEnd)
	}
	udpPortStart := d.Get("udp_start_port").(string)
	udpPortEnd := d.Get("udp_end_port").(string)
	startInt, sErr = strconv.Atoi(udpPortStart)
	endInt, eErr = strconv.Atoi(udpPortEnd)
	if sErr != nil {
		return sErr
	}
	if eErr != nil {
		return eErr
	}

	if endInt < startInt {
		return fmt.Errorf("`udp_start_port`:%s should not be greater than `udp_end_port`:%s.", udpPortStart, udpPortEnd)
	}
	request.TcpSportStart = &tcpPortStart
	request.TcpSportEnd = &tcpPortEnd
	request.UdpSportStart = &udpPortStart
	request.UdpSportEnd = &udpPortEnd

	minTcpPackageLen := d.Get("min_tcp_package_len").(string)
	maxTcpPackageLen := d.Get("max_tcp_package_len").(string)
	minTcpPackageLenInt, _ := strconv.Atoi(minTcpPackageLen)
	maxTcpPackageLenInt, _ := strconv.Atoi(maxTcpPackageLen)
	if maxTcpPackageLenInt < minTcpPackageLenInt {
		return fmt.Errorf("`min_tcp_package_len`:%s should not be greater than `max_tcp_package_len`:%s.", minTcpPackageLen, maxTcpPackageLen)
	}
	minUdpPackageLen := d.Get("min_udp_package_len").(string)
	maxUdpPackageLen := d.Get("max_udp_package_len").(string)
	minUdpPackageLenInt, _ := strconv.Atoi(minUdpPackageLen)
	maxUdpPackageLenInt, _ := strconv.Atoi(maxUdpPackageLen)
	if maxUdpPackageLenInt < minUdpPackageLenInt {
		return fmt.Errorf("`min_udp_package_len`:%s should not be greater than `max_udp_package_len`:%s.", minUdpPackageLen, maxUdpPackageLen)
	}
	request.MinTcpPackageLen = &minTcpPackageLen
	request.MaxTcpPackageLen = &maxTcpPackageLen
	request.MinUdpPackageLen = &minUdpPackageLen
	request.MaxUdpPackageLen = &maxTcpPackageLen

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.ModifyDdosPolicyCase(ctx, request)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudDayuDdosPolicyCaseRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyCaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_case.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDos policy")
	}
	resourceType := items[0]
	sceneId := items[1]

	dayuService := DayuService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := dayuService.DeleteDdosPolicyCase(ctx, resourceType, sceneId)

	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "ResourceInUse" {
				//if bind automatically, try to unbind policy first
				//get the automatically generated policy
				policies, dErr := dayuService.DescribeDdosPolicies(ctx, resourceType, "")
				if dErr != nil {
					err = dErr
				}
				bindPolicyId := ""
				bindResourceIds := []string{}
				for _, policy := range policies {
					if *policy.SceneId == sceneId {
						bindPolicyId = *policy.PolicyId
						for _, resources := range (*policy).BoundResources {
							bindResourceIds = append(bindResourceIds, *resources)
						}
					}
				}
				if bindPolicyId == "" || len(bindResourceIds) == 0 {
					return fmt.Errorf("the automatically generated policy of policy case %s can not be find", sceneId)
				}
				//unbind policy and resource
				for _, resourceId := range bindResourceIds {
					bErr := dayuService.UnbindDdosPolicy(ctx, resourceId, resourceType, bindPolicyId)
					if bErr != nil {
						err = bErr
					}
				}
			}
		}
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := dayuService.DeleteDdosPolicyCase(ctx, resourceType, sceneId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := dayuService.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
	if err != nil || has {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
			if err != nil {
				return retryError(err)
			}

			if has {
				err = fmt.Errorf("delete DDoS policy case fail, DDoS policy case still exist from sdk DescribeDDosPolicy")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete DDoS policy case fail, DDoS policy case still exist from sdk DescribeDDosPolicy")
	}
}
