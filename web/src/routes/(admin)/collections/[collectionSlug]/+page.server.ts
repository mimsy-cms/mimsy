import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';

type Resource = {
	id: number;
	slug: string;
};

export const load: PageServerLoad = async ({ params, fetch }) => {
	const collectionSlug = params.collectionSlug;

	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}`);
	const resources = (await response.json()) as Resource[];

	return {
		collectionSlug: collectionSlug,
		resources
	};
};
