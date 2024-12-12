{ packages }:
final: prev: {
  izu = packages.default;
  izuGenerate = formatter: hotkeys: packages.izuGenerate.override { inherit formatter hotkeys; };
}
