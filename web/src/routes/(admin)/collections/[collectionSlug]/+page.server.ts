import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';

export const load: PageServerLoad = async ({ params }) => {
	const collectionSlug = params.collectionSlug;

	const defRes = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/definition`);
	if (!defRes.ok) {
		throw new Error(`Failed to fetch collection definition for ${collectionSlug}`);
	}
	const collectionDef = await defRes.json();

	const items = [] as const;

	return {
		collectionSlug,
		collectionName: collectionDef.name,
		items
	};
};
