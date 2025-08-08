import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createProgram } from '../src/index'
import { Command } from 'commander'

describe('CLI integration tests', () => {
  let consoleSpy: ReturnType<typeof vi.spyOn>
  let processExitSpy: ReturnType<typeof vi.spyOn>
  let program: Command

  beforeEach(() => {
    consoleSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
    processExitSpy = vi.spyOn(process, 'exit').mockImplementation(() => {
      throw new Error('process.exit called')
    })
    program = createProgram()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('should have correct program metadata', () => {
    expect(program.name()).toBe('msy')
    expect(program.description()).toBe('A CLI tool for mimsy, the simple SvelteKit CMS')
    expect(program.version()).toBe('1.0.0')
  })

  it('should have all commands registered', () => {
    const commandNames = program.commands.map(cmd => cmd.name())
    expect(commandNames).toContain('test')
    expect(commandNames).toContain('export-schema')
  })

  it('should execute test command', async () => {
    program.exitOverride()
    
    try {
      await program.parseAsync(['test'], { from: 'user' })
    } catch (e) {
      // Ignore exit override errors
    }
    
    expect(consoleSpy).toHaveBeenCalledWith('Test command executed')
  })

  it('should show help when --help flag is used', () => {
    const writeSpy = vi.spyOn(process.stdout, 'write').mockImplementation(() => true)
    const helpProgram = createProgram()
    helpProgram.exitOverride()
    
    try {
      helpProgram.parse(['node', 'msy', '--help'], { from: 'user' })
    } catch (error) {
      // Expected to throw due to exitOverride
    }
    
    const output = writeSpy.mock.calls.map(call => call[0]).join('')
    expect(output).toContain('A CLI tool for mimsy')
    expect(output).toContain('test')
    expect(output).toContain('export-schema')
  })

  it('should show version with --version flag', () => {
    const writeSpy = vi.spyOn(process.stdout, 'write').mockImplementation(() => true)
    const versionProgram = createProgram()
    versionProgram.exitOverride()
    
    try {
      versionProgram.parse(['node', 'msy', '--version'], { from: 'user' })
    } catch (error) {
      // Expected to throw due to exitOverride
    }
    
    expect(writeSpy).toHaveBeenCalledWith(expect.stringContaining('1.0.0'))
  })

  it('should handle unknown commands gracefully', () => {
    const errorSpy = vi.spyOn(process.stderr, 'write').mockImplementation(() => true)
    const unknownProgram = createProgram()
    unknownProgram.exitOverride()
    
    expect(() => {
      unknownProgram.parse(['node', 'msy', 'unknown-command'], { from: 'user' })
    }).toThrow()
    
    const errorOutput = errorSpy.mock.calls.map(call => call[0]).join('')
    expect(errorOutput).toContain('unknown command')
  })
})