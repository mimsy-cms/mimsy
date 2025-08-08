import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const response = await fetch(`/api/v1/media/${params.id}`, {
		headers: {
			'Content-Type': 'application/json'
		}
	});

	const media = await response.json();

	return { media };
};
