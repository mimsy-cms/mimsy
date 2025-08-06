import { Command } from 'commander'

export function testCommandAction(): void {
  console.log('Test command executed')
}

export function testCommand(program: Command): Command {
  return program
    .command('test')
    .description('Testing command to bootstrap the cli')
    .action(testCommandAction)
}
