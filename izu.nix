{ pkgs, lib }:
pkgs.buildGoApplication rec {
  pname = "izu";
  version = lib.readFile ./VERSION;

  pwd = ./.;
  src = ./.;

  file = pkgs.fetchurl {
    url = "https://raw.githubusercontent.com/xkbcommon/libxkbcommon/master/include/xkbcommon/xkbcommon-keysyms.h";
    hash = "sha256-/f17XP7ASLZR6LgTS/tnP6y2+mJDld7nMUz8W+FJli8=";
  };

  preBuild = ''
    FILE="${file}" go generate ./...
  '';

  modules = ./gomod2nix.toml;
}
