/*
Provides a resource to create a cvm chc_deploy_vpc

Example Usage

```hcl
resource "tencentcloud_cvm_chc_deploy_vpc" "chc_deploy_vpc" {
  chc_ids =
  deploy_virtual_private_cloud {
		vpc_id = ""
		subnet_id = ""
		as_vpc_gateway =
		private_ip_addresses =
		ipv6_address_count =

  }
  deploy_security_group_ids =
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
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
			"chc_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "CHC host IDs.",
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
							Description: "Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;TRUE: yes;&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;FALSE: no&amp;lt;br&amp;gt;&amp;lt;br&amp;gt;Default: FALSE.",
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
		request  = cvm.NewConfigureChcDeployVpcRequest()
		response = cvm.NewConfigureChcDeployVpcResponse()
		chcId    string
		vpcId    string
		subnetId string
	)
	if v, ok := d.GetOk("chc_ids"); ok {
		chcIdsSet := v.(*schema.Set).List()
		for i := range chcIdsSet {
			chcIds := chcIdsSet[i].(string)
			request.ChcIds = append(request.ChcIds, &chcIds)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "deploy_virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			virtualPrivateCloud.VpcId = helper.String(v.(string))
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm chcDeployVpc failed, reason:%+v", logId, err)
		return err
	}

	chcId = *response.Response.ChcId
	d.SetId(strings.Join([]string{chcId, vpcId, subnetId}, FILED_SP))

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"PREPARED"}, 10*readRetryTimeout, time.Second, service.CvmChcDeployVpcStateRefreshFunc(d.Id(), []string{}))

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

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	chcId := idSplit[0]
	vpcId := idSplit[1]
	subnetId := idSplit[2]

	chcDeployVpc, err := service.DescribeCvmChcDeployVpcById(ctx, chcId, vpcId, subnetId)
	if err != nil {
		return err
	}

	if chcDeployVpc == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmChcDeployVpc` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if chcDeployVpc.ChcIds != nil {
		_ = d.Set("chc_ids", chcDeployVpc.ChcIds)
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	chcId := idSplit[0]
	vpcId := idSplit[1]
	subnetId := idSplit[2]

	if err := service.DeleteCvmChcDeployVpcById(ctx, chcId, vpcId, subnetId); err != nil {
		return err
	}

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"READY"}, 5*readRetryTimeout, time.Second, service.CvmChcDeployVpcStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
