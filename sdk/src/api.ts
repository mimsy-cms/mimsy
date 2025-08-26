import { MediaValue, UserValue } from "./builtins";
import {
  Collection,
  ObjectOf,
  ObjectOfCollection,
  postProcess,
  Schema,
} from "./collection";
import { CollectionOrBuiltin, UnfetchedRelation } from "./fields";
import { fetchRelation } from "./requests";

export class MimsyClient {
  constructor(private readonly baseUrl: string) {}

  fetch(url: string) {
    return fetch(`${this.baseUrl}${url}`);
  }

  with<T extends Collection<any>>(collection: T) {
    return new MimsyCollectionClient(this, collection);
  }

  async fetchRelation<T extends CollectionOrBuiltin<any>>(
    unfetchedRelation: UnfetchedRelation<T>,
  ): Promise<ObjectOfCollection<T>> {
    return fetchRelation(this, unfetchedRelation);
  }

  user() {
    return new MimsyUserClient(this);
  }

  media() {
    return new MimsyMediaClient(this);
  }
}

export class MimsyUserClient {
  constructor(private readonly client: MimsyClient) {}

  async all(): Promise<UserValue[]> {
    const response = await this.client.fetch("/v1/users");
    return response.json();
  }

  async get(id: string): Promise<UserValue> {
    const response = await this.client.fetch(`/v1/users/${id}`);
    return response.json();
  }
}

export class MimsyMediaClient {
  constructor(private readonly client: MimsyClient) {}

  async all(): Promise<MediaValue[]> {
    const response = await this.client.fetch("/v1/media");
    return response.json();
  }

  async get(id: string): Promise<MediaValue> {
    const response = await this.client.fetch(`/v1/media/${id}`);
    return response.json();
  }
}

export class MimsyCollectionClient<T extends Collection<any>> {
  constructor(
    private readonly client: MimsyClient,
    private readonly collection: T,
  ) {}

  async all(): Promise<ObjectOfCollection<T>[]> {
    const response = await this.client.fetch(`/v1/${this.collection.name}`);
    const data = await response.json();

    return data.map(
      (item: any) =>
        postProcess(this.collection, item) as ObjectOfCollection<T>,
    );
  }

  async get(id: string): Promise<ObjectOfCollection<T>> {
    const response = await this.client.fetch(
      `/v1/${this.collection.name}/${id}`,
    );
    const data = await response.json();

    return postProcess(this.collection, data) as ObjectOfCollection<T>;
  }
}
