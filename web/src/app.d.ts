// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
  namespace App {
    interface Locals {
      user?: {
        id: number;
        email: string;
        must_change_password: boolean;
        // add any other fields your user has
      };
    }

    // You can also extend PageData, Error, etc., if needed
    // interface PageData {}
    // interface Error {}
    // interface PageState {}
    // interface Platform {}
  }
}

export {};
