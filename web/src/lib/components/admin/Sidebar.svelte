<script lang="ts">
	import DatabaseIcon from '@lucide/svelte/icons/database';
	import GlobeIcon from '@lucide/svelte/icons/globe';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import LayoutDashboardIcon from '@lucide/svelte/icons/layout-dashboard';
	import Accordion from './Accordion.svelte';
	import { cn } from '$lib/cn';
	import { page } from '$app/state';

	const collections = [
		{ name: 'Posts', href: '/collections/posts' },
		{ name: 'Pages', href: '/collections/pages' },
		{ name: 'Events', href: '/collections/events' }
	];

	const globals = [
		{ name: 'Info', href: '/globals/info' },
		{ name: 'Services', href: '/globals/services' },
		{ name: 'Footer', href: '/globals/footer' }
	];
</script>

<div class="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col">
	<div class="flex flex-grow flex-col overflow-y-auto border-r border-gray-200 bg-white">
		<nav class="mt-4 flex-1 space-y-1 px-2">
			<a
				href="/"
				class={cn(
					'group flex w-full items-center rounded-md py-2 pl-2 pr-1 text-left text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900',
					{
						'bg-blue-200/80 text-blue-800': page.url.pathname === '/'
					}
				)}
			>
				<LayoutDashboardIcon class="mr-3 h-5 w-5 flex-shrink-0" />
				<span class="flex-1">Dashboard</span>
			</a>

			<Accordion text="Collections" icon={DatabaseIcon} items={collections} />

			<Accordion text="Globals" icon={GlobeIcon} items={globals} />
		</nav>

		<div class="p-2">
			<a
				href="/logout"
				data-sveltekit-preload-data="off"
				class="group flex items-center rounded-md px-2 py-2 text-sm text-gray-600 hover:bg-gray-50 hover:text-gray-900"
			>
				<LogOutIcon class="mr-3 h-5 w-5 flex-shrink-0" />
				Logout
			</a>
		</div>
	</div>
</div>
