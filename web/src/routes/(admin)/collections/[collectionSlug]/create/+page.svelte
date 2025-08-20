<script lang="ts">
  import { enhance } from '$app/forms';
  import ResourceForm from '$lib/components/admin/ResourceForm.svelte';
  export let data: import('../create/$types.js').PageData;
  export let form: any;
</script>

{#if form?.error}
    <div class="rounded border border-red-400 bg-red-100 px-4 py-3 mb-4 text-red-700">
        {form.error}
    </div>
{/if}

<form method="POST" use:enhance>
  <ResourceForm {data} slugEditable={true} mode="create" />
</form>