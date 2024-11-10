provider "aws" {
  region = "us-east-1a"
}


data "aws_availability_zones" "available" {} // This is a data source that fetches the availability zones in the region specified in the provider block.

data "aws_eks_cluster" "cluster" {
    name = module.eks.cluster_id
} // This is a data source that fetches the AWS Elastic Kubernetes Service (EKS) cluster details.

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
} // This is a data source that fetches the authentication details for the EKS cluster.

locals {
  cluster_name = "learnk8s"
} // This is a local block that defines a variable cluster_name with the value learnk8s.

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
} // This is a provider block that configures the Kubernetes provider with the EKS cluster details.

module "eks-kubeconfig" {
  source  = "hyperbadger/eks-kubeconfig/aws"
  version = "1.0.0"

  depends_on = [module.eks]
  cluster_id = module.eks.cluster_id
} // This is a module block that fetches the kubeconfig for the EKS cluster.

resource "local_file" "name" {
  content = module.eks-kubeconfig.kubeconfig
  filename = "kubeconfig_${local.cluster_name}"
} // This is a resource block that writes the kubeconfig details to a file on your local storage.

module "vpc" {
    source = "terraform-aws-modules/vpc/aws"
    version = "3.18.1"

    name = "k8s-vpc"
    cidr = "172.16.0.0/16"
    azs = data.aws_availability_zones.available.names
    private_subnets = ["172.16.1.0/24", "172.16.2.0/24","172.16.3.0/24"]
    public_subnets = ["172.16.4.0/24", "172.16.5.0/24","172.16.6.0/24"]
    enable_nat_gateway = true
    single_nat_gateway = true
    enable_dns_hostnames = true

    public_subnet_tags = {
        "kubernetes.io/cluster/${local.cluster_name}" = "shared"
        "kubernetes.io/role/elb" = "1"
    }
} // This is a module block that creates a VPC with public and private subnets, using the vpc module from the terraform-aws-modules GitHub repository. It specifies the VPC name, CIDR block, availability zones, subnets, NAT gateway, DNS hostnames, and subnet tags.

module "eks" {
  source = "terraform-aws-modules/eks/aws"
    version = "18.30.3"
   cluster_name = "${local.cluster_name}"
   cluster_version = "1.24"
   subnet_ids = module.vpc.private_subnets
   vpc_id = module.vpc.vpc_id

   eks_managed_node_groups = {
    first= {
        desired_capacity = 1
        max_capacity = 5
        min_capacity = 1

        instance_type = "t2.micro"
    }
   }
} // This is a module block that creates an EKS cluster with a managed node group, using the eks module from the terraform-aws-modules GitHub repository. It specifies the cluster name, version, subnet IDs, VPC ID, and node group details. which include the desired capacity, max capacity, min capacity, and instance type.

