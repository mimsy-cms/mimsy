// See https://svelte.dev/docs/kit/types#app.d.ts

import type { AuthUser } from '$lib/types/user';

// for information about these interfaces
declare global {
	namespace App {
		interface Locals {
			user?: AuthUser;
		}

		// You can also extend PageData, Error, etc., if needed
		// interface PageData {}
		// interface Error {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
