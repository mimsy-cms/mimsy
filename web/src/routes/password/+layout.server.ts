import type { Actions, PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const session = cookies.get('session');

  if (!session) {
    // Not logged in
    throw redirect(303, '/login');
  }
};
