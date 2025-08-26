<script lang="ts">
	import ResourceForm from '$lib/components/admin/ResourceForm.svelte';
	import FlashMessage from '$lib/components/FlashMessage.svelte';
	import { superForm } from 'sveltekit-superforms';
	import { getMessageType } from '$lib/utils/messageTypes';
	import { getFlash } from 'sveltekit-flash-message';
	import { page } from '$app/stores';

	const { data } = $props();
	const { form, message, enhance, submitting } = superForm(data.form);
	const flash = getFlash(page);
</script>

{#if $flash}
	<FlashMessage 
		message={$flash.message} 
		type={$flash.type} 
	/>
{/if}

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
		resource={data.resource}
		createdBy={data.createdBy}
		updatedBy={data.updatedBy}
		slugEditable={false}
		submitting={$submitting}
	/>
</form>