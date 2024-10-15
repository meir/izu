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
  version = "0.1";
  pwd = ./.;
  src = ./.;
  file = fetchurl {
    url = "https://raw.githubusercontent.com/xkbcommon/libxkbcommon/master/include/xkbcommon/xkbcommon-keysyms.h";
    hash = "sha256-uPDT22f98wWHAKcBP7QEsrDUP4mGKizzvLsEeWAZEjE=";
  };
  preBuild = ''
    FILE="${file}" go generate ./...
  '';
  modules = ./gomod2nix.toml;
}
