/*
Provides a resource to create a vpc eni_ipv4_address

Example Usage

```hcl
resource "tencentcloud_vpc_eni_ipv4_address" "eni_ipv4_address" {
  network_interface_id = ""
  private_ip_addresses {
		private_ip_address = ""
		primary =
		public_ip_address = ""
		address_id = ""
		description = ""
		is_wan_ip_blocked =
		state = ""
		qos_level = ""

  }
  secondary_private_ip_address_count =
  qos_level = ""
}
```

Import

vpc eni_ipv4_address can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_eni_ipv4_address.eni_ipv4_address eni_ipv4_address_id
```
*/
package vpc

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudEniIpv4Address() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEniIpv4AddressCreate,
		Read:   resourceTencentCloudEniIpv4AddressRead,
		Delete: resourceTencentCloudEniIpv4AddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the ENI instance, such as `eni-m6dyj72l`.",
			},

			"private_ip_addresses": {
				Optional:    true,
				Type:        schema.TypeSet,
				Computed: true,
				ForceNew:      true,
				ConflictsWith: []string{"secondary_private_ip_address_count", "qos_level"},
				Description: "The information on private IP addresses, of which you can specify a maximum of 10 at a time. You should provide either this parameter or SecondaryPrivateIpAddressCount, or both.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:      true,
							Description: "Private IP address.",
						},
						"primary": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether it is a primary IP.",
						},
						"public_ip_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public IP address.",
						},
						"address_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "EIP instance ID, such as `eip-11112222`.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Private IP description.",
						},
						"is_wan_ip_blocked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the public IP is blocked.",
						},
						"state": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP status: `PENDING`: Creating, `MIGRATING`: Migrating, `DELETING`: Deleting, `AVAILABLE`: Available.",
						},
						"qos_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP service level. Values: PT` (Gold), `AU` (Silver), `AG `(Bronze) and DEFAULT` (Default).",
						},
					},
				},
			},

			"secondary_private_ip_address_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"private_ip_addresses"},
				Description: "The number of newly-applied private IP addresses. You should provide either this parameter or PrivateIpAddresses, or both. The total number of private IP addresses cannot exceed the quota. For more information, see&amp;amp;lt;a href=&amp;amp;quot;/document/product/576/18527&amp;amp;quot;&amp;amp;gt;ENI Use Limits&amp;amp;lt;/a&amp;amp;gt;.",
			},

			"qos_level": {
				Optional:    true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"private_ip_addresses"},
				Type:        schema.TypeString,
				Description: "IP service level. It is used together with `SecondaryPrivateIpAddressCount`. Values: PT`(Gold), `AU`(Silver), `AG `(Bronze) and DEFAULT (Default).",
			},
		},
	}
}

func resourceTencentCloudEniIpv4AddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni_ipv4_address.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request            = vpc.NewAssignPrivateIpAddressesRequest()
		response           = vpc.NewAssignPrivateIpAddressesResponse()
		networkInterfaceId string
	)
	if v, ok := d.GetOk("network_interface_id"); ok {
		networkInterfaceId = v.(string)
		request.NetworkInterfaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("private_ip_addresses"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			privateIpAddressSpecification := vpc.PrivateIpAddressSpecification{}
			if v, ok := dMap["private_ip_address"]; ok {
				privateIpAddressSpecification.PrivateIpAddress = helper.String(v.(string))
			}
			if v, ok := dMap["primary"]; ok {
				privateIpAddressSpecification.Primary = helper.Bool(v.(bool))
			}
			if v, ok := dMap["public_ip_address"]; ok {
				privateIpAddressSpecification.PublicIpAddress = helper.String(v.(string))
			}
			if v, ok := dMap["address_id"]; ok {
				privateIpAddressSpecification.AddressId = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				privateIpAddressSpecification.Description = helper.String(v.(string))
			}
			if v, ok := dMap["is_wan_ip_blocked"]; ok {
				privateIpAddressSpecification.IsWanIpBlocked = helper.Bool(v.(bool))
			}
			if v, ok := dMap["state"]; ok {
				privateIpAddressSpecification.State = helper.String(v.(string))
			}
			if v, ok := dMap["qos_level"]; ok {
				privateIpAddressSpecification.QosLevel = helper.String(v.(string))
			}
			request.PrivateIpAddresses = append(request.PrivateIpAddresses, &privateIpAddressSpecification)
		}
	}

	if v, ok := d.GetOkExists("secondary_private_ip_address_count"); ok {
		request.SecondaryPrivateIpAddressCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("qos_level"); ok {
		request.QosLevel = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssignPrivateIpAddresses(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc eniIpv4Address failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(networkInterfaceId)

	return resourceTencentCloudEniIpv4AddressRead(d, meta)
}

func resourceTencentCloudEniIpv4AddressRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni_ipv4_address.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	networkInterfaceId := d.Id()

	enis, err := service.DescribeEniById(ctx, []string{networkInterfaceId})

	if err != nil {
		return err
	}

	if len(enis) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6EniAddress` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	eni := enis[0]

	if eni.NetworkInterfaceId != nil {
		_ = d.Set("network_interface_id", eni.NetworkInterfaceId)
	}

	ipv4s := make([]map[string]interface{}, 0, len(eni.PrivateIpAddressSet))
	for _, ipv4 := range eni.PrivateIpAddressSet {
		ipv4s = append(ipv4s, map[string]interface{}{
			"private_ip_address": ipv4.PrivateIpAddress,
			"primary":            ipv4.Primary,
			"public_ip_address":  ipv4.AddressId,
			"address_id":         ipv4.AddressId,
			"description":        ipv4.Description,
			"is_wan_ip_blocked":  ipv4.IsWanIpBlocked,
			"state":              ipv4.State,
			"qos_level":          ipv4.QosLevel,
		})
	}

	_ = d.Set("network_interface_id", networkInterfaceId)
	_ = d.Set("private_ip_addresses", ipv4s)
	_ = d.Set("secondary_private_ip_address_count", len(eni.PrivateIpAddressSet))
	if len(eni.PrivateIpAddressSet) > 0 {
		_ = d.Set("qos_level", eni.PrivateIpAddressSet[0].QosLevel)
	}

	return nil
}

func resourceTencentCloudEniIpv4AddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_ipv4_address.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	networkInterfaceId := d.Id()

	if err := service.DeleteVpcEniIpv4AddressById(ctx, networkInterfaceId); err != nil {
		return err
	}

	return nil
}
