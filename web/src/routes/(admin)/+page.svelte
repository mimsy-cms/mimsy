<script lang="ts">
	import Searchbar from '$lib/components/admin/Searchbar.svelte';

	let { data } = $props();

	let collectionSearch = $state('');
	let globalSearch = $state('');

	let collections = $state(data.collections);
	let globals = $state(data.globals);

	async function fetchCollections() {
		const response = await fetch(`/api/v1/collections?q=${collectionSearch}`);
		collections = await response.json();
	}

	async function fetchGlobals() {
		const response = await fetch(`/api/v1/collections/globals?q=${globalSearch}`);
		globals = await response.json();
	}

	$effect(() => {
		collectionSearch;
		fetchCollections();
	});

	$effect(() => {
		globalSearch;
		fetchGlobals();
	});
</script>

<div class="flex flex-col gap-6">
	<h2 class="text-4xl font-medium">Collections</h2>

	<Searchbar
		id="collection-search"
		name="collection-search"
		class="max-w-md"
		bind:value={collectionSearch}
	/>

	<ol class="flex min-h-32 gap-6">
		{#each collections as collection}
			<li class="contents">
				<a
					class="min-w-64 rounded-md border border-gray-200 bg-white px-3 py-2 hover:bg-gray-50"
					href={`/collections/${collection.slug}`}
				>
					<span class="text-xl font-medium">{collection.name}</span>
				</a>
			</li>
		{:else}
			<li class="flex items-center w-full text-gray-500 p-4">
				<span>No collections found</span>
			</li>
		{/each}
	</ol>

	<h2 class="text-4xl font-medium">Globals</h2>

	<Searchbar id="global-search" name="global-search" class="max-w-md" bind:value={globalSearch} />

	<ol class="flex min-h-32 gap-6">
		{#each globals as global}
			<li class="contents">
				<a
					class="min-w-64 rounded-md border border-gray-200 bg-white px-3 py-2 hover:bg-gray-50"
					href={`/collections/${global.slug}`}
				>
					<span class="text-xl font-medium">{global.name}</span>
				</a>
			</li>
		{:else}
			<li class="flex items-center w-full text-gray-500 p-4">
				<span>No globals found</span>
			</li>
		{/each}
	</ol>
</div>
