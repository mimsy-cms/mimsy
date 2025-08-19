import { expectTypeOf } from "vitest";
import { collection, Collection, toObject } from "$src/collection";
import * as fields from "$src/fields";
import * as builtins from "$src/builtins";

/**
 * This file contains tests the typescript types of the SDK.
 * Typecheck is enabled in the vitest config, so these tests will fail if the types are incorrect.
 */

test("simple collection", () => {
  const testCollection = collection("posts", {
    title: fields.shortString({
      description: "The title of the post",
    }),
  });

  expectTypeOf(testCollection).toMatchTypeOf<Collection<any>>();

  let test = toObject(testCollection);
  expectTypeOf(test).toEqualTypeOf({
    title: "",
  });
});

test("collections with builtins", () => {
  const testCollection = collection("posts", {
    title: fields.shortString({
      description: "The title of the post",
    }),
    coverImage: fields.media({
      description: "Test",
    }),
  });

  expectTypeOf(testCollection).toMatchTypeOf<Collection<any>>();

  let test = toObject(testCollection);
  expectTypeOf(test).toEqualTypeOf({
    title: "",
    coverImage: { id: "", name: "" },
  });
});

test("collection multiRelation", () => {
  const tags = collection("tags", {
    name: fields.shortString({
      description: "The name of the tag",
    }),
  });

  const testCollection = collection("posts", {
    title: fields.shortString({
      description: "The title of the post",
    }),
    tags: fields.multiRelation({
      relatesTo: tags,
    }),
  });

  expectTypeOf(testCollection).toMatchTypeOf<Collection<any>>();

  let test = toObject(testCollection);
  expectTypeOf(test).toEqualTypeOf({
    title: "",
    tags: [{ id: "", name: "" }],
  });
});

test("collection with user", () => {
  const tags = collection("tags", {
    name: fields.shortString({
      description: "The name of the tag",
    }),
  });

  const testCollection = collection("posts", {
    title: fields.shortString({
      description: "The title of the post",
    }),
    author: fields.relation({
      relatesTo: builtins.User,
    }),
  });

  expectTypeOf(testCollection).toMatchTypeOf<Collection<any>>();

  let test = toObject(testCollection);
  expectTypeOf(test).toEqualTypeOf({
    title: "",
    author: { id: "", name: "" },
  });
});
