{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix nix-buildproxy;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
        nix-buildproxy.overlays.default
      ];
    }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
  mkBuildproxy,
}:

buildGoApplication rec {
  pname = "izu";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  preBuild = ''
    source ${mkBuildproxy ./proxy_content.nix}
    go generate ./...
  '';
  modules = ./gomod2nix.toml;
}
