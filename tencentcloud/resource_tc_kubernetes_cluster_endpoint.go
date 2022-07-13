/*
Provide a resource to create a KubernetesClusterEndpoint.
This resource allows you to create an empty cluster first without any workers. Only all attached node depends create complete, cluster endpoint will finally be enabled.

~> **NOTE:** Recommend using `depends_on` to make sure endpoint create after node pools or workers does.

Example Usage

```hcl
resource "tencentcloud_kubernetes_node_pool" "pool1" {}

resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id = "cls-xxxxxxxx"
  cluster_internet = true
  cluster_intranet = true
  managed_cluster_internet_security_policies = [
    "192.168.0.0/24"
  ]
  cluster_intranet_subnet_id = "subnet-xxxxxxxx"
  depends_on = [
	tencentcloud_kubernetes_node_pool.pool1
  ]
}
```


Import

KubernetesClusterEndpoint instance can be imported by passing cluster id, e.g.
```
$ terraform import tencentcloud_kubernetes_cluster_endpoint.test cluster-id
```

*/
package tencentcloud

import (
	"context"
	"fmt"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTkeClusterEndpoint() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTkeClusterEndpointRead,
		Create: resourceTencentCloudTkeClusterEndpointCreate,
		Update: resourceTencentCloudTkeClusterEndpointUpdate,
		Delete: resourceTencentCloudTkeClusterEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specify cluster ID.",
			},
			"cluster_internet": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Open internet access or not.",
			},
			"cluster_intranet": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Open intranet access or not.",
			},
			"managed_cluster_internet_security_policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Security policies for managed cluster internet, like:'192.168.1.0/24' or '113.116.51.27', '0.0.0.0/0' means all." +
					" This field can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true." +
					" `managed_cluster_internet_security_policies` can not delete or empty once be set.",
			},
			"cluster_intranet_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Subnet id who can access this independent cluster, this field must and can only set  when `cluster_intranet` is true." +
					" `cluster_intranet_subnet_id` can not modify once be set.",
			},
			// Computed
			"cluster_deploy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster deploy type of `MANAGED_CLUSTER` or `INDEPENDENT_CLUSTER`.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User name of account.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Password of account.",
			},
			"certification_authority": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate used for access.",
			},
			"cluster_external_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External network address to access.",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain name for access.",
			},
			"pgw_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Intranet address used for access.",
			},
		},
	}
}

func resourceTencentCloudTkeClusterEndpointRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_endpoint.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client}

	id := d.Id()
	info, has, err := service.DescribeCluster(ctx, id)
	if err != nil {
		d.SetId("")
		return err
	}
	if !has {
		d.SetId("")
		return fmt.Errorf("cluster %s not found", id)
	}

	d.SetId(id)
	_ = d.Set("cluster_id", id)

	_ = d.Set("cluster_deploy_type", info.DeployType)

	response, err := service.DescribeClusterSecurity(ctx, id)

	if err != nil {
		return err
	}

	var (
		security        = response.Response
		clusterInternet = security.ClusterExternalEndpoint != nil && *security.ClusterExternalEndpoint != ""
		clusterIntranet = security.PgwEndpoint != nil && *security.PgwEndpoint != ""
	)

	_ = d.Set("cluster_internet", clusterInternet)
	_ = d.Set("cluster_intranet", clusterIntranet)
	_ = d.Set("user_name", security.UserName)
	_ = d.Set("password", security.Password)
	_ = d.Set("certification_authority", security.CertificationAuthority)
	_ = d.Set("cluster_external_endpoint", security.ClusterExternalEndpoint)
	_ = d.Set("domain", security.Domain)
	_ = d.Set("pgw_endpoint", security.PgwEndpoint)

	if len(security.SecurityPolicy) > 0 {
		_ = d.Set("managed_cluster_internet_security_policies", security.SecurityPolicy)
	}

	return nil
}

func resourceTencentCloudTkeClusterEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_endpoint.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client}

	id := d.Get("cluster_id").(string)
	var (
		err              error
		isManagedCluster bool
		securityPolicies []string
		clusterInternet  = d.Get("cluster_internet").(bool)
		clusterIntranet  = d.Get("cluster_intranet").(bool)
		intranetSubnetId = d.Get("cluster_intranet_subnet_id").(string)
	)

	clusterInfo, has, err := service.DescribeCluster(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("cluster %s not found", id)
	}
	isManagedCluster = clusterInfo.DeployType == TKE_DEPLOY_TYPE_MANAGED

	if v, ok := d.Get("managed_cluster_internet_security_policies").([]interface{}); ok && len(v) > 0 {
		securityPolicies = helper.InterfacesStrings(v)
	}

	if clusterIntranet && intranetSubnetId == "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` must set when `cluster_intranet` is true")
	}
	if !clusterIntranet && intranetSubnetId != "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` can only set when `cluster_intranet` is true")
	}

	if !(clusterInternet && isManagedCluster) && len(securityPolicies) > 0 {
		return fmt.Errorf("`managed_cluster_internet_security_policies` can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true.")
	}

	// Create Intranet(Private) Network
	if clusterIntranet {
		err := tencentCloudClusterIntranetSwitch(ctx, &service, id, intranetSubnetId, true)
		if err != nil {
			return err
		}
	}

	//TKE_DEPLOY_TYPE_INDEPENDENT Open the internet
	if clusterInternet {
		err := tencentCloudClusterInternetSwitch(ctx, &service, id, true, isManagedCluster, securityPolicies)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, true, isManagedCluster)
		if err != nil {
			return err
		}
	}

	d.SetId(id)

	return resourceTencentCloudTkeClusterEndpointRead(d, meta)
}

func resourceTencentCloudTkeClusterEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_endpoint.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client}
	id := d.Id()

	var (
		isManagedCluster bool
		err              error
	)
	if v, ok := d.GetOk("cluster_deploy_type"); ok {
		isManagedCluster = v.(string) == TKE_DEPLOY_TYPE_MANAGED
	}

	if d.HasChange("cluster_internet") {
		clusterInternet := d.Get("cluster_internet").(bool)
		policies := helper.InterfacesStrings(d.Get("managed_cluster_internet_security_policies").([]interface{}))
		err = tencentCloudClusterInternetSwitch(ctx, &service, id, clusterInternet, isManagedCluster, policies)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, clusterInternet, isManagedCluster)
		if err != nil {
			return err
		}
	}

	if d.HasChange("cluster_intranet") {
		clusterIntranet := d.Get("cluster_intranet").(bool)
		subnetId := d.Get("cluster_intranet_subnet_id").(string)
		err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, clusterIntranet)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudTkeClusterEndpointRead(d, meta)
}

func resourceTencentCloudTkeClusterEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_endpoint.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client}
	var (
		id               = d.Id()
		err              error
		isManagedCluster bool
	)
	if v, ok := d.GetOk("cluster_deploy_type"); ok {
		isManagedCluster = v.(string) == TKE_DEPLOY_TYPE_MANAGED
	}

	response, err := service.DescribeClusterSecurity(ctx, id)

	if err != nil {
		return err
	}

	var (
		security        = response.Response
		clusterInternet = security.ClusterExternalEndpoint != nil && *security.ClusterExternalEndpoint != ""
		clusterIntranet = security.PgwEndpoint != nil && *security.PgwEndpoint != ""
		errs            multierror.Error
	)

	if clusterInternet {
		err = tencentCloudClusterInternetSwitch(ctx, &service, id, false, isManagedCluster, nil)
		if err != nil {
			errs = *multierror.Append(err)
		} else {
			taskErr := waitForClusterEndpointFinish(ctx, &service, id, false, isManagedCluster)
			if taskErr != nil {
				errs = *multierror.Append(taskErr)
			}
		}
	}

	if clusterIntranet {
		err = tencentCloudClusterIntranetSwitch(ctx, &service, id, "", false)
		if err != nil {
			errs = *multierror.Append(err)
		}
	}

	return errs.ErrorOrNil()
}

func waitForClusterEndpointFinish(ctx context.Context, service *TkeService, id string, enabled bool, isManagedCluster bool) (err error) {
	return resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		var (
			status         string
			message        string
			inErr          error
			retryableState = TkeInternetStatusCreating
			finishStates   = []string{TkeInternetStatusNotfound, TkeInternetStatusCreated}
		)

		if !enabled {
			retryableState = TkeInternetStatusDeleting
			finishStates = []string{TkeInternetStatusNotfound, TkeInternetStatusDeleted}
		}

		if isManagedCluster {
			status, message, inErr = service.DescribeClusterEndpointVipStatus(ctx, id)
		} else {
			status, message, inErr = service.DescribeClusterEndpointStatus(ctx, id)
		}
		if inErr != nil {
			return retryError(inErr)
		}
		if status == retryableState {
			return resource.RetryableError(
				fmt.Errorf("%s create cluster internet endpoint status still is %s", id, status))
		}
		if IsContains(finishStates, status) {
			return nil
		}
		return resource.NonRetryableError(
			fmt.Errorf("%s create cluster internet endpoint error ,status is %s,message is %s", id, status, message))
	})
}

func tencentCloudClusterInternetSwitch(ctx context.Context, service *TkeService, id string, enable, isManagedCluster bool, policies []string) (err error) {
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if enable {
			if isManagedCluster {
				err = service.CreateClusterEndpointVip(ctx, id, policies)
			} else {
				err = service.CreateClusterEndpoint(ctx, id, "", true)
			}
			if err != nil {
				return retryError(err, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
			}
		} else {
			if isManagedCluster {
				err = service.DeleteClusterEndpointVip(ctx, id)
			} else {
				err = service.DeleteClusterEndpoint(ctx, id, true)
			}
			if err != nil {
				return retryError(err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func tencentCloudClusterIntranetSwitch(ctx context.Context, service *TkeService, id, subnetId string, enable bool) (err error) {
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if enable {
			err = service.CreateClusterEndpoint(ctx, id, subnetId, false)
			if err != nil {
				return retryError(err, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
			}
		} else {
			err = service.DeleteClusterEndpoint(ctx, id, false)
			if err != nil {
				return retryError(err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
