import { BuiltIn, BuiltInValue, Media, MediaValue } from "./builtins";
import { type Collection, type Schema } from "./collection";

// To make declaration of fields out of this file more difficult.
const fieldType: unique symbol = Symbol();

export type CollectionOrBuiltin<T extends Schema> = T extends BuiltInValue
  ? BuiltIn<T>
  : Collection<T>;

type FieldOptions<Constraints = {}> = {
  label?: string;
  description?: string;
  constraints?: Constraints & { required?: boolean };
};

type FieldGenerator<T, Constraints = {}> = (
  options?: FieldOptions<Constraints>,
) => Field<T>;

export type UnfetchedRelation<T extends CollectionOrBuiltin<any>> = {
  _collection: T;
  id: string;
};

export type UnfetchedMultiRelation<T extends CollectionOrBuiltin<any>> =
  UnfetchedRelation<T>[];

export type Field<T> = {
  _marker: typeof fieldType;
  type: string;
};

export const shortString: FieldGenerator<
  string,
  {
    minLength?: number;
    maxLength?: number;
  }
> = (options) => ({
  _marker: fieldType,
  type: "string",
  ...options,
});

export const richText: FieldGenerator<string> = (options?: FieldOptions) => ({
  _marker: fieldType,
  type: "rich_text",
  ...options,
});

export const checkbox: FieldGenerator<boolean> = (options?: FieldOptions) => ({
  _marker: fieldType,
  type: "checkbox",
  ...options,
});

export const dateTime: FieldGenerator<Date> = (options?: FieldOptions) => ({
  _marker: fieldType,
  type: "date_time",
  ...options,
});

export const number: FieldGenerator<
  number,
  {
    min?: number;
    max?: number;
  }
> = (options) => ({
  _marker: fieldType,
  type: "number",
  ...options,
});

export const email: FieldGenerator<string> = (options?: FieldOptions) => ({
  _marker: fieldType,
  type: "email",
  ...options,
});

export function relation<R extends CollectionOrBuiltin<any>>(
  options: {
    relatesTo: R;
  } & FieldOptions,
): Field<UnfetchedRelation<R>> {
  return {
    _marker: fieldType,
    type: "relation",
    ...options,
  };
}

export function multiRelation<T extends CollectionOrBuiltin<any>>(
  options?: {
    relatesTo: T;
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
