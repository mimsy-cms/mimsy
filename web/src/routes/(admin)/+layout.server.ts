import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';
import type { LayoutServerLoad } from './$types';
import type { Collection } from '$lib/collection/definition';

export const load: LayoutServerLoad = async ({ fetch, locals }) => {
	if (!locals.user) {
		throw redirect(303, '/login');
	}

	const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections`);
	const globalResponse = await fetch(`${env.PUBLIC_API_URL}/v1/collections/globals`);

	const collections = (await response.json()) as Collection[];
	const globals = (await globalResponse.json()) as Collection[];

	return {
		collections,
		globals
	};
};
