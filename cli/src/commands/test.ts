import { Command } from 'commander'

export function testCommand(program: Command) {
  program
    .command('test')
    .description('Testing command to bootstrap the cli')
    .action(() => {
      console.log('Test command executed')
    })
}
