import { MimsyClient } from "./api";
import { isBuiltIn as isBuiltin } from "./builtins";
import { ObjectOf, ObjectOfCollection, Schema } from "./collection";
import { CollectionOrBuiltin, Field, UnfetchedRelation } from "./fields";

export async function fetchRelation<T extends CollectionOrBuiltin<any>>(
  client: MimsyClient,
  unfetched: UnfetchedRelation<T>,
): Promise<ObjectOfCollection<T>> {
  const typeInfo = unfetched["_collection"];
  if (isBuiltin(typeInfo)) {
    // Handle different builtin types
    switch (typeInfo.name) {
      case "<builtins.user>":
        // I can confirm that the type is correct
        return (await client
          .user()
          .get(unfetched.id)) as unknown as ObjectOfCollection<T>;
      case "<builtins.media>":
        // I can confirm that the type is correct
        return (await client
          .media()
          .get(unfetched.id)) as unknown as ObjectOfCollection<T>;
      default:
        throw new Error(`Unknown builtin type: ${typeInfo.name}`);
    }
  } else {
    // Handle custom types
    const relation = await client.with(typeInfo).get(unfetched.id);
    return relation as unknown as ObjectOfCollection<T>;
  }
}
