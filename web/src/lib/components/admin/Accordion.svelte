<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import { slide } from 'svelte/transition';
	import type { Component } from 'svelte';
	import MenuItem from './MenuItem.svelte';

	type Item = {
		name: string;
		href: string;
	};

	type Props = {
		text: string;
		emptyText?: string;
		icon: Component;
		items: Item[];
		onNavigate?: () => void;
	};

	let { text, emptyText, icon: Icon, items, onNavigate }: Props = $props();

	let open = $state(items.length > 0);
</script>

<div class="space-y-1">
	<button
		onclick={() => (open = !open)}
		class="group flex w-full items-center rounded-md px-2 py-2 text-left text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900"
	>
		<Icon class="mr-3 h-5 w-5 flex-shrink-0" />
		<span class="flex-1">{text}</span>
		<ChevronDownIcon
			class="ml-3 size-4 transform transition-transform {open ? 'rotate-180' : ''}"
		/>
	</button>

	{#if open}
		<div transition:slide={{ duration: 80 }}>
			<div class="space-y-1">
				{#each items as item}
					<MenuItem href={item.href} {onNavigate}>{item.name}</MenuItem>
				{/each}

				{#if items.length === 0}
					<p class="px-2 text-sm text-gray-400">{emptyText}</p>
				{/if}
			</div>
		</div>
	{/if}
</div>
