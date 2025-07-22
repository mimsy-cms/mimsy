#!/usr/bin/env node
import { Command } from 'commander'
import { testCommand } from './commands/test'

export const program = new Command()

program
  .name('msy')
  .description('A CLI tool for mimsy, the simple SvelteKit CMS')
  .version('1.0.0')

testCommand(program)

program.parse(process.argv)
