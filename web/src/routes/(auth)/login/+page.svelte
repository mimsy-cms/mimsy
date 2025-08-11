<script lang="ts">
	import Error from '$lib/components/Error.svelte';
	import Input from '$lib/components/Input.svelte';
	import PasswordInput from '$lib/components/PasswordInput.svelte';
	import { superForm } from 'sveltekit-superforms';
	import { fa } from 'zod/v4/locales';

	let { data } = $props();

	let { form, errors, enhance, message } = superForm(
		data.form, {
			resetForm: false,
		});
</script>

<div class="flex min-h-screen flex-col items-center justify-center p-6">
	<h1 class="mb-6 text-4xl font-bold">Mimsy</h1>

	<form class="mb-4 flex w-full max-w-md flex-col gap-2" method="post" use:enhance novalidate>
		<div class="flex flex-col gap-2">
			<label for="email">Email</label>
			<Input
				type="email"
				id="email"
				name="email"
				bind:value={$form.email}
				error={!!$errors.email}
			/>
			<Error>{$errors.email}</Error>
		</div>

		<div class="flex flex-col gap-2">
			<label for="password">Password</label>
			<PasswordInput
				id="password"
				name="password"
				bind:value={$form.password}
				error={!!$errors.password}
			/>
			<Error>{$errors.password}</Error>
		</div>

		<button type="submit" class="btn mt-4">Login</button>
	</form>

	{#if $message}
		<p role="alert" class="text-red-600 mt-2 text-sm">{$message}</p>
	{/if}
</div>
