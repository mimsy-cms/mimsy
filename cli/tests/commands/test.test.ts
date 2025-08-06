import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { Command } from 'commander'
import { testCommand, testCommandAction } from '../../src/commands/test'

describe('test command', () => {
  let program: Command
  let consoleSpy: ReturnType<typeof vi.spyOn>

  beforeEach(() => {
    program = new Command()
    consoleSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('should register the test command', () => {
    const result = testCommand(program)
    
    expect(result).toBeInstanceOf(Command)
    const testCmd = program.commands.find(cmd => cmd.name() === 'test')
    expect(testCmd).toBeDefined()
    expect(testCmd?.description()).toBe('Testing command to bootstrap the cli')
  })

  it('should execute and log the correct message', async () => {
    testCommand(program)
    
    // Prevent Commander from calling process.exit
    program.exitOverride()
    
    try {
      await program.parseAsync(['test'], { from: 'user' })
    } catch (e) {
      // Ignore exit override errors
    }
    
    expect(consoleSpy).toHaveBeenCalledWith('Test command executed')
  })

  it('should execute action directly', () => {
    testCommandAction()
    
    expect(consoleSpy).toHaveBeenCalledWith('Test command executed')
  })
})