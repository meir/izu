{
  lib,
  pkgs,
  hotkeys ? [ ],
  formatter ? "sxhkd",
  izu,
}:
with lib;
let
  cfg = pkgs.writeScript "config" (concatStringsSep "\n\n" hotkeys);
in
pkgs.stdenv.mkDerivation {
  name = "izu";

  buildInputs = [ izu ];

  phases = "installPhase";

  installPhase = ''
    izu --config ${cfg} --formatter ${formatter} > "$out"
  '';

  meta = {
    description = "A unified hotkey config";
    homepage = "https://github.com/meir/izu";
  };
}
