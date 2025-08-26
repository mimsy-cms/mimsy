// To make declaration of fields out of this file more difficult.
const builtInType = Symbol("mimsy-builtin");

export type BuiltInValue = {
  _marker: typeof builtInType;
};

export type BuiltIn<T extends BuiltInValue> = {
  _marker: typeof builtInType;
  name: string;
};

export type UserValue = BuiltInValue & {
  id: string;
  email: string;
  is_admin: boolean;
  created_at: Date;
  updated_at: Date;
};

export type MediaValue = BuiltInValue & {
  id: string;
  uuid: string;
  name: string;
  content_type: string;
  created_at: Date;
  size: number;
  uploaded_by_id: string;
  url: string;
};

export const User: BuiltIn<UserValue> = {
  _marker: builtInType,
  name: "<builtins.user>",
};
export const Media: BuiltIn<MediaValue> = {
  _marker: builtInType,
  name: "<builtins.media>",
};

export function isBuiltIn(value: any): value is BuiltIn<any> {
  return (
    typeof value === "object" &&
    value !== null &&
    "_marker" in value &&
    value._marker === builtInType
  );
}
