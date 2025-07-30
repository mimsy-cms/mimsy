import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const collectionSlug = params.collectionSlug;

    const response = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/items`);
    if (!response.ok) {
        throw new Error(`Failed to fetch items: ${response.statusText}`);
    }

    const rawItems = await response.json();

    const items = rawItems.map((item: any) => ({
        id: item.id,
        resourceSlug: item.slug,
        ...item.data
    }));

    return {
        collectionSlug: collectionSlug,
        items
    };
};