# Creating Infrastructure

Linode is supported as the default provider because of the great [Linux Unplugged promotion](https://linode.com/unplugged), feel free to add support for other providers.

## Terraform

Ensure terraform is installed and in PATH

### Linode

```bash
cd terra/linode
make terraform-key
terraform apply
```

# Deployment

Currently, QCards only comes with built-in support for docker/ansible. Please feel free to add support for other deployment options. The below deployment also currently only exists as a deployment of multiple containers on a single host with no redundancy/backup. Improvements would be much appreciated.

## Ansible & Docker

Install ansible and replace intended host DNS/IP in `inventory.yaml`

```bash
# Set env vars
export SSH_KEY_FILE=./terra/linode/key # REPLACE ME
export DB_SECRET='secret' # REPLACE ME
export TOKEN_SIGNING_SECRET='secret' # REPLACE ME
export QPOKER_HOST='https://qcards.xyz' # REPLACE ME

# Provision Server
ansible-playbook -i ansible/inventory.yaml -e "ansible_ssh_private_key_file=$SSH_KEY_FILE" ansible/provision.yaml

# Deploy Postgres
ansible-playbook -i ansible/inventory.yaml -e "ansible_ssh_private_key_file=$SSH_KEY_FILE" -e "db_password=$DB_SECRET" ansible/deploy/deploy_pg.yaml

# Deploy app
ansible-playbook -i ansible/inventory.yaml -e "ansible_ssh_private_key_file=$SSH_KEY_FILE" -e "pg_connection=postgres://qcards:$DB_SECRET@qcards_pg:5432/qcards?sslmode=disable" -e "token_signing_value=$TOKEN_SIGNING_SECRET" -e "qpoker_host=$QPOKER_HOST" ansible/deploy/deploy_app.yaml
```
