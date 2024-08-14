{pkgs}: let
  pname = "cartservice";
in
  pkgs.buildGoApplication {
    # pname has to match the location (folder) where the main function is or use
    # subPackges to specify the file (e.g. subPackages = ["some/folder/main.go"];)
    inherit pname;
    name = "${pname}";
    pwd = ./.;
    src = ./.;
    modules = ./gomod2nix.toml;
    doCheck = false;
    CGO_ENABLED = 0;
  }
