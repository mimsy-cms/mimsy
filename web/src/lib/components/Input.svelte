<script lang="ts">
	import { cn } from '$lib/cn';
	import type { Snippet } from 'svelte';

	type Props = {
		id?: string;
		name?: string;
		value?: string | number;
		placeholder?: string;
		type?: 'text' | 'email' | 'password' | 'date' | 'url' | 'number';
		disabled?: boolean;
		required?: boolean;
		error?: boolean;
		class?: string;
		oninput?: (event: Event) => void;
	};

	let {
		id,
		name,
		value = $bindable(),
		placeholder,
		type = 'text',
		disabled = false,
		required = false,
		error = false,
		class: className,
		oninput
	}: Props = $props();
</script>

<input
	{id}
	{name}
	{type}
	{placeholder}
	{disabled}
	{required}
	aria-invalid={error ? 'true' : undefined}
	{oninput}
	bind:value
	class={cn(
		'rounded-md border border-gray-300 bg-white px-3 py-1.5 outline-none focus:border-blue-500 focus:ring-blue-500',
		{
			'border-red-500 text-red-500': error
		},
		className
	)}
/>
