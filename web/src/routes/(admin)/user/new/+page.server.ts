import type { Actions, PageServerLoad } from './$types';
import { z } from 'zod/v4';
import { zod4 } from 'sveltekit-superforms/adapters';
import { error, redirect } from '@sveltejs/kit';
import { fail, superValidate } from 'sveltekit-superforms';

const newUserSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
  isAdmin: z.coerce.boolean().optional()
});

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const session = cookies.get('session');
  console.log('[LOAD] Session:', session);

  if (!session) {
    console.log('[LOAD] No session, redirecting to login');
    throw redirect(303, '/login');
  }

  const res = await fetch('http://localhost:3000/v1/auth/me', {
    headers: {
      Authorization: `Bearer ${session}`,
    },
  });

  console.log('[LOAD] /me response status:', res.status);
  if (!res.ok) {
    console.log('[LOAD] Unauthorized session, redirecting to login');
    throw redirect(303, '/login');
  }

  const user = await res.json();
  console.log('[LOAD] Logged in user:', user);

  if (!user.is_admin) {
    console.log('[LOAD] User is not admin, redirecting to home');
    throw redirect(303, '/');
  }

  const form = await superValidate(zod4(newUserSchema));
  console.log('[LOAD] Loaded form:', form);

  return { form, session };
};

export const actions: Actions = {
  default: async ({ request, cookies, fetch }) => {
    const formData = await request.formData();

    for (const pair of formData.entries()) {
      console.log('[FORM DATA]', pair[0], '=', pair[1]);
    }

    const session = cookies.get('session');
    const form = await superValidate(formData, zod4(newUserSchema));

    console.log('[ACTION] Session:', session);

    if (!form.valid) {
      console.log('[ACTION] Form validation failed:', form);
      return fail(400, { form });
    }

    console.log('[ACTION] Form is valid:', form.data);

    const res = await fetch('http://localhost:3000/v1/auth/me', {
      headers: {
        Authorization: `Bearer ${session}`,
      },
    });

    console.log('[ACTION] /me response status:', res.status);
    if (!res.ok) {
      console.log('[ACTION] Session check failed, redirecting to login');
      throw redirect(303, '/login');
    }

    const user = await res.json();
    console.log('[ACTION] Logged in user:', user);

    if (!user.is_admin) {
      console.log('[ACTION] User is not admin, redirecting to home');
      throw redirect(303, '/');
    }

    const { email, password, isAdmin } = form.data;
    console.log('[ACTION] Registering new user with:', { email, password, isAdmin });

    const registerRes = await fetch('http://localhost:3000/v1/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${session}`,
      },
      body: JSON.stringify({ email, password, isAdmin }),
    });

    console.log('[ACTION] Register response status:', registerRes.status);

    if (!registerRes.ok) {
      const errorData = await registerRes.json();
      console.error('[ACTION] Failed to register user:', errorData);
      return fail(registerRes.status, {
        form,
        message: errorData.error || 'Failed to create user',
      });
    }

    console.log('[ACTION] User created successfully');
    return {
      form,
      message: 'User created successfully!',
    };
  }
};
