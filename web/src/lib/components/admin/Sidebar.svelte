<script lang="ts">
	import DatabaseIcon from '@lucide/svelte/icons/database';
	import GlobeIcon from '@lucide/svelte/icons/globe';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import LayoutDashboardIcon from '@lucide/svelte/icons/layout-dashboard';
	import Accordion from './Accordion.svelte';
	import { cn } from '$lib/cn';
	import MenuItem from './MenuItem.svelte';
	import { goto } from '$app/navigation';
	import AnvilIcon from '@lucide/svelte/icons/anvil';

	type Collection = {
		name: string;
		href: string;
	};

	type Props = {
		class?: string;
		onNavigate: () => void;
		collections: Collection[];
		globals: Collection[];
	};

	let { class: className, onNavigate, collections, globals }: Props = $props();

	const builtins = [
		{ name: 'Media', href: '/media' },
		{ name: 'Users', href: '/users' },
		{ name: 'Sync Status', href: '/sync' }
	];

	async function logout() {
		await fetch(`/logout`, { method: 'POST' });
		goto('/login');
	}
</script>

<div class={cn('lg:inset-y-0 lg:min-w-64', className)}>
	<div class="flex flex-1 flex-col border-r border-gray-200 bg-white">
		<nav class="mt-4 flex-1 space-y-1 overflow-auto px-2">
			<MenuItem href="/" class="font-medium" {onNavigate}>
				<LayoutDashboardIcon class="mr-3 h-5 w-5 flex-shrink-0" />
				<span class="flex-1">Dashboard</span>
			</MenuItem>

			<Accordion
				text="Collections"
				emptyText="No Collections"
				icon={DatabaseIcon}
				items={collections}
				{onNavigate}
			/>
			<Accordion
				text="Globals"
				emptyText="No Globals"
				icon={GlobeIcon}
				items={globals}
				{onNavigate}
			/>
			<Accordion text="Builtins" icon={AnvilIcon} items={builtins} {onNavigate} />
		</nav>

		<div class="p-2">
			<button
				onclick={logout}
				type="button"
				class="group flex w-full items-center rounded-md px-2 py-2 text-sm text-gray-600 hover:bg-gray-50 hover:text-gray-900"
			>
				<LogOutIcon class="mr-3 h-5 w-5 flex-shrink-0" />
				Logout
			</button>
		</div>
	</div>
</div>
