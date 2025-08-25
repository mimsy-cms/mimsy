<script lang="ts">
	import ResourceForm from '$lib/components/admin/ResourceForm.svelte';
	import { superForm } from 'sveltekit-superforms';

	const { data } = $props();

	const { form, message, enhance } = superForm(data.form);
</script>

<form method="POST" use:enhance>
	<ResourceForm {form} definition={data.definition} slugEditable={false} />
</form>
