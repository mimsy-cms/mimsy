import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';
import { redirect } from '@sveltejs/kit';
import type { CollectionDefinition, CollectionResource } from '$lib/collection/definition';
import type { User } from '$lib/types/user';

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
): Promise<CollectionResource> {
	const response = await fetch(
		`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/${resourceSlug}`
	);
	return response.json();
}

async function fetchUser(id: number, fetch: typeof globalThis.fetch): Promise<User> {
	const response = await fetch(`${env.PUBLIC_API_URL}/v1/users/${id}`);
	return response.json();
}

export const load: PageServerLoad = async ({ params, fetch }) => {
	if (params.resourceSlug === 'create') {
		throw redirect(307, `/collections/${params.collectionSlug}/create`);
	}

	const [definition, resource] = await Promise.all([
		fetchCollectionDefinition(params.collectionSlug, fetch),
		fetchResource(params.collectionSlug, params.resourceSlug, fetch)
	]);

	const [createdBy, updatedBy] = await Promise.all([
		fetchUser(resource.created_by, fetch),
		fetchUser(resource.updated_by, fetch)
	]);

	return {
		slug: params.collectionSlug,
		definition,
		resource,
		createdBy,
		updatedBy
	};
};
