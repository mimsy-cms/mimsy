import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';
import type { LayoutServerLoad } from './$types';

type Collection = {
	name: string;
	slug: string;
};

export const load: LayoutServerLoad = async ({ cookies, fetch, locals }) => {
	if (!locals.user) {
		throw redirect(303, '/login');
	}

	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections`);

	const collections = (await response.json()) as Collection[];

	return {
		collections
	};
};
