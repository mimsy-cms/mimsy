<script lang="ts">
	import ResourceForm from '$lib/components/admin/ResourceForm.svelte';
	import { superForm } from 'sveltekit-superforms';

	const { data } = $props();

	const { form, message, enhance } = superForm(data.form);
</script>

{#if $message}
	<div class="mb-4 rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700">
		{$message}
	</div>
{/if}

<form method="POST" use:enhance>
	<ResourceForm {form} definition={data.definition} slugEditable={true} />
</form>
