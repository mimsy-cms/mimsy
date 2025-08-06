import { Collection, Schema } from "./collection";
import { Field } from "./fields";

const collectionRegistry = new Map<string, Collection<any>>();

export function registerCollection<T extends Schema>(
  collection: Collection<T>,
): void {
  if (collectionRegistry.has(collection.name)) {
    console.warn(
      `[Mimsy SDK] Warning: A collection with the name "${collection.name}" is already registered. It will be overwritten.`,
    );
  }
  collectionRegistry.set(collection.name, collection);
}

export function getAllCollections(): Collection<any>[] {
  return Array.from(collectionRegistry.values());
}

export function clearRegistry(): void {
  collectionRegistry.clear();
}

export function getCollection(name: string): Collection<any> | undefined {
  return collectionRegistry.get(name);
}

// --- Serialization Logic ---

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
  schema: SerializedSchema;
};

function serializeField(field: Field<any>): SerializedField {
  const result: SerializedField = {
    type: field.type,
  };

  // Clone field to avoid modifying the original object
  const options = { ...field };

  // Remove internal properties that shouldn't be in the final schema
  delete (options as any)._marker;
  delete (options as any).type;

  if ((options as any).relatesTo) {
    // `relatesTo` can be a Collection or a BuiltIn, both have a `name` property.
    result.relatesTo = (options as any).relatesTo.name;
    delete (options as any).relatesTo;
  }

  if (Object.keys(options).length > 0) {
    result.options = options;
  }

  return result;
}

function serializeSchema(schema: Schema): SerializedSchema {
  if (typeof schema !== "object" || schema === null || "_marker" in schema) {
    // This case should ideally not happen for a collection's schema based on usage.
    // If it does, we return an empty schema.
    return {};
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
