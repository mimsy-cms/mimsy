import { describe, expect, test } from "vitest";
import {
  collection,
  global,
  postProcess,
  ObjectOf,
  FieldValue,
} from "$src/collection";
import { fields } from "$src/index";
import { User, Media } from "$src/builtins";
import type { UnfetchedRelation } from "$src/fields";

describe("collection and global reserved names", () => {
  test.each([
    "user",
    "media",
    "collection",
    "cron_locks",
    "session",
    "sync_status",
  ])("reserved collection name: %s", (name) => {
    expect(() => {
      collection(name, {
        name: fields.shortString(),
      });
    }).toThrowError(`Collection name '${name}' is reserved`);
  });

  test.each([
    "user",
    "media",
    "collection",
    "cron_locks",
    "session",
    "sync_status",
  ])("reserved global name: %s", (name) => {
    expect(() => {
      global(name, {
        name: fields.shortString(),
      });
    }).toThrowError(`Global name '${name}' is reserved`);
  });
});

describe("postProcess", () => {
  describe("simple fields", () => {
    const simpleCollection = collection("simple", {
      title: fields.shortString(),
      count: fields.number(),
      published: fields.checkbox(),
    });

    test("should process simple field types correctly", () => {
      const rawData = {
        title: "Test Title",
        count: 42,
        published: true,
      };

      const result = postProcess(simpleCollection, rawData);

      expect(result).toEqual({
        title: "Test Title",
        count: 42,
        published: true,
      });
    });

    test("should handle undefined fields", () => {
      const rawData = {
        title: "Test Title",
      };

      const result = postProcess(simpleCollection, rawData);

      expect(result).toEqual({
        title: "Test Title",
        count: undefined,
        published: undefined,
      });
    });
  });

  describe("relation fields", () => {
    const postCollection = collection("posts", {
      title: fields.shortString(),
      author: fields.relation({
        relatesTo: User,
      }),
      featuredImage: fields.relation({
        relatesTo: Media,
      }),
    });

    test("should convert relation fields to UnfetchedRelation", () => {
      const rawData = {
        title: "Post Title",
        author_id: 123,
        featuredImage_id: 456,
      };

      const result = postProcess(postCollection, rawData);

      expect(result.title).toBe("Post Title");

      expect(result.author).toEqual({
        _collection: User,
        id: "123",
      } satisfies UnfetchedRelation<typeof User>);

      expect(result.featuredImage).toEqual({
        _collection: Media,
        id: "456",
      } satisfies UnfetchedRelation<typeof Media>);
    });

    test("should handle missing relation IDs", () => {
      const rawData = {
        title: "Post Title",
      };

      const result = postProcess(postCollection, rawData);

      expect(result.title).toBe("Post Title");
      expect(result.author).toEqual({
        _collection: User,
        id: "undefined",
      });
      expect(result.featuredImage).toEqual({
        _collection: Media,
        id: "undefined",
      });
    });
  });

  describe("custom collection relations", () => {
    const categoryCollection = collection("categories", {
      name: fields.shortString(),
      slug: fields.shortString(),
    });

    const postWithCategoryCollection = collection("posts_with_category", {
      title: fields.shortString(),
      category: fields.relation({ relatesTo: categoryCollection }),
    });

    test("should handle custom collection relations", () => {
      const rawData = {
        title: "Post with Category",
        category_id: 789,
      };

      const result = postProcess(postWithCategoryCollection, rawData);

      expect(result.title).toBe("Post with Category");
      expect(result.category).toEqual({
        _collection: categoryCollection,
        id: "789",
      });
    });
  });

  describe("multi_relation fields", () => {
    const tagCollection = collection("tags", {
      name: fields.shortString(),
    });

    const postWithTagsCollection = collection("posts_with_tags", {
      title: fields.shortString(),
      tags: fields.multiRelation({ relatesTo: tagCollection }),
    });

    test("should mark multi_relation as unsupported", () => {
      const rawData = {
        title: "Post with Tags",
        tags: [1, 2, 3],
      };

      const result = postProcess(postWithTagsCollection, rawData);

      expect(result.title).toBe("Post with Tags");
      expect(result.tags).toBe("Unsupported type");
    });
  });

  describe("mixed field types", () => {
    const complexCollection = collection("complex", {
      title: fields.shortString(),
      description: fields.richText(),
      views: fields.number(),
      rating: fields.number(),
      isPublished: fields.checkbox(),
      publishedAt: fields.dateTime(),
      author: fields.relation({
        relatesTo: User,
      }),
      coverImage: fields.relation({
        relatesTo: Media,
      }),
    });

    test("should handle complex collections with mixed field types", () => {
      const rawData = {
        title: "Complex Post",
        description: "This is a longer description",
        views: 1000,
        rating: 4.5,
        isPublished: true,
        publishedAt: "2024-01-15T10:30:00Z",
        author_id: 555,
        coverImage_id: 999,
      };

      const result = postProcess(complexCollection, rawData);

      expect(result).toEqual({
        title: "Complex Post",
        description: "This is a longer description",
        views: 1000,
        rating: 4.5,
        isPublished: true,
        publishedAt: "2024-01-15T10:30:00Z",
        author: {
          _collection: User,
          id: "555",
        },
        coverImage: {
          _collection: Media,
          id: "999",
        },
      });
    });
  });

  describe("edge cases", () => {
    const testCollection = collection("test", {
      field1: fields.shortString(),
      field2: fields.number(),
    });

    test("should handle null data", () => {
      const result = postProcess(testCollection, null);

      expect(result).toEqual({
        field1: undefined,
        field2: undefined,
      });
    });

    test("should handle non-object data", () => {
      const result = postProcess(testCollection, "not an object");

      expect(result).toEqual({
        field1: undefined,
        field2: undefined,
      });
    });

    test("should handle empty object", () => {
      const result = postProcess(testCollection, {});

      expect(result).toEqual({
        field1: undefined,
        field2: undefined,
      });
    });

    test("should ignore extra fields in raw data", () => {
      const rawData = {
        field1: "value1",
        field2: 123,
        extraField: "should be ignored",
        anotherExtra: true,
      };

      const result = postProcess(testCollection, rawData);

      expect(result).toEqual({
        field1: "value1",
        field2: 123,
      });
      expect(result).not.toHaveProperty("extraField");
      expect(result).not.toHaveProperty("anotherExtra");
    });
  });
});
