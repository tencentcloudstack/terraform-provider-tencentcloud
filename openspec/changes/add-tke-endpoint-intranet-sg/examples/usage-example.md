# Usage Examples: TKE Cluster Endpoint with Intranet Security Group

## Example 1: Basic Intranet Endpoint with Security Group

```hcl
resource "tencentcloud_kubernetes_cluster_endpoint" "intranet_example" {
  cluster_id                       = "cls-xxxxxxxx"
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = "subnet-xxxxxxxx"
  cluster_intranet_security_group  = "sg-xxxxxxxx"  # ✅ New optional field
}
```

**Key Points**:
- `cluster_intranet_security_group` is an **optional** field for configuring intranet endpoint security
- Security group controls access to the intranet endpoint
- This is a **ForceNew** field - changing it will recreate the resource

---

## Example 2: Both Internet and Intranet Endpoints

```hcl
resource "tencentcloud_kubernetes_cluster_endpoint" "full_example" {
  cluster_id                       = "cls-xxxxxxxx"
  
  # Internet (external network) configuration
  cluster_internet                 = true
  cluster_internet_security_group  = "sg-internet-xxxx"
  cluster_internet_domain          = "my-cluster.example.com"
  
  # Intranet (internal network) configuration
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = "subnet-xxxxxxxx"
  cluster_intranet_security_group  = "sg-intranet-xxxx"  # ✅ New field
  cluster_intranet_domain          = "my-cluster.internal"
}
```

**Key Points**:
- Internet security group (`sg-internet-xxxx`) controls external access
- Intranet security group (`sg-intranet-xxxx`) controls internal VPC access
- Both can be configured simultaneously
- Different security groups allow fine-grained access control

---

## Example 3: With Security Group Resource

```hcl
# Create security group for intranet access
resource "tencentcloud_security_group" "intranet_sg" {
  name        = "tke-cluster-intranet-sg"
  description = "Security group for TKE cluster intranet endpoint"
}

resource "tencentcloud_security_group_rule_set" "intranet_rules" {
  security_group_id = tencentcloud_security_group.intranet_sg.id

  ingress {
    action      = "ACCEPT"
    cidr_block  = "10.0.0.0/16"  # Allow VPC CIDR
    protocol    = "TCP"
    port        = "443"
    description = "Allow HTTPS from VPC"
  }
}

# Use the security group for cluster endpoint
resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id                       = "cls-xxxxxxxx"
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = var.subnet_id
  cluster_intranet_security_group  = tencentcloud_security_group.intranet_sg.id
}
```

**Key Points**:
- Security group must be created before endpoint
- Configure appropriate ingress rules for cluster access (typically HTTPS/443)
- Use `depends_on` if rule creation timing is critical

---

## Example 4: Migration from Old Configuration

### Before (Without Security Group)

```hcl
# Old configuration - still works, backward compatible
resource "tencentcloud_kubernetes_cluster_endpoint" "old" {
  cluster_id                 = "cls-xxxxxxxx"
  cluster_intranet           = true
  cluster_intranet_subnet_id = "subnet-xxxxxxxx"
  # No cluster_intranet_security_group - OK, backward compatible
}
```

**Result**: Configuration continues to work (backward compatible)

### After (With Security Group - Recommended)

```hcl
resource "tencentcloud_kubernetes_cluster_endpoint" "new" {
  cluster_id                       = "cls-xxxxxxxx"
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = "subnet-xxxxxxxx"
  cluster_intranet_security_group  = "sg-xxxxxxxx"  # ✅ Added for better security
}
```

**Migration Steps**:
1. Identify existing security group or create new one
2. Add `cluster_intranet_security_group` to configuration (optional)
3. Apply changes (will recreate endpoint due to ForceNew)

---

## Example 5: Configuration Scenarios

### ✅ Valid Configurations

