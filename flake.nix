{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.inputs.flake-utils.follows = "flake-utils";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
      pre-commit-hooks,
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
      in
      rec {
        packages.default = pkgs.callPackage ./izu.nix { inherit pkgs; };
        packages.izuGenerate = pkgs.callPackage ./. {
          inherit pkgs;
          izu = packages.default;
          formatter = "sxhkd";
          hotkeys = [ ];
        };

        overlays.default = (
          final: prev: {
            izu = packages.default;
            izuGenerate = formatter: hotkeys: packages.izuGenerate.override { inherit formatter hotkeys; };
          }
        );

        devShells.default = pkgs.callPackage ./shell.nix { inherit pre-commit-hooks pkgs; };
      }
    ));
}
