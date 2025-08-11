import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';
import type { CollectionDefinition } from '$lib/collection/definition';

async function fetchCollectionDefinition(
	collectionSlug: string,
	fetch: typeof globalThis.fetch
): Promise<CollectionDefinition> {
	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/definition`);
	return response.json();
}

async function fetchResource(
	collectionSlug: string,
	resourceSlug: string,
	fetch: typeof globalThis.fetch
): Promise<any> {
	const response = await fetch(
		`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/${resourceSlug}`
	);
	if (!response.ok) {
		throw new Error(`Failed to fetch resource: ${response.statusText}`);
	}
	return response.json();
}

export const load: PageServerLoad = async ({ params, fetch }) => {
	const [definition, resource] = await Promise.all([
		fetchCollectionDefinition(params.collectionSlug, fetch),
		fetchResource(params.collectionSlug, params.resourceSlug, fetch)
	]);

	return {
		slug: params.collectionSlug,
		definition,
		resource
	};
};
