import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';

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
