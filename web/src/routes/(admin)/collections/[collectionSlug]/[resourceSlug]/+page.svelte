<script lang="ts">
  import ResourceForm from '$lib/components/admin/ResourceForm.svelte';
  export let data: import('../create/$types.js').PageData;
</script>

<ResourceForm {data} slugEditable={false} />
