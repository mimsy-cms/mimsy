import type { PageServerLoad, Actions } from './$types';
import { env } from '$env/dynamic/public';
import type { CollectionDefinition } from '$lib/collection/definition';
import { fail, redirect } from '@sveltejs/kit';
import { superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import z from 'zod';

async function fetchCollectionDefinition(
	collectionSlug: string,
	fetch: typeof globalThis.fetch
): Promise<CollectionDefinition> {
	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/definition`);
	return response.json();
}

export const load: PageServerLoad = async ({ params, fetch }) => {
	const form = await superValidate(zod(createSchema));
	const definition = await fetchCollectionDefinition(params.collectionSlug, fetch);

	return {
		slug: params.collectionSlug,
		form,
		definition
	};
};

const createSchema = z
	.record(z.string(), z.any())
	.refine((f) => /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(f.slug), { message: 'Invalid slug format' });

export const actions: Actions = {
	default: async ({ request, params, fetch }) => {
		const collectionSlug = params.collectionSlug;

		const form = await superValidate(request, zod(createSchema), { strict: false });
		if (!form.valid) {
			return fail(400, { form });
		}

		const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}`, {
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
				return fail(response.status, { error: `Failed to create resource: ${errorText}` });
			}
		}

		const createdResource = await response.json();

		throw redirect(303, `/collections/${collectionSlug}/${createdResource.slug}`);
	}
};
