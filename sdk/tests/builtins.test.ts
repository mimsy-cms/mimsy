import { describe, expect, test } from "vitest";
import { User, Media, isBuiltIn } from "$src/builtins";

describe("isBuiltIn", () => {
  test("should return true for User builtin", () => {
    expect(isBuiltIn(User)).toBe(true);
  });

  test("should return true for Media builtin", () => {
    expect(isBuiltIn(Media)).toBe(true);
  });

  test("should return false for non-builtin objects", () => {
    const notBuiltin = {
      name: "not a builtin",
    };
    expect(isBuiltIn(notBuiltin)).toBe(false);
  });

  test("should return false for objects with wrong marker", () => {
    const wrongMarker = {
      _marker: Symbol("wrong"),
      name: "<builtins.fake>",
    };
    expect(isBuiltIn(wrongMarker)).toBe(false);
  });

  test("should return false for null", () => {
    expect(isBuiltIn(null)).toBe(false);
  });

  test("should return false for undefined", () => {
    expect(isBuiltIn(undefined)).toBe(false);
  });

  test("should return false for primitives", () => {
    expect(isBuiltIn("string")).toBe(false);
    expect(isBuiltIn(123)).toBe(false);
    expect(isBuiltIn(true)).toBe(false);
  });

  test("should return false for arrays", () => {
    expect(isBuiltIn([])).toBe(false);
    expect(isBuiltIn([User])).toBe(false);
  });
});
