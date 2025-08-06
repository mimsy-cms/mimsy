#!/usr/bin/env node
import { Command } from 'commander'
import { testCommand } from './commands/test'
import { exportSchemaCommand } from './commands/export-schema'
import { applyCommand } from './commands/apply'

export function createProgram(): Command {
  const program = new Command()
  
  program
    .name('msy')
    .description('A CLI tool for mimsy, the simple SvelteKit CMS')
    .version('1.0.0')
  
  testCommand(program)
  exportSchemaCommand(program)
  applyCommand(program)
  
  return program
}

export const program = createProgram()

// Only parse if this is the main module
if (require.main === module) {
  program.parse(process.argv)
}
