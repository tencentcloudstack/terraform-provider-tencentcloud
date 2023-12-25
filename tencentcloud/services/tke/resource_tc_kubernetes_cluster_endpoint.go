package tke

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTkeClusterEndpoint() *schema.Resource {
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
			"cluster_internet_security_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify security group, NOTE: This argument must not be empty if cluster internet enabled.",
			},
			"managed_cluster_internet_security_policies": {
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "this argument was deprecated, use `cluster_internet_security_group` instead.",
				Description: "Security policies for managed cluster internet, like:'192.168.1.0/24' or '113.116.51.27', '0.0.0.0/0' means all." +
					" This field can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true." +
					" `managed_cluster_internet_security_policies` can not delete or empty once be set.",
			},
			"cluster_internet_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Domain name for cluster Kube-apiserver internet access. " +
					" Be careful if you modify value of this parameter, the cluster_external_endpoint value may be changed automatically too.",
			},
			"extensive_parameters": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The LB parameter. Only used for public network access.",
			},
			"cluster_intranet_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Domain name for cluster Kube-apiserver intranet access." +
					" Be careful if you modify value of this parameter, the pgw_endpoint value may be changed automatically too.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_endpoint.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
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
		security = response.Response
		//clusterInternet = security.ClusterExternalEndpoint != nil && *security.ClusterExternalEndpoint != ""
		//clusterIntranet = security.PgwEndpoint != nil && *security.PgwEndpoint != ""
	)

	//_ = d.Set("cluster_internet", clusterInternet)
	//_ = d.Set("cluster_intranet", clusterIntranet)
	_ = d.Set("user_name", security.UserName)
	_ = d.Set("password", security.Password)
	_ = d.Set("certification_authority", security.CertificationAuthority)
	_ = d.Set("cluster_external_endpoint", security.ClusterExternalEndpoint)
	_ = d.Set("domain", security.Domain)
	_ = d.Set("pgw_endpoint", security.PgwEndpoint)

	//if len(security.SecurityPolicy) > 0 {
	//	_ = d.Set("managed_cluster_internet_security_policies", security.SecurityPolicy)
	//}

	return nil
}

func resourceTencentCloudTkeClusterEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_endpoint.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TkeService{client}

	id := d.Get("cluster_id").(string)
	var (
		err                          error
		clusterInternet              = d.Get("cluster_internet").(bool)
		clusterIntranet              = d.Get("cluster_intranet").(bool)
		intranetSubnetId             = d.Get("cluster_intranet_subnet_id").(string)
		clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
		clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
		clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
		extensiveParameters          = d.Get("extensive_parameters").(string)
	)

	if err != nil {
		return err
	}

	if clusterIntranet && intranetSubnetId == "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` must set when `cluster_intranet` is true")
	}
	if !clusterIntranet && intranetSubnetId != "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` can only set when `cluster_intranet` is true")
	}

	if err := service.CheckOneOfClusterNodeReady(ctx, id, true); err != nil {
		return err
	}

	// Create Intranet(Private) Network
	if clusterIntranet {
		err := tencentCloudClusterIntranetSwitch(ctx, &service, id, intranetSubnetId, true, clusterIntranetDomain)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, true, false)
		if err != nil {
			return err
		}
	}

	//TKE_DEPLOY_TYPE_INDEPENDENT Open the internet
	if clusterInternet {
		err := tencentCloudClusterInternetSwitch(ctx, &service, id, true, clusterInternetSecurityGroup, clusterInternetDomain, extensiveParameters)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, true, true)
		if err != nil {
			return err
		}
	}

	d.SetId(id)

	return resourceTencentCloudTkeClusterEndpointRead(d, meta)
}

func resourceTencentCloudTkeClusterEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_endpoint.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TkeService{client}
	id := d.Id()

	var (
		clusterInternet              = d.Get("cluster_internet").(bool)
		clusterIntranet              = d.Get("cluster_intranet").(bool)
		clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
		clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
		clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
		subnetId                     = d.Get("cluster_intranet_subnet_id").(string)
		extensiveParameters          = d.Get("extensive_parameters").(string)
	)

	var (
		err error
	)

	if d.HasChange("cluster_internet_security_group") && !d.HasChange("cluster_internet") {
		if clusterInternet {
			err := service.ModifyClusterEndpointSG(ctx, id, clusterInternetSecurityGroup)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("cluster_internet") {
		err = tencentCloudClusterInternetSwitch(ctx, &service, id, clusterInternet, clusterInternetSecurityGroup, clusterInternetDomain, extensiveParameters)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, clusterInternet, true)
		if err != nil {
			return err
		}
	} else if clusterInternet && d.HasChange("cluster_internet_domain") {
		// only domain changed, need to close and reopen
		// close
		err = tencentCloudClusterInternetSwitch(ctx, &service, id, false, clusterInternetSecurityGroup, clusterInternetDomain, "")
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, false, true)
		if err != nil {
			return err
		}
		// reopen
		err = tencentCloudClusterInternetSwitch(ctx, &service, id, true, clusterInternetSecurityGroup, clusterInternetDomain, extensiveParameters)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, true, true)
		if err != nil {
			return err
		}
	}

	if d.HasChange("cluster_intranet") {
		err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, clusterIntranet, clusterIntranetDomain)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, clusterIntranet, false)
		if err != nil {
			return err
		}
	} else if clusterIntranet && d.HasChange("cluster_intranet_domain") {
		// only domain changed, need to close and reopen
		// close
		err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, false, clusterIntranetDomain)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, false, false)
		if err != nil {
			return err
		}
		// reopen
		err = tencentCloudClusterIntranetSwitch(ctx, &service, id, subnetId, true, clusterIntranetDomain)
		if err != nil {
			return err
		}
		err = waitForClusterEndpointFinish(ctx, &service, id, true, false)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudTkeClusterEndpointRead(d, meta)
}

func resourceTencentCloudTkeClusterEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_endpoint.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TkeService{client}
	var (
		id  = d.Id()
		err error
	)

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
		err = tencentCloudClusterInternetSwitch(ctx, &service, id, false, "", "", "")
		if err != nil {
			errs = *multierror.Append(err)
		} else {
			taskErr := waitForClusterEndpointFinish(ctx, &service, id, false, true)
			if taskErr != nil {
				errs = *multierror.Append(taskErr)
			}
		}
	}

	if clusterIntranet {
		err = tencentCloudClusterIntranetSwitch(ctx, &service, id, "", false, "")
		if err != nil {
			errs = *multierror.Append(err)
		}
	}

	return errs.ErrorOrNil()
}

func waitForClusterEndpointFinish(ctx context.Context, service *TkeService, id string, enabled bool, isInternet bool) (err error) {
	return resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

		status, message, inErr = service.DescribeClusterEndpointStatus(ctx, id, isInternet)

		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if status == retryableState {
			return resource.RetryableError(
				fmt.Errorf("%s create cluster internet endpoint status still is %s", id, status))
		}
		if tccommon.IsContains(finishStates, status) {
			return nil
		}
		return resource.NonRetryableError(
			fmt.Errorf("%s create cluster internet endpoint error ,status is %s,message is %s", id, status, message))
	})
}

func tencentCloudClusterInternetSwitch(ctx context.Context, service *TkeService, id string, enable bool, sg string, domain string, extensiveParameters string) (err error) {
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if enable {
			err = service.CreateClusterEndpoint(ctx, id, "", sg, true, domain, extensiveParameters)
			if err != nil {
				return tccommon.RetryError(err, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
			}
		} else {
			err = service.DeleteClusterEndpoint(ctx, id, true)
			if err != nil {
				return tccommon.RetryError(err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func tencentCloudClusterIntranetSwitch(ctx context.Context, service *TkeService, id, subnetId string, enable bool, domain string) (err error) {
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if enable {
			err = service.CreateClusterEndpoint(ctx, id, subnetId, "", false, domain, "")
			if err != nil {
				return tccommon.RetryError(err, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
			}
		} else {
			err = service.DeleteClusterEndpoint(ctx, id, false)
			if err != nil {
				return tccommon.RetryError(err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
