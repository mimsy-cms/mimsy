import { describe, expect, test, beforeEach } from "vitest";
import { collection, global } from "$src/collection";
import { fields } from "$src/index";
import { User, Media } from "$src/builtins";
import { clearRegistry, exportSchema } from "$src/registry";

describe("Serialization", () => {
  beforeEach(() => {
    clearRegistry();
  });

  describe("underscore field exclusion", () => {
    test("should exclude all underscore-prefixed fields from serialization", () => {
      const testCollection = collection("test", {
        publicField1: fields.shortString(),
        _privateField1: fields.number(),
        publicField2: fields.checkbox(),
        _privateField2: fields.richText(),
        _anotherPrivate: fields.dateTime(),
        lastPublicField: fields.number(),
      });

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(Object.keys(serialized.schema)).toEqual([
        "publicField1",
        "publicField2",
        "lastPublicField",
      ]);

      expect(serialized.schema).not.toHaveProperty("_privateField1");
      expect(serialized.schema).not.toHaveProperty("_privateField2");
      expect(serialized.schema).not.toHaveProperty("_anotherPrivate");
    });

    test("should handle fields with underscore in the middle or end", () => {
      const testCollection = collection("underscorePositions", {
        _startUnderscore: fields.shortString(),
        middle_underscore: fields.number(),
        endUnderscore_: fields.checkbox(),
        __doubleStart: fields.richText(),
        normal: fields.shortString(),
      });

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(serialized.schema).not.toHaveProperty("_startUnderscore");
      expect(serialized.schema).not.toHaveProperty("__doubleStart");
      expect(serialized.schema).toHaveProperty("middle_underscore");
      expect(serialized.schema).toHaveProperty("endUnderscore_");
      expect(serialized.schema).toHaveProperty("normal");
    });

    test("should exclude underscore fields with complex field types", () => {
      const categoryCollection = collection("category", {
        name: fields.shortString(),
      });

      const testCollection = collection("complexFields", {
        title: fields.shortString(),
        _privateRelation: fields.relation({ relatesTo: User }),
        publicRelation: fields.relation({ relatesTo: Media }),
        _privateMultiRelation: fields.multiRelation({
          relatesTo: categoryCollection,
        }),
        content: fields.richText(),
        _privateOptions: fields.shortString({
          constraints: { required: true, minLength: 1 },
        }),
      });

      const exported = exportSchema();
      const serialized = exported.collections.find(
        (c) => c.name === "complexFields",
      );

      expect(serialized?.schema).toHaveProperty("title");
      expect(serialized?.schema).toHaveProperty("publicRelation");
      expect(serialized?.schema).toHaveProperty("content");

      expect(serialized?.schema).not.toHaveProperty("_privateRelation");
      expect(serialized?.schema).not.toHaveProperty("_privateMultiRelation");
      expect(serialized?.schema).not.toHaveProperty("_privateOptions");
    });

    test("should handle collections with only underscore fields", () => {
      const privateOnlyCollection = collection("privateOnly", {
        _private1: fields.shortString(),
        _private2: fields.number(),
        _private3: fields.checkbox(),
        _private4: fields.richText(),
      });

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(serialized.name).toBe("privateOnly");
      expect(serialized.schema).toEqual({});
      expect(Object.keys(serialized.schema)).toHaveLength(0);
    });

    test("should handle empty collections", () => {
      const emptyCollection = collection("empty", {});

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(serialized.name).toBe("empty");
      expect(serialized.schema).toEqual({});
    });

    test("should exclude underscore fields in global collections", () => {
      const testGlobal = global("settings", {
        publicSetting: fields.shortString(),
        _privateSetting: fields.number(),
        anotherPublic: fields.checkbox(),
        _anotherPrivate: fields.richText(),
      });

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(serialized.isGlobal).toBe(true);
      expect(serialized.schema).toHaveProperty("publicSetting");
      expect(serialized.schema).toHaveProperty("anotherPublic");
      expect(serialized.schema).not.toHaveProperty("_privateSetting");
      expect(serialized.schema).not.toHaveProperty("_anotherPrivate");
    });
  });

  describe("field serialization correctness", () => {
    test("should correctly serialize simple field types", () => {
      const testCollection = collection("simpleTypes", {
        text: fields.shortString(),
        longText: fields.richText(),
        num: fields.number(),
        bool: fields.checkbox(),
        date: fields.dateTime(),
      });

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(serialized.schema.text).toEqual({ type: "string" });
      expect(serialized.schema.longText).toEqual({ type: "rich_text" });
      expect(serialized.schema.num).toEqual({ type: "number" });
      expect(serialized.schema.bool).toEqual({ type: "checkbox" });
      expect(serialized.schema.date).toEqual({ type: "date_time" });
    });

    test("should correctly serialize field options", () => {
      const testCollection = collection("withOptions", {
        requiredField: fields.shortString({
          constraints: { required: true, minLength: 5 },
        }),
        withOptions: fields.number({
          constraints: { min: 1, max: 100 },
        }),
        complexOptions: fields.shortString({
          label: "Complex Field",
          description: "A field with multiple options",
          constraints: {
            required: true,
            minLength: 5,
            maxLength: 100,
          },
        }),
      });

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(serialized.schema.requiredField).toEqual({
        type: "string",
        options: {
          constraints: { required: true, minLength: 5 },
        },
      });

      expect(serialized.schema.withOptions).toEqual({
        type: "number",
        options: {
          constraints: { min: 1, max: 100 },
        },
      });

      expect(serialized.schema.complexOptions).toEqual({
        type: "string",
        options: {
          label: "Complex Field",
          description: "A field with multiple options",
          constraints: {
            required: true,
            minLength: 5,
            maxLength: 100,
          },
        },
      });
    });

    test("should correctly serialize relation fields", () => {
      const categoryCollection = collection("categories", {
        name: fields.shortString(),
        _hidden: fields.number(),
      });

      const postCollection = collection("posts", {
        title: fields.shortString(),
        author: fields.relation({ relatesTo: User }),
        category: fields.relation({ relatesTo: categoryCollection }),
        thumbnail: fields.relation({ relatesTo: Media }),
        _privateRelation: fields.relation({ relatesTo: User }),
      });

      const exported = exportSchema();
      const serializedPosts = exported.collections.find(
        (c) => c.name === "posts",
      );

      expect(serializedPosts?.schema.author).toEqual({
        type: "relation",
        relatesTo: "<builtins.user>",
      });

      expect(serializedPosts?.schema.category).toEqual({
        type: "relation",
        relatesTo: "categories",
      });

      expect(serializedPosts?.schema.thumbnail).toEqual({
        type: "relation",
        relatesTo: "<builtins.media>",
      });

      expect(serializedPosts?.schema).not.toHaveProperty("_privateRelation");
    });

    test("should correctly serialize multi-relation fields", () => {
      const tagCollection = collection("tags", {
        name: fields.shortString(),
      });

      const postCollection = collection("posts", {
        title: fields.shortString(),
        tags: fields.multiRelation({ relatesTo: tagCollection }),
        _privateTags: fields.multiRelation({ relatesTo: tagCollection }),
      });

      const exported = exportSchema();
      const serializedPosts = exported.collections.find(
        (c) => c.name === "posts",
      );

      expect(serializedPosts?.schema.tags).toEqual({
        type: "multi_relation",
        relatesTo: "tags",
      });

      expect(serializedPosts?.schema).not.toHaveProperty("_privateTags");
    });
  });

  describe("multiple collections serialization", () => {
    test("should correctly serialize multiple collections with underscore fields", () => {
      const col1 = collection("collection1", {
        field1: fields.shortString(),
        _private1: fields.number(),
        field2: fields.checkbox(),
      });

      const col2 = collection("collection2", {
        _private2: fields.shortString(),
        field3: fields.richText(),
        _private3: fields.number(),
      });

      const col3 = global("globalSettings", {
        setting1: fields.shortString(),
        _privateSetting: fields.checkbox(),
        setting2: fields.number(),
      });

      const exported = exportSchema();

      expect(exported.collections).toHaveLength(3);

      const serialized1 = exported.collections.find(
        (c) => c.name === "collection1",
      );
      expect(Object.keys(serialized1!.schema)).toEqual(["field1", "field2"]);

      const serialized2 = exported.collections.find(
        (c) => c.name === "collection2",
      );
      expect(Object.keys(serialized2!.schema)).toEqual(["field3"]);

      const serialized3 = exported.collections.find(
        (c) => c.name === "globalSettings",
      );
      expect(Object.keys(serialized3!.schema)).toEqual([
        "setting1",
        "setting2",
      ]);
      expect(serialized3!.isGlobal).toBe(true);
    });

    test("should handle mixed collections with and without underscore fields", () => {
      const noUnderscoreCol = collection("noUnderscore", {
        field1: fields.shortString(),
        field2: fields.number(),
      });

      const withUnderscoreCol = collection("withUnderscore", {
        publicField: fields.shortString(),
        _privateField: fields.number(),
      });

      const onlyUnderscoreCol = collection("onlyUnderscore", {
        _private1: fields.shortString(),
        _private2: fields.number(),
      });

      const exported = exportSchema();

      const noUnderscoreSerialized = exported.collections.find(
        (c) => c.name === "noUnderscore",
      );
      expect(Object.keys(noUnderscoreSerialized!.schema)).toEqual([
        "field1",
        "field2",
      ]);

      const withUnderscoreSerialized = exported.collections.find(
        (c) => c.name === "withUnderscore",
      );
      expect(Object.keys(withUnderscoreSerialized!.schema)).toEqual([
        "publicField",
      ]);

      const onlyUnderscoreSerialized = exported.collections.find(
        (c) => c.name === "onlyUnderscore",
      );
      expect(Object.keys(onlyUnderscoreSerialized!.schema)).toEqual([]);
    });
  });

  describe("export metadata", () => {
    test("should include correct metadata in export", () => {
      collection("test", {
        field: fields.shortString(),
        _private: fields.number(),
      });

      const exported = exportSchema();

      expect(exported).toHaveProperty("collections");
      expect(exported).toHaveProperty("generatedAt");
      expect(Array.isArray(exported.collections)).toBe(true);
      expect(typeof exported.generatedAt).toBe("string");

      const date = new Date(exported.generatedAt);
      expect(date.toString()).not.toBe("Invalid Date");
    });

    test("should maintain collection metadata during serialization", () => {
      const testCol = collection("test", {
        field: fields.shortString(),
        _private: fields.number(),
      });

      const testGlobal = global("settings", {
        setting: fields.shortString(),
        _privateSetting: fields.number(),
      });

      const exported = exportSchema();

      const serializedCol = exported.collections.find((c) => c.name === "test");
      expect(serializedCol?.name).toBe("test");
      expect(serializedCol?.isGlobal).toBe(false);

      const serializedGlobal = exported.collections.find(
        (c) => c.name === "settings",
      );
      expect(serializedGlobal?.name).toBe("settings");
      expect(serializedGlobal?.isGlobal).toBe(true);
    });
  });

  describe("edge cases", () => {
    test("should handle very long underscore prefixes", () => {
      const testCollection = collection("longUnderscore", {
        normal: fields.shortString(),
        ___tripleUnderscore: fields.number(),
        ____quadUnderscore: fields.checkbox(),
        _____fiveUnderscore: fields.richText(),
      });

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(Object.keys(serialized.schema)).toEqual(["normal"]);
      expect(serialized.schema).not.toHaveProperty("___tripleUnderscore");
      expect(serialized.schema).not.toHaveProperty("____quadUnderscore");
      expect(serialized.schema).not.toHaveProperty("_____fiveUnderscore");
    });

    test("should handle special characters with underscore", () => {
      const testCollection = collection("special", {
        _$private: fields.shortString(),
        $public: fields.number(),
        _123numeric: fields.checkbox(),
        "456numeric": fields.richText(),
        "_with-dash": fields.shortString(),
        "with-dash": fields.number(),
      });

      const exported = exportSchema();
      const serialized = exported.collections[0];

      expect(serialized.schema).not.toHaveProperty("_$private");
      expect(serialized.schema).not.toHaveProperty("_123numeric");
      expect(serialized.schema).not.toHaveProperty("_with-dash");

      expect(serialized.schema).toHaveProperty("$public");
      expect(serialized.schema).toHaveProperty("456numeric");
      expect(serialized.schema).toHaveProperty("with-dash");
    });

    test("should handle null/undefined gracefully", () => {
      const testCollection = collection("nullish", {
        normalField: fields.shortString(),
        _privateField: fields.number(),
      });

      const exported = exportSchema();
      expect(exported.collections).toHaveLength(1);

      const serialized = exported.collections[0];
      expect(serialized.schema.normalField).toBeDefined();
      expect(serialized.schema._privateField).toBeUndefined();
    });
  });
});
