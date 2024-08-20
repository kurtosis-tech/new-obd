# enforces a minimum Tilt version
version_settings(constraint='>=0.22.2')

# This allows Tilt to build and push images directly to Minikube's Docker daemon
local('eval $(minikube docker-env)')

# CART SERVICE
docker_build(
    'kurtosistech/cartservice',
    context='./src/cartservice',
    dockerfile='./src/cartservice/Dockerfile',
)

# PRODUCT CATALOG SERVICE
docker_build(
    'kurtosistech/productcatalogservice',
    context='./src/productcatalogservice',
    dockerfile='./src/productcatalogservice/Dockerfile',
)

# FRONTEND
docker_build(
    'kurtosistech/frontend',
    context='./src',
    ignore='./src/cartservice',
    dockerfile='./src/frontend.dockerfile',
)

k8s_yaml('./kubernetes-manifests/trace-router.yaml')

apply_command = ['python3', 'kontrol-get-yaml.py']

k8s_yaml(local(apply_command))

k8s_resource( new_name='istio-ingressgateway',objects=['istio-ingressgateway:Service:istio-system'], port_forwards=[8080, 8443])

