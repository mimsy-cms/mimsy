{inputs, ...}: {
  imports = [
    inputs.process-compose-flake.flakeModule
  ];

  perSystem = {
    config,
    lib,
    pkgs,
    ...
  }: {
    process-compose.dev = {
      imports = [
        inputs.services-flake.processComposeModules.default
      ];

      services.postgres."pg1" = {
        enable = true;

        superuser = "mimsy";

        extensions = exts: [
          exts.system_stats
        ];

        initialDatabases = [
          {
            name = "mimsy";
          }
        ];
      };

      services.pgadmin."pgad1" = {
        enable = true;
        initialEmail = "test@runelabs.xyz";
        initialPassword = "password";
      };
      settings.processes = {
        go-deps = {
          command = ''
            cd api
            ${pkgs.go}/bin/go mod download
            echo "Go dependencies downloaded successfully"
          '';
        };

        pnpm-deps = {
          command = ''
            ${pkgs.pnpm}/bin/pnpm install
            echo "Node dependencies installed successfully"
          '';
        };

        api = {
          command = ''
            cd api
            ${pkgs.air}/bin/air
          '';
          depends_on = {
            go-deps.condition = "process_completed_successfully";
            pg1.condition = "process_healthy";
          };
          environment = {
            POSTGRES_HOST = "localhost";
            POSTGRES_PORT = "5432";
            POSTGRES_USER = "mimsy";
            POSTGRES_PASSWORD = "mimsy";
            POSTGRES_DB = "mimsy";
          };
        };

        web = {
          command = ''
            cd web
            ${pkgs.pnpm}/bin/pnpm dev
          '';
          depends_on = {
            pnpm-deps.condition = "process_completed_successfully";
            api.condition = "process_started";
          };
        };

        landing = {
          command = ''
            cd landing
            ${pkgs.pnpm}/bin/pnpm dev
          '';
          depends_on = {
            pnpm-deps.condition = "process_completed_successfully";
          };
        };
      };
    };
  };
}
