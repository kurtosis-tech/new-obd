# enforces a minimum Tilt version
version_settings(constraint='>=0.22.2')

# This allows Tilt to build and push images directly to Minikube's Docker daemon
local('eval $(minikube docker-env)')

# POSTGRES
k8s_yaml('./kubernetes-manifests/postgres.yaml')

# CART SERVICE
docker_build(
    'cartservice',
    context='./src/cartservice',
    dockerfile='./src/cartservice/Dockerfile',
)

k8s_yaml('./kubernetes-manifests/cartservice.yaml')

# PRODUCT CATALOG SERVICE
docker_build(
    'productcatalogservice',
    context='./src/productcatalogservice',
    dockerfile='./src/productcatalogservice/Dockerfile',
)

k8s_yaml('./kubernetes-manifests/productcatalogservice.yaml')

# FRONTEND
docker_build(
    'frontend',
    context='./src',
    dockerfile='./src/frontend.dockerfile',
)

k8s_yaml('./kubernetes-manifests/frontend.yaml')
