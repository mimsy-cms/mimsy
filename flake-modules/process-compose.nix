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
    };
  };
}
