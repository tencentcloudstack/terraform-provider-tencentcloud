resource "tencentcloud_tcr_instance" "example" {
  name        = "testacctcrinstance"
  instance_type = "standard"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_token" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  description       = "test"
  enable   = false
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  name        = "test"
  is_public   = false
}

resource "tencentcloud_tcr_repository" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  namespace_name        = tencentcloud_tcr_namespace.example.name
  name = "test"
  brief_desc = "example"
  description = "long example"
}

data "tencentcloud_tcr_instances" "example" {
  name = tencentcloud_tcr_instance.example.name
}

data "tencentcloud_tcr_tokens" "example" {
  instance_id = tencentcloud_tcr_token.example.instance_id
}

data "tencentcloud_tcr_namespaces" "example" {
  instance_id = tencentcloud_tcr_namespace.example.instance_id
}

data "tencentcloud_tcr_repositories" "example" {
  instance_id = tencentcloud_tcr_repository.example.instance_id
  namespace_name = tencentcloud_tcr_namespace.example.name
}

# get internal vpc access
resource "tencentcloud_vpc" "example" {
  name       = "example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "example" {
  availability_zone = var.availability_zone
  name              = "example"
  vpc_id            = tencentcloud_vpc.example.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_tcr_vpc_attachment" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  vpc_id = tencentcloud_vpc.example.id
  subnet_id = tencentcloud_subnet.example.id
}

data "tencentcloud_tcr_vpc_attachments" "example" {
  instance_id = tencentcloud_tcr_vpc_attachment.example.instance_id
}