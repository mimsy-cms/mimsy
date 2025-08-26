import type { PageServerLoad, Actions } from './$types';
import { env } from '$env/dynamic/public';
import type { CollectionDefinition } from '$lib/collection/definition';
import { fail } from '@sveltejs/kit';
import { superValidate, message } from 'sveltekit-superforms';
import { redirect } from 'sveltekit-flash-message/server';
import { zod } from 'sveltekit-superforms/adapters';
import z from 'zod';

async function fetchCollectionDefinition(
	collectionSlug: string,
	fetch: typeof globalThis.fetch
): Promise<CollectionDefinition> {
	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/definition`);
	return response.json();
}

const createSchema = z
	.record(z.string(), z.any())
	.refine((f) => /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(f.slug), { message: 'Invalid slug format' });

export const load: PageServerLoad = async ({ params, fetch }) => {
	const form = await superValidate(zod(createSchema));
	const definition = await fetchCollectionDefinition(params.collectionSlug, fetch);

	return {
		slug: params.collectionSlug,
		form,
		definition
	};
};

export const actions: Actions = {
	default: async ({ request, params, fetch, cookies }) => {
		const collectionSlug = params.collectionSlug;

		const form = await superValidate(request, zod(createSchema), { strict: false });

		if (!form.valid) {
			return message(form, 'Invalid form data', { status: 400 });
		}

		try {
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

				let errorMessage: string;
				if (response.status === 409) {
					errorMessage = 'Resource with that slug already exists';
				} else if (response.status === 401) {
					errorMessage = 'You are not authorized to create resources';
				} else if (response.status === 404) {
					errorMessage = 'Collection not found';
				} else {
					errorMessage = `Failed to create resource: ${errorText}`;
				}

				return message(form, errorMessage, { status: response.status as 400 | 404 | 409 | 401 });
			}

			const createdResource = await response.json();

			redirect(`/collections/${collectionSlug}/${createdResource.slug}`, { type: 'success', message: 'Resource created successfully' }, cookies);
		} catch (error) {
			if (error instanceof Error) {
				console.error('Error creating resource:', error);
				return message(form, 'Failed to create resource', { status: 500 });
			}
			throw error;
		}
	}
};
