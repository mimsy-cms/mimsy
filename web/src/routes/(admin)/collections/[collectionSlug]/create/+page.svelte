<script lang="ts">
  import ResourceForm from '$lib/components/admin/ResourceForm.svelte';
  export let data: import('../create/$types.js').PageData;
</script>

<form method="POST" use:enhance>
  <ResourceForm {data} slugEditable={true} mode="create" />
</form>