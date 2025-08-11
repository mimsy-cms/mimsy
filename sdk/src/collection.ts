import { BuiltInValue } from "./builtins";
import { Field } from "./fields";
import { registerCollection, registerGlobal } from "./registry";

export type Schema =
  | BuiltInValue
  | {
      [key: string]: Field<any>;
    };

export type Collection<T extends Schema> = {
  name: string;
  schema: T;
  isGlobal: false;
};

export type Global<T extends Schema> = {
  name: string;
  schema: T;
  isGlobal: true;
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
    isGlobal: false,
  } satisfies Collection<T>;
  registerCollection(coll);
  return coll;
}

export function global<T extends Schema>(
  name: string,
  schema: T,
): Global<T> {
  const coll = {
    name,
    schema,
    isGlobal: true,
  } satisfies Global<T>;
  registerGlobal(coll);
  return coll;
}
