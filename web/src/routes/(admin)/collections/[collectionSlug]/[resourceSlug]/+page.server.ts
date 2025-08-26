import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';
import type { CollectionDefinition, CollectionResource } from '$lib/collection/definition';
import type { User } from '$lib/types/user';
import z from 'zod';
import { superValidate, message } from 'sveltekit-superforms';
import { getFlash } from 'sveltekit-flash-message';
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

export const load: PageServerLoad = async ({ params, fetch, cookies }) => {
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
	default: async ({ request, params, fetch, cookies }) => {
		const form = await superValidate(request, zod(updateSchema));
		if (!form.valid) {
			return message(form, 'Invalid form data', { status: 400 });
		}

		try {
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

				let errorMessage: string;
				if (response.status === 409) {
					errorMessage = 'Resource with that slug already exists';
				} else if (response.status === 401) {
					errorMessage = 'You are not authorized to create resources';
				} else if (response.status === 404) {
					errorMessage = 'Collection not found';
				} else {
					errorMessage = `Failed to update resource: ${errorText}`;
				}

				return message(form, errorMessage, { status: response.status as 400 | 404 | 409 | 401 });
			}

			return message(form, 'Resource updated successfully');
		} catch (error) {
			console.error('Error updating resource:', error);
			return message(form, 'Failed to update resource', { status: 500 });
		}
	}
};
