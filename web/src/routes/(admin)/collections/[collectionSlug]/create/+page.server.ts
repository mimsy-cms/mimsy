import type { PageServerLoad, Actions } from './$types';
import { env } from '$env/dynamic/public';
import type { CollectionDefinition } from '$lib/collection/definition';
import { fail, redirect } from '@sveltejs/kit';

async function fetchCollectionDefinition(
	collectionSlug: string,
	fetch: typeof globalThis.fetch
): Promise<CollectionDefinition> {
	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/definition`);
	return response.json();
}

export const load: PageServerLoad = async ({ params, fetch }) => {
	const [definition] = await Promise.all([fetchCollectionDefinition(params.collectionSlug, fetch)]);

	return {
		slug: params.collectionSlug,
		definition
	};
};

export const actions: Actions = {
	create: async ({ request, params, fetch }) => {
		const collectionSlug = params.collectionSlug;
		const data = await request.formData();
		const slug = data.get('slug') as string;

		if (!slug || slug.trim() === '') {
			return { error: 'Slug is required' };
		}

		const slugPattern = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;
		if (!slugPattern.test(slug.trim())) {
			return fail(400, {
				error: 'Invalid slug format. Use lowercase letters, numbers, and hyphens only.'
			});
		}

		const body: Record<string, any> = {};
		for (const [key, value] of data.entries()) {
			body[key] = value;
		}

		console.log(body);

		const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(body)
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
