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
    mkGoEnv ? pkgs.mkGoEnv,
    gomod2nix ? pkgs.gomod2nix,
    pkgsUnstable ? (import (fetchTree flakeLock.nodes.nixpkgs-unstable.locked) {}),
  }: let
    goEnv = mkGoEnv {pwd = ./.;};
  in
    pkgs.mkShell {
      name = "speaker-shell";
      packages = with pkgs; [
        goEnv
        gomod2nix
        pkgsUnstable.delve
        pkgsUnstable.go
        go-tools
        mpg123
        libxkbcommon
        xorg.libX11.dev
        xorg.libXtst
        pkgsUnstable.translate-shell
        pkgsUnstable.python312Packages.gtts
      ];

      postShellHook = '''';
    }
