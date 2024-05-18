module "database" {
  source = "./modules/database"

  db_name = "customers"
}

module "secret" {
  source = "./modules/secret"
}

module "auth" {
  source = "./modules/auth"

  lambda_name = "login"
  vpc_name    = var.vpc_name

  sign_key = module.secret.sign_key

  security_group_id = module.database.security_group_id

  depends_on = [
    module.secret
  ]
}
