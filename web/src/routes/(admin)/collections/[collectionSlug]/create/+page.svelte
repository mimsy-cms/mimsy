<script lang="ts">
	import ResourceForm from '$lib/components/admin/ResourceForm.svelte';
	import FlashMessage from '$lib/components/FlashMessage.svelte';
	import { superForm } from 'sveltekit-superforms';
	import { getMessageType } from '$lib/utils/messageTypes';
	
	const { data } = $props();
	const { form, message, enhance, submitting } = superForm(data.form);
</script>

{#if $message}
	<FlashMessage 
		message={$message} 
		type={getMessageType($form, $message)} 
	/>
{/if}

<form method="POST" use:enhance>
	<ResourceForm 
		{form} 
		definition={data.definition} 
		slugEditable={true}
		submitting={$submitting}
	/>
</form>