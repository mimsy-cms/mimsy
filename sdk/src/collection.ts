import { BuiltIn, BuiltInValue } from "./builtins";
import { CollectionOrBuiltin, Field, UnfetchedRelation } from "./fields";
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

export type FieldValue<T> = T extends Field<infer U> ? U : never;

export type ObjectOfCollection<T extends CollectionOrBuiltin<any>> =
  T extends BuiltIn<infer V>
    ? V
    : T extends Collection<infer S>
      ? ObjectOf<S>
      : never;

export type ObjectOf<S extends Schema> =
  S extends BuiltIn<infer V>
    ? V
    : {
        [key in keyof S]: FieldValue<S[key]>;
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

export function global<T extends Schema>(name: string, schema: T): Global<T> {
  const coll = {
    name,
    schema,
    isGlobal: true,
  } satisfies Global<T>;
  registerGlobal(coll);
  return coll;
}

export function postProcess<T extends Exclude<Schema, BuiltInValue>>(
  collection: Collection<T>,
  data: unknown,
): ObjectOf<T> {
  // Iterate for fields
  const processedData = {} as ObjectOf<T>;
  for (const key in collection.schema) {
    const field = collection.schema[key];
    switch (field.type) {
      case "multi_relation":
        processedData[key] = "Unsupported type" as any;
        break;
      case "relation":
        // TODO: Get the name of the field + "_id"
        // And then put this with the schema info into a unfetched relation array
        const relationIdValue = get<number>(data, `${key}_id`);
        const relationId = relationIdValue
          ? relationIdValue.toString()
          : "undefined";
        processedData[key] = {
          _collection: (field as any).relatesTo,
          id: relationId,
        } satisfies UnfetchedRelation<any> as any;
        break;
      default:
        processedData[key] = get(data, key);
        break;
    }
  }
  return processedData;
}

function get<T>(from: unknown, key: string): T {
  if (from && typeof from === "object" && key in from) {
    return (from as any)[key] as T;
  }
  return undefined as T;
}
