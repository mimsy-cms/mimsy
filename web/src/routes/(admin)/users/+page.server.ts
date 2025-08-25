import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const res = await fetch('/api/v1/users');
	if (!res.ok) {
		throw new Error('Failed to fetch users');
	}

	const users = await res.json();

	return { users };
};
