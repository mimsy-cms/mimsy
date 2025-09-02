import { describe, expect, test, vi } from "vitest";
import {
  MimsyClient,
  MimsyUserClient,
  MimsyMediaClient,
  MimsyCollectionClient,
} from "$src/api";
import { collection } from "$src/collection";
import { fields } from "$src/index";
import { User, Media } from "$src/builtins";

describe("MimsyClient", () => {
  const baseUrl = "https://api.example.com";
  let client: MimsyClient;

  beforeEach(() => {
    client = new MimsyClient(baseUrl);
  });

  describe("fetch", () => {
    test("should append baseUrl to the request path", async () => {
      global.fetch = vi.fn().mockResolvedValue(new Response("{}"));

      await client.fetch("/v1/test");

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/test`);
    });
  });

  describe("with", () => {
    test("should return a MimsyCollectionClient instance", () => {
      const testCollection = collection("posts", {
        title: fields.shortString(),
        content: fields.richText(),
      });

      const collectionClient = client.with(testCollection);

      expect(collectionClient).toBeInstanceOf(MimsyCollectionClient);
    });
  });

  describe("user", () => {
    test("should return a MimsyUserClient instance", () => {
      const userClient = client.user();

      expect(userClient).toBeInstanceOf(MimsyUserClient);
    });
  });

  describe("media", () => {
    test("should return a MimsyMediaClient instance", () => {
      const mediaClient = client.media();

      expect(mediaClient).toBeInstanceOf(MimsyMediaClient);
    });
  });

  describe("fetchRelation", () => {
    test("should call fetchRelation with correct parameters", async () => {
      const mockRelation = {
        _collection: User,
        id: "123",
      };

      global.fetch = vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            id: "123",
            email: "test@example.com",
            is_admin: false,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          }),
        ),
      );

      const result = await client.fetchRelation(mockRelation);

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/users/123`);
      expect(result).toHaveProperty("id", "123");
    });
  });
});

describe("MimsyUserClient", () => {
  const baseUrl = "https://api.example.com";
  let client: MimsyClient;
  let userClient: MimsyUserClient;

  beforeEach(() => {
    client = new MimsyClient(baseUrl);
    userClient = client.user();
  });

  describe("all", () => {
    test("should fetch all users", async () => {
      const mockUsers = [
        {
          id: "1",
          email: "user1@example.com",
          is_admin: false,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: "2",
          email: "user2@example.com",
          is_admin: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
      ];

      global.fetch = vi
        .fn()
        .mockResolvedValue(new Response(JSON.stringify(mockUsers)));

      const users = await userClient.all();

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/users`);
      expect(users).toEqual(mockUsers);
    });
  });

  describe("get", () => {
    test("should fetch a specific user by id", async () => {
      const mockUser = {
        id: "123",
        email: "user@example.com",
        is_admin: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      global.fetch = vi
        .fn()
        .mockResolvedValue(new Response(JSON.stringify(mockUser)));

      const user = await userClient.get("123");

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/users/123`);
      expect(user).toEqual(mockUser);
    });
  });
});

describe("MimsyMediaClient", () => {
  const baseUrl = "https://api.example.com";
  let client: MimsyClient;
  let mediaClient: MimsyMediaClient;

  beforeEach(() => {
    client = new MimsyClient(baseUrl);
    mediaClient = client.media();
  });

  describe("all", () => {
    test("should fetch all media", async () => {
      const mockMedia = [
        {
          id: "1",
          uuid: "uuid-1",
          name: "image1.jpg",
          content_type: "image/jpeg",
          created_at: new Date().toISOString(),
          size: 1024,
          uploaded_by_id: "user1",
          url: "https://example.com/media/1",
        },
        {
          id: "2",
          uuid: "uuid-2",
          name: "image2.png",
          content_type: "image/png",
          created_at: new Date().toISOString(),
          size: 2048,
          uploaded_by_id: "user2",
          url: "https://example.com/media/2",
        },
      ];

      global.fetch = vi
        .fn()
        .mockResolvedValue(new Response(JSON.stringify(mockMedia)));

      const media = await mediaClient.all();

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/media`);
      expect(media).toEqual(mockMedia);
    });
  });

  describe("get", () => {
    test("should fetch a specific media by id", async () => {
      const mockMedia = {
        id: "456",
        uuid: "uuid-456",
        name: "document.pdf",
        content_type: "application/pdf",
        created_at: new Date().toISOString(),
        size: 5000,
        uploaded_by_id: "user1",
        url: "https://example.com/media/456",
      };

      global.fetch = vi
        .fn()
        .mockResolvedValue(new Response(JSON.stringify(mockMedia)));

      const media = await mediaClient.get("456");

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/media/456`);
      expect(media).toEqual(mockMedia);
    });
  });
});

describe("MimsyCollectionClient", () => {
  const baseUrl = "https://api.example.com";
  let client: MimsyClient;

  const testCollection = collection("posts", {
    title: fields.shortString(),
    content: fields.richText(),
    author: fields.relation({ relatesTo: User }),
    category_id: fields.number(),
  });

  beforeEach(() => {
    client = new MimsyClient(baseUrl);
  });

  describe("all", () => {
    test("should fetch all items from a collection", async () => {
      const mockPosts = [
        {
          title: "Post 1",
          content: "Content 1",
          author_id: 1,
          category_id: 10,
        },
        {
          title: "Post 2",
          content: "Content 2",
          author_id: 2,
          category_id: 20,
        },
      ];

      global.fetch = vi
        .fn()
        .mockResolvedValue(new Response(JSON.stringify(mockPosts)));

      const collectionClient = client.with(testCollection);
      const posts = await collectionClient.all();

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/collections/posts`);
      expect(posts).toHaveLength(2);
      expect(posts[0]).toHaveProperty("title", "Post 1");
      expect(posts[0]).toHaveProperty("author");
      expect(posts[0].author).toHaveProperty("id", "1");
      expect(posts[0].author).toHaveProperty("_collection", User);
    });
  });

  describe("get", () => {
    test("should fetch a specific item from a collection by id", async () => {
      const mockPost = {
        title: "Post Title",
        content: "Post Content",
        author_id: 5,
        category_id: 15,
      };

      global.fetch = vi
        .fn()
        .mockResolvedValue(new Response(JSON.stringify(mockPost)));

      const collectionClient = client.with(testCollection);
      const post = await collectionClient.get("789");

      expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/v1/collections/posts/789`);
      expect(post).toHaveProperty("title", "Post Title");
      expect(post).toHaveProperty("author");
      expect(post.author).toHaveProperty("id", "5");
      expect(post.author).toHaveProperty("_collection", User);
    });
  });
});
