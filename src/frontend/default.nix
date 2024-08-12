{pkgs}: let
  pname = "frontend";
  staticFiles = pkgs.stdenv.mkDerivation {
    name = "${pname}-static";
    src = ./static;
    installPhase = ''
      mkdir -p $out/static
      cp -r * $out/static/
    '';
  };
  templateFiles = pkgs.stdenv.mkDerivation {
    name = "${pname}-static";
    src = ./templates;
    installPhase = ''
      mkdir -p $out/templates
      cp -r * $out/templates/
    '';
  };
  app = pkgs.buildGoApplication {
    # pname has to match the location (folder) where the main function is or use
    # subPackges to specify the file (e.g. subPackages = ["some/folder/main.go"];)
    inherit pname;
    name = "${pname}-app";
    pwd = ./.;
    src = ./.;
    modules = ./gomod2nix.toml;
    doCheck = false;
    CGO_ENABLED = 0;
  };
in
  pkgs.stdenv.mkDerivation {
    name = "${pname}";
    phases = ["installPhase"];
    installPhase = ''
      mkdir -p $out
      ln -s ${staticFiles}/static $out/static
      ln -s ${templateFiles}/templates $out/templates
      ln -s ${app}/bin/${pname} $out/${pname}
    '';
  }
