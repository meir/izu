{
  pkgs,
  mkGoEnv,
  gomod2nix,
  pre-commit-hooks,
}:

let
  goEnv = mkGoEnv { pwd = ./.; };
  pre-commit-check = pre-commit-hooks.lib.${pkgs.system}.run {
    src = ./.;
    hooks = {
      gofmt.enable = true;
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
  ];
}
