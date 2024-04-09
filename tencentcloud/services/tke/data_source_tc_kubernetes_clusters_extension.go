package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKubernetesClustersReadPreRequest0(ctx context.Context, req *tke.DescribeClustersRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	var clusterID string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterID = v.(string)
	}
	var clusterName string
	if v, ok := d.GetOk("cluster_name"); ok {
		clusterName = v.(string)
	}
	if clusterID != "" && clusterName != "" {
		return fmt.Errorf("cluster_id, cluster_name only one can be set one")
	}
	if clusterID != "" {
		req.ClusterIds = []*string{&clusterID}
	}
	if clusterName != "" {
		filter := &tke.Filter{
			Name:   helper.String("ClusterName"),
			Values: []*string{&clusterName},
		}
		req.Filters = []*tke.Filter{filter}
	}

	return nil
}

func dataSourceTencentCloudKubernetesClustersReadPreHandleResponse0(ctx context.Context, req map[string]interface{}, resp *[]*tke.Cluster) error {
	if resp == nil {
		return nil
	}
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	tags := helper.GetTags(d, "tags")
	if len(tags) == 0 {
		return nil
	}

	// 通过标签过滤
	var newResp []*tke.Cluster
clsLoop:
	for _, cls := range *resp {
		clsTags := map[string]string{}
		// 建立 cls tag 的索引
		if len(cls.TagSpecification) > 0 {
			for _, tag := range cls.TagSpecification[0].Tags {
				clsTags[*tag.Key] = *tag.Value
			}
		}
		// 比对是否包含指定的 tag ，以及值是否一致
		for k, v := range tags {
			if clsTags[k] != v {
				continue clsLoop
			}
		}

		newResp = append(newResp, cls)
	}
	*resp = newResp
	return nil
}

func dataSourceTencentCloudKubernetesClustersReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *[]*tke.Cluster) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}
	logID := tccommon.GetLogId(ctx)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	list, ok := d.GetOk("list")
	if !ok {
		return nil
	}
	listSlice, ok := list.([]interface{})
	if !ok {
		return nil
	}
	for i, v := range listSlice {
		clsData, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		clsData["cluster_os"] = tkeToShowClusterOs(clsData["cluster_os"].(string))
		clsData["cluster_deploy_type"] = strings.ToUpper(clsData["cluster_deploy_type"].(string))

		clusterID := clsData["cluster_id"].(string)

		// 获取集群节点信息
		_, nodes, err := service.DescribeClusterInstances(ctx, clusterID)
		if err != nil {
			_, nodes, err = service.DescribeClusterInstances(ctx, clusterID)
		}
		if err != nil {
			log.Printf(
				"[CRITAL]%s tencentcloud_kubernetes_clusters DescribeClusterInstances fail, reason:%s\n ",
				logID, err.Error(),
			)
			return err
		}
		workerInstancesList := make([]map[string]interface{}, 0, len(nodes))
		for _, node := range nodes {
			tempMap := make(map[string]interface{})
			tempMap["instance_id"] = node.InstanceId
			tempMap["instance_role"] = node.InstanceRole
			tempMap["instance_state"] = node.InstanceState
			tempMap["failed_reason"] = node.FailedReason
			tempMap["lan_ip"] = node.LanIp
			workerInstancesList = append(workerInstancesList, tempMap)
		}
		clsData["worker_instances_list"] = workerInstancesList

		// 获取集群安全信息
		securityRet, err := service.DescribeClusterSecurity(ctx, clusterID)
		if err != nil {
			securityRet, err = service.DescribeClusterSecurity(ctx, clusterID)
		}
		if err != nil {
			log.Printf(
				"[CRITAL]%s tencentcloud_kubernetes_clusters DescribeClusterSecurity fail, reason:%s\n ",
				logID, err.Error(),
			)
			return err
		}
		policies := make([]string, 0, len(securityRet.Response.SecurityPolicy))
		for _, v := range securityRet.Response.SecurityPolicy {
			policies = append(policies, *v)
		}
		clsData["user_name"] = helper.PString(securityRet.Response.UserName)
		clsData["password"] = helper.PString(securityRet.Response.Password)
		clsData["certification_authority"] = helper.PString(securityRet.Response.CertificationAuthority)
		clsData["cluster_external_endpoint"] = helper.PString(securityRet.Response.ClusterExternalEndpoint)
		clsData["domain"] = helper.PString(securityRet.Response.Domain)
		clsData["pgw_endpoint"] = helper.PString(securityRet.Response.PgwEndpoint)
		clsData["security_policy"] = policies

		// 获取集群连接凭证
		config, err := service.DescribeClusterConfig(ctx, clusterID, true)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				config, err = service.DescribeClusterConfig(ctx, d.Id(), true)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			log.Printf(
				"[CRITAL]%s tencentcloud_kubernetes_clusters DescribeClusterInstances fail, reason:%s\n ",
				logID, err.Error(),
			)
			return err
		}
		intranetConfig, err := service.DescribeClusterConfig(ctx, clusterID, false)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				config, err = service.DescribeClusterConfig(ctx, d.Id(), false)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			log.Printf(
				"[CRITAL]%s tencentcloud_kubernetes_clusters DescribeClusterInstances fail, reason:%s\n ",
				logID, err.Error(),
			)
			return err
		}
		clsData["kube_config"] = config
		clsData["kube_config_intranet"] = intranetConfig
		clusterInternet, err := getClusterNetworkStatus(ctx, &service, clusterID, true)
		if err != nil {
			log.Printf(
				"[CRITAL]%s tencentcloud_kubernetes_clusters get cluster internet status fail, reason:%s\n ",
				logID, err.Error(),
			)
			return err
		}
		clusterIntranet, err := getClusterNetworkStatus(ctx, &service, clusterID, false)
		if err != nil {
			log.Printf(
				"[CRITAL]%s tencentcloud_kubernetes_clusters get cluster intranet status fail, reason:%s\n ",
				logID, err.Error(),
			)
			return err
		}
		if v, ok := d.GetOk("kube_config_file_prefix"); ok {
			kubeConfigFilePrefix := v.(string)
			if clusterInternet {
				kubeConfigFile := kubeConfigFilePrefix + fmt.Sprintf("-%s-kubeconfig", clsData)
				if err = tccommon.WriteToFile(kubeConfigFile, config); err != nil {
					return err
				}
			}
			if clusterIntranet {
				kubeConfigIntranetFile := kubeConfigFilePrefix + fmt.Sprintf("-%s-kubeconfig-intranet", clusterID)
				if err = tccommon.WriteToFile(kubeConfigIntranetFile, intranetConfig); err != nil {
					return err
				}
			}
		}

		if len(*resp) <= i {
			continue
		}
		cls := (*resp)[i]

		// 获取 vpc_cni_type
		property, err := helper.JsonToMap(*cls.Property)
		if err != nil {
			return err
		}
		if property["VpcCniType"] != nil {
			clsData["vpc_cni_type"] = property["VpcCniType"].(string)
		}

		// 获取 tags
		if len(cls.TagSpecification) > 0 {
			tags := map[string]string{}
			for _, tag := range cls.TagSpecification[0].Tags {
				tags[*tag.Key] = *tag.Value
			}
			clsData["tags"] = tags
		}

		// cluster_extra_args
		clsData["cluster_extra_args"] = []map[string]interface{}{{
			"kube_apiserver":          []string(nil),
			"kube_controller_manager": []string(nil),
			"kube_scheduler":          []string(nil),
		}}
	}

	_ = d.Set("list", listSlice)

	return nil
}

func getClusterNetworkStatus(ctx context.Context, service *TkeService, clusterId string, isInternet bool) (networkStatus bool, err error) {
	var status string
	var isOpened bool
	var errRet error
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		status, _, errRet = service.DescribeClusterEndpointStatus(ctx, clusterId, isInternet)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if status == TkeInternetStatusCreating || status == TkeInternetStatusDeleting {
			return resource.RetryableError(
				fmt.Errorf("%s create cluster internet endpoint status still is %s", clusterId, status))
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	if status == TkeInternetStatusNotfound || status == TkeInternetStatusDeleted {
		isOpened = false
	}
	if status == TkeInternetStatusCreated {
		isOpened = true
	}
	networkStatus = isOpened

	return networkStatus, nil
}

func dataSourceTencentCloudKubernetesClustersReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return []interface{}(nil)
	}
	return d.Get("list")
}
