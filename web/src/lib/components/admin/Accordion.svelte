<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import { slide } from 'svelte/transition';
	import type { Component } from 'svelte';
	import { cn } from '$lib/cn';
	import { page } from '$app/state';

	type Item = {
		name: string;
		href: string;
	};

	type Props = {
		text: string;
		icon: Component;
		items: Item[];
	};

	let { text, icon: Icon, items }: Props = $props();

	let open = $state(true);
</script>

<div class="space-y-1">
	<button
		onclick={() => (open = !open)}
		class="group flex w-full items-center rounded-md py-2 pl-2 pr-1 text-left text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900"
	>
		<Icon class="mr-3 h-5 w-5 flex-shrink-0" />
		<span class="flex-1">{text}</span>
		<ChevronDownIcon
			class="ml-3 h-4 w-4 transform transition-transform {open ? 'rotate-180' : ''}"
		/>
	</button>

	{#if open}
		<div transition:slide>
			<div class="space-y-1">
				{#each items as item}
					<a
						href={item.href}
						class={cn(
							'group flex w-full items-center rounded-md py-2 pl-2 pr-1 text-left text-sm text-gray-600 hover:bg-gray-50 hover:text-gray-900',
							{
								'bg-blue-200/80 text-blue-800': page.url.pathname === item.href
							}
						)}
					>
						{item.name}
					</a>
				{/each}
			</div>
		</div>
	{/if}
</div>
