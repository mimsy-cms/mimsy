<script lang="ts">
	import Sidebar from '$lib/components/admin/Sidebar.svelte';

	let { children } = $props();
</script>

<div class="flex min-h-screen">
	<Sidebar />

	<main class="flex-1 py-6">
		<div class="px-4 sm:px-6 lg:px-8">
			{@render children()}
		</div>
	</main>
</div>
