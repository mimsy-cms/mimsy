import { collection, global } from "$src/collection";
import { fields } from "$src/index";
import { test } from "vitest";

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
