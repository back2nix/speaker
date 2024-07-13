{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
        ];
      }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
}:
buildGoApplication {
  pname = "speaker";
  version = "1.0.0";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;

  buildInputs = with pkgs; [
    translate-shell
    python310Packages.gtts
    mpg123
    libxkbcommon
    xorg.libX11.dev
    xorg.libXtst
  ];
}
