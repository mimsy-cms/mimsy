import { type Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('session');

	if (token) {
		const res = await event.fetch(`${env.PUBLIC_API_URL}/v1/auth/me`);

		if (res.ok) {
			const user = await res.json();
			event.locals.user = user;
		}
	}

	return resolve(event);
};
