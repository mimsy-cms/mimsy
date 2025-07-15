<script lang="ts">
	import SearchIcon from '@lucide/svelte/icons/search';
	import Input from '../Input.svelte';
	import { cn } from '$lib/cn';

	type Props = {
		id: string;
		name: string;
		value?: string;
		placeholder?: string;
		class?: string;
	};

	let { id, name, value = $bindable(), placeholder, class: className }: Props = $props();
</script>

<div
	class={cn(
		'flex gap-2 rounded-md border border-gray-300 bg-white px-3 py-1.5 outline-none focus-within:border-blue-500 focus-within:ring-blue-500',
		className
	)}
>
	<button type="submit" class="text-gray-400 hover:text-gray-600">
		<SearchIcon class="size-5" />
	</button>
	<Input
		{id}
		{name}
		type="search"
		bind:value
		{placeholder}
		class="flex-1 border-none bg-transparent p-0"
	/>
</div>
