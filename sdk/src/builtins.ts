// To make declaration of fields out of this file more difficult.
const builtInType = Symbol("mimsy-builtin");

export type BuiltInValue = {
  _marker: typeof builtInType;
};

export type BuiltIn<T extends BuiltInValue> = {
  name: string;
};

export type UserValue = BuiltInValue & {
  id: string;
  email: string;
  name: string;
};

export type MediaValue = BuiltInValue & {
  id: string;
  url: string;
};

export const User: BuiltIn<UserValue> = {
  name: "User",
};
export const Media: BuiltIn<MediaValue> = {
  name: "Media",
};
