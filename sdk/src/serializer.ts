import { Field } from "./fields";
import { Collection, Schema } from "./collection";
import { BuiltInValue } from "./builtins";
import { getAllCollections } from "./registry";

type SerializedField = {
  type: string;
  options?: any;
  relatesTo?: string; // collection name for relations
};

type SerializedSchema = {
  [key: string]: SerializedField;
};

type SerializedCollection = {
  name: string;
  schema: SerializedSchema | { type: "builtin"; name: string };
};

function serializeField(field: Field<any>): SerializedField {
  const result: SerializedField = {
    type: field.type,
  };

  // Create a copy of the field to avoid modifying the original object
  const options = { ...field } as any;

  // Remove internal properties that shouldn't be in the final schema
  delete options._marker;
  delete options.type;

  // Handle relation fields specifically
  if (options.relatesTo) {
    result.relatesTo = options.relatesTo.name;
    delete options.relatesTo;
  }

  // Add remaining properties as options
  if (Object.keys(options).length > 0) {
    result.options = options;
  }

  return result;
}

function serializeSchema(
  schema: Schema,
): SerializedSchema | { type: "builtin"; name: string } {
  const asBuiltIn = schema as BuiltInValue;
  if (asBuiltIn._marker) {
    // This is a built-in type. We need a way to get its name.
    // Assuming the built-in object has a 'name' property.
    return { type: "builtin", name: (schema as any).name };
  }

  const serialized: SerializedSchema = {};
  for (const [key, field] of Object.entries(schema)) {
    serialized[key] = serializeField(field as Field<any>);
  }
  return serialized;
}

export function exportSchema(): {
  collections: SerializedCollection[];
  generatedAt: string;
} {
  const collections = getAllCollections();

  return {
    collections: collections.map((coll) => ({
      name: coll.name,
      schema: serializeSchema(coll.schema),
    })),
    generatedAt: new Date().toISOString(),
  };
}
