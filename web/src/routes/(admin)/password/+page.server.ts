import { z } from 'zod/v4';
import type { Actions, PageServerLoad } from './$types';
import { fail, superValidate } from 'sveltekit-superforms';
import { zod4 } from 'sveltekit-superforms/adapters';
import { error, redirect } from '@sveltejs/kit';

const passwordSchema = z
  .object({
    old_password: z.string().min(1, { message: 'Current password is required' }),
    new_password: z.string().min(8, { message: 'New password must be at least 8 characters long' }),
    confirm_password: z.string().min(1, { message: 'Please confirm your new password' })
  })
  .refine(data => data.new_password === data.confirm_password, {
    path: ['confirm_password'],
    message: 'Passwords do not match'
  });

export const load: PageServerLoad = async () => {
    const form = await superValidate(zod4(passwordSchema));

    return {
        form
    };
};

export const actions: Actions = {
  default: async ({ request, fetch }) => {
    const form = await superValidate(request, zod4(passwordSchema));
    if (!form.valid) return fail(400, { form });

    const { old_password, new_password } = form.data;

    const res = await fetch('http://localhost:3000/v1/auth/password', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        old_password,
        new_password
      })
    });

    if (!res.ok) {
      const errorText = await res.text();
      form.message = `Failed to change password: ${errorText}`;
      return fail(res.status, { form });
    }

    throw redirect(303, '/');
  }
};

