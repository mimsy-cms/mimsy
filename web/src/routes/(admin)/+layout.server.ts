import type { Actions, PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const session = cookies.get('session');
  
    if (!session) {
      throw redirect(303, '/login');
    }

    const res = await fetch(`${env.PUBLIC_API_URL}/v1/auth/me`, {
      headers: {
        Authorization: `Bearer ${session}`,
      },
    });
  
    if (!res.ok) {
      throw redirect(303, '/login');
    }
  
    const user = await res.json();
  
    if (user.must_change_password) {
      throw redirect(303, '/password');
    }
  
    return { user };
};
