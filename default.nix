{
  lib,
  pkgs,
  hotkeys ? [ ],
  formatter ? "sxhkd",
}:
with lib;
let
  cfg = pkgs.writeScript "config" (concatStringsSep "\n\n" hotkeys);
  izu = pkgs.buildGoApplication rec {
    pname = "izu";
    version = "0.2.1";

    pwd = ./.;
    src = ./.;

    file = pkgs.fetchurl {
      url = "https://raw.githubusercontent.com/xkbcommon/libxkbcommon/master/include/xkbcommon/xkbcommon-keysyms.h";
      hash = "sha256-U5ibymrhoq+glsoB1gDIdgpMaoBp8ySccah7bUfojYc=";
    };

    preBuild = ''
      FILE="${file}" go generate ./...
    '';

    modules = ./gomod2nix.toml;
  };

in
pkgs.stdenv.mkDerivation {
  name = "izu";

  buildInputs = [ izu ];

  phases = "installPhase";

  installPhase = ''
    izu --config ${cfg} --formatter ${formatter} > "$out"
  '';

  meta = {
    description = "A unified hotkey config";
    homepage = "https://github.com/meir/izu";
  };
}
