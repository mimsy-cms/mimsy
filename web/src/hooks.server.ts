import { redirect, type Handle, type RequestEvent } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

async function fetchAuthUser(event: RequestEvent) {
	const response = await event.fetch(`${env.PUBLIC_API_URL}/v1/auth/me`);

	if (response.ok) {
		return await response.json();
	}

	return null;
}

export const handle: Handle = async ({ event, resolve }) => {
	const user = await fetchAuthUser(event);
	const isAdminRoute = event.route.id?.startsWith('/(admin)');

	// We redirect the user to the login page if they are not authenticated
	// and trying to access an admin route
	if (isAdminRoute && !user) {
		redirect(302, '/login');
	}

	event.locals.user = user;

	return resolve(event);
};
