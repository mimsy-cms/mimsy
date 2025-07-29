import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const collectionSlug = params.collectionSlug;

    const response = await fetch(`http://localhost:3000/v1/collections/${collectionSlug}/items`);
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