import { env } from '$env/dynamic/public';
import type { LayoutServerLoad } from './$types';
import type { Collection } from '$lib/collection/definition';

export const load: LayoutServerLoad = async ({ fetch }) => {
	const [collections, globals]: [Collection[], Collection[]] = await Promise.all([
		fetch(`${env.PUBLIC_API_URL}/v1/collections`).then((res) => res.json()),
		fetch(`${env.PUBLIC_API_URL}/v1/collections/globals`).then((res) => res.json())
	]);

	return {
		collections,
		globals
	};
};
