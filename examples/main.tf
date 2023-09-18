terraform {
  required_providers {
    kraud = {
      source = "kraud.cloud/app/kraud"
    }
  }
}

provider "kraud" {
  auth_token = "0a8c7134-c51c-4c82-87c5-b8f4f254ea3c:27a0e0e7f7b8a2ef6a84e49640699023062442940a0642fa4a13f72c5ea35b88"
}

data "kraud_volumes" "app" {}

output "app_volumes" {
  value = data.kraud_volumes.app
}
