import { vi } from 'vitest'

export function mockConsole() {
  return {
    log: vi.spyOn(console, 'log').mockImplementation(() => {}),
    error: vi.spyOn(console, 'error').mockImplementation(() => {}),
    warn: vi.spyOn(console, 'warn').mockImplementation(() => {}),
  }
}

export function mockProcess() {
  return {
    exit: vi.spyOn(process, 'exit').mockImplementation((code?: number) => {
      throw new Error(`process.exit called with code ${code}`)
    }),
    stdout: {
      write: vi.spyOn(process.stdout, 'write').mockImplementation(() => true)
    },
    stderr: {
      write: vi.spyOn(process.stderr, 'write').mockImplementation(() => true)
    }
  }
}

export function createMockFileSystem() {
  const files = new Map<string, string>()
  
  return {
    files,
    writeFileSync: vi.fn((path: string, content: string) => {
      files.set(path, content)
    }),
    readFileSync: vi.fn((path: string) => {
      if (!files.has(path)) {
        throw new Error(`ENOENT: no such file or directory, open '${path}'`)
      }
      return files.get(path)
    }),
    existsSync: vi.fn((path: string) => files.has(path)),
  }
}

export function waitForAsync(ms: number = 0): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}