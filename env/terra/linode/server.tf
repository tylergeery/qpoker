provider "linode" {
    token = var.api_token
}

resource "linode_instance" "linode-qpoker" {
    label = "linode-qpoker"
    image = "linode/ubuntu18.04"
    region = "us-west"
    type = "g6-nanode-1"
    authorized_keys = [var.authorized_key]
}
