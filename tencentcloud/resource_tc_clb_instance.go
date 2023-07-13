/*
Provides a resource to create a CLB instance.

Example Usage

INTERNAL CLB

```hcl
resource "tencentcloud_clb_instance" "internal_clb" {
  network_type = "INTERNAL"
  clb_name     = "myclb"
  project_id   = 0
  vpc_id       = "vpc-7007ll7q"
  subnet_id    = "subnet-12rastkr"

  tags = {
    test = "tf"
  }
}
```

OPEN CLB

```hcl
resource "tencentcloud_clb_instance" "open_clb" {
  network_type              = "OPEN"
  clb_name                  = "myclb"
  project_id                = 0
  vpc_id                    = "vpc-da7ffa61"
  security_groups           = ["sg-o0ek7r93"]
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "vpc-da7ffa61"

  tags = {
    test = "tf"
  }
}
```

Dynamic Vip Instance

```hcl
resource "tencentcloud_security_group" "foo" {
  name = "clb-instance-open-sg"
}

resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-open-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb_open" {
  network_type              = "OPEN"
  clb_name                  = "clb-instance-open"
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.foo.id
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id
  security_groups           = [tencentcloud_security_group.foo.id]

  dynamic_vip = true

  tags = {
    test = "tf"
  }
}

output "domain" {
  value = tencentcloud_clb_instance.clb_open.domain
}
```

Default enable

```hcl
resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-1"
  name              = "sdk-feature-test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_security_group" "sglab" {
  name        = "sg_o0ek7r93"
  description = "favourite sg"
  project_id  = 0
}

resource "tencentcloud_vpc" "foo" {
  name         = "for-my-open-clb"
  cidr_block   = "10.0.0.0/16"

  tags = {
    "test" = "mytest"
  }
}

resource "tencentcloud_clb_instance" "open_clb" {
  network_type                 = "OPEN"
  clb_name                     = "my-open-clb"
  project_id                   = 0
  vpc_id                       = tencentcloud_vpc.foo.id
  load_balancer_pass_to_target = true

  security_groups              = [tencentcloud_security_group.sglab.id]
  target_region_info_region    = "ap-guangzhou"
  target_region_info_vpc_id    = tencentcloud_vpc.foo.id

  tags = {
    test = "open"
  }
}
```

CREATE multiple instance

```hcl
resource "tencentcloud_clb_instance" "open_clb1" {
  network_type              = "OPEN"
  clb_name = "hello"
  master_zone_id = "ap-guangzhou-3"
}
```

CREATE instance with log
```hcl
resource "tencentcloud_vpc" "vpc_test" {
  name = "clb-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "rtb_test" {
  name = "clb-test"
  vpc_id = "${tencentcloud_vpc.vpc_test.id}"
}

resource "tencentcloud_subnet" "subnet_test" {
  name = "clb-test"
  cidr_block = "10.0.1.0/24"
  availability_zone = "ap-guangzhou-3"
  vpc_id = "${tencentcloud_vpc.vpc_test.id}"
  route_table_id = "${tencentcloud_route_table.rtb_test.id}"
}

resource "tencentcloud_clb_log_set" "set" {
  period = 7
}

resource "tencentcloud_clb_log_topic" "topic" {
  log_set_id = "${tencentcloud_clb_log_set.set.id}"
  topic_name = "clb-topic"
}

resource "tencentcloud_clb_instance" "internal_clb" {
  network_type = "INTERNAL"
  clb_name = "myclb"
  project_id = 0
  vpc_id = "${tencentcloud_vpc.vpc_test.id}"
  subnet_id = "${tencentcloud_subnet.subnet_test.id}"
  load_balancer_pass_to_target = true
  log_set_id = "${tencentcloud_clb_log_set.set.id}"
  log_topic_id = "${tencentcloud_clb_log_topic.topic.id}"

  tags = {
    test = "tf"
  }
}

```

Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.foo lb-7a0t6zqb
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var clbActionMu = &sync.Mutex{}

func resourceTencentCloudClbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceCreate,
		Read:   resourceTencentCloudClbInstanceRead,
		Update: resourceTencentCloudClbInstanceUpdate,
		Delete: resourceTencentCloudClbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CLB_NETWORK_TYPE),
				Description:  "Type of CLB instance. Valid values: `OPEN` and `INTERNAL`.",
			},
			"clb_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the CLB. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'.",
			},
			"clb_vips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The virtual service address table of the CLB.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "ID of the project within the CLB instance, `0` - Default Project.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "VPC ID of the CLB.",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
				Description:  "Subnet ID of the CLB. Effective only for CLB within the VPC. Only supports `INTERNAL` CLBs. Default is `ipv4`.",
			},
			"address_ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IP version, only applicable to open CLB. Valid values are `ipv4`, `ipv6` and `IPv6FullChain`.",
			},
			"internet_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bandwidth package id. If set, the `internet_charge_type` must be `BANDWIDTH_PACKAGE`.",
			},
			"internet_bandwidth_max_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is MB.",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security groups of the CLB instance. Supports both `OPEN` and `INTERNAL` CLBs.",
			},
			"target_region_info_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.",
			},
			"target_region_info_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Vpc information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.",
			},
			"snat_pro": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether Binding IPs of other VPCs feature switch.",
			},
			"snat_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Snat Ip List, required with `snat_pro=true`. NOTE: This argument cannot be read and modified here because dynamic ip is untraceable, please import resource `tencentcloud_clb_snat_ip` to handle fixed ips.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Snat IP address, If set to empty will auto allocated.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Snat subnet ID.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The available tags within this CLB.",
			},
			"vip_isp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network operator, only applicable to open CLB. Valid values are `CMCC`(China Mobile), `CTCC`(Telecom), `CUCC`(China Unicom) and `BGP`. If this ISP is specified, network billing method can only use the bandwidth package billing (BANDWIDTH_PACKAGE).",
			},
			"load_balancer_pass_to_target": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the target allow flow come from clb. If value is true, only check security group of clb, or check both clb and backend instance security group.",
			},
			"master_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Setting master zone id of cross available zone disaster recovery, only applicable to open CLB.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Available zone id, only applicable to open CLB.",
			},
			"slave_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Setting slave zone id of cross available zone disaster recovery, only applicable to open CLB. this zone will undertake traffic when the master is down.",
			},
			"log_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of log set.",
			},
			"log_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of log topic.",
			},
			"dynamic_vip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If create dynamic vip CLB instance, `true` or `false`.",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain name of the CLB instance.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)

	networkType := d.Get("network_type").(string)
	clbName := d.Get("clb_name").(string)
	flag, e := checkSameName(clbName, meta)
	if e != nil {
		return e
	}
	if flag {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: Same CLB name %s exists!", clbName)
	}
	targetRegionInfoRegion := ""
	targetRegionInfoVpcId := ""
	if v, ok := d.GetOk("target_region_info_region"); ok {
		targetRegionInfoRegion = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with target_region_info")
		}
	}
	if v, ok := d.GetOk("target_region_info_vpc_id"); ok {
		targetRegionInfoVpcId = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with target_region_info")
		}
	}
	if (targetRegionInfoRegion != "" && targetRegionInfoVpcId == "") || (targetRegionInfoRegion == "" && targetRegionInfoVpcId != "") {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: region and vpc_id must be set at same time")
	}

	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerType = helper.String(networkType)
	request.LoadBalancerName = helper.String(clbName)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		if networkType == CLB_NETWORK_TYPE_OPEN {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: OPEN network_type do not support this operation with subnet_id")
		}
		request.SubnetId = helper.String(v.(string))
	}

	//vip
	if v, ok := d.GetOk("vip_isp"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support vip ISP setting")
		}
		request.VipIsp = helper.String(v.(string))
	}

	//ip version
	if v, ok := d.GetOk("address_ip_version"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support IP version setting")
		}
		request.AddressIPVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snat_pro"); ok {
		request.SnatPro = helper.Bool(v.(bool))
	}

	if v, ok := d.Get("snat_ips").([]interface{}); ok && len(v) > 0 {
		for i := range v {
			item := v[i].(map[string]interface{})
			subnetId := item["subnet_id"].(string)
			snatIp := &clb.SnatIp{
				SubnetId: &subnetId,
			}
			if v, ok := item["ip"].(string); ok && v != "" {
				snatIp.Ip = &v
			}
			request.SnatIps = append(request.SnatIps, snatIp)
		}
	}

	v, ok := d.GetOk("internet_charge_type")
	bv, bok := d.GetOk("internet_bandwidth_max_out")
	pv, pok := d.GetOk("bandwidth_package_id")

	chargeType := v.(string)

	//internet charge type
	if ok || bok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support internet charge type setting")
		}
		request.InternetAccessible = &clb.InternetAccessible{}
		if ok {
			request.InternetAccessible.InternetChargeType = helper.String(chargeType)
		}
		if bok {
			request.InternetAccessible.InternetMaxBandwidthOut = helper.IntInt64(bv.(int))
		}
		if pok {
			if pok && chargeType != INTERNET_CHARGE_TYPE_BANDWIDTH_PACKAGE {
				return fmt.Errorf("[CHECK][CLB instance][Create] check: internet_charge_type must `BANDWIDTH_PACKAGE` when bandwidth_package_id was set")
			}
			request.BandwidthPackageId = helper.String(pv.(string))
		} else if chargeType == INTERNET_CHARGE_TYPE_BANDWIDTH_PACKAGE {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: the `bandwidth_package_id` must be specified if internet_charge_type was `BANDWIDTH_PACKAGE`")
		}
	}

	if v, ok := d.GetOk("master_zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support master zone id setting")
		}
		request.MasterZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support zone id setting")
		}
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("slave_zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support slave zone id setting")
		}
		request.SlaveZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_pass_to_target"); ok {
		request.LoadBalancerPassToTarget = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("dynamic_vip"); ok {
		request.DynamicVip = helper.Bool(v.(bool))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			tmpKey := k
			tmpValue := v
			request.Tags = append(request.Tags, &clb.TagInfo{
				TagKey:   &tmpKey,
				TagValue: &tmpValue,
			})
		}
	}

	clbId := ""
	var response *clb.CreateLoadBalancerResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().CreateLoadBalancer(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return retryError(errors.WithStack(retryErr))
			}
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
		return err
	}
	if len(response.Response.LoadBalancerIds) < 1 {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: response error, load balancer id is nil")
	}
	d.SetId(*response.Response.LoadBalancerIds[0])
	clbId = *response.Response.LoadBalancerIds[0]

	if v, ok := d.GetOk("security_groups"); ok {
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = helper.String(clbId)
		securityGroups := v.([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			if securityGroups[i] != nil {
				securityGroup := securityGroups[i].(string)
				sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
			}
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			sgResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId

				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create CLB instance security_groups failed, reason:%+v", logId, err)
			return err
		}
	}

	if v, ok := d.GetOk("log_set_id"); ok {
		if u, ok := d.GetOk("log_topic_id"); ok {
			logRequest := clb.NewSetLoadBalancerClsLogRequest()
			logRequest.LoadBalancerId = helper.String(clbId)
			logRequest.LogSetId = helper.String(v.(string))
			logRequest.LogTopicId = helper.String(u.(string))
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				logResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerClsLog(logRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, logRequest.GetAction(), logRequest.ToJsonString(), logResponse.ToJsonString())
					requestId := *logResponse.Response.RequestId

					retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
					if retryErr != nil {
						return retryError(errors.WithStack(retryErr))
					}
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s set CLB instance log failed, reason:%+v", logId, err)
				return err
			}

		} else {
			return fmt.Errorf("log_topic_id and log_set_id must be set together.")
		}
	}

	if targetRegionInfoRegion != "" {
		isLoadBalancePassToTgt := d.Get("load_balancer_pass_to_target").(bool)
		targetRegionInfo := clb.TargetRegionInfo{
			Region: &targetRegionInfoRegion,
			VpcId:  &targetRegionInfoVpcId,
		}
		mRequest := clb.NewModifyLoadBalancerAttributesRequest()
		mRequest.LoadBalancerId = helper.String(clbId)
		mRequest.TargetRegionInfo = &targetRegionInfo
		mRequest.LoadBalancerPassToTarget = &isLoadBalancePassToTgt
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			mResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerAttributes(mRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
				requestId := *mResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClbInstanceRead(d, meta)
}

func resourceTencentCloudClbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *clb.LoadBalancer
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeLoadBalancerById(ctx, clbId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB instance failed, reason:%+v", logId, err)
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("network_type", instance.LoadBalancerType)
	_ = d.Set("clb_name", instance.LoadBalancerName)
	_ = d.Set("clb_vips", helper.StringsInterfaces(instance.LoadBalancerVips))
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("target_region_info_region", instance.TargetRegionInfo.Region)
	_ = d.Set("target_region_info_vpc_id", instance.TargetRegionInfo.VpcId)
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecureGroups))
	_ = d.Set("domain", instance.LoadBalancerDomain)

	if instance.VipIsp != nil {
		_ = d.Set("vip_isp", instance.VipIsp)
	}
	if instance.AddressIPVersion != nil {
		_ = d.Set("address_ip_version", instance.AddressIPVersion)
	}
	if instance.NetworkAttributes != nil {
		_ = d.Set("internet_bandwidth_max_out", instance.NetworkAttributes.InternetMaxBandwidthOut)
		_ = d.Set("internet_charge_type", instance.NetworkAttributes.InternetChargeType)
	}

	_ = d.Set("load_balancer_pass_to_target", instance.LoadBalancerPassToTarget)
	//_ = d.Set("master_zone_id", instance.MasterZone.ZoneId)
	//_ = d.Set("zone_id", instance.Zones)
	//_ = d.Set("slave_zone_id", instance.MasterZone)
	_ = d.Set("log_set_id", instance.LogSetId)
	_ = d.Set("log_topic_id", instance.LogTopicId)

	if _, ok := d.GetOk("snat_pro"); ok {
		_ = d.Set("snat_pro", instance.SnatPro)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "clb", "clb", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudClbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	d.Partial(true)

	clbId := d.Id()
	request := clb.NewModifyLoadBalancerAttributesRequest()
	request.LoadBalancerId = helper.String(clbId)
	clbName := ""
	targetRegionInfo := clb.TargetRegionInfo{}
	internet := clb.InternetAccessible{}
	changed := false
	isLoadBalancerPassToTgt := false
	snatPro := d.Get("snat_pro").(bool)

	if d.HasChange("clb_name") {
		changed = true
		clbName = d.Get("clb_name").(string)
		flag, err := checkSameName(clbName, meta)
		if err != nil {
			return err
		}
		if flag {
			return fmt.Errorf("[CHECK][CLB instance][Update] check: Same CLB name %s exists!", clbName)
		}
		request.LoadBalancerName = helper.String(clbName)
	}

	if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc_id") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance %s][Update] check: INTERNAL network_type do not support this operation with target_region_info", clbId)
		}
		changed = true
		region := d.Get("target_region_info_region").(string)
		vpcId := d.Get("target_region_info_vpc_id").(string)
		targetRegionInfo = clb.TargetRegionInfo{
			Region: &region,
			VpcId:  &vpcId,
		}
		request.TargetRegionInfo = &targetRegionInfo
	}

	if d.HasChange("internet_charge_type") || d.HasChange("internet_bandwidth_max_out") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance %s][Update] check: INTERNAL network_type do not support this operation with internet setting", clbId)
		}
		changed = true
		chargeType := d.Get("internet_charge_type").(string)
		bandwidth := d.Get("internet_bandwidth_max_out").(int)
		if chargeType != "" {
			internet.InternetChargeType = &chargeType
		}
		if bandwidth > 0 {
			internet.InternetMaxBandwidthOut = helper.IntInt64(bandwidth)
		}
		request.InternetChargeInfo = &internet
	}

	if d.HasChange("load_balancer_pass_to_target") {
		changed = true
		isLoadBalancerPassToTgt = d.Get("load_balancer_pass_to_target").(bool)
		request.LoadBalancerPassToTarget = &isLoadBalancerPassToTgt
	}

	if d.HasChange("snat_pro") {
		changed = true
		request.SnatPro = &snatPro
	}

	immutableArgs := []string{"snat_ips", "dynamic_vip"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if changed {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerAttributes(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(retryErr)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("security_groups") {

		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = helper.String(clbId)
		securityGroups := d.Get("security_groups").([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			sgResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance security_group failed, reason:%+v", logId, err)
			return err
		}

	}

	if d.HasChange("log_set_id") || d.HasChange("log_topic_id") {
		logSetId := d.Get("log_set_id")
		logTopicId := d.Get("log_topic_id")
		logRequest := clb.NewSetLoadBalancerClsLogRequest()
		logRequest.LoadBalancerId = helper.String(clbId)
		logRequest.LogSetId = helper.String(logSetId.(string))
		logRequest.LogTopicId = helper.String(logTopicId.(string))
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			logResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerClsLog(logRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, logRequest.GetAction(), logRequest.ToJsonString(), logResponse.ToJsonString())
				requestId := *logResponse.Response.RequestId

				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s set CLB instance log failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("clb", "clb", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}
	d.Partial(false)

	return nil
}

func resourceTencentCloudClbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteLoadBalancerById(ctx, clbId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB instance failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func checkSameName(name string, meta interface{}) (flag bool, errRet error) {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	flag = false
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	params := make(map[string]interface{})
	params["clb_name"] = name
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		clbs, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		if len(clbs) > 0 {
			//this describe function is a fuzzy query
			// so take a further check
			for _, clb := range clbs {
				if *clb.LoadBalancerName == name {
					flag = true
					return nil
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB instance failed, reason:%+v", logId, err)
	}
	errRet = err
	return
}
