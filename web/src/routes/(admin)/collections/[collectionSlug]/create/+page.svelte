<script lang="ts">
	import { enhance } from '$app/forms';
	import ResourceForm from '$lib/components/admin/ResourceForm.svelte';

	const { data, form } = $props();
</script>

{#if form?.error}
	<div class="mb-4 rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700">
		{form.error}
	</div>
{/if}

<form method="POST" use:enhance>
	<ResourceForm {data} slugEditable={true} />
</form>
