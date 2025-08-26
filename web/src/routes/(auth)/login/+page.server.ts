import { z } from 'zod/v4';
import type { Actions, PageServerLoad } from './$types';
import { fail, superValidate } from 'sveltekit-superforms';
import { zod4 } from 'sveltekit-superforms/adapters';
import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

const loginSchema = z.object({
	email: z.email(),
	password: z.string().min(1, { error: 'This field is required' })
});

export const load: PageServerLoad = async () => {
	const form = await superValidate(zod4(loginSchema));

	return {
		form
	};
};

export const actions: Actions = {
	default: async ({ request, cookies, fetch }) => {
		const form = await superValidate(request, zod4(loginSchema));

		if (!form.valid) {
			form.data.password = '';
			return fail(400, { form });
		}

		const res = await fetch(`${env.PUBLIC_API_URL}/v1/auth/login`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(form.data)
		});

		if (!res.ok) {
			const msg = await res.text();
			form.message = msg ?? 'Invalid credentials';
			return fail(res.status, { form });
		}

		const data = await res.json();
		cookies.set('session', data.session, {
			httpOnly: true,
			path: '/',
			maxAge: 60 * 60 * 24 * 7, // 7 days
			sameSite: 'strict',
			secure: true
		});

		if (data.mustChangePassword) {
			redirect(303, '/password');
		}

		redirect(302, '/');
	}
};
