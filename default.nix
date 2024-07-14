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
  lib,
}: let
in
  buildGoApplication {
    pname = "speaker";
    version = "1.0.0";
    pwd = ./.;
    src = ./.;
    modules = ./gomod2nix.toml;

    buildInputs = with pkgs; [
      libxkbcommon
      xorg.libX11
      xorg.libXtst
      translate-shell
      python310Packages.gtts
      mpg123
    ];

    nativeBuildInputs = with pkgs; [makeWrapper];

    postInstall = ''
      cp -r sound $out/;

      wrapProgram "$out/bin/speaker" \
      --prefix PATH : ${lib.makeBinPath [pkgs.mpg123 pkgs.translate-shell pkgs.python310Packages.gtts pkgs.xorg.libXtst pkgs.libxkbcommon pkgs.xorg.libX11]}:$out \
    '';
  }
