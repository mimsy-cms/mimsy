import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const response = await fetch('/api/v1/collections/media', {
		headers: {
			'Content-Type': 'application/json'
		}
	});

	const data = await response.json();

	return { media: data };
};
