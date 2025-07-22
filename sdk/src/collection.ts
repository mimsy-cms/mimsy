import { BuiltInValue } from "./builtins";
import { Field } from "./fields";
import { registerCollection } from "./registry";

export type Schema =
  | BuiltInValue
  | {
      [key: string]: Field<any>;
    };

export type Collection<T extends Schema> = {
  name: string;
  schema: T;
};

type ObjectOf<S extends Schema> = S extends BuiltInValue
  ? BuiltInValue
  : {
      [key in keyof S]: S[key] extends Field<infer U> ? U : never;
    };

export function toObject<S extends Schema>(
  collection: Collection<S>,
): ObjectOf<S> {
  // TODO: Implement toObject function
  return {} as ObjectOf<S>;
}

export function collection<T extends Schema>(
  name: string,
  schema: T,
): Collection<T> {
  const coll = {
    name,
    schema,
  };
  return coll;
}
