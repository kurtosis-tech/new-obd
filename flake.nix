{
  description = "Go development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
    unstable.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.inputs.flake-utils.follows = "flake-utils";
  };
  outputs = {
    self,
    nixpkgs,
    flake-utils,
    unstable,
    gomod2nix,
    ...
  }:
    flake-utils.lib.eachDefaultSystem
    (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            (import "${gomod2nix}/overlay.nix")
          ];
        };

        version = toString (self.ref or self.shortRev or self.dirtyShortRev or self.lastModified or "unknown");
        tag-name = let
          branch-name = import ./branch-name.nix;
          tag = toString (
            if isNull branch-name
            then "${version}"
            else "${toString branch-name}"
          );
        in
          builtins.replaceStrings ["/"] ["_"] tag;

        service_names = [
          "cartservice"
          "frontend"
          "metrics"
          "productcatalogservice"
        ];
        architectures = ["amd64" "arm64"];
        imageRegistry = "kurtosistech";

        matchingContainerArch =
          if builtins.match "aarch64-.*" system != null
          then "arm64"
          else if builtins.match "x86_64-.*" system != null
          then "amd64"
          else throw "Unsupported system type: ${system}";

        mergeContainerPackages = acc: service:
          pkgs.lib.recursiveUpdate acc {
            packages."${service}-container" = self.containers.${system}.${service}.${matchingContainerArch};
          };

        multiPlatformDockerPusher = acc: service:
          pkgs.lib.recursiveUpdate acc {
            packages."publish-${service}-container" = let
              name = "${imageRegistry}/${service}";
              tagBase = tag-name;
              images =
                map (
                  arch: rec {
                    inherit arch;
                    image = self.containers.${system}.${service}.${arch};
                    tag = "${tagBase}-${arch}";
                  }
                )
                architectures;
              loadAndPush = builtins.concatStringsSep "\n" (pkgs.lib.concatMap
                ({
                  arch,
                  image,
                  tag,
                }: [
                  "$docker load -i ${image}"
                  "$docker push ${name}:${tag}"
                ])
                images);
              imageNames =
                builtins.concatStringsSep " "
                (map ({
                  arch,
                  image,
                  tag,
                }: "${name}:${tag}")
                images);
            in
              pkgs.writeTextFile {
                inherit name;
                text = ''
                  #!${pkgs.stdenv.shell}
                  set -euxo pipefail
                  docker=${pkgs.docker}/bin/docker
                  ${loadAndPush}
                  $docker manifest create --amend ${name}:${tagBase} ${imageNames}
                  $docker manifest push ${name}:${tagBase}
                '';
                executable = true;
                destination = "/bin/push";
              };
          };

        systemOutput = rec {
          devShells.default = let
            go-tidy-all = pkgs.writeShellApplication {
              name = "go-tidy-all";
              runtimeInputs = with pkgs; [go git gomod2nix];
              text = ''
                root_dirpath=$(git rev-parse --show-toplevel)
                find "$root_dirpath" -type f -name 'go.mod' -exec sh -c 'dir=$(dirname "$1") && cd "$dir" && echo "$dir" && go mod tidy && gomod2nix' shell {} \;
              '';
            };

            tag-branch = pkgs.writeShellApplication {
              name = "tag-branch";
              runtimeInputs = with pkgs; [git];
              text = ''
                GIT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
                CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
                echo "\"$CURRENT_BRANCH\"" >"$GIT_ROOT/branch-name.nix"
                echo "Branch name \"$CURRENT_BRANCH\" written to $GIT_ROOT/branch-name.nix"
              '';
            };

            tag-branch-version = pkgs.writeShellApplication {
              name = "tag-branch-version";
              runtimeInputs = with pkgs; [git];
              text = ''
                GIT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
                CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
                SHORT_COMMIT_HASH=$(git rev-parse --short HEAD)
                echo "\"$CURRENT_BRANCH-$SHORT_COMMIT_HASH\"" >"$GIT_ROOT/branch-name.nix"
                echo "Branch name \"$CURRENT_BRANCH-$SHORT_COMMIT_HASH\" written to $GIT_ROOT/branch-name.nix"
              '';
            };
          in
            pkgs.mkShell {
              buildInputs = with pkgs; [
                go
                gopls
                go-tools
                golangci-lint
                pkgs.gomod2nix
                skaffold
                minikube
                dive
                go-tidy-all
                tag-branch
                tag-branch-version
              ];

              shellHook = ''
                echo "Go development environment loaded"
                go version
              '';
            };

          packages.cartservice = pkgs.callPackage ./src/cartservice {
            inherit pkgs;
          };

          packages.metrics = pkgs.callPackage ./src/metrics {
            inherit pkgs;
          };

          packages.frontend = pkgs.callPackage ./src/frontend {
            inherit pkgs;
          };

          packages.productcatalogservice = pkgs.callPackage ./src/productcatalogservice {
            inherit pkgs;
          };

          containers = let
            os = "linux";
            all =
              pkgs.lib.mapCartesianProduct ({
                arch,
                service_name,
              }: {
                "${service_name}" = {
                  "${toString arch}" = let
                    nix_arch =
                      builtins.replaceStrings
                      ["arm64" "amd64"] ["aarch64" "x86_64"]
                      arch;

                    container_pkgs = import nixpkgs {
                      system = "${nix_arch}-${os}";
                    };

                    # if running from linux no cross-compilation is needed to palce the service in a container
                    needsCrossCompilation =
                      "${nix_arch}-${os}"
                      != system;

                    service =
                      if !needsCrossCompilation
                      then
                        packages.${service_name}.overrideAttrs
                        (old: old // {doCheck = false;})
                      else
                        packages.${service_name}.overrideAttrs (old:
                          old
                          // {
                            GOOS = os;
                            GOARCH = arch;
                            # CGO_ENABLED = disabled breaks the CLI compilation
                            # CGO_ENABLED = 0;
                            doCheck = false;
                          });
                  in
                    pkgs.dockerTools.buildImage {
                      name = "${imageRegistry}/${service_name}";
                      tag = "${tag-name}-${arch}";
                      # tag = commit_hash;
                      created = "now";
                      copyToRoot = pkgs.buildEnv {
                        name = "image-root";
                        paths = [
                          service
                          container_pkgs.bashInteractive
                          container_pkgs.nettools
                          container_pkgs.gnugrep
                          container_pkgs.coreutils
                          container_pkgs.cacert
                        ];
                        pathsToLink = ["/bin"];
                      };
                      architecture = arch;
                      config.Cmd =
                        if !needsCrossCompilation
                        then ["${service}/bin/${service.pname}"]
                        else ["${service}/bin/${os}_${arch}/${service.pname}"];
                      config.Env = ["SSL_CERT_FILE=${container_pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"];
                    };
                };
              }) {
                arch = architectures;
                service_name = service_names;
              };
          in
            pkgs.lib.foldl' (set: acc: pkgs.lib.recursiveUpdate acc set) {}
            all;
        };
        # Add containers matching architecture with local system as toplevel packages
        # this means calling `nix build .#<SERVICE_NAME>-container` will build the container matching the local system.
        # For cross-compilation use the containers attribute directly: `nix build .containers.<LOCAL_SYSTEM>.<SERVICE_NAME>.<ARCH>`
        outputWithContaniers = pkgs.lib.foldl' mergeContainerPackages systemOutput service_names;
        outputWithContainersAndPushers = pkgs.lib.foldl' multiPlatformDockerPusher outputWithContaniers service_names;
      in
        outputWithContainersAndPushers
    );
}
