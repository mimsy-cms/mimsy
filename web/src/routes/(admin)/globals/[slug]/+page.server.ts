import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';
import type { CollectionDefinition, CollectionResource } from '$lib/collection/definition';

async function fetchCollectionDefinition(
	collectionSlug: string,
	fetch: typeof globalThis.fetch
): Promise<CollectionDefinition> {
	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/definition`);
	return response.json();
}

async function fetchResource(
	collectionSlug: string,
	fetch: typeof globalThis.fetch
): Promise<CollectionResource> {
	const response = await fetch(
		`${env.PUBLIC_API_URL}/v1/globals/${collectionSlug}`
	);
	return response.json();
}

export const load: PageServerLoad = async ({ params, fetch }) => {
	const [definition, resource] = await Promise.all([
		fetchCollectionDefinition(params.slug, fetch),
		fetchResource(params.slug, fetch)
	]);

	return {
		slug: params.slug,
		definition,
		resource
	};
};
