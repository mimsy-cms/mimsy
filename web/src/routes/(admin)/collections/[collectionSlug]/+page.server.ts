import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const collectionSlug = params.collectionSlug;

	const items = [] as const;

	return {
		collectionSlug: collectionSlug,
		items
	};
};
