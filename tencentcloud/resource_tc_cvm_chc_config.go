/*
Provides a resource to create a cvm chc_config

Example Usage

```hcl
resource "tencentcloud_cvm_chc_config" "chc_config" {
  chc_id = "chc-xxxxxx"
  instance_name = "xxxxxx"
  bmc_user = "admin"
  password = "xxxxxx"
    bmc_virtual_private_cloud {
		vpc_id = "vpc-xxxxxx"
		subnet_id = "subnet-xxxxxx"

  }
  bmc_security_group_ids = ["sg-xxxxxx"]

  deploy_virtual_private_cloud {
		vpc_id = "vpc-xxxxxx"
		subnet_id = "subnet-xxxxxx"
  }
  deploy_security_group_ids = ["sg-xxxxxx"]
}
```

Import

cvm chc_config can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_chc_config.chc_config chc_config_id
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

func resourceTencentCloudCvmChcConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmChcConfigCreate,
		Update: resourceTencentCloudCvmChcConfigUpdate,
		Read:   resourceTencentCloudCvmChcConfigRead,
		Delete: resourceTencentCloudCvmChcConfigDelete,
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

			"instance_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "CHC host name.",
			},

			"device_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Server type.",
			},

			"bmc_user": {
				Optional:     true,
				RequiredWith: []string{"password"},
				Type:         schema.TypeString,
				Description:  "Valid characters: Letters, numbers, hyphens and underscores. Only set when update password.",
			},

			"password": {
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"bmc_user"},
				Type:         schema.TypeString,

				Description: "The password can contain 8 to 16 characters, including letters, numbers and special symbols (()`~!@#$%^&amp;amp;*-+=_|{}).",
			},

			"bmc_virtual_private_cloud": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Out-of-band network information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.",
						},
						"private_ip_addresses": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Number of IPv6 addresses randomly generated for the ENI.",
						},
					},
				},
			},

			"bmc_security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				RequiredWith: []string{"bmc_virtual_private_cloud"},
				Description:  "Out-of-band network security group list.",
			},

			"deploy_virtual_private_cloud": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Deployment network information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.",
						},
						"private_ip_addresses": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Number of IPv6 addresses randomly generated for the ENI.",
						},
					},
				},
			},

			"deploy_security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				RequiredWith: []string{"deploy_virtual_private_cloud"},
				Description:  "Deployment network security group list.",
			},
		},
	}
}

func resourceTencentCloudCvmChcConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		assistChange  bool
		deployChange  bool
		chcId         string
		vpcId         string
		assistRequest = cvm.NewConfigureChcAssistVpcRequest()
		deployRequest = cvm.NewConfigureChcDeployVpcRequest()
	)
	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	if v, ok := d.GetOk("chc_id"); ok {
		chcId = v.(string)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.InstanceName = helper.String(v.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}

	if v, ok := d.GetOk("device_type"); ok {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.DeviceType = helper.String(v.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}
	bmcUser, bmcUserok := d.GetOk("bmc_user")
	password, passwordOk := d.GetOk("password")
	if bmcUserok && passwordOk {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.BmcUser = helper.String(bmcUser.(string))
		attributeRequest.Password = helper.String(password.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "bmc_virtual_private_cloud"); ok {
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
			privateIpAddresses := v.([]interface{})
			for i := range privateIpAddresses {
				privateIpAddresses := privateIpAddresses[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		assistChange = true
		assistRequest.BmcVirtualPrivateCloud = &virtualPrivateCloud
	}

	if v, ok := d.GetOk("bmc_security_group_ids"); ok {
		bmcSecurityGroupIds := v.([]interface{})
		for i := range bmcSecurityGroupIds {
			bmcSecurityGroupIds := bmcSecurityGroupIds[i].(string)
			assistRequest.BmcSecurityGroupIds = append(assistRequest.BmcSecurityGroupIds, &bmcSecurityGroupIds)
		}
		assistChange = true
	}

	if assistChange {
		assistRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ConfigureChcAssistVpc(assistRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, assistRequest.GetAction(), assistRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cvm chcAssistVpc failed, reason:%+v", logId, err)
			return err
		}
		conf := BuildStateChangeConf([]string{}, []string{"READY"}, 20*readRetryTimeout, time.Second, service.CvmChcInstanceStateRefreshFunc(chcId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

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
			privateIpAddresses := v.([]interface{})
			for i := range privateIpAddresses {
				privateIpAddresses := privateIpAddresses[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		deployRequest.DeployVirtualPrivateCloud = &virtualPrivateCloud
		deployChange = true
	}

	if v, ok := d.GetOk("deploy_security_group_ids"); ok {
		deploySecurityGroupIds := v.([]interface{})
		for i := range deploySecurityGroupIds {
			deploySecurityGroupIds := deploySecurityGroupIds[i].(string)
			deployRequest.DeploySecurityGroupIds = append(deployRequest.DeploySecurityGroupIds, &deploySecurityGroupIds)
		}
		deployChange = true
	}

	if deployChange {
		deployRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ConfigureChcDeployVpc(deployRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deployRequest.GetAction(), deployRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cvm chcDeployVpc failed, reason:%+v", logId, err)
			return err
		}

		conf := BuildStateChangeConf([]string{}, []string{vpcId}, 10*readRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(chcId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	d.SetId(chcId)

	return resourceTencentCloudCvmChcConfigRead(d, meta)
}

func resourceTencentCloudCvmChcConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	chcId := d.Id()

	if d.HasChange("instance_name") {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.ChcIds = []*string{&chcId}
		if v, ok := d.GetOk("instance_name"); ok {
			attributeRequest.InstanceName = helper.String(v.(string))
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}
	if d.HasChange("device_type") {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.ChcIds = []*string{&chcId}
		if v, ok := d.GetOk("device_type"); ok {
			attributeRequest.DeviceType = helper.String(v.(string))
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}
	if d.HasChange("bmc_user") || d.HasChange("password") {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.ChcIds = []*string{&chcId}
		if v, ok := d.GetOk("bmc_user"); ok {
			attributeRequest.BmcUser = helper.String(v.(string))
		}

		if v, ok := d.GetOk("password"); ok {
			attributeRequest.Password = helper.String(v.(string))
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudCvmChcConfigRead(d, meta)
}
func resourceTencentCloudCvmChcConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	chcId := d.Id()

	params := map[string]interface{}{
		"chc_ids": []string{chcId},
	}
	chcHosts, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}

	if len(chcHosts) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmChcAssistVpc` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	chcHost := chcHosts[0]
	if chcHost.ChcId != nil {
		_ = d.Set("chc_id", chcHost.ChcId)
	}
	_ = d.Set("instance_name", chcHost.InstanceName)
	_ = d.Set("device_type", chcHost.DeviceType)
	if chcHost.BmcVirtualPrivateCloud != nil {
		bmcVirtualPrivateCloudMap := map[string]interface{}{}

		if chcHost.BmcVirtualPrivateCloud.VpcId != nil {
			bmcVirtualPrivateCloudMap["vpc_id"] = chcHost.BmcVirtualPrivateCloud.VpcId
		}

		if chcHost.BmcVirtualPrivateCloud.SubnetId != nil {
			bmcVirtualPrivateCloudMap["subnet_id"] = chcHost.BmcVirtualPrivateCloud.SubnetId
		}

		if chcHost.BmcVirtualPrivateCloud.AsVpcGateway != nil {
			bmcVirtualPrivateCloudMap["as_vpc_gateway"] = chcHost.BmcVirtualPrivateCloud.AsVpcGateway
		}

		if chcHost.BmcVirtualPrivateCloud.PrivateIpAddresses != nil {
			privateIpAddresses := make([]string, 0)
			for _, p := range chcHost.BmcVirtualPrivateCloud.PrivateIpAddresses {
				privateIpAddresses = append(privateIpAddresses, *p)
			}
			bmcVirtualPrivateCloudMap["private_ip_addresses"] = privateIpAddresses
		}

		if chcHost.BmcVirtualPrivateCloud.Ipv6AddressCount != nil {
			bmcVirtualPrivateCloudMap["ipv6_address_count"] = chcHost.BmcVirtualPrivateCloud.Ipv6AddressCount
		}

		_ = d.Set("bmc_virtual_private_cloud", []interface{}{bmcVirtualPrivateCloudMap})
	}

	if chcHost.BmcSecurityGroupIds != nil {
		bmcSecurityGroupIds := make([]string, 0)
		for _, sgId := range chcHost.BmcSecurityGroupIds {
			bmcSecurityGroupIds = append(bmcSecurityGroupIds, *sgId)
		}
		_ = d.Set("bmc_security_group_ids", bmcSecurityGroupIds)
	}

	if chcHost.DeployVirtualPrivateCloud != nil {
		deployVirtualPrivateCloudMap := map[string]interface{}{}

		if chcHost.DeployVirtualPrivateCloud.VpcId != nil {
			deployVirtualPrivateCloudMap["vpc_id"] = chcHost.DeployVirtualPrivateCloud.VpcId
		}

		if chcHost.DeployVirtualPrivateCloud.SubnetId != nil {
			deployVirtualPrivateCloudMap["subnet_id"] = chcHost.DeployVirtualPrivateCloud.SubnetId
		}

		if chcHost.DeployVirtualPrivateCloud.AsVpcGateway != nil {
			deployVirtualPrivateCloudMap["as_vpc_gateway"] = chcHost.DeployVirtualPrivateCloud.AsVpcGateway
		}

		if chcHost.DeployVirtualPrivateCloud.PrivateIpAddresses != nil {
			privateIpAddresses := make([]string, 0)
			for _, p := range chcHost.DeployVirtualPrivateCloud.PrivateIpAddresses {
				privateIpAddresses = append(privateIpAddresses, *p)
			}
			deployVirtualPrivateCloudMap["private_ip_addresses"] = privateIpAddresses
		}

		if chcHost.DeployVirtualPrivateCloud.Ipv6AddressCount != nil {
			deployVirtualPrivateCloudMap["ipv6_address_count"] = chcHost.DeployVirtualPrivateCloud.Ipv6AddressCount
		}

		_ = d.Set("deploy_virtual_private_cloud", []interface{}{deployVirtualPrivateCloudMap})
	}

	if chcHost.DeploySecurityGroupIds != nil {
		deploySecurityGroupIds := make([]string, 0)
		for _, sgId := range chcHost.DeploySecurityGroupIds {
			deploySecurityGroupIds = append(deploySecurityGroupIds, *sgId)
		}
		_ = d.Set("deploy_security_group_ids", deploySecurityGroupIds)
	}

	return nil
}

func resourceTencentCloudCvmChcConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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

	params := map[string]interface{}{
		"chc_ids": []string{chcId},
	}
	chcHosts, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}
	if len(chcHosts) > 0 && *chcHosts[0].InstanceState == "INIT" {
		return nil
	}

	if err := service.DeleteCvmChcAssistVpcById(ctx, chcId); err != nil {
		return err
	}

	conf = BuildStateChangeConf([]string{}, []string{"INIT"}, 10*readRetryTimeout, time.Second, service.CvmChcInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
