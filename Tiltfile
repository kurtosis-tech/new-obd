# constants
FRONTEND_SERVICE_NAME='frontend'
CARTSERVICE_SERVICE_NAME='cartservice'
PRODUCTCATALOG_SERVICE_NAME='productcatalogservice'
DEV_IMAGE_SUFFIX='-dev-image'
FRONTEND_REF=FRONTEND_SERVICE_NAME + DEV_IMAGE_SUFFIX
CARTSERVICE_REF=CARTSERVICE_SERVICE_NAME + DEV_IMAGE_SUFFIX
PRODUCTCATALOG_REF=PRODUCTCATALOG_SERVICE_NAME + DEV_IMAGE_SUFFIX

# enforces a minimum Tilt version
version_settings(constraint='>=0.22.2')

# check if Kardinal is installed
load('./dev/require-tool', 'require_tool')
require_tool('kardinal', 'Kardinal CLI is not installed, visit https://github.com/kurtosis-tech/kardinal#installation to install it')

# define and get the arguments
config.define_string_list('build')
arguments = config.parse()

# get the 'build' argument
build_arg = []
if 'build' in arguments:
    build_arg = arguments['build']

build_arg_len = len(build_arg)

# validate the 'build' argument
valid_build_arguments = [FRONTEND_SERVICE_NAME, CARTSERVICE_SERVICE_NAME, PRODUCTCATALOG_SERVICE_NAME]
for bu_arg in build_arg:
    if bu_arg not in valid_build_arguments:
        fail('build argument {} is not valid. Valid arguments are: {}'.format(bu_arg, valid_build_arguments))

# enforces to use this Tilt with Minikube or Docker desktop contexts to avoid deploying this stuff in cloud servers
allow_k8s_contexts(['minikube', 'docker-desktop'])

if k8s_context() == 'minikube':
    #This allows Tilt to build and push images directly to Minikube's Docker daemon
    print('Building container images with Minikube')
    local('eval $(minikube docker-env)')

# set the env var to use the local Kontrol
local('export KARDINAL_CLI_DEV_MODE=TRUE')

# clean current flows before creating new ones, it's just for the tilt up cmd
if config.tilt_subcommand == 'up':
    local(['kardinal', 'flow', 'delete', '', '--all'])

    # executing kardinal deploy
    local(['kardinal', 'deploy', '-k', './release/obd-kardinal.yaml'])

    if build_arg_len > 0:
        # execute the kardinal flow create for the services we want to build
        cmd_list = ['kardinal', 'flow', 'create']

        for index, value in enumerate(build_arg):
            service_name = value
            service_image = value + DEV_IMAGE_SUFFIX
            if index == 0:
                cmd_list.extend([service_name, service_image])
            else:
                cmd_list.append('-s')
                add_sv_flag = '{}={}'.format(service_name, service_image)
                cmd_list.append(add_sv_flag)

        # now execute the create flow command
        local(cmd_list)

# CART SERVICE
if 'cartservice' in build_arg:
    docker_build(
        CARTSERVICE_REF,
        context='./src/cartservice',
        dockerfile='./src/cartservice/Dockerfile',
    )

# PRODUCT CATALOG SERVICE
if 'productcatalogservice' in build_arg:
    docker_build(
        PRODUCTCATALOG_REF,
        context='./src/productcatalogservice',
        dockerfile='./src/productcatalogservice/Dockerfile',
    )

# FRONTEND
if 'frontend' in build_arg:
    docker_build(
        FRONTEND_REF,
        context='./src',
        dockerfile='./src/frontend.dockerfile',
    )

kardinal_topology_yaml = local(['kardinal', 'topology', 'print-manifest', '--add-trace-router'], quiet=True)
kardinal_topology_yaml_str = str(kardinal_topology_yaml)

if kardinal_topology_yaml_str != '':
    k8s_yaml(kardinal_topology_yaml, allow_duplicates = True)

local_resource(
    name='ingress-gateway-port-forward',
    serve_cmd=['kubectl', 'port-forward', 'service/istio-ingressgateway', '80:80', '-n', 'istio-system']
)
