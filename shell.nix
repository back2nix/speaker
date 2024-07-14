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
  mkGoEnv ? pkgs.mkGoEnv,
  gomod2nix ? pkgs.gomod2nix,
  pkgsUnstable,
}: let
  goEnv = mkGoEnv {pwd = ./.;};
in
  pkgs.mkShell {
    name = "speaker-shell";
    packages = with pkgs; [
      goEnv
      gomod2nix
      go-tools
      pkgsUnstable.translate-shell
      pkgsUnstable.python312Packages.gtts
      mpg123
      libxkbcommon
      xorg.libX11.dev
      xorg.libXtst
    ];

    postShellHook = ''
    '';
  }
