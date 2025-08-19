import type { Actions, PageServerLoad } from './$types';
import { z } from 'zod/v4';
import { zod4 } from 'sveltekit-superforms/adapters';
import { error, redirect } from '@sveltejs/kit';
import { fail, superValidate } from 'sveltekit-superforms';
import { env } from '$env/dynamic/public';

const newUserSchema = z.object({
	email: z.string().email(),
	password: z.string().min(8),
	isAdmin: z.coerce.boolean().optional()
});

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const session = cookies.get('session');

	if (!session) {
		throw redirect(303, '/login');
	}

	const res = await fetch(`${env.PUBLIC_API_URL}/v1/auth/me`, {
		headers: {
			Authorization: `Bearer ${session}`
		}
	});

	if (!res.ok) {
		throw redirect(303, '/login');
	}

	const user = await res.json();

	if (!user.is_admin) {
		throw redirect(303, '/');
	}

	const form = await superValidate(zod4(newUserSchema));

	return { form, session };
};

export const actions: Actions = {
	default: async ({ request, cookies, fetch }) => {
		const formData = await request.formData();

		const session = cookies.get('session');
		const form = await superValidate(formData, zod4(newUserSchema));

		if (!form.valid) {
			return fail(400, { form });
		}

		const res = await fetch(`${env.PUBLIC_API_URL}/v1/auth/me`, {
			headers: {
				Authorization: `Bearer ${session}`
			}
		});

		if (!res.ok) {
			throw redirect(303, '/login');
		}

		const user = await res.json();

		if (!user.is_admin) {
			throw redirect(303, '/');
		}

		const { email, password, isAdmin } = form.data;

		const registerRes = await fetch(`${env.PUBLIC_API_URL}/v1/auth/register`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${session}`
			},
			body: JSON.stringify({ email, password, isAdmin })
		});

		if (!registerRes.ok) {
			const errorData = await registerRes.json();
			return fail(registerRes.status, {
				form,
				message: errorData.error || 'Failed to create user'
			});
		}

		return {
			form,
			message: 'User created successfully!'
		};
	}
};
