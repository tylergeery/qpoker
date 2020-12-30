terraform-key:
	ssh-keygen -f ./key -t rsa -b 4096

terraform-apply:
	eval $(egrep -v '^#' .env | xargs) terraform apply

terraform-plan:
	eval "$(egrep -v '^#' .env | xargs)" terraform plan
