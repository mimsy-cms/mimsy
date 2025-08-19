{
  description = "Mimsy development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    process-compose-flake.url = "github:Platonic-Systems/process-compose-flake";
    services-flake.url = "github:juspay/services-flake";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];

      imports = [
        ./flake-modules/process-compose.nix
      ];

      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: {
        # Development shell
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go development
            go
            gopls
            go-tools
            golangci-lint
            air
            delve
            pgroll
            mockgen

            # Node.js development
            nodejs_22
            pnpm
            nodePackages.typescript
            nodePackages.svelte-language-server
            nodePackages.vscode-langservers-extracted

            # Development utilities
            git
            curl
            jq
            httpie
            just
            watchexec

            # Nix tools
            nil
            nixpkgs-fmt
            gh
          ];

          shellHook = ''
            echo "🚀 Mimsy development environment"
            echo ""
            echo "Available commands:"
            echo "  cd api && go run .     - Run API server"
            echo "  cd landing && pnpm dev - Run landing page"
            echo "  cd web && pnpm dev     - Run web app"
            echo ""
            echo "Services will be available at:"
            echo "  API:     http://localhost:8080"
            echo "  Landing: http://localhost:5173"
            echo "  Web:     http://localhost:5174"
            echo ""
          '';

          # Environment variables
        };
      };
    };
}
