<script lang="ts">
	import { page } from '$app/state';
	import { cn } from '$lib/cn';
	import type { Snippet } from 'svelte';

	type Props = {
		href: string;
		class?: string;
		children: Snippet;
		onNavigate?: () => void;
	};

	let { href, class: className, children, onNavigate }: Props = $props();
</script>

<a
	{href}
	onclick={onNavigate}
	class={cn(
		'group flex w-full items-center rounded-md py-2 pl-2 pr-1 text-left text-sm text-gray-600 hover:bg-gray-50 hover:text-gray-900',
		{
			'bg-blue-200/80 text-blue-800 hover:bg-blue-200 hover:text-blue-800':
				page.url.pathname === href
		},
		className
	)}
>
	{@render children()}
</a>
