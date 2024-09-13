{
  description = "A basic gomod2nix flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";
  inputs.pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  inputs.nix-buildproxy.url = "github:polygon/nix-buildproxy";

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
      pre-commit-hooks,
      nix-buildproxy,
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        gomod2nixPkgs = gomod2nix.legacyPackages.${system};
        buildproxy = (nix-buildproxy.overlays.default nixpkgs nixpkgs);

        # The current default sdk for macOS fails to compile go projects, so we use a newer one for now.
        # This has no effect on other platforms.
        callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;
      in
      {
        packages.default = callPackage ./. {
          inherit (gomod2nixPkgs) buildGoApplication;
          inherit (buildproxy.lib) mkBuildproxy;
        };
        devShells.default = callPackage ./shell.nix {
          inherit (gomod2nixPkgs) mkGoEnv gomod2nix;
          inherit pre-commit-hooks;
          inherit (buildproxy) buildproxy-capture;
        };
      }
    ));
}
