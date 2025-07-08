<script lang="ts">
	import Sidebar from '$lib/components/admin/Sidebar.svelte';

	let { children } = $props();
</script>

<div class="min-h-screen">
	<Sidebar />

	<main class="py-6">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			{@render children()}
		</div>
	</main>
</div>
