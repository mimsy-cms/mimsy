<script lang="ts">
	import Searchbar from '$lib/components/admin/Searchbar.svelte';

	let { data } = $props();

	const builtins = [
		{ name: 'Media', href: '/media' },
		{ name: 'Users', href: '/users' }
	];
</script>

<div class="flex flex-col gap-6">
	<h2 class="text-4xl font-medium">Collections</h2>

	<Searchbar id="collection-search" name="collection-search" class="max-w-md" />

	<ol class="flex min-h-32 gap-6">
		{#each data.collections as collection}
			<li class="contents">
				<a
					class="min-w-64 rounded-md border border-gray-200 bg-white px-3 py-2 hover:bg-gray-50"
					href={`/collections/${collection.slug}`}
				>
					<span class="text-xl font-medium">{collection.name}</span>
				</a>
			</li>
		{/each}
	</ol>

	<h2 class="text-4xl font-medium">Globals</h2>

	<Searchbar id="global-search" name="global-search" class="max-w-md" />

	<ol class="flex min-h-32 gap-6">
		{#each data.globals as global}
			<li class="contents">
				<a
					class="min-w-64 rounded-md border border-gray-200 bg-white px-3 py-2 hover:bg-gray-50"
					href={`/collections/${global.slug}`}
				>
					<span class="text-xl font-medium">{global.name}</span>
				</a>
			</li>
		{/each}
	</ol>

	<Searchbar id="builtin-search" name="builtin-search" class="max-w-md" />

	<ol class="flex min-h-32 gap-6">
		{#each builtins as builtin}
			<li class="contents">
				<a
					class="min-w-64 rounded-md border border-gray-200 bg-white px-3 py-2 hover:bg-gray-50"
					href={builtin.href}
				>
					<span class="text-xl font-medium">{builtin.name}</span>
				</a>
			</li>
		{/each}
	</ol>
</div>
