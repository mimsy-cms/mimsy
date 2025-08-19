import { Media, MediaValue } from "./builtins";
import { type Collection, type Schema } from "./collection";

// To make declaration of fields out of this file more difficult.
const fieldType: unique symbol = Symbol();

type FieldOptions = any;

type FieldGenerator<T> = (options?: FieldOptions) => Field<T>;

type UnfetchedRelation<T extends Schema> = {
  id: string;
  name: string;
};

type UnfetchedMultiRelation<T extends Schema> = UnfetchedRelation<T>[];

export type Field<T> = {
  _marker: typeof fieldType;
  type: string;
};

export const shortString: FieldGenerator<string> = (
  options?: FieldOptions,
) => ({
  _marker: fieldType,
  type: "string",
  ...options,
});

export const richText: FieldGenerator<string> = (options?: FieldOptions) => ({
  _marker: fieldType,
  type: "rich_text",
  ...options,
});

export function relation<T extends Schema>(
  options?: {
    relatesTo: Collection<T>;
  } & FieldOptions,
): Field<UnfetchedRelation<T>> {
  return {
    _marker: fieldType,
    type: "relation",
    ...options,
  };
}

export function multiRelation<T extends Schema>(
  options?: {
    relatesTo: Collection<T>;
  } & FieldOptions,
): Field<UnfetchedMultiRelation<T>> {
  return {
    _marker: fieldType,
    type: "multi_relation",
    ...options,
  };
}

export function media(
  options?: FieldOptions,
): Field<UnfetchedRelation<MediaValue>> {
  return relation({
    relatesTo: Media,
    ...options,
  });
}
