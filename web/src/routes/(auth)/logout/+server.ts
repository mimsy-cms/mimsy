import { redirect, json } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

export async function POST({ fetch, cookies }) {
	const res = await fetch(`${env.PUBLIC_API_URL}/v1/auth/logout`, {
		method: 'POST',
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json'
		}
	});

	if (!res.ok) {
		const errorText = await res.text();
		return new Response(errorText, { status: res.status });
	}

	cookies.delete('session', { path: '/' });
	return json({ success: true });
}
