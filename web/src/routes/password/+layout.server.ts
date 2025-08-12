import type { LayoutServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: LayoutServerLoad = async ({ cookies }) => {
	const session = cookies.get('session');

	if (!session) {
		// Not logged in
		throw redirect(303, '/login');
	}
};
