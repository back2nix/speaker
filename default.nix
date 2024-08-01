let
  inherit (builtins) fetchTree fromJSON readFile;
  flakeLock = fromJSON (readFile ./flake.lock);
in
  {
    pkgs ? (import (fetchTree flakeLock.nodes.nixpkgs.locked) {
      overlays = [
        (import "${fetchTree flakeLock.nodes.gomod2nix.locked}/overlay.nix")
      ];
    }),
    buildGoApplication ? pkgs.buildGoApplication,
    pkgsUnstable ? (import (fetchTree flakeLock.nodes.nixpkgs-unstable.locked) {}),
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
        pkgsUnstable.go
        pkgsUnstable.delve
      ];

      nativeBuildInputs = with pkgs; [makeWrapper];

      postInstall = ''
        cp -r sound $out/;

        wrapProgram "$out/bin/speaker" \
        --prefix PATH : ${lib.makeBinPath [
          pkgs.mpg123
          pkgs.xorg.libXtst
          pkgs.libxkbcommon
          pkgs.xorg.libX11
          pkgsUnstable.translate-shell
          pkgsUnstable.python312Packages.gtts
        ]}:$out
      '';
    }
