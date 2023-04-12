/*
Provides a resource to create a cvm chc_deploy_vpc

Example Usage

```hcl
resource "tencentcloud_cvm_chc_deploy_vpc" "chc_deploy_vpc" {
  chc_id = "chc-xxxxx"
  deploy_virtual_private_cloud {
		vpc_id = "vpc-xxxxx"
		subnet_id = "subnet-xxxxx"
  }
  deploy_security_group_ids = ["sg-xxxxx"]
}
```

Import

cvm chc_deploy_vpc can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_chc_deploy_vpc.chc_deploy_vpc chc_deploy_vpc_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCvmChcDeployVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmChcDeployVpcCreate,
		Read:   resourceTencentCloudCvmChcDeployVpcRead,
		Delete: resourceTencentCloudCvmChcDeployVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"chc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CHC host ID.",
			},

			"deploy_virtual_private_cloud": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Deployment network information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.",
						},
						"private_ip_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of IPv6 addresses randomly generated for the ENI.",
						},
					},
				},
			},

			"deploy_security_group_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Deployment network security group list.",
			},
		},
	}
}

func resourceTencentCloudCvmChcDeployVpcCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_deploy_vpc.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = cvm.NewConfigureChcDeployVpcRequest()
		chcId   string
		vpcId   string
	)
	chcId = d.Get("chc_id").(string)
	request.ChcIds = []*string{&chcId}

	if dMap, ok := helper.InterfacesHeadMap(d, "deploy_virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			vpcId = v.(string)
			virtualPrivateCloud.VpcId = helper.String(vpcId)
		}
		if v, ok := dMap["subnet_id"]; ok {
			virtualPrivateCloud.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["as_vpc_gateway"]; ok {
			virtualPrivateCloud.AsVpcGateway = helper.Bool(v.(bool))
		}
		if v, ok := dMap["private_ip_addresses"]; ok {
			privateIpAddressesSet := v.(*schema.Set).List()
			for i := range privateIpAddressesSet {
				privateIpAddresses := privateIpAddressesSet[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		request.DeployVirtualPrivateCloud = &virtualPrivateCloud
	}

	if v, ok := d.GetOk("deploy_security_group_ids"); ok {
		deploySecurityGroupIdsSet := v.(*schema.Set).List()
		for i := range deploySecurityGroupIdsSet {
			deploySecurityGroupIds := deploySecurityGroupIdsSet[i].(string)
			request.DeploySecurityGroupIds = append(request.DeploySecurityGroupIds, &deploySecurityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ConfigureChcDeployVpc(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm chcDeployVpc failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(chcId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{vpcId}, 10*readRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCvmChcDeployVpcRead(d, meta)
}

func resourceTencentCloudCvmChcDeployVpcRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_deploy_vpc.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	chcId := d.Id()

	params := map[string]interface{}{
		"chc_ids": []string{chcId},
	}
	chcDeployVpcs, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}

	if len(chcDeployVpcs) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmChcDeployVpc` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	chcDeployVpc := chcDeployVpcs[0]
	if chcDeployVpc.ChcId != nil {
		_ = d.Set("chc_id", chcDeployVpc.ChcId)
	}

	if chcDeployVpc.DeployVirtualPrivateCloud != nil {
		deployVirtualPrivateCloudMap := map[string]interface{}{}

		if chcDeployVpc.DeployVirtualPrivateCloud.VpcId != nil {
			deployVirtualPrivateCloudMap["vpc_id"] = chcDeployVpc.DeployVirtualPrivateCloud.VpcId
		}

		if chcDeployVpc.DeployVirtualPrivateCloud.SubnetId != nil {
			deployVirtualPrivateCloudMap["subnet_id"] = chcDeployVpc.DeployVirtualPrivateCloud.SubnetId
		}

		if chcDeployVpc.DeployVirtualPrivateCloud.AsVpcGateway != nil {
			deployVirtualPrivateCloudMap["as_vpc_gateway"] = chcDeployVpc.DeployVirtualPrivateCloud.AsVpcGateway
		}

		if chcDeployVpc.DeployVirtualPrivateCloud.PrivateIpAddresses != nil {
			deployVirtualPrivateCloudMap["private_ip_addresses"] = chcDeployVpc.DeployVirtualPrivateCloud.PrivateIpAddresses
		}

		if chcDeployVpc.DeployVirtualPrivateCloud.Ipv6AddressCount != nil {
			deployVirtualPrivateCloudMap["ipv6_address_count"] = chcDeployVpc.DeployVirtualPrivateCloud.Ipv6AddressCount
		}

		_ = d.Set("deploy_virtual_private_cloud", []interface{}{deployVirtualPrivateCloudMap})
	}

	if chcDeployVpc.DeploySecurityGroupIds != nil {
		_ = d.Set("deploy_security_group_ids", chcDeployVpc.DeploySecurityGroupIds)
	}

	return nil
}

func resourceTencentCloudCvmChcDeployVpcDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_deploy_vpc.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	chcId := d.Id()

	request := cvm.NewRemoveChcDeployVpcRequest()
	request.ChcIds = []*string{&chcId}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().RemoveChcDeployVpc(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s remove Chc deploy vpc failed, reason:%+v", logId, err)
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{""}, 5*readRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
