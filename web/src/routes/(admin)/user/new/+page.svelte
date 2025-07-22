<script lang="ts">
  import { superForm } from 'sveltekit-superforms';
  import Input from '$lib/components/Input.svelte';
  import PasswordInput from '$lib/components/PasswordInput.svelte';
  import Error from '$lib/components/Error.svelte';
  import SaveIcon from '@lucide/svelte/icons/save';
  import UndoIcon from '@lucide/svelte/icons/undo-2';

  export let data;

  const { form, errors, enhance, message, reset } = superForm(data.form, {
    resetForm: true
  });

  function handleCancel() {
    reset();
  }
</script>

<div class="flex flex-col gap-6 p-6"
     style="max-width: calc(100vw - 280px - 48px);">
  <h1 class="text-4xl font-medium">Create user</h1>

  <div class="flex justify-end gap-2">
    <button
      type="button"
      on:click={handleCancel}
      class="flex items-center border border-gray-300 text-gray-700 px-2 py-1 rounded-md hover:bg-gray-100"
    >
      <UndoIcon class="mr-3 h-5 w-5 flex-shrink-0" />
      Cancel
    </button>
    <button
      type="submit"
      form="createUserForm"
      class="btn flex items-center"
    >
      <SaveIcon class="mr-3 h-5 w-5 flex-shrink-0" />
      Save
    </button>
  </div>

  <div class="rounded-md border border-gray-300 bg-white p-6 w-full">
    <form id="createUserForm" class="flex flex-col gap-4" method="POST" use:enhance>
      <div class="flex flex-col gap-2">
        <label for="email">Email</label>
        <Input
          type="email"
          id="email"
          name="email"
          bind:value={$form.email}
          error={!!$errors.email}
        />
        <Error>{$errors.email}</Error>
      </div>

      <div class="flex flex-col gap-2">
        <label for="password">Temporary password</label>
        <PasswordInput
          id="password"
          name="password"
          bind:value={$form.password}
          error={!!$errors.password}
        />
        <Error>{$errors.password}</Error>
      </div>

      <label class="flex items-center gap-2">
        <input type="checkbox" name="isAdmin" bind:checked={$form.isAdmin} />
        User management
      </label>

      <input type="hidden" name="session" value={data.session} />
    </form>

    {#if $message}
      <p class="text-red-600 mt-4 text-sm">{$message}</p>
    {/if}
  </div>
</div>
