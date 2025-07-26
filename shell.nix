{
  # Чистая функция. Все аргументы передаются из flake.nix.
  pkgs,
  pkgsUnstable,
}: let
  goEnv = pkgsUnstable.mkGoEnv { pwd = ./.; };
in
  pkgs.mkShell {
    name = "speaker-shell";
    packages =
      # Стабильные пакеты
      (with pkgs; [
        libxkbcommon
        xorg.libX11.dev
        xorg.libXtst
        mpg123
      ])
      # Нестабильные пакеты
      ++ (with pkgsUnstable; [
        go
        go-tools
        delve
        gomod2nix
        translate-shell
        python312Packages.gtts
        goEnv
      ]);

    postShellHook = '''';
  }
