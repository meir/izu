{
  lib,
  pkgs,
  hotkeys ? [ ],
  formatter ? "sxhkd",
}:
with lib;
let
  cfg = pkgs.writeScript "config" (concatStringsSep "\n\n" hotkeys);
in
buildGoApplication rec {
  pname = "izu";
  version = "0.2.0";

  pwd = ./.;
  src = ./.;

  file = pkgs.fetchurl {
    url = "https://raw.githubusercontent.com/xkbcommon/libxkbcommon/master/include/xkbcommon/xkbcommon-keysyms.h";
    hash = "sha256-U5ibymrhoq+glsoB1gDIdgpMaoBp8ySccah7bUfojYc=";
  };

  phases = "installPhase";

  preBuild = ''
    FILE="${file}" go generate ./...
  '';

  installPhase = ''
    izu --config ${cfg} --formatter ${formatter} > "$out"
  '';

  modules = ./gomod2nix.toml;
}