```hcl
# Valid: Intranet with security group (recommended)
resource "tencentcloud_kubernetes_cluster_endpoint" "valid1" {
  cluster_id                       = "cls-xxx"
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = "subnet-xxx"
  cluster_intranet_security_group  = "sg-xxx"
}

# Valid: Intranet without security group (backward compatible)
resource "tencentcloud_kubernetes_cluster_endpoint" "valid2" {
  cluster_id                 = "cls-xxx"
  cluster_intranet           = true
  cluster_intranet_subnet_id = "subnet-xxx"
  # No security group - OK
}

# Valid: Internet with security group
resource "tencentcloud_kubernetes_cluster_endpoint" "valid3" {
  cluster_id                       = "cls-xxx"
  cluster_internet                 = true
  cluster_internet_security_group  = "sg-xxx"
}

# Valid: Both intranet and internet
resource "tencentcloud_kubernetes_cluster_endpoint" "valid4" {
  cluster_id                       = "cls-xxx"
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = "subnet-xxx"
  cluster_intranet_security_group  = "sg-intranet"
  cluster_internet                 = true
  cluster_internet_security_group  = "sg-internet"
}
```

---

## Example 6: ForceNew Behavior

```hcl
# Initial configuration
resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id                       = "cls-xxxxxxxx"
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = "subnet-xxxxxxxx"
  cluster_intranet_security_group  = "sg-old-xxxx"
}
```

**Change security group**:
```hcl
resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id                       = "cls-xxxxxxxx"
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = "subnet-xxxxxxxx"
  cluster_intranet_security_group  = "sg-new-xxxx"  # ⚠️ Changed
}
```

**Terraform Plan Output**:
```
# tencentcloud_kubernetes_cluster_endpoint.example must be replaced
-/+ resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
      ~ cluster_intranet_security_group = "sg-old-xxxx" -> "sg-new-xxxx" # forces replacement
        # ... other attributes
    }
```

**Warning**: Changing `cluster_intranet_security_group` will **destroy and recreate** the endpoint, which may cause brief service disruption.

---

## Example 7: Data Source Reference

```hcl
# Query existing security group
data "tencentcloud_security_groups" "tke_sg" {
  name = "tke-cluster-sg"
}

# Use in cluster endpoint
resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id                       = "cls-xxxxxxxx"
  cluster_intranet                 = true
  cluster_intranet_subnet_id       = var.subnet_id
  cluster_intranet_security_group  = data.tencentcloud_security_groups.tke_sg.security_groups[0].security_group_id
}
```

---

## Comparison: Internet vs Intranet Security Groups

| Feature | `cluster_internet_security_group` | `cluster_intranet_security_group` |
|---------|-----------------------------------|-----------------------------------|
| **Network Type** | Internet (External) | Intranet (Internal VPC) |
| **Required When** | `cluster_internet` = true | `cluster_intranet` = true |
| **ForceNew** | ❌ No (can update) | ✅ Yes (must recreate) |
| **Update Support** | ✅ Yes (via `ModifyClusterEndpointSG`) | ❌ No (API limitation) |
| **Typical Usage** | Public access control | VPC internal access control |

---

## Best Practices

1. **Security Group Design**:
   - Use separate security groups for internet and intranet
   - Restrict intranet security group to VPC CIDR only
   - Follow principle of least privilege

2. **ForceNew Consideration**:
   - Plan security group changes during maintenance windows
   - Use `lifecycle` block to prevent accidental recreation:
     ```hcl
     lifecycle {
       prevent_destroy = true
     }
     ```

3. **Naming Convention**:
   ```hcl
   resource "tencentcloud_security_group" "tke_intranet" {
     name = "${var.cluster_name}-intranet-sg"
   }
   
   resource "tencentcloud_security_group" "tke_internet" {
     name = "${var.cluster_name}-internet-sg"
   }
   ```

4. **Validation**:
   - Always test in non-production environment first
   - Verify security group rules before applying
   - Check cluster connectivity after changes

---

## Troubleshooting

### Error: "cluster_intranet_security_group must be set"

**Cause**: Intranet enabled without security group.

**Solution**: Add security group to configuration:
```hcl
cluster_intranet_security_group = "sg-xxxxxxxx"
```

### Error: "cluster_intranet_security_group can only be set when cluster_intranet is true"

**Cause**: Security group configured but intranet disabled.

**Solution**: Either enable intranet or remove security group:
```hcl
cluster_intranet = true  # Enable intranet
```

### ForceNew Warning

**Cause**: Changing `cluster_intranet_security_group` triggers recreation.

**Solution**: Plan for brief downtime or use blue-green deployment strategy.
