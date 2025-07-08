<script lang="ts">
	import Searchbar from '$lib/components/admin/Searchbar.svelte';

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

<div class="flex flex-col gap-6">
	<h2 class="text-4xl font-medium">Collections</h2>

	<Searchbar id="collection-search" name="collection-search" class="max-w-md" />

	<ol class="flex min-h-32 gap-6">
		{#each collections as collection}
			<li class="contents">
				<a
					class="min-w-64 rounded-md border border-gray-200 bg-white px-3 py-2 hover:bg-gray-50"
					href={collection.href}
				>
					<span class="text-xl font-medium">{collection.name}</span>
				</a>
			</li>
		{/each}
	</ol>

	<h2 class="text-4xl font-medium">Globals</h2>

	<Searchbar id="global-search" name="global-search" class="max-w-md" />

	<ol class="flex min-h-32 gap-6">
		{#each globals as global}
			<li class="contents">
				<a
					class="min-w-64 rounded-md border border-gray-200 bg-white px-3 py-2 hover:bg-gray-50"
					href={global.href}
				>
					<span class="text-xl font-medium">{global.name}</span>
				</a>
			</li>
		{/each}
	</ol>
</div>
