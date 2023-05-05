/*
Provide a datasource to query EKS cluster credential info (offlined).

~> **NOTE:**  This resource was offline no longer supported.

Example Usage

```hcl
data "tencentcloud_eks_cluster_credential" "foo" {
  cluster_id = "cls-xxxxxxxx"
}

# example outputs
output "addresses" {
  value = data.tencentcloud_eks_cluster_credential.cred.addresses
}

output "ca_cert" {
  value = data.tencentcloud_eks_cluster_credential.cred.credential.ca_cert
}

output "token" {
  value = data.tencentcloud_eks_cluster_credential.cred.credential.token
}

output "public_lb_param" {
  value = data.tencentcloud_eks_cluster_credential.cred.public_lb.0.extra_param
}

output "internal_lb_subnet" {
  value = data.tencentcloud_eks_cluster_credential.cred.internal_lb.0.subnet_id
}

```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

func datasourceTencentCloudEksClusterCredential() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource was offline no longer supported.",
		Read:               datasourceTencentCloudEksClusterCredentialRead,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EKS Cluster ID.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save result.",
			},
			"addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of IP Address information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of IP, can be `advertise`, `public`, etc.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP Address.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port.",
						},
					},
				},
			},
			"credential": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Credential info. Format `{ ca_cert: String, token: String }`.",
			},
			"public_lb": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster public access LoadBalancer info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,

							Description: "Indicates weather the public access LB enabled.",
						},
						"allow_from_cidrs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of CIDRs which allowed to access.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"security_policies": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of security allow IP or CIDRs, default deny all.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"extra_param": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Extra param text json.",
						},
						"security_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group.",
						},
					},
				},
			},
			"internal_lb": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster internal access LoadBalancer info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates weather the internal access LB enabled.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of subnet which related to Internal LB.",
						},
					},
				},
			},
			"proxy_lb": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the new internal/public network function.",
			},
			"kube_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EKS cluster kubeconfig.",
			},
		},
	}
}

func datasourceTencentCloudEksClusterCredentialRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("datasource.tencentcloud_eks_cluster_credential.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := EksService{client: client}

	clusterId := d.Get("cluster_id").(string)

	request := tke.NewDescribeEKSClusterCredentialRequest()
	request.ClusterId = &clusterId

	info, err := service.DescribeEKSClusterCredential(ctx, request)

	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId("eks-cluster-credential-" + clusterId)

	_ = d.Set("proxy_lb", info.ProxyLB)

	_ = d.Set("kube_config", info.KubeConfig)

	addresses := make([]map[string]interface{}, 0)

	for i := range info.Addresses {
		item := info.Addresses[i]
		addr := map[string]interface{}{
			"type": item.Type,
			"ip":   item.Ip,
			"port": item.Port,
		}
		addresses = append(addresses, addr)
	}
	_ = d.Set("addresses", addresses)

	credential := make(map[string]interface{})
	if info.Credential != nil {
		credential = map[string]interface{}{
			"token":   info.Credential.Token,
			"ca_cert": info.Credential.CACert,
		}
		_ = d.Set("credential", credential)
	}

	internalLB := make([]map[string]interface{}, 0)
	if info.InternalLB != nil {
		lb := map[string]interface{}{
			"enabled":   info.InternalLB.Enabled,
			"subnet_id": info.InternalLB.SubnetId,
		}
		internalLB = append(internalLB, lb)
		_ = d.Set("internal_lb", internalLB)
	}

	publicLB := make([]map[string]interface{}, 0)
	if info.PublicLB != nil {
		lb := map[string]interface{}{
			"enabled":           info.PublicLB.Enabled,
			"extra_param":       info.PublicLB.ExtraParam,
			"allow_from_cidrs":  info.PublicLB.AllowFromCidrs,
			"security_group":    info.PublicLB.SecurityGroup,
			"security_policies": info.PublicLB.SecurityPolicies,
		}
		publicLB = append(publicLB, lb)
		_ = d.Set("public_lb", publicLB)
	}

	result := map[string]interface{}{
		"credential":  credential,
		"addresses":   addresses,
		"public_lb":   publicLB,
		"internal_lb": internalLB,
		"proxy_lb":    info.ProxyLB,
		"kube_config": info.KubeConfig,
	}

	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), result)
	}

	return nil
}
