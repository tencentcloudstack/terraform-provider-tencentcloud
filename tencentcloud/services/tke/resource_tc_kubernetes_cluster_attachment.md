Provide a resource to attach an existing  cvm to kubernetes cluster.

Example Usage

```hcl

	variable "availability_zone" {
	  default = "ap-guangzhou-3"
	}

	variable "cluster_cidr" {
	  default = "172.16.0.0/16"
	}

	variable "default_instance_type" {
	  default = "S1.SMALL1"
	}

	data "tencentcloud_images" "default" {
	  image_type = ["PUBLIC_IMAGE"]
	  os_name    = "centos"
	}

	data "tencentcloud_vpc_subnets" "vpc" {
	  is_default        = true
	  availability_zone = var.availability_zone
	}

	data "tencentcloud_instance_types" "default" {
	  filter {
	    name   = "instance-family"
	    values = ["SA2"]
	  }

	  cpu_core_count = 8
	  memory_size    = 16
	}

	resource "tencentcloud_instance" "foo" {
	  instance_name     = "tf-auto-test-1-1"
	  availability_zone = var.availability_zone
	  image_id          = data.tencentcloud_images.default.images.0.image_id
	  instance_type     = var.default_instance_type
	  system_disk_type  = "CLOUD_PREMIUM"
	  system_disk_size  = 50
	}

	resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
	  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
	  cluster_cidr            = "10.1.0.0/16"
	  cluster_max_pod_num     = 32
	  cluster_name            = "keep"
	  cluster_desc            = "test cluster desc"
	  cluster_max_service_num = 32

	  worker_config {
	    count                      = 1
	    availability_zone          = var.availability_zone
	    instance_type              = var.default_instance_type
	    system_disk_type           = "CLOUD_SSD"
	    system_disk_size           = 60
	    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
	    internet_max_bandwidth_out = 100
	    public_ip_assigned         = true
	    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

	    data_disk {
	      disk_type = "CLOUD_PREMIUM"
	      disk_size = 50
	    }

	    enhanced_security_service = false
	    enhanced_monitor_service  = false
	    user_data                 = "dGVzdA=="
	    password                  = "ZZXXccvv1212"
	  }

	  cluster_deploy_type = "MANAGED_CLUSTER"
	}

	resource "tencentcloud_kubernetes_cluster_attachment" "test_attach" {
	  cluster_id  = tencentcloud_kubernetes_cluster.managed_cluster.id
	  instance_id = tencentcloud_instance.foo.id
	  password    = "Lo4wbdit"

	  labels = {
	    "test1" = "test1",
	    "test2" = "test2",
	  }

	  worker_config_overrides {
	    desired_pod_num = 8
	  }
	}

```