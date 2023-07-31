data "tencentcloud_instance_types" "cvm4c8m" {
	exclude_sold_out=true
	cpu_core_count=4
	memory_size=8
    filter {
      name   = "instance-charge-type"
      values = ["POSTPAID_BY_HOUR"]
    }
    filter {
    name   = "zone"
    values = [var.availability_zone]
  }
}

resource "tencentcloud_vpc" "emr_vpc" {
  name       = "emr-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "emr_subnet" {
  availability_zone = var.availability_zone
  name              = "emr-subnets"
  vpc_id            = tencentcloud_vpc.emr_vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_security_group" "emr_sg" {
  name        = "emr-sg"
  description = "emr sg"
  project_id  = 0
}

resource "tencentcloud_emr_cluster" "emr_cluster" {
	product_id=4
	display_strategy="clusterList"
	vpc_settings={
	  vpc_id=tencentcloud_vpc.emr_vpc.id
      subnet_id=tencentcloud_subnet.emr_subnet.id
	}
	softwares=[
	  "zookeeper-3.6.1",
    ]
	support_ha=0
	instance_name="emr-cluster-test"
	resource_spec {
	  master_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
	  }
	  master_count=1
	  core_count=2
	}
	login_settings={
	  password="Tencent@cloud123"
	}
	time_span=3600
	time_unit="s"
	pay_mode=0
	placement={
	  zone=var.availability_zone
	  project_id=0
	}
	sg_id=tencentcloud_security_group.emr_sg.id
}