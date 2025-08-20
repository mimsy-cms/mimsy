import type { PageServerLoad, Actions } from './$types';
import { env } from '$env/dynamic/public';
import { fail, redirect, type Redirect } from '@sveltejs/kit';

type Resource = {
	id: number;
	slug: string;
	created_by_email: string;
	updated_at: string;
};

export const load: PageServerLoad = async ({ params, fetch }) => {
	const collectionSlug = params.collectionSlug;

	const defRes = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/definition`);
	if (!defRes.ok) {
		throw new Error(`Failed to fetch collection definition for ${collectionSlug}`);
	}
	const collectionDef = await defRes.json();

	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}`);
	const resources = (await response.json()) as Resource[];

	return {
		collectionSlug: collectionSlug,
		collectionName: collectionDef.name,
		resources
	};
};

export const actions: Actions = {
	create: async ({ request, params, fetch, cookies }) => {
		const collectionSlug = params.collectionSlug;
		const data = await request.formData();
		const slug = data.get('slug') as string;

		if (!slug || slug.trim() === '') {
			return { error: 'Slug is required' };
		}

		const slugPattern = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;
		if (!slugPattern.test(slug.trim())) {
			return fail(400, { error: 'Invalid slug format. Use lowercase letters, numbers, and hyphens only.' });
		}

		const cookieHeader = cookies.getAll()
			.map(cookie => `${cookie.name}=${cookie.value}`)
			.join('; ');

		const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Cookie': cookieHeader
			},
			body: JSON.stringify({ slug: slug.trim() })
		});

		if (!response.ok) {
			const errorText = await response.text();
			console.error('API Error:', response.status, errorText);

			if (response.status === 409) {
				return fail(409, { error: 'Resource already exists' });
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

