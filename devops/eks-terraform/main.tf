module "dev_cluster" {
  source = "./cluster"
  cluster_name = "dev"
}

module "staging_cluster" {
  source = "./cluster"
  cluster_name = "staging"
}

module "prod_cluster" {
  source = "./cluster"
  cluster_name = "production"
}


