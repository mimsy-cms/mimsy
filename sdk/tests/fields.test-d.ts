import { User, UserValue } from "$src/builtins";
import { collection, FieldValue, ObjectOf } from "$src/collection";
import { fields } from "$src/index";

describe("ObjectOf type", () => {
  test("should correctly type simple collections", () => {
    const simpleCollection = collection("simple", {
      title: fields.shortString(),
      count: fields.number(),
    });

    type SimpleObject = ObjectOf<typeof simpleCollection.schema>;

    expectTypeOf<SimpleObject>().toEqualTypeOf<{
      title: string;
      count: number;
    }>();
  });

  test("should correctly type BuiltInValue", () => {
    type UserObject = ObjectOf<typeof User>;

    expectTypeOf<UserObject>().toEqualTypeOf<UserValue>();
  });
});

describe("FieldValue type", () => {
  test("should extract correct value types from fields", () => {
    type StringValue = FieldValue<ReturnType<typeof fields.shortString>>;
    type LongStringValue = FieldValue<ReturnType<typeof fields.longString>>;
    type NumberValue = FieldValue<ReturnType<typeof fields.number>>;
    type BooleanValue = FieldValue<ReturnType<typeof fields.checkbox>>;

    expectTypeOf<StringValue>().toEqualTypeOf<string>();
    expectTypeOf<LongStringValue>().toEqualTypeOf<string>();
    expectTypeOf<NumberValue>().toEqualTypeOf<number>();
    expectTypeOf<BooleanValue>().toEqualTypeOf<boolean>();
  });
});
