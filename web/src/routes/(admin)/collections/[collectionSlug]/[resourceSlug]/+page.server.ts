import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';
import type { CollectionDefinition, CollectionResource } from '$lib/collection/definition';
import type { User } from '$lib/types/user';
import z from 'zod';
import { superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import type { Actions } from './$types';
import { fail } from '@sveltejs/kit';

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

const updateSchema = z.record(z.string(), z.any());

export const load: PageServerLoad = async ({ params, fetch }) => {
	const [definition, resource] = await Promise.all([
		fetchCollectionDefinition(params.collectionSlug, fetch),
		fetchResource(params.collectionSlug, params.resourceSlug, fetch)
	]);

	const form = await superValidate(zod(updateSchema), {
		defaults: resource
	});

	const [createdBy, updatedBy] = await Promise.all([
		fetchUser(resource.created_by, fetch),
		fetchUser(resource.updated_by, fetch)
	]);

	return {
		slug: params.collectionSlug,
		form,
		definition,
		resource,
		createdBy,
		updatedBy
	};
};

export const actions: Actions = {
	default: async ({ request, params, fetch }) => {
		const form = await superValidate(request, zod(updateSchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const response = await fetch(
			`${env.PUBLIC_API_URL}/v1/collections/${params.collectionSlug}/${params.resourceSlug}`,
			{
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(form.data)
			}
		);

		if (!response.ok) {
			const errorText = await response.text();
			console.error('API Error:', response.status, errorText);

			if (response.status === 409) {
				return fail(409, { error: 'Resource with that slug already exists' });
			} else if (response.status === 401) {
				return fail(401, { error: 'You are not authorized to create resources' });
			} else if (response.status === 404) {
				return fail(404, { error: 'Collection not found' });
			} else {
				return fail(response.status, { error: `Failed to update resource: ${errorText}` });
			}
		}
	}
};
