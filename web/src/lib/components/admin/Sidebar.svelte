<script lang="ts">
	import DatabaseIcon from '@lucide/svelte/icons/database';
	import GlobeIcon from '@lucide/svelte/icons/globe';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import LayoutDashboardIcon from '@lucide/svelte/icons/layout-dashboard';
	import Accordion from './Accordion.svelte';
	import { cn } from '$lib/cn';
	import MenuItem from './MenuItem.svelte';
	import { goto } from '$app/navigation';

	type Collection = {
		name: string;
		slug: string;
	};

	type Props = {
		class?: string;
		onNavigate: () => void;
		collections: Collection[];
	};

	let { class: className, onNavigate, collections }: Props = $props();

	const collectionLinks = collections.map((collection) => ({
		name: collection.name,
		href: `/collections/${collection.slug}`
	}));

	const globals = [
		{ name: 'Info', href: '/globals/info' },
		{ name: 'Services', href: '/globals/services' },
		{ name: 'Footer', href: '/globals/footer' }
	];

	async function logout() {
		const res = await fetch(`/logout`, {
			method: 'POST'
		});

		if (res.ok) {
			goto('/login');
		} else {
			const errorText = await res.text();
			alert('Logout failed: ' + errorText);
		}
	}
</script>

<div class={cn('lg:inset-y-0 lg:min-w-64', className)}>
	<div class="flex flex-1 flex-col border-r border-gray-200 bg-white">
		<nav class="mt-4 flex-1 space-y-1 overflow-auto px-2">
			<MenuItem href="/" class="font-medium" {onNavigate}>
				<LayoutDashboardIcon class="mr-3 h-5 w-5 flex-shrink-0" />
				<span class="flex-1">Dashboard</span>
			</MenuItem>

			<Accordion text="Collections" icon={DatabaseIcon} items={collectionLinks} {onNavigate} />

			<Accordion text="Globals" icon={GlobeIcon} items={globals} {onNavigate} />
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
