import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/public';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { collectionSlug } = params;

    const definitionResponse = await fetch(`${env.PUBLIC_API_URL}/v1/collections/${collectionSlug}/definition`);
    if (!definitionResponse.ok) {
        throw new Error(`Failed to fetch collection definition: ${definitionResponse.statusText}`);
    }

    const definition = await definitionResponse.json();

    return {
        slug: collectionSlug,
        definition
    };
};