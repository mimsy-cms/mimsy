#!/usr/bin/env node
import { Command } from "commander";
import { updateCommand } from "./commands/update";
import { initCommand } from "./commands/init";

export function createProgram(): Command {
  const program = new Command();

  program
    .name("msy")
    .description("A CLI tool for mimsy, the simple SvelteKit CMS")
    .version("1.0.0");

  initCommand(program);
  updateCommand(program);

  return program;
}

export const program = createProgram();

// Only parse if this is the main module
if (require.main === module) {
  program.parse(process.argv);
}
