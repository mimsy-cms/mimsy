<script lang="ts">
	import { cn } from '$lib/cn';
	import MobileMenu from '$lib/components/admin/MobileMenu.svelte';
	import Sidebar from '$lib/components/admin/Sidebar.svelte';

	let { children, data } = $props();

	let mobileNavOpen = $state(false);

	let collections = $derived(
		(data.collections ?? []).map((collection) => ({
			name: collection.name,
			href: `/collections/${collection.slug}`
		}))
	);

	let globals = $derived(
		(data.globals ?? []).map((global) => ({
			name: global.name,
			href: `/collections/${global.slug}`
		}))
	);
</script>

<div class="flex min-h-screen flex-col lg:flex-row">
	<MobileMenu
		class="sticky top-0 lg:hidden"
		onToggle={() => {
			mobileNavOpen = !mobileNavOpen;
		}}
	/>

	<Sidebar
		onNavigate={() => (mobileNavOpen = false)}
		{collections}
		{globals}
		class={cn('mt-13 absolute bottom-0 left-0 right-0 top-0 pt-[1px] lg:static lg:mt-0 lg:flex', {
			hidden: !mobileNavOpen
		})}
	/>

	<main class="flex flex-1 overflow-x-hidden">
		<div class="flex-1 overflow-auto px-4 py-6 sm:px-6 lg:px-8">
			{@render children()}
		</div>
	</main>
</div>
