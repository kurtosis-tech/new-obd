{pkgs}: let
  pname = "productcatalogservice";
  dataFiles = pkgs.stdenv.mkDerivation {
    name = "${pname}-static";
    src = ./data;
    installPhase = ''
      mkdir -p $out/data
      cp -r * $out/data/
    '';
  };
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
    postInstall = ''
      # find relative path of the binary file
      PNAME_PATH=$(find $out -name "${pname}" | head -n 1)
      PNAME_FOLDER=$(dirname "$PNAME_PATH")

      # link static files to the same folder as the binary
      mkdir -p "$PNAME_FOLDER"
      cp -R "${dataFiles}/data" "$PNAME_FOLDER/"

      # wrapper script to start the binary in the same folder as the static files
      mv $PNAME_PATH $PNAME_FOLDER/bin-${pname}
      echo "#!/bin/bash" > $PNAME_PATH
      echo "cd $PNAME_FOLDER && $PNAME_FOLDER/bin-${pname}" >> $PNAME_PATH
      chmod +x $PNAME_PATH
    '';
  }
