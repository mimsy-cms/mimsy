<script lang="ts">
	import ResourceForm from '$lib/components/admin/ResourceForm.svelte';
	const { data } = $props();
</script>

<ResourceForm {data} slugEditable={false} />
