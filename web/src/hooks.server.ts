import { redirect, type Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('session');

	if (token) {
		const res = await fetch(`${env.PUBLIC_API_URL}/v1/auth/me`, {
			headers: {
				Authorization: `Bearer ${token}`
			}
		});

		if (res.ok) {
			const user = await res.json();

			if (user.must_change_password && event.url.pathname !== '/password') {
				redirect(303, '/password');
			}

			event.locals.user = user;
		}
	}

	return resolve(event);
};
