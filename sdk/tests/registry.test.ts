import { describe, expect, test, beforeEach } from "vitest";
import {
  registerCollection,
  registerGlobal,
  getAllCollections,
  clearRegistry,
  getCollection,
  exportSchema,
} from "$src/registry";
import { collection, global } from "$src/collection";
import { fields } from "$src/index";
import { User, Media } from "$src/builtins";

describe("Registry", () => {
  beforeEach(() => {
    clearRegistry();
  });

  describe("collection registration", () => {
    test("should register a collection successfully", () => {
      const testCollection = collection("test", {
        title: fields.shortString(),
        content: fields.richText(),
      });

      const allCollections = getAllCollections();
      expect(allCollections).toHaveLength(1);
      expect(allCollections[0].name).toBe("test");
    });

    test("should warn when overwriting an existing collection", () => {
      const consoleSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
      
      collection("duplicate", {
        field1: fields.shortString(),
      });
      
      collection("duplicate", {
        field2: fields.number(),
      });

      expect(consoleSpy).toHaveBeenCalledWith(
        '[Mimsy SDK] Warning: A collection with the name "duplicate" is already registered. It will be overwritten.'
      );
      
      const registered = getCollection("duplicate");
      expect(registered?.schema).toHaveProperty("field2");
      expect(registered?.schema).not.toHaveProperty("field1");
      
      consoleSpy.mockRestore();
    });

    test("should throw error for reserved names", () => {
      expect(() => {
        collection("user", { name: fields.shortString() });
      }).toThrow("Collection name 'user' is reserved");
    });

    test("should retrieve a specific collection by name", () => {
      const testCol = collection("retrievable", {
        field: fields.shortString(),
      });

      const retrieved = getCollection("retrievable");
      expect(retrieved).toBeDefined();
      expect(retrieved?.name).toBe("retrievable");
      expect(retrieved?.isGlobal).toBe(false);
    });

    test("should return undefined for non-existent collection", () => {
      const result = getCollection("nonexistent");
      expect(result).toBeUndefined();
    });
  });

  describe("global registration", () => {
    test("should register a global successfully", () => {
      const testGlobal = global("settings", {
        siteTitle: fields.shortString(),
        maintenance: fields.checkbox(),
      });

      const allCollections = getAllCollections();
      expect(allCollections).toHaveLength(1);
      expect(allCollections[0].name).toBe("settings");
      expect(allCollections[0].isGlobal).toBe(true);
    });

    test("should warn when overwriting an existing global", () => {
      const consoleSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
      
      global("config", {
        field1: fields.shortString(),
      });
      
      global("config", {
        field2: fields.number(),
      });

      expect(consoleSpy).toHaveBeenCalledWith(
        '[Mimsy SDK] Warning: A global with the name "config" is already registered. It will be overwritten.'
      );
      
      consoleSpy.mockRestore();
    });

    test("should throw error for reserved global names", () => {
      expect(() => {
        global("media", { name: fields.shortString() });
      }).toThrow("Global name 'media' is reserved");
    });
  });

  describe("registry management", () => {
    test("should clear all registered collections", () => {
      collection("test1", { field: fields.shortString() });
      collection("test2", { field: fields.number() });
      global("test3", { field: fields.checkbox() });

      expect(getAllCollections()).toHaveLength(3);
      
      clearRegistry();
      
      expect(getAllCollections()).toHaveLength(0);
    });

    test("should return all collections including globals", () => {
      const col1 = collection("posts", {
        title: fields.shortString(),
      });
      
      const col2 = global("settings", {
        siteName: fields.shortString(),
      });
      
      const col3 = collection("pages", {
        content: fields.richText(),
      });

      const all = getAllCollections();
      expect(all).toHaveLength(3);
      expect(all.map(c => c.name)).toContain("posts");
      expect(all.map(c => c.name)).toContain("settings");
      expect(all.map(c => c.name)).toContain("pages");
    });
  });

  describe("schema serialization with underscore fields", () => {
    test("should ignore fields starting with underscore during serialization", () => {
      const testCollection = collection("withUnderscore", {
        publicField: fields.shortString(),
        _privateField: fields.number(),
        anotherPublic: fields.checkbox(),
        _anotherPrivate: fields.richText(),
      });

      const exported = exportSchema();
      const serializedCollection = exported.collections.find(
        c => c.name === "withUnderscore"
      );

      expect(serializedCollection).toBeDefined();
      expect(serializedCollection?.schema).toHaveProperty("publicField");
      expect(serializedCollection?.schema).toHaveProperty("anotherPublic");
      expect(serializedCollection?.schema).not.toHaveProperty("_privateField");
      expect(serializedCollection?.schema).not.toHaveProperty("_anotherPrivate");
    });

    test("should handle collections with only underscore fields", () => {
      const testCollection = collection("onlyPrivate", {
        _private1: fields.shortString(),
        _private2: fields.number(),
        _private3: fields.checkbox(),
      });

      const exported = exportSchema();
      const serializedCollection = exported.collections.find(
        c => c.name === "onlyPrivate"
      );

      expect(serializedCollection).toBeDefined();
      expect(serializedCollection?.schema).toEqual({});
    });

    test("should handle mixed underscore and regular fields in complex schema", () => {
      const complexCollection = collection("complex", {
        title: fields.shortString(),
        _internalId: fields.number(),
        content: fields.richText(),
        _metadata: fields.shortString(),
        author: fields.relation({ relatesTo: User }),
        _cachedData: fields.richText(),
        published: fields.checkbox(),
      });

      const exported = exportSchema();
      const serializedCollection = exported.collections.find(
        c => c.name === "complex"
      );

      expect(serializedCollection?.schema).toHaveProperty("title");
      expect(serializedCollection?.schema).toHaveProperty("content");
      expect(serializedCollection?.schema).toHaveProperty("author");
      expect(serializedCollection?.schema).toHaveProperty("published");
      expect(serializedCollection?.schema).not.toHaveProperty("_internalId");
      expect(serializedCollection?.schema).not.toHaveProperty("_metadata");
      expect(serializedCollection?.schema).not.toHaveProperty("_cachedData");
    });

    test("should preserve underscore fields in non-serialized collection object", () => {
      const testCollection = collection("preserveUnderscore", {
        publicField: fields.shortString(),
        _privateField: fields.number(),
      });

      const retrieved = getCollection("preserveUnderscore");
      expect(retrieved?.schema).toHaveProperty("publicField");
      expect(retrieved?.schema).toHaveProperty("_privateField");
    });

    test("should correctly serialize relation fields while ignoring underscore fields", () => {
      const categoryCollection = collection("categories", {
        name: fields.shortString(),
        _internalSlug: fields.shortString(),
      });

      const postCollection = collection("posts", {
        title: fields.shortString(),
        _draft: fields.checkbox(),
        category: fields.relation({ relatesTo: categoryCollection }),
        author: fields.relation({ relatesTo: User }),
        _tempData: fields.richText(),
      });

      const exported = exportSchema();
      
      const serializedPosts = exported.collections.find(c => c.name === "posts");
      expect(serializedPosts?.schema).toEqual({
        title: { type: "string" },
        category: {
          type: "relation",
          relatesTo: "categories",
        },
        author: {
          type: "relation",
          relatesTo: "<builtins.user>",
        },
      });

      const serializedCategories = exported.collections.find(c => c.name === "categories");
      expect(serializedCategories?.schema).toEqual({
        name: { type: "string" },
      });
    });

    test("should handle global collections with underscore fields", () => {
      const settingsGlobal = global("settings", {
        siteName: fields.shortString(),
        _apiKey: fields.shortString(),
        maintenanceMode: fields.checkbox(),
        _internalConfig: fields.richText(),
      });

      const exported = exportSchema();
      const serializedGlobal = exported.collections.find(
        c => c.name === "settings"
      );

      expect(serializedGlobal).toBeDefined();
      expect(serializedGlobal?.isGlobal).toBe(true);
      expect(serializedGlobal?.schema).toHaveProperty("siteName");
      expect(serializedGlobal?.schema).toHaveProperty("maintenanceMode");
      expect(serializedGlobal?.schema).not.toHaveProperty("_apiKey");
      expect(serializedGlobal?.schema).not.toHaveProperty("_internalConfig");
    });

    test("should correctly serialize field options while ignoring underscore fields", () => {
      const testCollection = collection("withOptions", {
        requiredField: fields.shortString({ 
          constraints: { required: true, minLength: 5 }
        }),
        _privateRequired: fields.shortString({ 
          constraints: { required: true }
        }),
        optionalField: fields.number({ 
          constraints: { min: 1, max: 100 }
        }),
        _privateWithDefault: fields.number({ 
          constraints: { min: 0 }
        }),
      });

      const exported = exportSchema();
      const serializedCollection = exported.collections.find(
        c => c.name === "withOptions"
      );

      expect(serializedCollection?.schema).toEqual({
        requiredField: {
          type: "string",
          options: { 
            constraints: { required: true, minLength: 5 }
          },
        },
        optionalField: {
          type: "number",
          options: { 
            constraints: { min: 1, max: 100 }
          },
        },
      });
    });

    test("should handle multiple collections with underscore fields", () => {
      collection("collection1", {
        field1: fields.shortString(),
        _private1: fields.number(),
      });

      collection("collection2", {
        _private2: fields.shortString(),
        field2: fields.checkbox(),
      });

      global("globalSettings", {
        setting1: fields.richText(),
        _privateSetting: fields.shortString(),
      });

      const exported = exportSchema();
      
      expect(exported.collections).toHaveLength(3);
      
      const col1 = exported.collections.find(c => c.name === "collection1");
      expect(col1?.schema).toHaveProperty("field1");
      expect(col1?.schema).not.toHaveProperty("_private1");
      
      const col2 = exported.collections.find(c => c.name === "collection2");
      expect(col2?.schema).toHaveProperty("field2");
      expect(col2?.schema).not.toHaveProperty("_private2");
      
      const globalCol = exported.collections.find(c => c.name === "globalSettings");
      expect(globalCol?.schema).toHaveProperty("setting1");
      expect(globalCol?.schema).not.toHaveProperty("_privateSetting");
    });

    test("should include generatedAt timestamp in export", () => {
      collection("test", {
        field: fields.shortString(),
      });

      const exported = exportSchema();
      
      expect(exported).toHaveProperty("generatedAt");
      expect(new Date(exported.generatedAt)).toBeInstanceOf(Date);
      expect(exported.generatedAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/);
    });

    test("should handle empty registry export", () => {
      const exported = exportSchema();
      
      expect(exported.collections).toEqual([]);
      expect(exported).toHaveProperty("generatedAt");
    });
  });

  describe("edge cases", () => {
    test("should handle collection with no fields", () => {
      const emptyCollection = collection("empty", {});
      
      const exported = exportSchema();
      const serialized = exported.collections.find(c => c.name === "empty");
      
      expect(serialized).toBeDefined();
      expect(serialized?.schema).toEqual({});
    });

    test("should handle special characters in collection names", () => {
      const specialCollection = collection("special-name_123", {
        field: fields.shortString(),
      });
      
      const retrieved = getCollection("special-name_123");
      expect(retrieved).toBeDefined();
      expect(retrieved?.name).toBe("special-name_123");
    });

    test("should maintain order of registration", () => {
      collection("first", { field: fields.shortString() });
      collection("second", { field: fields.number() });
      collection("third", { field: fields.checkbox() });
      
      const all = getAllCollections();
      expect(all[0].name).toBe("first");
      expect(all[1].name).toBe("second");
      expect(all[2].name).toBe("third");
    });
  });
});