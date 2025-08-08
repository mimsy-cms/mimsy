import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const response = await fetch('/api/v1/media', {
		headers: {
			'Content-Type': 'application/json'
		}
	});

	if (!response.ok) {
		console.error('Failed to fetch media:', response.status, await response.text());
		return { media: [] };
	}

	const data = await response.json();

	return { media: data };
};
