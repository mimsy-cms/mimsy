import type { Actions, PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';
import type { CollectionDefinition, CollectionResource } from '$lib/collection/definition';
import { fail, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import z from 'zod';

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
	const response = await fetch(`${env.PUBLIC_API_URL}/v1/globals/${collectionSlug}`);
	return response.json();
}

const updateSchema = z.record(z.string(), z.any());

export const load: PageServerLoad = async ({ params, fetch }) => {
	const [definition, resource] = await Promise.all([
		fetchCollectionDefinition(params.slug, fetch),
		fetchResource(params.slug, fetch)
	]);

	const form = await superValidate(zod(updateSchema), {
		defaults: resource
	});

	return {
		slug: params.slug,
		form,
		definition,
		resource
	};
};

export const actions: Actions = {
	default: async ({ request, params, fetch }) => {
		const form = await superValidate(request, zod(updateSchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${params.slug}`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(form.data)
		});

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
