<script lang="ts">
  import Error from '$lib/components/Error.svelte';
  import PasswordInput from '$lib/components/PasswordInput.svelte';
  import { superForm } from 'sveltekit-superforms';
  import SaveIcon from '@lucide/svelte/icons/save';
  import UndoIcon from '@lucide/svelte/icons/undo-2';

  export let data;

  const { form, errors, enhance, message, reset } = superForm(data.form, {
    resetForm: true,
  });

  function handleCancel() {
    reset();
  }
</script>

<div class="flex flex-col gap-6 p-6" style="max-width: calc(100vw - 280px - 48px);">
  <h1 class="text-4xl font-medium">Reset Password</h1>

  <div class="flex justify-end gap-2">
    <button
      type="button"
      on:click={handleCancel}
      class="flex items-center border border-gray-300 text-gray-700 px-2 py-1 rounded-md hover:bg-gray-100"
    >
      <UndoIcon class="mr-3 h-5 w-5 flex-shrink-0" />
      Cancel
    </button>
    <button type="submit" form="resetPasswordForm" class="btn flex items-center">
      <SaveIcon class="mr-3 h-5 w-5 flex-shrink-0" />
      Change
    </button>
  </div>

  <div class="rounded-md border border-gray-300 bg-white p-6 w-full">
    <form
      id="resetPasswordForm"
      class="flex flex-col gap-4"
      method="POST"
      use:enhance
      novalidate
    >
      <div class="flex flex-col gap-2">
        <label for="old_password">Current Password</label>
        <PasswordInput
          id="old_password"
          name="old_password"
          bind:value={$form.old_password}
          error={!!$errors.old_password}
        />
        <Error>{$errors.old_password}</Error>
      </div>

      <div class="flex flex-col gap-2">
        <label for="new_password">New Password</label>
        <PasswordInput
          id="new_password"
          name="new_password"
          bind:value={$form.new_password}
          error={!!$errors.new_password}
        />
        <Error>{$errors.new_password}</Error>
      </div>

      <div class="flex flex-col gap-2">
        <label for="confirm_password">Confirm New Password</label>
        <PasswordInput
          id="confirm_password"
          name="confirm_password"
          bind:value={$form.confirm_password}
          error={!!$errors.confirm_password}
        />
        <Error>{$errors.confirm_password}</Error>
      </div>
    </form>

    {#if $message}
      <p class="text-red-600 mt-4 text-sm">{$message}</p>
    {/if}
  </div>
</div>
