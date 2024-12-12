{ pkgs }:
pkgs.buildGoApplication rec {
  pname = "izu";
  version = "0.2.2";

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
}
