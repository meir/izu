{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [ (import "${fetchTree gomod2nix.locked}/overlay.nix") ];
    }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
  fetchurl ? pkgs.fetchurl,
}:

buildGoApplication rec {
  pname = "izu";
  version = "0.2.0";
  pwd = ./.;
  src = ./.;
  file = fetchurl {
    url = "https://raw.githubusercontent.com/xkbcommon/libxkbcommon/master/include/xkbcommon/xkbcommon-keysyms.h";
    hash = "sha256-U5ibymrhoq+glsoB1gDIdgpMaoBp8ySccah7bUfojYc=";
  };
  preBuild = ''
    FILE="${file}" go generate ./...
  '';
  modules = ./gomod2nix.toml;
}
