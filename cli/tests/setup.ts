import { vi } from 'vitest'

// Mock the SDK module globally
vi.mock('@mimsy/sdk', () => ({
  exportSchema: vi.fn(() => ({
    collections: [],
    generatedAt: new Date().toISOString()
  })),
  clearRegistry: vi.fn(),
  collection: vi.fn(),
  fields: {
    shortString: vi.fn(),
    relation: vi.fn(),
    multiRelation: vi.fn(),
    media: vi.fn()
  },
  builtins: {
    User: {}
  }
}))