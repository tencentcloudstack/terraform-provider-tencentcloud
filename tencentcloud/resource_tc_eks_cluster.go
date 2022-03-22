/*
Provides an elastic kubernetes cluster resource.

Example Usage

```
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-eks-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "sub" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}
resource "tencentcloud_subnet" "sub2" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.10.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_eks_cluster" "foo" {
  cluster_name = "tf-test-eks"
  k8s_version = "1.18.4"
  vpc_id = tencentcloud_vpc.vpc.id
  subnet_ids = [
    tencentcloud_subnet.sub.id,
    tencentcloud_subnet.sub2.id,
  ]
  cluster_desc = "test eks cluster created by terraform"
  service_subnet_id =     tencentcloud_subnet.sub.id
  dns_servers {
    domain = "www.example1.com"
    servers = ["1.1.1.1:8080", "1.1.1.1:8081", "1.1.1.1:8082"]
  }
  enable_vpc_core_dns = true
  need_delete_cbs = true
  tags = {
    hello = "world"
  }
}
```

Import

```
terraform import tencentcloud_eks_cluster.foo cluster-id
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudEksCluster() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentcloudEKSClusterRead,
		Create: resourceTencentcloudEKSClusterCreate,
		Update: resourceTencentcloudEKSClusterUpdate,
		Delete: resourceTencentcloudEKSClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of EKS cluster.",
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Kubernetes version of EKS cluster.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Vpc Id of EKS cluster.",
			},
			"subnet_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Subnet Ids for EKS cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cluster_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of EKS cluster.",
			},
			"service_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet id of service.",
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of cluster custom DNS Server info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Optional:    true,
							Type:        schema.TypeString,
							Description: "DNS Server domain. Empty indicates all domain.",
						},
						"servers": {
							Optional:    true,
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of DNS Server IP address, pattern: \"ip[:port]\".",
						},
					},
				},
			},
			"extra_param": {
				Type:        schema.TypeMap,
				ForceNew:    true,
				Optional:    true,
				Description: "Extend parameters.",
			},
			"enable_vpc_core_dns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Indicates whether to enable dns in user cluster, default value is `true`.",
			},
			// update after create
			"need_delete_cbs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Delete CBS after EKS cluster remove.",
			},
			"public_lb": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Cluster public access LoadBalancer info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates weather the public access LB enabled.",
						},
						"allow_from_cidrs": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of CIDRs which allowed to access.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"security_policies": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of security allow IP or CIDRs, default deny all.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"extra_param": {
							Type:        schema.TypeString,
							Optional:    true,
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
				MaxItems:    1,
				Optional:    true,
				Description: "Cluster internal access LoadBalancer info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates weather the internal access LB enabled.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of subnet which related to Internal LB.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of EKS cluster.",
			},
		},
	}
}

func resourceTencentcloudEKSClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eks_cluster.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EksService{client: meta.(*TencentCloudClient).apiV3Conn}

	cluster, has, err := service.DescribeEksCluster(ctx, d.Id())
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("cluster %s not found", d.Id())
	}

	_ = d.Set("cluster_name", cluster.ClusterName)
	_ = d.Set("cluster_desc", cluster.ClusterDesc)
	_ = d.Set("k8s_version", cluster.K8SVersion)
	_ = d.Set("vpc_id", cluster.VpcId)
	_ = d.Set("service_subnet_id", cluster.ServiceSubnetId)
	_ = d.Set("subnet_ids", cluster.SubnetIds)
	_ = d.Set("dns_servers", cluster.DnsServers)
	_ = d.Set("tags", cluster.Tags)
	_ = d.Set("subnet_ids", cluster.SubnetIds)
	_ = d.Set("need_delete_cbs", cluster.NeedDeleteCbs)
	_ = d.Set("enable_vpc_core_dns", cluster.EnableVpcCoreDNS)

	info, err := service.DescribeEKSClusterCredentialById(ctx, d.Id())

	if err != nil {
		return err
	}

	if _, ok := d.GetOk("internal_lb"); ok && info.InternalLB != nil {
		internalLB := make([]map[string]interface{}, 0)
		lb := map[string]interface{}{
			"enabled":   info.InternalLB.Enabled,
			"subnet_id": info.InternalLB.SubnetId,
		}
		internalLB = append(internalLB, lb)
		err = d.Set("internal_lb", internalLB)
		if err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("public_lb"); ok && info.PublicLB != nil {
		publicLB := make([]map[string]interface{}, 0)
		var (
			cidrs    []*string
			policies []*string
		)
		// Avoid empty string "" write to data
		for _, v := range info.PublicLB.AllowFromCidrs {
			if *v != "" {
				cidrs = append(cidrs, v)
			}
		}
		for _, v := range info.PublicLB.SecurityPolicies {
			if *v != "" {
				policies = append(cidrs, v)
			}
		}
		lb := map[string]interface{}{
			"enabled":           info.PublicLB.Enabled,
			"extra_param":       info.PublicLB.ExtraParam,
			"allow_from_cidrs":  cidrs,
			"security_group":    info.PublicLB.SecurityGroup,
			"security_policies": policies,
		}
		publicLB = append(publicLB, lb)
		_ = d.Set("public_lb", publicLB)
	}

	return nil
}

func resourceTencentcloudEKSClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eks_cluster.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		client           = meta.(*TencentCloudClient).apiV3Conn
		service          = EksService{client: client}
		tagService       = TagService{client: client}
		k8sVersion       = d.Get("k8s_version").(string)
		clusterName      = d.Get("cluster_name").(string)
		vpcId            = d.Get("vpc_id").(string)
		clusterDesc      = d.Get("cluster_desc").(string)
		enableVpcCoreDns = d.Get("enable_vpc_core_dns").(bool)
		tags             = helper.GetTags(d, "tags")
		subnetIds        []*string
		dnsServers       []*tke.DnsServerConf
	)

	request := tke.NewCreateEKSClusterRequest()
	request.K8SVersion = helper.String(k8sVersion)
	request.ClusterName = helper.String(clusterName)
	request.VpcId = helper.String(vpcId)
	request.EnableVpcCoreDNS = helper.Bool(enableVpcCoreDns)

	for k, v := range tags {
		if len(request.TagSpecification) == 0 {
			request.TagSpecification = []*tke.TagSpecification{{
				ResourceType: helper.String("cluster"),
			}}
		}

		request.TagSpecification[0].Tags = append(request.TagSpecification[0].Tags, &tke.Tag{
			Key:   helper.String(k),
			Value: helper.String(v),
		})
	}
	if clusterDesc != "" {
		request.ClusterDesc = helper.String(clusterDesc)
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		for _, id := range v.([]interface{}) {
			subnetIds = append(subnetIds, helper.String(id.(string)))
		}
		request.SubnetIds = subnetIds
	}

	if v, ok := d.GetOk("dns_servers"); ok {
		for _, i := range v.([]interface{}) {
			conf := i.(map[string]interface{})
			domain := conf["domain"].(string)
			dnsConf := &tke.DnsServerConf{
				Domain: helper.String(domain),
			}
			for _, server := range conf["servers"].([]interface{}) {
				dnsConf.DnsServers = append(dnsConf.DnsServers, helper.String(server.(string)))
			}
			dnsServers = append(dnsServers, dnsConf)
		}

		request.DnsServers = dnsServers
	}

	if extraParamRaw, ok := d.GetOk("extra_param"); ok {
		param, err := json.Marshal(extraParamRaw.(map[string]string))
		if err != nil {
			return err
		}
		request.ExtraParam = helper.String(string(param))
	}

	id, err := service.CreateEksCluster(ctx, request)

	if err != nil {
		return err
	}

	d.SetId(id)

	err = resource.Retry(readRetryTimeout*3, func() *resource.RetryError {
		cluster, _, err := service.DescribeEksCluster(ctx, id)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if cluster.Status == "Initializing" {
			return resource.RetryableError(fmt.Errorf("eks cluster is %s, retry", cluster.Status))
		}

		return nil
	})

	if err != nil {
		return err
	}

	upgradeRequest := tke.NewUpdateEKSClusterRequest()

	if v, ok := d.GetOk("need_delete_cbs"); ok {
		upgradeRequest.ClusterId = helper.String(id)
		upgradeRequest.NeedDeleteCbs = helper.Bool(v.(bool))
	}

	enablePublic := false
	enableInternal := false

	if lb, ok := helper.InterfacesHeadMap(d, "internal_lb"); ok {
		upgradeRequest.ClusterId = helper.String(id)
		enabled := lb["enabled"].(bool)
		upgradeRequest.InternalLB = &tke.ClusterInternalLB{
			Enabled: &enabled,
		}
		if v, ok := lb["subnet_id"].(string); ok && v != "" {
			upgradeRequest.InternalLB.SubnetId = &v
		}
		enableInternal = enabled
	}

	if lb, ok := helper.InterfacesHeadMap(d, "public_lb"); ok {
		upgradeRequest.ClusterId = helper.String(id)
		enabled := lb["enabled"].(bool)
		upgradeRequest.PublicLB = &tke.ClusterPublicLB{
			Enabled: &enabled,
		}
		if v, ok := lb["security_group"].(string); ok && v != "" {
			upgradeRequest.PublicLB.SecurityGroup = &v
		}
		if v, ok := lb["allow_from_cidrs"].([]interface{}); ok && len(v) > 0 {
			upgradeRequest.PublicLB.AllowFromCidrs = helper.InterfacesStringsPoint(v)
		}
		if v, ok := lb["security_policies"].([]interface{}); ok && len(v) > 0 {
			upgradeRequest.PublicLB.SecurityPolicies = helper.InterfacesStringsPoint(v)
		}
		enablePublic = enabled
	}

	if upgradeRequest.ClusterId != nil {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := service.UpdateEksCluster(ctx, upgradeRequest)
			if inErr != nil {
				return retryError(err, InternalError)
			}
			return nil
		})

		if err != nil {
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, inErr := service.DescribeEKSClusterCredentialById(ctx, id)
			if inErr != nil {
				return retryError(inErr)
			}
			if info.InternalLB == nil || *info.InternalLB.Enabled != enableInternal {
				return resource.RetryableError(fmt.Errorf("waiting for internal lb upgrade"))
			}

			if info.PublicLB == nil || *info.PublicLB.Enabled != enablePublic {
				return resource.RetryableError(fmt.Errorf("waiting for public lb upgrade"))
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		region := client.Region
		resourceName := BuildTagResourceName("ccs", "cluster", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, tags, []string{}); err != nil {
			fmt.Printf("[WARN]: update tags failed: %s", err.Error())
		}
	}

	return resourceTencentcloudEKSClusterRead(d, meta)
}

func resourceTencentcloudEKSClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eks_cluster.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	id := d.Id()
	request := tke.NewUpdateEKSClusterRequest()
	request.ClusterId = helper.String(id)
	client := meta.(*TencentCloudClient).apiV3Conn

	var (
		service     = EksService{client: client}
		tagService  = TagService{client: client}
		clusterName = d.Get("cluster_name").(string)
		clusterDesc = d.Get("cluster_desc").(string)
		updateAttrs []string
	)

	if d.HasChange("cluster_name") {
		updateAttrs = append(updateAttrs, "cluster_name")
		request.ClusterName = helper.String(clusterName)
	}

	if d.HasChange("cluster_desc") {
		updateAttrs = append(updateAttrs, "cluster_desc")
		request.ClusterDesc = helper.String(clusterDesc)
	}

	if d.HasChange("subnet_ids") {
		updateAttrs = append(updateAttrs, "subnet_ids")
		ids := d.Get("subnet_ids").([]interface{})
		for _, id := range ids {
			request.SubnetIds = append(request.SubnetIds, helper.String(id.(string)))
		}
	}

	if d.HasChange("dns_servers") {
		updateAttrs = append(updateAttrs, "dns_servers")
		var dnsServers []*tke.DnsServerConf
		v := d.Get("dns_servers")
		servers := v.([]interface{})
		if len(servers) > 0 {
			for _, i := range v.([]interface{}) {
				conf := i.(map[string]interface{})
				domain := conf["domain"].(string)
				dnsConf := &tke.DnsServerConf{
					Domain: helper.String(domain),
				}
				for _, server := range conf["servers"].([]interface{}) {
					dnsConf.DnsServers = append(dnsConf.DnsServers, helper.String(server.(string)))
				}
				dnsServers = append(dnsServers, dnsConf)
			}

			request.DnsServers = dnsServers
		} else {
			request.ClearDnsServer = helper.String("1")
		}
	}

	if d.HasChange("need_delete_cbs") {
		updateAttrs = append(updateAttrs, "need_delete_cbs")
		needDelete := d.Get("need_delete_cbs").(bool)
		request.NeedDeleteCbs = helper.Bool(needDelete)
	}

	enableInternal := false
	enablePublic := false
	if d.HasChange("internal_lb") {
		updateAttrs = append(updateAttrs, "internal_lb")
		if v, ok := d.GetOk("internal_lb"); ok {
			lb := v.([]map[string]interface{})[0]
			enabled := lb["enabled"].(bool)
			request.InternalLB = &tke.ClusterInternalLB{
				Enabled: &enabled,
			}
			if v, ok := lb["subnet_id"].(string); ok && v != "" {
				request.InternalLB.SubnetId = &v
			}
			enableInternal = enabled
		} else {
			request.InternalLB = &tke.ClusterInternalLB{
				Enabled: helper.Bool(false),
			}
		}
	}

	if d.HasChange("public_lb") {
		updateAttrs = append(updateAttrs, "public_lb")
		if v, ok := d.GetOk("public_lb"); ok {
			lb := v.([]map[string]interface{})[0]
			enabled := lb["enabled"].(bool)
			request.PublicLB = &tke.ClusterPublicLB{
				Enabled: &enabled,
			}
			if v, ok := lb["security_group"].(string); ok && v != "" {
				request.PublicLB.SecurityGroup = &v
			}
			if v, ok := lb["allow_from_cidrs"].([]interface{}); ok && len(v) > 0 {
				request.PublicLB.AllowFromCidrs = helper.InterfacesStringsPoint(v)
			}
			if v, ok := lb["security_policies"].([]interface{}); ok && len(v) > 0 {
				request.PublicLB.SecurityPolicies = helper.InterfacesStringsPoint(v)
			}
			enablePublic = enabled
		} else {
			request.PublicLB = &tke.ClusterPublicLB{
				Enabled: helper.Bool(false),
			}
		}

		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, inErr := service.DescribeEKSClusterCredentialById(ctx, id)
			if inErr != nil {
				return retryError(inErr)
			}
			if info.InternalLB == nil || *info.InternalLB.Enabled != enableInternal {
				return resource.RetryableError(fmt.Errorf("waiting for internal lb upgrade"))
			}

			if info.PublicLB == nil || *info.PublicLB.Enabled != enablePublic {
				return resource.RetryableError(fmt.Errorf("waiting for public lb upgrade"))
			}
			return nil
		})

		if err != nil {
			return err
		}

	}

	if len(updateAttrs) > 0 {
		err := service.UpdateEksCluster(ctx, request)
		if err != nil {
			return err
		}
		for _, attr := range updateAttrs {
			d.SetPartial(attr)
		}
	}

	if d.HasChange("tags") {
		region := client.Region
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("ccs", "cluster", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	return resourceTencentcloudEKSClusterRead(d, meta)
}

func resourceTencentcloudEKSClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eks_cluster.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := EksService{client: meta.(*TencentCloudClient).apiV3Conn}

	request := tke.NewDeleteEKSClusterRequest()
	request.ClusterId = helper.String(id)

	if err := service.DeleteEksCluster(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		info, has, err := service.DescribeEksCluster(ctx, d.Id())
		if has && info.Status == "Terminating" {
			return resource.RetryableError(fmt.Errorf("cluster %s terminating, retrying", d.Id()))
		}
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			log.Printf("Describe Error Code: %s", e.GetCode())
			if e.GetCode() == "ResourceNotFound" {
				return nil
			}
		}
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
