{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs =
    {
      self,
      nixpkgs,
      gomod2nix,
      pre-commit-hooks,
    }:
    let
      forAllSystems =
        function:
        nixpkgs.lib.genAttrs
          [
            "x86_64-linux"
            "aarch64-linux"
            "aarch64-darwin"

          ]
          (
            system:
            function (
              import nixpkgs {
                inherit system;
                overlays = [ gomod2nix.overlays.default ];
              }
            )
          );
    in
    rec {
      packages = forAllSystems (pkgs: rec {
        default = pkgs.callPackage ./izu.nix { };
        izuGenerate = pkgs.callPackage ./izu-generate.nix { izu = default; };
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.callPackage ./shell.nix { inherit pre-commit-hooks pkgs; };
      });
    };

}
