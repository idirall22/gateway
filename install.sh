helm repo update
helm install nginx-ingress stable/nginx-ingress --set controller.publishService.enabled=true