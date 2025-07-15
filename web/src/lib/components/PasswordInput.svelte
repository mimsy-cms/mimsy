<script lang="ts">
	import EyeOffIcon from '@lucide/svelte/icons/eye-off';
	import EyeIcon from '@lucide/svelte/icons/eye';
	import Input from './Input.svelte';
	import { cn } from '$lib/cn';

	type Props = {
		id: string;
		name: string;
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		error?: boolean;
		class?: string;
	};

	let {
		id,
		name,
		value = $bindable(),
		placeholder,
		disabled = false,
		error = false,
		class: className
	}: Props = $props();

	let showPassword = $state(false);
</script>

<div
	class="flex gap-2 rounded-md border border-gray-300 bg-white px-3 py-1.5 outline-none focus-within:border-blue-500 focus-within:ring-blue-500"
>
	<Input
		{id}
		{name}
		type={showPassword ? 'text' : 'password'}
		bind:value
		{placeholder}
		{disabled}
		{error}
		class={cn('flex-1 border-none bg-transparent p-0', className)}
	/>
	<button
		type="button"
		class="text-gray-400 hover:text-gray-600"
		onclick={() => (showPassword = !showPassword)}
	>
		{#if showPassword}
			<EyeOffIcon class="size-5" />
		{:else}
			<EyeIcon class="size-5" />
		{/if}
	</button>
</div>
