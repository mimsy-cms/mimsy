import { redirect, json } from '@sveltejs/kit';

export async function POST({ fetch, cookies }) {
	const res = await fetch('http://localhost:3000/v1/auth/logout', {
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
