import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { collectionSlug } = params;

    const definitionResponse = await fetch(`http://localhost:3000/v1/collections/${collectionSlug}/definition`);
    if (!definitionResponse.ok) {
        throw new Error(`Failed to fetch collection definition: ${definitionResponse.statusText}`);
    }

    const definition = await definitionResponse.json();

    return {
        slug: collectionSlug,
        definition
    };
};