import { z } from 'zod/v4';
import type { Actions, PageServerLoad } from './$types';
import { fail, superValidate } from 'sveltekit-superforms';
import { zod4 } from 'sveltekit-superforms/adapters';
import { redirect } from '@sveltejs/kit';

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
	default: async ({ request }) => {
		const form = await superValidate(request, zod4(loginSchema));

		if (!form.valid) {
			return fail(400, { form });
		}

		// TODO: Login user

		redirect(302, '/');
	}
};
