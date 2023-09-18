terraform {
  required_providers {
    kraud = {
      source = "kraud.cloud/app/kraud"
    }
  }
}

provider "kraud" {
  auth_token = "0a8c7134-c51c-4c82-87c5-b8f4f254ea3c:e8984daa2af7bf20c8dea853c4ab65ccdcfd5c0ead096bc62e0332c188f2ddd2"
}

data "kraud_volumes" "example" {}

# output "edu_coffees" {
#   value = data.hashicups_coffees.edu
# }
