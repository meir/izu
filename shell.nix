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
  mkGoEnv ? pkgs.mkGoEnv,
  gomod2nix ? pkgs.gomod2nix,
  pre-commit-hooks,
  buildproxy-capture,
}:

let
  goEnv = mkGoEnv { pwd = ./.; };
  pre-commit-check = pre-commit-hooks.lib.${pkgs.system}.run {
    src = ./.;
    hooks = {
      gofmt.enable = true;
      # broken
      # golangci-lint = {
      #   enable = true;
      #   name = "golangci-lint";
      #   description = "Lint my golang code";
      #   files = "\.go$";
      #   entry = "${pkgs.golangci-lint}/bin/golangci-lint run --new-from-rev HEAD --fix";
      #   require_serial = true;
      #   pass_filenames = false;
      # };
      goimports = {
        enable = true;
        name = "goimports";
        description = "Format my golang code";
        files = "\.go$";
        entry =
          let
            script = pkgs.writeShellScript "precommit-goimports" ''
              set -e
              failed=false
              for file in "$@"; do
                  # redirect stderr so that violations and summaries are properly interleaved.
                  if ! ${pkgs.gotools}/bin/goimports -l -d "$file" 2>&1
                  then
                      failed=true
                  fi
              done
              if [[ $failed == "true" ]]; then
                  exit 1
              fi
            '';
          in
          builtins.toString script;
      };
    };
  };
in
pkgs.mkShell {
  inherit (pre-commit-check) shellHook;
  packages = [
    goEnv
    gomod2nix
    pkgs.go_1_22
    pkgs.gotools
    pkgs.go-junit-report
    pkgs.go-task
    pkgs.delve
    buildproxy-capture
  ];
}
