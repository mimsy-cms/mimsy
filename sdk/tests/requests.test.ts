import { describe, expect, test, vi } from "vitest";
import { fetchRelation } from "$src/requests";
import { MimsyClient } from "$src/api";
import { User, Media } from "$src/builtins";
import { collection } from "$src/collection";
import { fields } from "$src/index";
import { UnfetchedRelation } from "$src/fields";

describe("fetchRelation", () => {
  const baseUrl = "https://api.example.com";
  let client: MimsyClient;

  beforeEach(() => {
    client = new MimsyClient(baseUrl);
    vi.clearAllMocks();
  });

  describe("builtin types", () => {
    test("should fetch a User relation", async () => {
      const mockUser = {
        id: "user-123",
        email: "test@example.com",
        is_admin: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      global.fetch = vi.fn().mockResolvedValue(
        new Response(JSON.stringify(mockUser))
      );

      const unfetchedUser: UnfetchedRelation<typeof User> = {
        _collection: User,
        id: "user-123",
      };

      const result = await fetchRelation(client, unfetchedUser);

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/users/user-123`);
      expect(result).toEqual(mockUser);
    });

    test("should fetch a Media relation", async () => {
      const mockMedia = {
        id: "media-456",
        uuid: "uuid-456",
        name: "image.jpg",
        content_type: "image/jpeg",
        created_at: new Date().toISOString(),
        size: 2048,
        uploaded_by_id: "user-1",
        url: "https://example.com/media/456",
      };

      global.fetch = vi.fn().mockResolvedValue(
        new Response(JSON.stringify(mockMedia))
      );

      const unfetchedMedia: UnfetchedRelation<typeof Media> = {
        _collection: Media,
        id: "media-456",
      };

      const result = await fetchRelation(client, unfetchedMedia);

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/media/media-456`);
      expect(result).toEqual(mockMedia);
    });

    test("should throw error for unknown builtin type", async () => {
      // Create a fake builtin with the same marker as real builtins
      const unknownBuiltin = {
        _marker: User._marker, // Use the same marker as real builtins
        name: "<builtins.unknown>",
      };

      const unfetchedUnknown = {
        _collection: unknownBuiltin,
        id: "123",
      } as any;

      await expect(fetchRelation(client, unfetchedUnknown)).rejects.toThrow(
        "Unknown builtin type: <builtins.unknown>"
      );
    });
  });

  describe("custom collection types", () => {
    test("should fetch a custom collection relation", async () => {
      const postCollection = collection("posts", {
        title: fields.shortString(),
        content: fields.richText(),
      });

      const mockPost = {
        title: "Test Post",
        content: "This is a test post content",
      };

      global.fetch = vi.fn().mockResolvedValue(
        new Response(JSON.stringify(mockPost))
      );

      const unfetchedPost = {
        _collection: postCollection,
        id: "post-789",
      };

      const result = await fetchRelation(client, unfetchedPost);

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/collections/posts/post-789`);
      expect(result).toHaveProperty("title", "Test Post");
      expect(result).toHaveProperty("content", "This is a test post content");
    });

    test("should handle nested relations in custom collections", async () => {
      const categoryCollection = collection("categories", {
        name: fields.shortString(),
        description: fields.richText(),
      });

      const postWithCategoryCollection = collection("posts_with_category", {
        title: fields.shortString(),
        content: fields.richText(),
        category: fields.relation({ relatesTo: categoryCollection }),
      });

      const mockPostData = {
        title: "Post with Category",
        content: "Content here",
        category_id: 5,
      };

      global.fetch = vi.fn().mockResolvedValue(
        new Response(JSON.stringify(mockPostData))
      );

      const unfetchedPost = {
        _collection: postWithCategoryCollection,
        id: "post-999",
      };

      const result = await fetchRelation(client, unfetchedPost);

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/collections/posts_with_category/post-999`);
      expect(result).toHaveProperty("title", "Post with Category");
      expect(result).toHaveProperty("category");
      expect((result as any).category).toHaveProperty("id", "5");
      expect((result as any).category).toHaveProperty("_collection", categoryCollection);
    });
  });

  describe("error handling", () => {
    test("should handle fetch errors", async () => {
      global.fetch = vi.fn().mockRejectedValue(new Error("Network error"));

      const unfetchedUser: UnfetchedRelation<typeof User> = {
        _collection: User,
        id: "user-error",
      };

      await expect(fetchRelation(client, unfetchedUser)).rejects.toThrow("Network error");
    });

    test("should handle invalid JSON response", async () => {
      global.fetch = vi.fn().mockResolvedValue(
        new Response("Invalid JSON")
      );

      const unfetchedUser: UnfetchedRelation<typeof User> = {
        _collection: User,
        id: "user-invalid",
      };

      await expect(fetchRelation(client, unfetchedUser)).rejects.toThrow();
    });
  });
});